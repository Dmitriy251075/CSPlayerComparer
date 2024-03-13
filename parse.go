package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"

	"github.com/mholt/archiver/v3"
)

var isCanceling bool = false

const maxDemosUnziping uint32 = 5
var currentDemosUnziping uint32

var totalDemoFiles uint32
var currentCompletedDemoFiles uint32
var usedDemoFiles uint32
var errorDemoFiles uint32

var demofiles []string

var PlrsStats []*PlrStats

func uncompress(path string, name string) string {
	for currentDemosUnziping >= maxDemosUnziping {
		time.Sleep(time.Millisecond * 500)
	}
	atomic.AddUint32(&currentDemosUnziping, 1)

	tmpname := createTmpName()

	if isCanceling {
		atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
		return "isCanceling"
	}

	err := archiver.DecompressFile(path, filepath.Join(os.TempDir(), name+tmpname))
	if err != nil {
		log.Println("failed to decompress file: ", err)
		atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
		return ""
	}

	atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
	return filepath.Join(os.TempDir(), name+tmpname)
}

func demPrepare(path string, name string) {
	defer wgDem.Done()

	ext := filepath.Ext(path)
	if ext == ".bz2" || ext == ".gz" {
		if isCanceling {
			return
		}

		decompressed := uncompress(path, name)
		if decompressed == "" {
			atomic.AddUint32(&errorDemoFiles, 1)
			return
		} else if decompressed == "isCanceling" {
			return
		}
		log.Println("file decompressed: ", name)

		demParse(decompressed)

		log.Println("file parsed: ", name)
		atomic.AddUint32(&currentCompletedDemoFiles, 1)

		os.Remove(decompressed)
	} else if ext == ".dem" {
		if isCanceling {
			return
		}

		demParse(path)
		log.Println("file parsed: ", name)
		atomic.AddUint32(&currentCompletedDemoFiles, 1)
	}
}

func demParse(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("failed to open demo file: ", err)
		atomic.AddUint32(&errorDemoFiles, 1)
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

	p.RegisterEventHandler(func(e events.Kill) {
		for _, plrstat := range PlrsStats {
			plrstat.appendStatKills(&e)
		}
	})

	p.RegisterEventHandler(func(e events.AnnouncementWinPanelMatch) {
		gs := p.GameState()

		ct := gs.TeamCounterTerrorists()
		t := gs.TeamTerrorists()

		// Create a demo ID to remove duplicates from stats
		demoid := strconv.Itoa(gs.TotalRoundsPlayed()) + "." + ct.ClanName() + "." + t.ClanName() + "-" + strconv.Itoa(gs.IngameTick()) + "-" + strconv.Itoa(ct.Score()) + "." + strconv.Itoa(t.Score())

		demofound := false
		for _, demo := range demofiles {
			if demo == demoid {
				demofound = true
			}
		}
		if demofound {
			atomic.AddUint32(&errorDemoFiles, 1)
			return
		}
		demofiles = append(demofiles, demoid)

		var Plrs []*common.Player

		// CT
		for _, plr := range ct.Members() {
			for _, plrstat := range PlrsStats {
				if plr.SteamID64 == plrstat.SteamID64 {
					Plrs = append(Plrs, plr)
				}
			}
		}

		// T
		for _, plr := range t.Members() {
			for _, plrstat := range PlrsStats {
				if plr.SteamID64 == plrstat.SteamID64 {
					Plrs = append(Plrs, plr)
				}
			}
		}

		foundPlrs := ""
		for _, plr := range Plrs {
			foundPlrs += " " + plr.Name
		}
		log.Println("found players:", foundPlrs)

		if len(Plrs) == len(PlrsStats) {
			for _, plr := range Plrs {
				for _, plrstat := range PlrsStats {
					if plr.SteamID64 == plrstat.SteamID64 {
						plrstat.appendStats(plr)
					}
				}
			}
			atomic.AddUint32(&usedDemoFiles, 1)
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		log.Println("failed to parse demo: ", err)
		atomic.AddUint32(&errorDemoFiles, 1)
	}
}

func createTmpName() string {
	tmpf, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		log.Println("failed to create temp file: ", err)
		return ""
	}
	tmpname := filepath.Base(tmpf.Name())
	tmpf.Close()
	err = os.Remove(tmpf.Name())
	if err != nil {
		log.Println("failed to delete temp file: ", err)
		return ""
	}
	return tmpname
}
