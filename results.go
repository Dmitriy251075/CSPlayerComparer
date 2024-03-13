package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
)

const SteamID64Line = 0
const NameLine = 1
const ScoreLine = 2
const TotalDamageLine = 3
const KillsLine = 4
const AssistsLine = 5
const DeathsLine = 6
const MVPsLine = 7
const PingLine = 8
const KillsPenetratedObjectsLine = 9
const KillsHeadshotLine = 10
const KillsAssistedFlashsLine = 11
const KillsAttackerBlindsLine = 12
const KillsNoScopeLine = 13
const KillsThroughSmokeLine = 14

func PrintResults() {
	fmt.Println("Avg Stats")
	fmt.Println("Result from demofiles: ", *dir)

	result := [][]string {
		{"Players SteamID64:"}, // SteamID64Line 0
		{"Players Name:"}, // NameLine 1
		{"Avg Score:"}, // ScoreLine 2
		{"Avg TotalDamage:"}, // TotalDamageLine 3
		{"Avg Kills:"}, // KillsLine 4
		{"Avg Assists:"}, // AssistsLine 5
		{"Avg Deaths:"}, // DeathsLine 6
		{"Avg MVPs:"}, // MVPsLine 7
		{"Avg Ping:"}, // PingLine 8
		{"Avg Kills Penetrated Objects:"}, // KillsPenetratedObjectsLine 9
		{"Avg Kills Headshot:"}, // KillsHeadshotLine 10
		{"Avg Kills Assisted Flashs:"}, // KillsAssistedFlashsLine 11
		{"Avg Kills Attacker Blinds:"}, // KillsAttackerBlindsLine 12
		{"Avg Kills No Scope:"}, // KillsNoScopeLine 13
		{"Avg Kills Through Smoke:"}, // KillsThroughSmokeLine 14
	}

	for _, plr := range PlrsStats {
		result = PrintPlrStat(plr, result)
	}

	tw := table.NewWriter()
	for line := 0; line < len(result); line++ {
		row := table.Row{result[line][0]}
		for col := 1; col < len(result[line]); col++ {
			row = append(row, result[line][col])
		}
		tw.AppendRow(row)
	}

	str := tw.Render()

	fmt.Print(str)

	if *resultFile != "" {
		err := os.WriteFile(*resultFile, []byte(str), 0644)
		if err != nil {
			log.Panicln("failed to write result file: ", err)
		}
	}
}

func PrintPlrStat(plr *PlrStats, result [][]string) [][]string {
	var statScorePlr float64 = 0
	var statDamagePlr float64 = 0
	var statKillsPlr float64 = 0
	var statAssistsPlr float64 = 0
	var statDeathsPlr float64 = 0
	var statMVPsPlr float64 = 0
	var statPingPlr float64 = 0
	var statPenetratedObjectsPlr float64 = 0
	var statHeadShotsPlr float64 = 0
	var statAssistedFlashsPlr float64 = 0
	var statAttackerBlindsPlr float64 = 0
	var statNoScopesPlr float64 = 0
	var statThroughSmokesPlr float64 = 0

	statScorePlr = float64(plr.statsScore) / float64(usedDemoFiles)
	statDamagePlr = float64(plr.statsDamage) / float64(usedDemoFiles)
	statKillsPlr = float64(plr.statsKills) / float64(usedDemoFiles)
	statAssistsPlr = float64(plr.statsAssists) / float64(usedDemoFiles)
	statDeathsPlr = float64(plr.statsDeaths) / float64(usedDemoFiles)
	statMVPsPlr = float64(plr.statsMVPs) / float64(usedDemoFiles)
	statPingPlr = float64(plr.statsPing) / float64(usedDemoFiles)
	statPenetratedObjectsPlr = float64(plr.statsPenetratedObjects) / float64(usedDemoFiles)
	statHeadShotsPlr = float64(plr.statsHeadShots) / float64(usedDemoFiles)
	statAssistedFlashsPlr = float64(plr.statsAssistedFlashs) / float64(usedDemoFiles)
	statAttackerBlindsPlr = float64(plr.statsAttackerBlinds) / float64(usedDemoFiles)
	statNoScopesPlr = float64(plr.statsNoScopes) / float64(usedDemoFiles)
	statThroughSmokesPlr = float64(plr.statsThroughSmokes) / float64(usedDemoFiles)

	for line := 0; line < len(result); line++ {
		if (line == SteamID64Line) {
			result[line] = append(result[line], strconv.FormatUint(plr.SteamID64, 10))
		} else if (line == NameLine) {
			result[line] = append(result[line], plr.Name)
		} else if (line == ScoreLine) {
			result[line] = append(result[line], strconv.FormatFloat(statScorePlr, 'f', 4, 64))
		} else if (line == TotalDamageLine) {
			result[line] = append(result[line], strconv.FormatFloat(statDamagePlr, 'f', 4, 64))
		} else if (line == KillsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statKillsPlr, 'f', 4, 64))
		} else if (line == AssistsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statAssistsPlr, 'f', 4, 64))
		} else if (line == DeathsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statDeathsPlr, 'f', 4, 64))
		} else if (line == MVPsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statMVPsPlr, 'f', 4, 64))
		} else if (line == PingLine) {
			result[line] = append(result[line], strconv.FormatFloat(statPingPlr, 'f', 4, 64))
		} else if (line == KillsPenetratedObjectsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statPenetratedObjectsPlr, 'f', 4, 64))
		} else if (line == KillsHeadshotLine) {
			result[line] = append(result[line], strconv.FormatFloat(statHeadShotsPlr, 'f', 4, 64))
		} else if (line == KillsAssistedFlashsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statAssistedFlashsPlr, 'f', 4, 64))
		} else if (line == KillsAttackerBlindsLine) {
			result[line] = append(result[line], strconv.FormatFloat(statAttackerBlindsPlr, 'f', 4, 64))
		} else if (line == KillsNoScopeLine) {
			result[line] = append(result[line], strconv.FormatFloat(statNoScopesPlr, 'f', 4, 64))
		} else if (line == KillsThroughSmokeLine) {
			result[line] = append(result[line], strconv.FormatFloat(statThroughSmokesPlr, 'f', 4, 64))
		}
	}

	return result
}
