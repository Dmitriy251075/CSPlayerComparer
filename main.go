package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
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

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

var recurse *bool

var resultFile *string

var plrID1 *uint64
var plrID2 *uint64

var wgDem sync.WaitGroup

var isCanceling bool = false

const maxDemosUnziping uint32 = 5
var currentDemosUnziping uint32

var totalDemoFiles uint32
var currentCompletedDemoFiles uint32
var usedDemoFiles uint32
var errorDemoFiles uint32

var demofiles []string

var statsScorePlr1 uint64
var statsDamagePlr1 uint64
var statsKillsPlr1 uint64
var statsAssistsPlr1 uint64
var statsDeathsPlr1 uint64
var statsMVPsPlr1 uint64
var statsPingPlr1 uint64
var statsPenetratedObjectsPlr1 uint64
var statsHeadShotsPlr1 uint64
var statsAssistedFlashsPlr1 uint64
var statsAttackerBlindsPlr1 uint64
var statsNoScopesPlr1 uint64
var statsThroughSmokesPlr1 uint64

var statsScorePlr2 uint64
var statsDamagePlr2 uint64
var statsKillsPlr2 uint64
var statsAssistsPlr2 uint64
var statsDeathsPlr2 uint64
var statsMVPsPlr2 uint64
var statsPingPlr2 uint64
var statsPenetratedObjectsPlr2 uint64
var statsHeadShotsPlr2 uint64
var statsAssistedFlashsPlr2 uint64
var statsAttackerBlindsPlr2 uint64
var statsNoScopesPlr2 uint64
var statsThroughSmokesPlr2 uint64

func main() {
	plrID1 = flag.Uint64("p1", 0, "SteamID64 of first player")
	plrID2 = flag.Uint64("p2", 0, "SteamID64 of second player (not required)")

	dir := flag.String("dir", "", "directory containing demo files")

	recurse = flag.Bool("r", false, "recurse into subdirectories")

	resultFile = flag.String("f", "", "result file")

	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
        sig := <-sigs

        fmt.Println()
        fmt.Println(sig)
		log.Println("canceling...")
        isCanceling = true
    }()

	log.Println("dir: ", *dir+"/")

	if *dir == "" {
		log.Panicln("dir of demofiles not set")
	}
	if *plrID1 == 0 {
		log.Panicln("plrID1 not set")
	}

	dirParse(filepath.Dir(*dir + "/"))

	for currentDemosUnziping != 0 && currentCompletedDemoFiles < totalDemoFiles {
		PrintProgress()
		time.Sleep(time.Millisecond * 5000)
	}
	wgDem.Wait()
	PrintProgress()

	var statScorePlr1 float64 = 0
	var statDamagePlr1 float64 = 0
	var statKillsPlr1 float64 = 0
	var statAssistsPlr1 float64 = 0
	var statDeathsPlr1 float64 = 0
	var statMVPsPlr1 float64 = 0
	var statPingPlr1 float64 = 0
	var statPenetratedObjectsPlr1 float64 = 0
	var statHeadShotsPlr1 float64 = 0
	var statAssistedFlashsPlr1 float64 = 0
	var statAttackerBlindsPlr1 float64 = 0
	var statNoScopesPlr1 float64 = 0
	var statThroughSmokesPlr1 float64 = 0

	var statScorePlr2 float64 = 0
	var statDamagePlr2 float64 = 0
	var statKillsPlr2 float64 = 0
	var statAssistsPlr2 float64 = 0
	var statDeathsPlr2 float64 = 0
	var statMVPsPlr2 float64 = 0
	var statPingPlr2 float64 = 0
	var statPenetratedObjectsPlr2 float64 = 0
	var statHeadShotsPlr2 float64 = 0
	var statAssistedFlashsPlr2 float64 = 0
	var statAttackerBlindsPlr2 float64 = 0
	var statNoScopesPlr2 float64 = 0
	var statThroughSmokesPlr2 float64 = 0

	statScorePlr1 = float64(statsScorePlr1) / float64(usedDemoFiles)
	statDamagePlr1 = float64(statsDamagePlr1) / float64(usedDemoFiles)
	statKillsPlr1 = float64(statsKillsPlr1) / float64(usedDemoFiles)
	statAssistsPlr1 = float64(statsAssistsPlr1) / float64(usedDemoFiles)
	statDeathsPlr1 = float64(statsDeathsPlr1) / float64(usedDemoFiles)
	statMVPsPlr1 = float64(statsMVPsPlr1) / float64(usedDemoFiles)
	statPingPlr1 = float64(statsPingPlr1) / float64(usedDemoFiles)
	statPenetratedObjectsPlr1 = float64(statsPenetratedObjectsPlr1) / float64(usedDemoFiles)
	statHeadShotsPlr1 = float64(statsHeadShotsPlr1) / float64(usedDemoFiles)
	statAssistedFlashsPlr1 = float64(statsAssistedFlashsPlr1) / float64(usedDemoFiles)
	statAttackerBlindsPlr1 = float64(statsAttackerBlindsPlr1) / float64(usedDemoFiles)
	statNoScopesPlr1 = float64(statsNoScopesPlr1) / float64(usedDemoFiles)
	statThroughSmokesPlr1 = float64(statsThroughSmokesPlr1) / float64(usedDemoFiles)

	if *plrID2 != 0 {
		statScorePlr2 = float64(statsScorePlr2) / float64(usedDemoFiles)
		statDamagePlr2 = float64(statsDamagePlr2) / float64(usedDemoFiles)
		statKillsPlr2 = float64(statsKillsPlr2) / float64(usedDemoFiles)
		statAssistsPlr2 = float64(statsAssistsPlr2) / float64(usedDemoFiles)
		statDeathsPlr2 = float64(statsDeathsPlr2) / float64(usedDemoFiles)
		statMVPsPlr2 = float64(statsMVPsPlr2) / float64(usedDemoFiles)
		statPingPlr2 = float64(statsPingPlr2) / float64(usedDemoFiles)
		statPenetratedObjectsPlr2 = float64(statsPenetratedObjectsPlr2) / float64(usedDemoFiles)
		statHeadShotsPlr2 = float64(statsHeadShotsPlr2) / float64(usedDemoFiles)
		statAssistedFlashsPlr2 = float64(statsAssistedFlashsPlr2) / float64(usedDemoFiles)
		statAttackerBlindsPlr2 = float64(statsAttackerBlindsPlr2) / float64(usedDemoFiles)
		statNoScopesPlr2 = float64(statsNoScopesPlr2) / float64(usedDemoFiles)
		statThroughSmokesPlr2 = float64(statsThroughSmokesPlr2) / float64(usedDemoFiles)
	}

	result := "Result from demofiles: " + *dir + "\n"
	result += "SteamID64 of first player: " + strconv.FormatUint(*plrID1, 10) + "\n"
	if *plrID2 != 0 {
		result += "SteamID64 of second player: " + strconv.FormatUint(*plrID2, 10) + "\n"
	}
	result += "First Player Avg Score: " + strconv.FormatFloat(statScorePlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Score: " + strconv.FormatFloat(statScorePlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg TotalDamage: " + strconv.FormatFloat(statDamagePlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg TotalDamage: " + strconv.FormatFloat(statDamagePlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Kills: " + strconv.FormatFloat(statKillsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Kills: " + strconv.FormatFloat(statKillsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Assists: " + strconv.FormatFloat(statAssistsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Assists: " + strconv.FormatFloat(statAssistsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Deaths: " + strconv.FormatFloat(statDeathsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Deaths: " + strconv.FormatFloat(statDeathsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg MVPs: " + strconv.FormatFloat(statMVPsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg MVPs: " + strconv.FormatFloat(statMVPsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Ping: " + strconv.FormatFloat(statPingPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Ping: " + strconv.FormatFloat(statPingPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Penetrated Objects: " + strconv.FormatFloat(statPenetratedObjectsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Penetrated Objects: " + strconv.FormatFloat(statPenetratedObjectsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Kills Headshot: " + strconv.FormatFloat(statHeadShotsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Kills Headshot: " + strconv.FormatFloat(statHeadShotsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Assisted Flashs: " + strconv.FormatFloat(statAssistedFlashsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Assisted Flashs: " + strconv.FormatFloat(statAssistedFlashsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Attacker Blinds: " + strconv.FormatFloat(statAttackerBlindsPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Attacker Blinds: " + strconv.FormatFloat(statAttackerBlindsPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Kills No Scope: " + strconv.FormatFloat(statNoScopesPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Kills No Scope: " + strconv.FormatFloat(statNoScopesPlr2, 'f', 4, 64) + "\n"
	}
	result += "First Player Avg Kills Through Smoke: " + strconv.FormatFloat(statThroughSmokesPlr1, 'f', 4, 64) + "\n"
	if *plrID2 != 0 {
		result += "Second Player Avg Kills Through Smoke: " + strconv.FormatFloat(statThroughSmokesPlr2, 'f', 4, 64) + "\n"
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
	osDir, err := os.ReadDir(path)
	check(err)

	for _, entry := range osDir {
		if entry.IsDir() && *recurse {
			dirParse(filepath.Join(path, entry.Name() + "/"))
		} else {
			wgDem.Add(1)
			go demPrepare(filepath.Join(path, entry.Name()), entry.Name())
			totalDemoFiles++
		}
		time.Sleep(time.Microsecond * 100)
	}
}

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

func PrintProgress() {
	str := "\n"
	str += "Progress: " + strconv.Itoa(int(currentCompletedDemoFiles)) + " / " + strconv.Itoa(int(totalDemoFiles)) + " %" + strconv.FormatFloat(float64(currentCompletedDemoFiles) / float64(totalDemoFiles) * 100, 'f', 4, 64) + "\n"
	str += "Total demos: " + strconv.Itoa(int(totalDemoFiles)) + "\n"
	str += "Current parsed: " + strconv.Itoa(int(currentCompletedDemoFiles)) + "\n"
	str += "Current used for stats: " + strconv.Itoa(int(usedDemoFiles)) + "\n"
	str += "Errors or duplicates: " + strconv.Itoa(int(errorDemoFiles)) + "\n"

	fmt.Println(str)
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

func appendStatsPlr(plr *common.Player) {
	if plr.SteamID64 == *plrID1 {
		atomic.AddUint64(&statsScorePlr1, uint64(plr.Score()))
		atomic.AddUint64(&statsDamagePlr1, uint64(plr.TotalDamage()))
		atomic.AddUint64(&statsKillsPlr1, uint64(plr.Kills()))
		atomic.AddUint64(&statsAssistsPlr1, uint64(plr.Assists()))
		atomic.AddUint64(&statsDeathsPlr1, uint64(plr.Deaths()))
		atomic.AddUint64(&statsMVPsPlr1, uint64(plr.MVPs()))
		atomic.AddUint64(&statsPingPlr1, uint64(plr.Ping()))
	} else if plr.SteamID64 == *plrID2 {
		atomic.AddUint64(&statsScorePlr2, uint64(plr.Score()))
		atomic.AddUint64(&statsDamagePlr2, uint64(plr.TotalDamage()))
		atomic.AddUint64(&statsKillsPlr2, uint64(plr.Kills()))
		atomic.AddUint64(&statsAssistsPlr2, uint64(plr.Assists()))
		atomic.AddUint64(&statsDeathsPlr2, uint64(plr.Deaths()))
		atomic.AddUint64(&statsMVPsPlr2, uint64(plr.MVPs()))
		atomic.AddUint64(&statsPingPlr2, uint64(plr.Ping()))
	}
}

func appendStatKillsPlr(e *events.Kill) {
	if e.Killer != nil && e.Killer.SteamID64 == *plrID1 {
		atomic.AddUint64(&statsPenetratedObjectsPlr1, uint64(e.PenetratedObjects))
		if e.IsHeadshot {
			atomic.AddUint64(&statsHeadShotsPlr1, 1)
		}
		if e.AttackerBlind {
			atomic.AddUint64(&statsAttackerBlindsPlr1, 1)
		}
		if e.NoScope {
			atomic.AddUint64(&statsNoScopesPlr1, 1)
		}
		if e.ThroughSmoke {
			atomic.AddUint64(&statsThroughSmokesPlr1, 1)
		}
	} else if e.Assister != nil && e.Assister.SteamID64 == *plrID1 {
		if e.AssistedFlash {
			atomic.AddUint64(&statsAssistedFlashsPlr1, 1)
		}
	} else if e.Killer != nil && e.Killer.SteamID64 == *plrID2 {
		atomic.AddUint64(&statsPenetratedObjectsPlr2, uint64(e.PenetratedObjects))
		if e.IsHeadshot {
			atomic.AddUint64(&statsHeadShotsPlr2, 1)
		}
		if e.AttackerBlind {
			atomic.AddUint64(&statsAttackerBlindsPlr2, 1)
		}
		if e.NoScope {
			atomic.AddUint64(&statsNoScopesPlr2, 1)
		}
		if e.ThroughSmoke {
			atomic.AddUint64(&statsThroughSmokesPlr2, 1)
		}
	} else if e.Assister != nil && e.Assister.SteamID64 == *plrID2 {
		if e.AssistedFlash {
			atomic.AddUint64(&statsAssistedFlashsPlr2, 1)
		}
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
		gs := p.GameState()

		ct := gs.TeamCounterTerrorists()
		t := gs.TeamTerrorists()



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

		if plr1 != nil && plr2 != nil {
			if e.Killer != nil && e.Killer.SteamID64 == *plrID1 || e.Assister != nil && e.Assister.SteamID64 == *plrID1 ||e.Killer != nil && e.Killer.SteamID64 == *plrID2 || e.Assister != nil && e.Assister.SteamID64 == *plrID2 {
				appendStatKillsPlr(&e)
			}
		} else if plr1 != nil && *plrID2 == 0 {
			if e.Killer != nil && e.Killer.SteamID64 == *plrID1 || e.Assister != nil && e.Assister.SteamID64 == *plrID1 ||e.Killer != nil && e.Killer.SteamID64 == *plrID2 || e.Assister != nil && e.Assister.SteamID64 == *plrID2 {
				appendStatKillsPlr(&e)
			}
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

		if plr1 != nil && plr2 != nil {
			appendStatsPlr(plr1)
			appendStatsPlr(plr2)
			atomic.AddUint32(&usedDemoFiles, 1)
		} else if plr1 != nil && *plrID2 == 0 {
			appendStatsPlr(plr1)
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
