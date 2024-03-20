package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
	"github.com/mholt/archiver/v3"
)

var isCanceling bool = false

const maxDemosUnziping int32 = 6
var currentDemosUnziping int32

const maxDemosUnziped int32 = 6
var currentDemosUnziped int32

const maxDemosParseing int32 = maxDemosUnziped + 2
var currentDemosParseing int32

var totalDemoFiles uint32
var currentCompletedDemoFiles uint32
var usedDemoFiles uint32
var currentCachedDemoFiles uint32
var errorDemoFiles uint32

var wgDem sync.WaitGroup

var PlrsStats []*PlrStats

var useStatsMatchmaking bool = false
var useStatsWingman bool = false
var useStatsOther bool = false

func uncompress(path string) *bytes.Buffer {
	for currentDemosUnziping >= maxDemosUnziping {
		time.Sleep(time.Millisecond * 100)
	}
	atomic.AddInt32(&currentDemosUnziping, 1)

	if isCanceling {
		atomic.AddInt32(&currentDemosUnziping, -1)
		return nil
	}

	// Open file
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		log.Println("failed to open file: ", err)
		atomic.AddInt32(&currentDemosUnziping, -1)
		return nil
	}
	defer file.Close()

	// Create interface
	iface, err := archiver.ByExtension(path)
	if err != nil {
		log.Println("failed to create interface: ", err)
		atomic.AddInt32(&currentDemosUnziping, -1)
		return nil
	}
	decomp := iface.(archiver.Decompressor)

	// Create buffer and decompress to buffer
	buf := new(bytes.Buffer)
	err = decomp.Decompress(file, buf)
	if err != nil {
		log.Println("failed to decompress file: ", err)
		atomic.AddInt32(&currentDemosUnziping, -1)
		return nil
	}

	atomic.AddInt32(&currentDemosUnziping, -1)
	atomic.AddInt32(&currentDemosUnziped, 1)
	return buf
}

func demPrepare(path string, name string) {
	defer wgDem.Done()

	dem := findDemoInMemByName(name)
	if dem != nil {
		log.Println("demo already parsed: ", path)
		atomic.AddUint32(&errorDemoFiles, 1)
		return
	}

	for currentDemosParseing >= maxDemosParseing {
		time.Sleep(time.Millisecond * 100)
	}
	atomic.AddInt32(&currentDemosParseing, 1)

	if isCanceling {
		atomic.AddInt32(&currentDemosParseing, -1)
		atomic.AddUint32(&errorDemoFiles, 1)
		return
	}

	cache := loadDemoCache(name)
	if cache != nil {
		var AllPlrsStats []*PlrStats
		AllPlrsStats = append(AllPlrsStats, cache.PlrsStats...)

		addTargetsStats(AllPlrsStats[:], cache.Gamemode)
			
		log.Println("file cached: ", name)
		atomic.AddInt32(&currentDemosParseing, -1)
		atomic.AddUint32(&currentCachedDemoFiles, 1)
		return
	}
	atomic.AddInt32(&currentDemosParseing, -1)

	ext := filepath.Ext(path)
	if ext == ".bz2" || ext == ".gz" {
		for currentDemosUnziped >= maxDemosUnziped {
			time.Sleep(time.Millisecond * 100)
		}

		// Decompress file to buffer
		decompressed := uncompress(path)
		if decompressed == nil {
			atomic.AddUint32(&errorDemoFiles, 1)
			return
		}
		log.Println("file decompressed: ", name)

		// Parse buffer
		IsOK := demParse(decompressed, path)
		if !IsOK {
			atomic.AddInt32(&currentDemosUnziped, -1)
			return
		}
		atomic.AddInt32(&currentDemosUnziped, -1)

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
		IsOK := demParse(file, path)
		if !IsOK {
			return
		}

		log.Println("file parsed: ", name)
		atomic.AddUint32(&currentCompletedDemoFiles, 1)
	}
}

func demParse(reader io.Reader, path string) bool {
	for currentDemosParseing >= maxDemosParseing {
		time.Sleep(time.Millisecond * 100)
	}
	atomic.AddInt32(&currentDemosParseing, 1)

	p := demoinfocs.NewParser(reader)
	defer p.Close()

	var AllPlrsStats []*PlrStats

	RegAllPlrsStats := func() {
		gs := p.GameState()

		for _, plr := range gs.Participants().Playing() {
			found := false
			for _, plrstat := range AllPlrsStats {
				if plr.SteamID64 == plrstat.SteamID64 {
					found = true
					plrstat.setStats(plr)
				}
			}
			if (!found) {
				pstats := PlrStats{SteamID64: plr.SteamID64}
				pstats.setStats(plr)
				AllPlrsStats = append(AllPlrsStats, &pstats)
			}
		}
	}

	p.RegisterEventHandler(func(e events.AnnouncementMatchStarted) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.AnnouncementLastRoundHalf) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.AnnouncementFinalRound) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.AnnouncementWinPanelMatch) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.RoundStart) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.Kill) {
		for _, plrstat := range AllPlrsStats {
			plrstat.appendStatKills(&e)
		}
	})

	p.RegisterEventHandler(func(e events.OtherDeath) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.PlayerHurt) {
		RegAllPlrsStats()
	})

	// Parse to end
	err := p.ParseToEnd()
	if err != nil {
		log.Println("failed to parse demo: ", err)
		atomic.AddUint32(&errorDemoFiles, 1)
		atomic.AddInt32(&currentDemosParseing, -1)
		return false
	}

	// Check demoid
	demCache := findDemoInMemByDemoID(p)
	if demCache != nil {
		log.Println("demo already parsed: ", path)
		atomic.AddUint32(&errorDemoFiles, 1)
		atomic.AddInt32(&currentDemosParseing, -1)
		return false
	}

	// Create demo cache
	cache := createDemoCache(path, p, AllPlrsStats[:])
	cache.saveToDisk(true)

	addTargetsStats(AllPlrsStats[:], 0)
	
	atomic.AddInt32(&currentDemosParseing, -1)
	return true
}

func addTargetsStats(AllPlrsStats []*PlrStats, gamemode int) {
	foundPlrs := ""
	for _, plr := range AllPlrsStats {
		foundPlrs += "\t" + plr.Name
	}
	log.Println("found players:", foundPlrs)

	var PlrsResults []*PlrStats
	for _, plrstat := range PlrsStats {
		plr := PlrStats{SteamID64: plrstat.SteamID64}
		PlrsResults = append(PlrsResults, &plr)
	}

	var countTargetInDemo = 0
	for _, plrstat := range AllPlrsStats {
		for _, plr := range PlrsResults {
			if plr.SteamID64 == plrstat.SteamID64 {
				countTargetInDemo++
			}
		}
	}

	if countTargetInDemo == len(PlrsResults) {
		for _, plr := range AllPlrsStats {
			for _, plrstat := range PlrsResults {
				if plr.SteamID64 == plrstat.SteamID64 {
					plrstat.appendStatsFromPlrStats(plr)
				}
			}
		}
	} else {
		return
	}

	Merge := func() {
		for _, plr := range PlrsResults {
			for _, plrstat := range PlrsStats {
				if plr.SteamID64 == plrstat.SteamID64 {
					plrstat.appendStatsFromPlrStats(plr)
				}
			}
		}
		atomic.AddUint32(&usedDemoFiles, 1)
	}

	if gamemode == 0 {
		gamemode = getGameMode(len(AllPlrsStats))
	}

	if useStatsMatchmaking && gamemode == Matchmaking {
		Merge()
	} else if useStatsWingman && gamemode == Wingman {
		Merge()
	} else if useStatsOther && gamemode == Other {
		Merge()
	}
}
