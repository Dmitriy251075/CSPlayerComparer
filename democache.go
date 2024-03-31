package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
)

var CacheMutex sync.Mutex
var Cache []*DemoCache

type DemoCache struct {
	Name      string
	demoID    string
	Gamemode  int
	PlrsStats []*PlrStats
}

const (
	Matchmaking = iota + 1
	Wingman
	Other
)

func getGameMode(PlrStatsLength int) int {
	var mode = Other // Default other
	if PlrStatsLength >= 8 && PlrStatsLength <= 11 {
		mode = Matchmaking
	} else if PlrStatsLength >= 3 && PlrStatsLength <= 5 {
		mode = Wingman
	} else if PlrStatsLength >= 0 {
		mode = Other
	}
	return mode
}

// Use this feature after completing the demo analysis.
func createDemoCache(path string, parser demoinfocs.Parser, Plrs []*PlrStats, force bool) *DemoCache {
	dem := DemoCache{
		Name: filepath.Base(path),
	}

	demoid := createDemoID(parser)

	if !force && findDemoInMemByDemoID(demoid) != nil {
		return nil
	}

	dem.demoID = demoid
	dem.Gamemode = getGameMode(len(Plrs))
	dem.PlrsStats = append(dem.PlrsStats, Plrs...)

	CacheMutex.Lock()
	defer CacheMutex.Unlock()

	Cache = append(Cache, &dem)
	return &dem
}

const StrfmtDemoCacheEnd = "\n"
const StrfmtDemoCacheDemoID = "demoid=\t"
const StrfmtDemoCacheDemoIDEnd = StrfmtDemoCacheEnd
const StrfmtDemoCacheMode = "mode=\t"
const StrfmtDemoCacheModeEnd = StrfmtDemoCacheEnd
const StrfmtDemoCachePlr = "plr=\t"
const StrfmtDemoCachePlrEnd = StrfmtDemoCacheEnd

// Returns demoCache and true if found in mem or false if not found in mem
func loadDemoCache(demoName string) (*DemoCache, bool) {
	if *useCache {
		demMem := findDemoInMemByName(demoName)
		if demMem != nil {
			return demMem, true
		}

		data, err := os.ReadFile(*cacheDir + "/" + filepath.Base(demoName) + ".txt")
		if err != nil {
			log.Println("failed to read cache file: ", err)
			return nil, false
		}

		dem := DemoCache{
			Name: filepath.Base(demoName),
		}

		// Parse demoID
		SplitDemoID := strings.SplitN(string(data), StrfmtDemoCacheDemoID, 2)
		if len(SplitDemoID) < 2 {
			log.Println("failed to parse " + filepath.Base(demoName) + ".txt" + ": demoID not found")
			return nil, false
		}
		DemoID := strings.SplitN(SplitDemoID[1], StrfmtDemoCacheDemoIDEnd, 2)[0]
		if DemoID != "" {
			dem.demoID = DemoID
		} else {
			log.Println("failed to parse " + filepath.Base(demoName) + ".txt" + ": demoID not found")
			return nil, false
		}

		demMem = findDemoInMemByDemoID(DemoID)
		if demMem != nil {
			return demMem, true
		}

		// Parse gamemode
		SplitMode := strings.SplitN(string(data), StrfmtDemoCacheMode, 2)
		if len(SplitMode) < 2 {
			log.Println("failed to parse " + filepath.Base(demoName) + ".txt" + ": gamemode not found")
			return nil, false
		}
		Mode := strings.SplitN(SplitMode[1], StrfmtDemoCacheModeEnd, 2)[0]
		if Mode != "" {
			gamemode, err := strconv.ParseInt(Mode, 10, 32)
			if err != nil {
				log.Println("failed to parse " + filepath.Base(demoName) + ".txt" + ": ", err)
				return nil, false
			}
			dem.Gamemode = int(gamemode)
		} else {
			log.Println("failed to parse " + filepath.Base(demoName) + ".txt" + ": gamemode not found")
			return nil, false
		}

		// Parse players
		plrSplit := strings.Split(string(data), StrfmtDemoCachePlr)

		for plrI := 1; plrI < len(plrSplit); plrI++ {
			plr := strings.SplitN(plrSplit[plrI], StrfmtDemoCachePlrEnd, 2)[0]

			p := PlrStats{}
				
			IsOK := p.fromString(plr)
			if !IsOK {
				log.Println("failed to parse " + filepath.Base(demoName) + ".txt" + ": player not found: " + plr)
				return nil, false
			}

			dem.PlrsStats = append(dem.PlrsStats, &p)
		}

		CacheMutex.Lock()
		defer CacheMutex.Unlock()

		demMem = findDemoInMemByDemoID(DemoID)
		if demMem != nil {
			return demMem, true
		}

		Cache = append(Cache, &dem)
		return &dem, false
	}
	return nil, false
}

func (dem *DemoCache) saveToDisk(overwrite bool) {
	d, err :=os.ReadDir(*cacheDir + "/")
	if err != nil {
		log.Println("failed to read cache directory: ", err)
		if os.IsNotExist(err) {
			err := os.Mkdir(*cacheDir, os.ModePerm)
			if err != nil {
				log.Println("failed to create cache directory: ", err)
			}
		}
	}

	var IsExists bool = false
	for _, file := range d {
		if file.Name() == dem.Name {
			IsExists = true
			log.Println("cache file already exists: ", dem.Name)
		}
	}

	if !IsExists || overwrite {
		CacheMutex.Lock()
		defer CacheMutex.Unlock()

		file, err := os.Create(*cacheDir + "/" + dem.Name + ".txt")
		if err != nil {
			log.Println("failed to create cache file: ", err)
		}
		defer file.Close()

		str := StrfmtDemoCacheDemoID + dem.demoID + StrfmtDemoCacheDemoIDEnd
		str += StrfmtDemoCacheMode + strconv.FormatUint(uint64(dem.Gamemode), 10) + StrfmtDemoCacheModeEnd
		for _, plr := range dem.PlrsStats {
			str += StrfmtDemoCachePlr + plr.toString() + StrfmtDemoCachePlrEnd
		}

		file.WriteString(str)
	}
}

// If it returns not nil, then the demo has already been parsed.
func findDemoInMemByName(demoName string) *DemoCache {
	for _, demo := range Cache {
		if demo.Name == demoName {
			return demo
		}
	}
	return nil
}

// If it returns not nil, then the demo has already been parsed.
func findDemoInMemByDemoID(demoID string) *DemoCache {
	for _, demo := range Cache {
		if demo.demoID == demoID {
			return demo
		}
	}
	return nil
}

// Use this feature after completing the demo analysis.
// If it returns not nil, then the demo has already been parsed.
func findDemoInMemByDemoIDThroughParser(parser demoinfocs.Parser) *DemoCache {
	demoid := createDemoID(parser)

	for _, demo := range Cache {
		if demo.demoID == demoid {
			return demo
		}
	}
	return nil
}

// Use this feature after completing the demo analysis.
func createDemoID(parser demoinfocs.Parser) string {
	gs := parser.GameState()
	header := parser.Header()

	ct := gs.TeamCounterTerrorists()
	t := gs.TeamTerrorists()

	demoid := strconv.Itoa(gs.TotalRoundsPlayed()) + "-" + strconv.Itoa(gs.OvertimeCount()) + "-" + strconv.Itoa(ct.Score()) + "-" + strconv.Itoa(t.Score()) + "_" + strconv.Itoa(gs.IngameTick()) + "\t" + ct.ClanName() + "\t" + t.ClanName() + "\t-\t" + header.ServerName + "\t-\t" + header.MapName + "\t-\t" + header.PlaybackTime.Abs().String()

	return demoid
}
