package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var vertResult *bool
var dir *string
var recurse *bool
var resultFile *string

var cacheDir *string
var useCache *bool

var gamemode *string

func main() {
	plrID1 := flag.Uint64("p1", 0, "SteamID64 of first player")
	plrID2 := flag.Uint64("p2", 0, "SteamID64 of second player (not required)")
	plrID3 := flag.Uint64("p3", 0, "SteamID64 of third player (not required)")
	plrID4 := flag.Uint64("p4", 0, "SteamID64 of fourth player (not required)")
	plrID5 := flag.Uint64("p5", 0, "SteamID64 of fifth player (not required)")

	vertResult = flag.Bool("v", false, "display the result horizontally (compact and perhaps not beautiful)")

	dir = flag.String("dir", "", "directory containing demo files")
	recurse = flag.Bool("r", false, "recursion into subdirectories")
	resultFile = flag.String("f", "", "results file")

	cacheDir = flag.String("cd", "cache", "cache directory")
	useCache = flag.Bool("c", true, "use cache")

	gamemode = flag.String("gm", "", "game mode (m - matchmaking, w - wingman, o - other) can be specified separated by commas (format -gm=m or -gm='m,w,o')")

	flag.Parse()

	if *plrID1 != 0 {
		p := PlrStats{SteamID64: *plrID1}
		PlrsStats = append(PlrsStats, &p)
	}
	if *plrID2 != 0 {
		p := PlrStats{SteamID64: *plrID2}
		PlrsStats = append(PlrsStats, &p)
	}
	if *plrID3 != 0 {
		p := PlrStats{SteamID64: *plrID3}
		PlrsStats = append(PlrsStats, &p)
	}
	if *plrID4 != 0 {
		p := PlrStats{SteamID64: *plrID4}
		PlrsStats = append(PlrsStats, &p)
	}
	if *plrID5 != 0 {
		p := PlrStats{SteamID64: *plrID5}
		PlrsStats = append(PlrsStats, &p)
	}

	if *plrID1 == 0 && *plrID2 == 0 && *plrID3 == 0 && *plrID4 == 0 && *plrID5 == 0 {
		log.Println("one of -p1, -p2, -p3, -p4, -p5 is required")
		flag.PrintDefaults()
		return
	}

	log.Println("dir: ", *dir+"/")

	if *dir == "" {
		log.Println("-dir is required")
		flag.PrintDefaults()
		return
	}

	gmSplit := strings.Split(*gamemode, ",")

	for _, gm := range gmSplit {
		if gm == "m" {
			useStatsMatchmaking = true
		} else if gm == "w" {
			useStatsWingman = true
		} else if gm == "o" {
			useStatsOther = true
		}
	}

	if !useStatsMatchmaking && !useStatsWingman && !useStatsOther {
		log.Println("-gm is not set, used -gm=m,w,o")
		useStatsMatchmaking = true
		useStatsWingman = true
		useStatsOther = true
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
        sig := <-sigs

        fmt.Println()
        fmt.Println(sig)
		log.Println("canceling...")
        isCanceling = true
    }()

	dirParse(filepath.Dir(*dir + "/"))

	for currentDemosUnziping != 0 && currentCompletedDemoFiles < totalDemoFiles {
		PrintProgress()
		time.Sleep(time.Millisecond * 5000)
	}
	wgDem.Wait()
	PrintProgress()

	PrintResults()
}

func dirParse(path string) {
	osDir, err := os.ReadDir(path)
	if err != nil {
		log.Println("failed to read dir: ", err)
		return
	}

	for _, entry := range osDir {
		if entry.IsDir() && *recurse {
			dirParse(filepath.Join(path, entry.Name() + "/"))
		} else {
			wgDem.Add(1)
			go demPrepare(filepath.Join(path, entry.Name()), entry.Name())
			totalDemoFiles++
		}
		time.Sleep(time.Microsecond * 25)
	}
}

func PrintProgress() {
	str := "\n"
	str += "Progress: " + strconv.Itoa(int(currentCompletedDemoFiles + currentCachedDemoFiles + errorDemoFiles)) + " / " + strconv.Itoa(int(totalDemoFiles)) + " %" + strconv.FormatFloat(float64(currentCompletedDemoFiles + currentCachedDemoFiles + errorDemoFiles) / float64(totalDemoFiles) * 100, 'f', 4, 64) + "\n"
	str += "Total demos: " + strconv.Itoa(int(totalDemoFiles)) + "\n"
	str += "Current parsed: " + strconv.Itoa(int(currentCompletedDemoFiles)) + "\n"
	str += "Current used for stats: " + strconv.Itoa(int(usedDemoFiles)) + "\n"
	str += "Current cached: " + strconv.Itoa(int(currentCachedDemoFiles)) + "\n"
	str += "Errors, duplicates or skips: " + strconv.Itoa(int(errorDemoFiles)) + "\n"

	fmt.Println(str)
}