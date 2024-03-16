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
)

func check(e error) {
	if e != nil {
		log.Panic(e)
	}
}

var vertResult *bool
var dir *string
var recurse *bool
var resultFile *string

var cacheDir *string
var useCache *bool

var wgDem sync.WaitGroup

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
		flag.PrintDefaults()
		return
		//log.Panicln("one of -p1, -p2, -p3, -p4, -p5 is required")
	}

	log.Println("dir: ", *dir+"/")

	if *dir == "" {
		flag.PrintDefaults()
		return
		//log.Panicln("dir of demofiles not set")
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

func PrintProgress() {
	str := "\n"
	str += "Progress: " + strconv.Itoa(int(currentCompletedDemoFiles + errorDemoFiles)) + " / " + strconv.Itoa(int(totalDemoFiles)) + " %" + strconv.FormatFloat(float64(currentCompletedDemoFiles + errorDemoFiles) / float64(totalDemoFiles) * 100, 'f', 4, 64) + "\n"
	str += "Total demos: " + strconv.Itoa(int(totalDemoFiles)) + "\n"
	str += "Current parsed: " + strconv.Itoa(int(currentCompletedDemoFiles)) + "\n"
	str += "Current used for stats: " + strconv.Itoa(int(usedDemoFiles)) + "\n"
	str += "Errors, duplicates or skips: " + strconv.Itoa(int(errorDemoFiles)) + "\n"

	fmt.Println(str)
}