package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/mholt/archiver/v3"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

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

func avgUint64(a []uint64) float64 {
	var sum float64 = 0
	for _, b := range a {
		sum += float64(b)
	}
	sum /= float64(len(a))
	return sum
}

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

var recurse *bool

var resultFile *string

var plrID1 *uint64
var plrID2 *uint64

//var wgParse sync.WaitGroup
var wgDem sync.WaitGroup

const maxDemosUnziping = 10

var currentDemosUnziping = 0

var demofiles []string

var statsScorePlr1 []uint64
var statsDamagePlr1 []uint64
var statsKillsPlr1 []uint64
var statsAssistsPlr1 []uint64
var statsDeathsPlr1 []uint64
var statsMVPsPlr1 []uint64
var statsPingPlr1 []uint64

var statsScorePlr2 []uint64
var statsDamagePlr2 []uint64
var statsKillsPlr2 []uint64
var statsAssistsPlr2 []uint64
var statsDeathsPlr2 []uint64
var statsMVPsPlr2 []uint64
var statsPingPlr2 []uint64

func main() {
	plrID1 = flag.Uint64("p1", 0, "SteamID64 of first player")
	plrID2 = flag.Uint64("p2", 0, "SteamID64 of second player (not required)")

	dir := flag.String("dir", "", "directory containing demo files")

	recurse = flag.Bool("recurse", false, "recurse into subdirectories")

	resultFile = flag.String("f", "", "result file")

	flag.Parse()

	log.Println("dir: ", *dir+"/")

	if *dir == "" {
		log.Panicln("dir of demofiles not set")
	}
	if *plrID1 == 0 {
		log.Panicln("plrID1 not set")
	}
	//if *plrID2 == 0 {
	//	log.Panicln("plrID2 not set")
	//}

	//wgParse.Add(1)
	dirParse(filepath.Dir(*dir + "/"))

	//wgParse.Wait()
	wgDem.Wait()

	var statScorePlr1 float64 = 0
	var statDamagePlr1 float64 = 0
	var statKillsPlr1 float64 = 0
	var statAssistsPlr1 float64 = 0
	var statDeathsPlr1 float64 = 0
	var statMVPsPlr1 float64 = 0
	var statPingPlr1 float64 = 0

	var statScorePlr2 float64 = 0
	var statDamagePlr2 float64 = 0
	var statKillsPlr2 float64 = 0
	var statAssistsPlr2 float64 = 0
	var statDeathsPlr2 float64 = 0
	var statMVPsPlr2 float64 = 0
	var statPingPlr2 float64 = 0

	statScorePlr1 = avgUint64(statsScorePlr1)
	statDamagePlr1 = avgUint64(statsDamagePlr1)
	statKillsPlr1 = avgUint64(statsKillsPlr1)
	statAssistsPlr1 = avgUint64(statsAssistsPlr1)
	statDeathsPlr1 = avgUint64(statsDeathsPlr1)
	statMVPsPlr1 = avgUint64(statsMVPsPlr1)
	statPingPlr1 = avgUint64(statsPingPlr1)

	if (*plrID2 != 0) {
		statScorePlr2 = avgUint64(statsScorePlr2)
		statDamagePlr2 = avgUint64(statsDamagePlr2)
		statKillsPlr2 = avgUint64(statsKillsPlr2)
		statAssistsPlr2 = avgUint64(statsAssistsPlr2)
		statDeathsPlr2 = avgUint64(statsDeathsPlr2)
		statMVPsPlr2 = avgUint64(statsMVPsPlr2)
		statPingPlr2 = avgUint64(statsPingPlr2)
	}

	result := "Result from demofiles: " + *dir + "\n"
	result = result + "SteamID64 of first player: " + strconv.FormatUint(*plrID1, 10) + "\n"
	if (*plrID2 != 0) {
		result = result + "SteamID64 of second player: " + strconv.FormatUint(*plrID2, 10) + "\n"
	}
	result = result + "First Player Avg Score: " + strconv.FormatFloat(statScorePlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg Score: " + strconv.FormatFloat(statScorePlr2, 'f', 4, 64) + "\n"
	}
	result = result + "First Player Avg TotalDamage: " + strconv.FormatFloat(statDamagePlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg TotalDamage: " + strconv.FormatFloat(statDamagePlr2, 'f', 4, 64) + "\n"
	}
	result = result + "First Player Avg Kills: " + strconv.FormatFloat(statKillsPlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg Kills: " + strconv.FormatFloat(statKillsPlr2, 'f', 4, 64) + "\n"
	}
	result = result + "First Player Avg Assists: " + strconv.FormatFloat(statAssistsPlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg Assists: " + strconv.FormatFloat(statAssistsPlr2, 'f', 4, 64) + "\n"
	}
	result = result + "First Player Avg Deaths: " + strconv.FormatFloat(statDeathsPlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg Deaths: " + strconv.FormatFloat(statDeathsPlr2, 'f', 4, 64) + "\n"
	}
	result = result + "First Player Avg MVPs: " + strconv.FormatFloat(statMVPsPlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg MVPs: " + strconv.FormatFloat(statMVPsPlr2, 'f', 4, 64) + "\n"
	}
	result = result + "First Player Avg Ping: " + strconv.FormatFloat(statPingPlr1, 'f', 4, 64) + "\n"
	if (*plrID2 != 0) {
		result = result + "Second Player Avg Ping: " + strconv.FormatFloat(statPingPlr2, 'f', 4, 64) + "\n"
	}

	fmt.Println(result)

	if *resultFile != "" {
		err := os.WriteFile(*resultFile, []byte(result), os.ModePerm)
		if err != nil {
			log.Panicln("failed to write result file: ", err)
		}
	}
}

func dirParse(path string) {
	//defer wgParse.Done()
	osDir, err := os.ReadDir(path)
	check(err)

	for _, entry := range osDir {
		if entry.IsDir() {
			if *recurse {
				//wgParse.Add(1)
				//go dirParse(filepath.Dir(filepath.Join(path, entry.Name())))
				dirParse(filepath.Dir(filepath.Join(path, entry.Name())))
			}
		} else {
			wgDem.Add(1)
			go demPrepare(filepath.Join(path, entry.Name()), entry.Name())
		}
		time.Sleep(time.Millisecond * 100)
	}
}

//var mutexCompress sync.Mutex

func uncompress(path string, name string) string {
	//mutexCompress.Lock()
	tmpname := createTmpName()

	for currentDemosUnziping >= maxDemosUnziping {
		time.Sleep(time.Millisecond * 500)
	}

	currentDemosUnziping += 1
	err := archiver.DecompressFile(path, filepath.Join(os.TempDir(), name+tmpname))
	if err != nil {
		log.Println("failed to decompress file: ", err)
		currentDemosUnziping -= 1
		return ""
	}
	currentDemosUnziping -= 1
	//mutexCompress.Unlock()

	return filepath.Join(os.TempDir(), name+tmpname)
}

func demPrepare(path string, name string) {
	defer wgDem.Done()

	ext := filepath.Ext(path)
	if ext == ".bz2" || ext == ".gz" {
		decompressed := uncompress(path, name)
		if decompressed == "" {
			return
		}
		log.Println("file decompressed: ", name)

		demParse(decompressed)

		log.Println("file parsed: ", name)

		os.Remove(decompressed)
	} else if ext == ".dem" {
		demParse(path)
		log.Println("file parsed: ", name)
	}
}

func appendStatsPlr1(plr1 *common.Player) {
	statsScorePlr1 = append(statsScorePlr1, uint64(plr1.Score()))
	statsDamagePlr1 = append(statsDamagePlr1, uint64(plr1.TotalDamage()))
	statsKillsPlr1 = append(statsKillsPlr1, uint64(plr1.Kills()))
	statsAssistsPlr1 = append(statsAssistsPlr1, uint64(plr1.Assists()))
	statsDeathsPlr1 = append(statsDeathsPlr1, uint64(plr1.Deaths()))
	statsMVPsPlr1 = append(statsMVPsPlr1, uint64(plr1.MVPs()))
	statsPingPlr1 = append(statsPingPlr1, uint64(plr1.Ping()))
}

func appendStatsPlr2(plr2 *common.Player) {
	statsScorePlr2 = append(statsScorePlr2, uint64(plr2.Score()))
	statsDamagePlr2 = append(statsDamagePlr2, uint64(plr2.TotalDamage()))
	statsKillsPlr2 = append(statsKillsPlr2, uint64(plr2.Kills()))
	statsAssistsPlr2 = append(statsAssistsPlr2, uint64(plr2.Assists()))
	statsDeathsPlr2 = append(statsDeathsPlr2, uint64(plr2.Deaths()))
	statsMVPsPlr2 = append(statsMVPsPlr2, uint64(plr2.MVPs()))
	statsPingPlr2 = append(statsPingPlr2, uint64(plr2.Ping()))
}

func demParse(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("failed to open demo file: ", err)
	}
	defer f.Close()

	p := demoinfocs.NewParser(f)
	defer p.Close()

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
			return
		}
		demofiles = append(demofiles, demoid)



		var plr1 *common.Player
		var plr2 *common.Player

		// CT
		for _, plr := range ct.Members() {
			if plr.SteamID64 == *plrID1 {
				plr1 = plr
			} else if plr.SteamID64 == *plrID2 {
				plr2 = plr
			}
		}

		// T
		for _, plr := range t.Members() {
			if plr.SteamID64 == *plrID1 {
				plr1 = plr
			} else if plr.SteamID64 == *plrID2 {
				plr2 = plr
			}
		}
		
		log.Println("found players: ", plr1, plr2)

		if plr1 != nil && plr2 == nil {
			appendStatsPlr1(plr1)
		}
		
		if plr1 != nil && plr2 != nil {
			appendStatsPlr1(plr1)
			appendStatsPlr2(plr2)
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		log.Println("failed to parse demo: ", err)
	}
}
