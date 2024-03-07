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

/*func avgArrayUint64(a []uint64) float64 {
	var sum float64 = 0
	for _, b := range a {
		sum += float64(b)
	}
	sum /= float64(len(a))
	return sum
}*/

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

var isCanceling bool = false

const maxDemosUnziping = 5
var currentDemosUnziping = 0

var totalDemoFiles = 0
var currentCompletedDemoFiles = 0
var usedDemoFiles = 0
var errorDemoFiles = 0

var demofiles []string

/*var statsScorePlr1 []uint64
var statsDamagePlr1 []uint64
var statsKillsPlr1 []uint64
var statsAssistsPlr1 []uint64
var statsDeathsPlr1 []uint64
var statsMVPsPlr1 []uint64
var statsPingPlr1 []uint64*/

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

/*var statsScorePlr2 []uint64
var statsDamagePlr2 []uint64
var statsKillsPlr2 []uint64
var statsAssistsPlr2 []uint64
var statsDeathsPlr2 []uint64
var statsMVPsPlr2 []uint64
var statsPingPlr2 []uint64*/

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

	recurse = flag.Bool("recurse", false, "recurse into subdirectories")

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
	//if *plrID2 == 0 {
	//	log.Panicln("plrID2 not set")
	//}

	//wgParse.Add(1)
	dirParse(filepath.Dir(*dir + "/"))

	for currentDemosUnziping != 0 && currentCompletedDemoFiles < totalDemoFiles {
		PrintProgress()
		time.Sleep(time.Millisecond * 5000)
	}
	wgDem.Wait()
	PrintProgress()

	//wgParse.Wait()
	//wgDem.Wait()

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

	/*statScorePlr1 = avgArrayUint64(statsScorePlr1)
	statDamagePlr1 = avgArrayUint64(statsDamagePlr1)
	statKillsPlr1 = avgArrayUint64(statsKillsPlr1)
	statAssistsPlr1 = avgArrayUint64(statsAssistsPlr1)
	statDeathsPlr1 = avgArrayUint64(statsDeathsPlr1)
	statMVPsPlr1 = avgArrayUint64(statsMVPsPlr1)
	statPingPlr1 = avgArrayUint64(statsPingPlr1)*/

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
		/*statScorePlr2 = avgArrayUint64(statsScorePlr2)
		statDamagePlr2 = avgArrayUint64(statsDamagePlr2)
		statKillsPlr2 = avgArrayUint64(statsKillsPlr2)
		statAssistsPlr2 = avgArrayUint64(statsAssistsPlr2)
		statDeathsPlr2 = avgArrayUint64(statsDeathsPlr2)
		statMVPsPlr2 = avgArrayUint64(statsMVPsPlr2)
		statPingPlr2 = avgArrayUint64(statsPingPlr2)*/

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
			totalDemoFiles++
		}
		time.Sleep(time.Millisecond)
	}
}

//var mutexCompress sync.Mutex

func uncompress(path string, name string) string {
	//mutexCompress.Lock()
	tmpname := createTmpName()

	for currentDemosUnziping >= maxDemosUnziping {
		time.Sleep(time.Millisecond * 500)
	}

	if isCanceling {
		return "isCanceling"
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

func PrintProgress() {
	str := "\n"
	str += "Progress: " + strconv.Itoa(currentCompletedDemoFiles) + " / " + strconv.Itoa(totalDemoFiles) + " %" + strconv.FormatFloat(float64(currentCompletedDemoFiles) / float64(totalDemoFiles) * 100, 'f', 4, 64) + "\n"
	str += "Total demos: " + strconv.Itoa(totalDemoFiles) + "\n"
	str += "Current parsed: " + strconv.Itoa(currentCompletedDemoFiles) + "\n"
	str += "Current used for stats: " + strconv.Itoa(usedDemoFiles) + "\n"
	str += "Errors or duplicates: " + strconv.Itoa(errorDemoFiles) + "\n"
	
	//if (currentCompletedDemoFiles == totalDemoFiles) {
	//	fmt.Print(str)
	//} else if (currentCompletedDemoFiles % 5 == 0) {
	//	fmt.Print(str)
	//}
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
			errorDemoFiles++
			return
		} else if decompressed == "isCanceling" {
			return
		}
		log.Println("file decompressed: ", name)

		demParse(decompressed)

		log.Println("file parsed: ", name)
		currentCompletedDemoFiles++

		os.Remove(decompressed)
	} else if ext == ".dem" {
		if isCanceling {
			return
		}

		demParse(path)
		log.Println("file parsed: ", name)
		currentCompletedDemoFiles++
	}
}

func appendStatsPlr(plr *common.Player) {
	if plr.SteamID64 == *plrID1 {
		/*statsScorePlr1 = append(statsScorePlr1, uint64(plr1.Score()))
		statsDamagePlr1 = append(statsDamagePlr1, uint64(plr1.TotalDamage()))
		statsKillsPlr1 = append(statsKillsPlr1, uint64(plr1.Kills()))
		statsAssistsPlr1 = append(statsAssistsPlr1, uint64(plr1.Assists()))
		statsDeathsPlr1 = append(statsDeathsPlr1, uint64(plr1.Deaths()))
		statsMVPsPlr1 = append(statsMVPsPlr1, uint64(plr1.MVPs()))
		statsPingPlr1 = append(statsPingPlr1, uint64(plr1.Ping()))*/

		statsScorePlr1 += uint64(plr.Score())
		statsDamagePlr1 += uint64(plr.TotalDamage())
		statsKillsPlr1 += uint64(plr.Kills())
		statsAssistsPlr1 += uint64(plr.Assists())
		statsDeathsPlr1 += uint64(plr.Deaths())
		statsMVPsPlr1 += uint64(plr.MVPs())
		statsPingPlr1 += uint64(plr.Ping())
	} else if plr.SteamID64 == *plrID2 {
		/*statsScorePlr2 = append(statsScorePlr2, uint64(plr2.Score()))
		statsDamagePlr2 = append(statsDamagePlr2, uint64(plr2.TotalDamage()))
		statsKillsPlr2 = append(statsKillsPlr2, uint64(plr2.Kills()))
		statsAssistsPlr2 = append(statsAssistsPlr2, uint64(plr2.Assists()))
		statsDeathsPlr2 = append(statsDeathsPlr2, uint64(plr2.Deaths()))
		statsMVPsPlr2 = append(statsMVPsPlr2, uint64(plr2.MVPs()))
		statsPingPlr2 = append(statsPingPlr2, uint64(plr2.Ping()))*/

		statsScorePlr2 += uint64(plr.Score())
		statsDamagePlr2 += uint64(plr.TotalDamage())
		statsKillsPlr2 += uint64(plr.Kills())
		statsAssistsPlr2 += uint64(plr.Assists())
		statsDeathsPlr2 += uint64(plr.Deaths())
		statsMVPsPlr2 += uint64(plr.MVPs())
		statsPingPlr2 += uint64(plr.Ping())
	}
}

func appendStatKillsPlr(e *events.Kill) {
	if e.Killer != nil && e.Killer.SteamID64 == *plrID1 {
		statsPenetratedObjectsPlr1 += uint64(e.PenetratedObjects) 
		if e.IsHeadshot {
			statsHeadShotsPlr1++
		}
		if e.AttackerBlind {
			statsAttackerBlindsPlr1++
		}
		if e.NoScope {
			statsNoScopesPlr1++
		}
		if e.ThroughSmoke {
			statsThroughSmokesPlr1++
		}
	} else if e.Assister != nil && e.Assister.SteamID64 == *plrID1 {
		if e.AssistedFlash {
			statsAssistedFlashsPlr1++
		}
	} else if e.Killer != nil && e.Killer.SteamID64 == *plrID2 {
		statsPenetratedObjectsPlr2 += uint64(e.PenetratedObjects) 
		if e.IsHeadshot {
			statsHeadShotsPlr2++
		}
		if e.AttackerBlind {
			statsAttackerBlindsPlr2++
		}
		if e.NoScope {
			statsNoScopesPlr2++
		}
		if e.ThroughSmoke {
			statsThroughSmokesPlr2++
		}
	} else if e.Assister != nil && e.Assister.SteamID64 == *plrID2 {
		if e.AssistedFlash {
			statsAssistedFlashsPlr2++
		}
	}
}

func demParse(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("failed to open demo file: ", err)
		errorDemoFiles++
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
			errorDemoFiles++
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
			usedDemoFiles++
		} else if plr1 != nil && *plrID2 == 0 {
			appendStatsPlr(plr1)
			usedDemoFiles++
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		log.Println("failed to parse demo: ", err)
		errorDemoFiles++
	}
}
