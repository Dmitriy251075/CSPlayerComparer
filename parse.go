package main

import (
	"bytes"
	"io"
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

const maxDemosUnziping uint32 = 4 // 1 archive for approximately ~1 GB of RAM.
var currentDemosUnziping uint32

var totalDemoFiles uint32
var currentCompletedDemoFiles uint32
var usedDemoFiles uint32
var errorDemoFiles uint32

var demofiles []string

var PlrsStats []*PlrStats

func uncompress(path string) *bytes.Buffer {
	for currentDemosUnziping >= maxDemosUnziping {
		time.Sleep(time.Millisecond * 500)
	}
	atomic.AddUint32(&currentDemosUnziping, 1)

	if isCanceling {
		atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
		return nil
	}

	// Open file
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		log.Println("failed to open file: ", err)
		atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
		return nil
	}
	defer file.Close()

	// Create interface
	iface, err := archiver.ByExtension(path)
	if err != nil {
		log.Println("failed to create interface: ", err)
		atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
		return nil
	}
	decomp := iface.(archiver.Decompressor)

	// Create buffer and decompress to buffer
	buf := new(bytes.Buffer)
	err = decomp.Decompress(file, buf)
	if err != nil {
		log.Println("failed to decompress file: ", err)
		atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
		return nil
	}

	atomic.StoreUint32(&currentDemosUnziping, currentDemosUnziping - 1)
	return buf
}

func demPrepare(path string, name string) {
	defer wgDem.Done()

	ext := filepath.Ext(path)
	if ext == ".bz2" || ext == ".gz" {
		if isCanceling {
			return
		}

		// Decompress file to buffer
		decompressed := uncompress(path)
		if decompressed == nil {
			atomic.AddUint32(&errorDemoFiles, 1)
			return
		}
		log.Println("file decompressed: ", name)

		// Parse buffer
		demParse(decompressed)

		log.Println("file parsed: ", name)
		atomic.AddUint32(&currentCompletedDemoFiles, 1)
	} else if ext == ".dem" {
		if isCanceling {
			return
		}

		// Open file
		file, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			log.Println("failed to open file: ", err)
			return
		}
		defer file.Close()

		// Parse file
		demParse(file)

		log.Println("file parsed: ", name)
		atomic.AddUint32(&currentCompletedDemoFiles, 1)
	}
}

func demParse(reader io.Reader) {
	p := demoinfocs.NewParser(reader)
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
	err := p.ParseToEnd()
	if err != nil {
		log.Println("failed to parse demo: ", err)
		atomic.AddUint32(&errorDemoFiles, 1)
	}
}
