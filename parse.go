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
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
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

var totalFiles uint32
var currentCompletedDemoFiles uint32
var usedDemoFiles uint32
var currentCachedDemoFiles uint32
var errorDemoFiles uint32

var wgDem sync.WaitGroup
var wgCache sync.WaitGroup

var PlrsStats []*PlrStats

var useStatsMatchmaking bool = false
var useStatsWingman bool = false
var useStatsOther bool = false

var usedDemoFileNamesMutex sync.Mutex
var usedDemoFileNames []string

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

// Returns 0 if success parsed demo, 1 if found in cache, 2 if error, 3 if canceled
func cacheParse(path string, name string) int {
	defer wgCache.Done()

	dem := findDemoInMemByName(name)
	if dem != nil {
		log.Println("demo already parsed: ", path)
		atomic.AddUint32(&errorDemoFiles, 1)
		return 1
	}

	for currentDemosParseing >= maxDemosParseing {
		time.Sleep(time.Millisecond * 100)
	}
	atomic.AddInt32(&currentDemosParseing, 1)

	if isCanceling {
		atomic.AddInt32(&currentDemosParseing, -1)
		atomic.AddUint32(&errorDemoFiles, 1)
		return 3
	}

	cache, IsDupe := loadDemoCache(name)
	if cache != nil {
		var AllPlrsStats []*PlrStats
		AllPlrsStats = append(AllPlrsStats, cache.PlrsStats...)

		addTargetsStats(AllPlrsStats[:], cache.Gamemode, cache.Name)
			
		log.Println("file cached: ", name)
		atomic.AddInt32(&currentDemosParseing, -1)
		atomic.AddUint32(&currentCachedDemoFiles, 1)
		return 0
	} else if IsDupe {
		log.Println("demo already parsed: ", path)
		atomic.AddInt32(&currentDemosParseing, -1)
		atomic.AddUint32(&errorDemoFiles, 1)
		return 1
	}
	atomic.AddInt32(&currentDemosParseing, -1)
	return 2
}

func demPrepare(path string, name string) {
	defer wgDem.Done()

	wgCache.Add(1)
	cacheResult := cacheParse(path, name)
	if cacheResult == 0 {
		return
	} else if cacheResult == 1 {
		return
	} else if cacheResult == 3 {
		return
	}

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
	} else {
		atomic.AddUint32(&errorDemoFiles, 1)
	}
}

// Returns true if success
func demParse(reader io.Reader, path string) bool {
	for currentDemosParseing >= maxDemosParseing {
		time.Sleep(time.Millisecond * 100)
	}
	atomic.AddInt32(&currentDemosParseing, 1)

	if isCanceling {
		atomic.AddUint32(&errorDemoFiles, 1)
		atomic.AddInt32(&currentDemosParseing, -1)
		return false
	}

	p := demoinfocs.NewParser(reader)
	defer p.Close()

	var AllPlrsStats []*PlrStats
	var AllPlrsStatsMap = make(map[uint64]*PlrStats)
	var PlrsPing = make(map[uint64]uint64)
	var PlrsPingChecks = make(map[uint64]uint64)

	RegAllPlrsStats := func() {
		gs := p.GameState()

		SetStats := func (plr *common.Player)  {
			AllPlrsStatsMap[plr.SteamID64].setStats(plr)
			PlrsPing[plr.SteamID64] += uint64(plr.Ping())
			PlrsPingChecks[plr.SteamID64] += 1
		}

		for _, plr := range gs.Participants().Playing() {
			if AllPlrsStatsMap[plr.SteamID64] != nil {
				SetStats(plr)
				continue
			}

			pstats := PlrStats{SteamID64: plr.SteamID64}
			AllPlrsStats = append(AllPlrsStats, &pstats)
			AllPlrsStatsMap[plr.SteamID64] = &pstats
			SetStats(plr)
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

	p.RegisterEventHandler(func(e events.RoundEnd) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.PlayerConnect) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.PlayerDisconnected) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.PlayerFlashed) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.PlayerJump) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.BotConnect) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.BotTakenOver) {
		RegAllPlrsStats()
	})

	p.RegisterEventHandler(func(e events.Kill) {
		RegAllPlrsStats()

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
	demCache := findDemoInMemByDemoIDThroughParser(p)
	if demCache != nil {
		log.Println("demo already parsed: ", path)
		atomic.AddUint32(&errorDemoFiles, 1)
		atomic.AddInt32(&currentDemosParseing, -1)
		return false
	}

	// Set Avg Ping
	for _, plrstat := range AllPlrsStats {
		if PlrsPing[plrstat.SteamID64] == 0 {
			continue
		}
		plrstat.statsPing = uint64(PlrsPing[plrstat.SteamID64] / PlrsPingChecks[plrstat.SteamID64])
	}

	// Create demo cache
	cache := createDemoCache(path, p, AllPlrsStats[:])
	cache.saveToDisk(true)

	addTargetsStats(AllPlrsStats[:], 0, filepath.Base(path))
	
	atomic.AddInt32(&currentDemosParseing, -1)
	return true
}

func addTargetsStats(AllPlrsStats []*PlrStats, gamemode int, nameDemo string) {
	foundPlrs := ""
	for _, plr := range AllPlrsStats {
		foundPlrs += "\t" + plr.Name
	}
	log.Println("found players:", foundPlrs)

	var PlrsResults = make(map[uint64]*PlrStats)
	for _, plrstat := range PlrsStats {
		plr := PlrStats{SteamID64: plrstat.SteamID64}
		PlrsResults[plrstat.SteamID64] = &plr
	}

	Merge := func() {
		if len(PlrsStats) == len(PlrsResults) {
			for _, plr := range PlrsStats {
				plr.appendStatsFromPlrStats(PlrsResults[plr.SteamID64])
			}
			atomic.AddUint32(&usedDemoFiles, 1)
	
			usedDemoFileNamesMutex.Lock()
			usedDemoFileNames = append(usedDemoFileNames, nameDemo)
			usedDemoFileNamesMutex.Unlock()
		}
	}

	var countTargetInDemo = 0
	for _, plrstat := range AllPlrsStats {
		if PlrsResults[plrstat.SteamID64] != nil {
			countTargetInDemo++
			PlrsResults[plrstat.SteamID64].appendStatsFromPlrStats(plrstat)
		}
	}

	if countTargetInDemo == len(PlrsResults) {
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
	} else {
		return
	}
}
