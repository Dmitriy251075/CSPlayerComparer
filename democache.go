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
	PlrsStats []*PlrStats
}

// Use this feature after completing the demo analysis.
func createDemoCache(path string, parser demoinfocs.Parser, Plrs []*PlrStats) *DemoCache {
	dem := DemoCache{
		Name: filepath.Base(path),
	}

	demoid := createDemoID(parser)

	for _, demo := range Cache {
		if demo.demoID == demoid {
			return nil
		}
	}
	CacheMutex.Lock()
	defer CacheMutex.Unlock()
	
	dem.demoID = demoid
	Cache = append(Cache, &dem)

	dem.PlrsStats = append(dem.PlrsStats, Plrs...)

	return &dem
}

const StrfmtDemoCacheDemoID = "demoid=\t"
const StrfmtDemoCacheDemoIDEnd = "\n"
const StrfmtDemoCachePlr = "plr=\t"
const StrfmtDemoCachePlrEnd = "\n"

func loadDemoCache(demoName string) *DemoCache {
	if *useCache {
		dem := findDemoInMemByName(demoName)
		if dem != nil {
			return dem
		}

		d, err :=os.ReadDir("cache/")
		if err != nil {
			log.Println("failed to read cache directory: ", err)
			if os.IsNotExist(err) {
				err := os.Mkdir("cache", os.ModePerm)
				if err != nil {
					log.Println("failed to create cache directory: ", err)
				}
			}
		}

		var IsExists bool = false
		for _, file := range d {
			if file.Name() == filepath.Base(demoName) + ".txt" {
				IsExists = true
			}
		}

		if IsExists {
			CacheMutex.Lock()
			defer CacheMutex.Unlock()

			data, err := os.ReadFile("cache/" + filepath.Base(demoName) + ".txt")
			if err != nil {
				log.Println("failed to read cache file: ", err)
			}

			dem := DemoCache{
				Name: filepath.Base(demoName),
			}

			
			demoidAndPlrs := strings.SplitN(strings.SplitN(string(data), StrfmtDemoCacheDemoID, 2)[1], StrfmtDemoCacheDemoIDEnd, 2)

			dem.demoID = demoidAndPlrs[0]

			plrSplit := strings.Split(demoidAndPlrs[1], StrfmtDemoCachePlr)

			for plrI := 1; plrI < len(plrSplit); plrI++ {
				plr := strings.Split(plrSplit[plrI], StrfmtDemoCachePlrEnd)[0]

				p := PlrStats{}

				p.fromString(plr)

				dem.PlrsStats = append(dem.PlrsStats, &p)
			}

			Cache = append(Cache, &dem)
			return &dem
		}
		return nil
	}
	return nil
}

func (dem *DemoCache) saveToDisk() {
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

	if !IsExists {
		CacheMutex.Lock()
		defer CacheMutex.Unlock()

		file, err := os.Create(*cacheDir + "/" + dem.Name + ".txt")
		if err != nil {
			log.Println("failed to create cache file: ", err)
		}
		defer file.Close()

		str := StrfmtDemoCacheDemoID + dem.demoID + StrfmtDemoCacheDemoIDEnd
		for _, plr := range dem.PlrsStats {
			str += StrfmtDemoCachePlr + plr.toString() + StrfmtDemoCachePlrEnd
		}

		file.WriteString(str)
	}
}

// Use this feature after completing the demo analysis.
// If it returns not nil, then the demo has already been parsed.
func findDemoInMemByName(demoName string) *DemoCache {
	for _, demo := range Cache {
		if demo.Name == demoName {
			return demo
		}
	}
	return nil
}

// Use this feature after completing the demo analysis.
// If it returns not nil, then the demo has already been parsed.
func findDemoInMemByDemoID(parser demoinfocs.Parser) *DemoCache {
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
