package main

import (
	"log"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/events"
)

type PlrStats struct {
	Name                   string
	SteamID64              uint64
	statsScore             uint64
	statsDamage            uint64
	statsKills             uint64
	statsAssists           uint64
	statsDeaths            uint64
	statsMVPs              uint64
	statsPing              uint64
	statsPenetratedObjects uint64
	statsHeadShots         uint64
	statsAssistedFlashs    uint64
	statsAttackerBlinds    uint64
	statsNoScopes          uint64
	statsThroughSmokes     uint64
}

func (p *PlrStats) appendStatKills(e *events.Kill) {
	if e.Killer != nil && e.Killer.SteamID64 == p.SteamID64 {
		atomic.AddUint64(&p.statsPenetratedObjects, uint64(e.PenetratedObjects))
		if e.IsHeadshot {
			atomic.AddUint64(&p.statsHeadShots, 1)
		}
		if e.AttackerBlind {
			atomic.AddUint64(&p.statsAttackerBlinds, 1)
		}
		if e.NoScope {
			atomic.AddUint64(&p.statsNoScopes, 1)
		}
		if e.ThroughSmoke {
			atomic.AddUint64(&p.statsThroughSmokes, 1)
		}
	} else if e.Assister != nil && e.Assister.SteamID64 == p.SteamID64 {
		if e.AssistedFlash {
			atomic.AddUint64(&p.statsAssistedFlashs, 1)
		}
	}
}

func (p *PlrStats) setStats(plr *common.Player) {
	if p.SteamID64 == plr.SteamID64 {
		if p.Name == "" {
			p.Name = plr.Name
		}
	
		atomic.StoreUint64(&p.statsScore, uint64(plr.Score()))
		atomic.StoreUint64(&p.statsDamage, uint64(plr.TotalDamage()))
		atomic.StoreUint64(&p.statsKills, uint64(plr.Kills()))
		atomic.StoreUint64(&p.statsAssists, uint64(plr.Assists()))
		atomic.StoreUint64(&p.statsDeaths, uint64(plr.Deaths()))
		atomic.StoreUint64(&p.statsMVPs, uint64(plr.MVPs()))
		atomic.StoreUint64(&p.statsPing, uint64(plr.Ping()))
	}
}

func (p *PlrStats) appendStatsFromPlrStats(otherPlrStats *PlrStats) {
	if (p.SteamID64 == otherPlrStats.SteamID64) {
		if p.Name == "" {
			p.Name = otherPlrStats.Name
		}

		atomic.AddUint64(&p.statsScore, otherPlrStats.statsScore)
		atomic.AddUint64(&p.statsDamage, otherPlrStats.statsDamage)
		atomic.AddUint64(&p.statsKills, otherPlrStats.statsKills)
		atomic.AddUint64(&p.statsAssists, otherPlrStats.statsAssists)
		atomic.AddUint64(&p.statsDeaths, otherPlrStats.statsDeaths)
		atomic.AddUint64(&p.statsMVPs, otherPlrStats.statsMVPs)
		atomic.AddUint64(&p.statsPing, otherPlrStats.statsPing)

		atomic.AddUint64(&p.statsPenetratedObjects, otherPlrStats.statsPenetratedObjects)
		atomic.AddUint64(&p.statsHeadShots, otherPlrStats.statsHeadShots)
		atomic.AddUint64(&p.statsAttackerBlinds, otherPlrStats.statsAttackerBlinds)
		atomic.AddUint64(&p.statsNoScopes, otherPlrStats.statsNoScopes)
		atomic.AddUint64(&p.statsThroughSmokes, otherPlrStats.statsThroughSmokes)
		atomic.AddUint64(&p.statsAssistedFlashs, otherPlrStats.statsAssistedFlashs)
	}
}

const StrfmtPlrStatsEnd = "\t"
const StrfmtPlrStatsName = "Name='"
const StrfmtPlrStatsNameEnd = "'" + StrfmtPlrStatsEnd
const StrfmtPlrStatsSteamID64 = "SteamID64="
const StrfmtPlrStatsSteamID64End = StrfmtPlrStatsEnd
const StrfmtPlrStatsScore = "Score="
const StrfmtPlrStatsScoreEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsDamage = "Damage="
const StrfmtPlrStatsDamageEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsKills = "Kills="
const StrfmtPlrStatsKillsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsAssists = "Assists="
const StrfmtPlrStatsAssistsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsDeaths = "Deaths="
const StrfmtPlrStatsDeathsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsMVPs = "MVPs="
const StrfmtPlrStatsMVPsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsPing = "Ping="
const StrfmtPlrStatsPingEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsPenetratedObjects = "PenetratedObjects="
const StrfmtPlrStatsPenetratedObjectsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsHeadShots = "HeadShots="
const StrfmtPlrStatsHeadShotsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsAttackerBlinds = "AttackerBlinds="
const StrfmtPlrStatsAttackerBlindsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsAssistedFlashs = "AssistedFlashs="
const StrfmtPlrStatsAssistedFlashsEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsNoScopes = "NoScopes="
const StrfmtPlrStatsNoScopesEnd = StrfmtPlrStatsEnd
const StrfmtPlrStatsThroughSmokes = "ThroughSmokes="
const StrfmtPlrStatsThroughSmokesEnd = StrfmtPlrStatsEnd

func (p *PlrStats) toString() string {
	str := StrfmtPlrStatsName + p.Name + StrfmtPlrStatsNameEnd
	str += StrfmtPlrStatsSteamID64 + strconv.FormatUint(p.SteamID64, 10) + StrfmtPlrStatsSteamID64End
	str += StrfmtPlrStatsScore + strconv.FormatUint(p.statsScore, 10) + StrfmtPlrStatsScoreEnd
	str += StrfmtPlrStatsDamage + strconv.FormatUint(p.statsDamage, 10) + StrfmtPlrStatsDamageEnd
	str += StrfmtPlrStatsKills + strconv.FormatUint(p.statsKills, 10) + StrfmtPlrStatsKillsEnd
	str += StrfmtPlrStatsAssists + strconv.FormatUint(p.statsAssists, 10) + StrfmtPlrStatsAssistsEnd
	str += StrfmtPlrStatsDeaths + strconv.FormatUint(p.statsDeaths, 10) + StrfmtPlrStatsDeathsEnd
	str += StrfmtPlrStatsMVPs + strconv.FormatUint(p.statsMVPs, 10) + StrfmtPlrStatsMVPsEnd
	str += StrfmtPlrStatsPing + strconv.FormatUint(p.statsPing, 10) + StrfmtPlrStatsPingEnd
	str += StrfmtPlrStatsPenetratedObjects + strconv.FormatUint(p.statsPenetratedObjects, 10) + StrfmtPlrStatsPenetratedObjectsEnd
	str += StrfmtPlrStatsHeadShots + strconv.FormatUint(p.statsHeadShots, 10) + StrfmtPlrStatsHeadShotsEnd
	str += StrfmtPlrStatsAttackerBlinds + strconv.FormatUint(p.statsAttackerBlinds, 10) + StrfmtPlrStatsAttackerBlindsEnd
	str += StrfmtPlrStatsAssistedFlashs + strconv.FormatUint(p.statsAssistedFlashs, 10) + StrfmtPlrStatsAssistedFlashsEnd
	str += StrfmtPlrStatsNoScopes + strconv.FormatUint(p.statsNoScopes, 10) + StrfmtPlrStatsNoScopesEnd
	str += StrfmtPlrStatsThroughSmokes + strconv.FormatUint(p.statsThroughSmokes, 10) + StrfmtPlrStatsThroughSmokesEnd

	return str
}

// Returns true if successful
func (p *PlrStats) fromString(str string) bool {
	var err error
	var Split []string

	var IsError bool = false

	Split = strings.Split(str, StrfmtPlrStatsName)
	if len(Split) >= 2 {
		p.Name = strings.Split(Split[1], StrfmtPlrStatsNameEnd)[0]
	} else {
		log.Println("failed to parse Name: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsSteamID64)
	if len(Split) >= 2 {
		p.SteamID64, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsSteamID64End)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse SteamID64: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse SteamID64: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsScore)
	if len(Split) >= 2 {
		p.statsScore, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsScoreEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsScore: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsScore: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsDamage)
	if len(Split) >= 2 {
		p.statsDamage, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsDamageEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsDamage: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsDamage: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsKills)
	if len(Split) >= 2 {
		p.statsKills, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsKillsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsKills: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsKills: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsAssists)
	if len(Split) >= 2 {
		p.statsAssists, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsAssistsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsAssists: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsAssists: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsDeaths)
	if len(Split) >= 2 {
		p.statsDeaths, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsDeathsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsDeaths: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsDeaths: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsMVPs)
	if len(Split) >= 2 {
		p.statsMVPs, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsMVPsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsMVPs: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsMVPs: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsPing)
	if len(Split) >= 2 {
		p.statsPing, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsPingEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsPing: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsPing: not found")
		IsError = true
	}
	
	Split = strings.Split(str, StrfmtPlrStatsPenetratedObjects)
	if len(Split) >= 2 {
		p.statsPenetratedObjects, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsPenetratedObjectsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsPenetratedObjects: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsPenetratedObjects: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsHeadShots)
	if len(Split) >= 2 {
		p.statsHeadShots, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsHeadShotsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsHeadShots: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsHeadShots: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsAssistedFlashs)
	if len(Split) >= 2 {
		p.statsAssistedFlashs, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsAssistedFlashsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsAssistedFlashs: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsAssistedFlashs: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsAttackerBlinds)
	if len(Split) >= 2 {
		p.statsAttackerBlinds, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsAttackerBlindsEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsAttackerBlinds: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsAttackerBlinds: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsNoScopes)
	if len(Split) >= 2 {
		p.statsNoScopes, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsNoScopesEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsNoScopes: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsNoScopes: not found")
		IsError = true
	}

	Split = strings.Split(str, StrfmtPlrStatsThroughSmokes)
	if len(Split) >= 2 {
		p.statsThroughSmokes, err = strconv.ParseUint(strings.Split(Split[1], StrfmtPlrStatsThroughSmokesEnd)[0], 10, 64)
		if err != nil {
			log.Println("failed to parse statsThroughSmokes: ", err)
			IsError = true
		}
	} else {
		log.Println("failed to parse statsThroughSmokes: not found")
		IsError = true
	}

	return !IsError
}