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

// Maybe not needed
/*func (p *PlrStats) appendStats(plr *common.Player) {
	if p.SteamID64 == plr.SteamID64 {
		if p.Name == "" {
			p.Name = plr.Name
		}
	
		atomic.AddUint64(&p.statsScore, uint64(plr.Score()))
		atomic.AddUint64(&p.statsDamage, uint64(plr.TotalDamage()))
		atomic.AddUint64(&p.statsKills, uint64(plr.Kills()))
		atomic.AddUint64(&p.statsAssists, uint64(plr.Assists()))
		atomic.AddUint64(&p.statsDeaths, uint64(plr.Deaths()))
		atomic.AddUint64(&p.statsMVPs, uint64(plr.MVPs()))
		atomic.AddUint64(&p.statsPing, uint64(plr.Ping()))
	}
}*/

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

const StrfmtPlrStatsName = "Name='"
const StrfmtPlrStatsNameEnd = "'\t"
const StrfmtPlrStatsSteamID64 = "SteamID64="
const StrfmtPlrStatsSteamID64End = "\t"
const StrfmtPlrStatsScore = "Score='"
const StrfmtPlrStatsScoreEnd = "\t"
const StrfmtPlrStatsDamage = "Damage='"
const StrfmtPlrStatsDamageEnd = "\t"
const StrfmtPlrStatsKills = "Kills='"
const StrfmtPlrStatsKillsEnd = "\t"
const StrfmtPlrStatsAssists = "Assists='"
const StrfmtPlrStatsAssistsEnd = "\t"
const StrfmtPlrStatsDeaths = "Deaths='"
const StrfmtPlrStatsDeathsEnd = "\t"
const StrfmtPlrStatsMVPs = "MVPs='"
const StrfmtPlrStatsMVPsEnd = "\t"
const StrfmtPlrStatsPing = "Ping='"
const StrfmtPlrStatsPingEnd = "\t"
const StrfmtPlrStatsPenetratedObjects = "PenetratedObjects='"
const StrfmtPlrStatsPenetratedObjectsEnd = "\t"
const StrfmtPlrStatsHeadShots = "HeadShots='"
const StrfmtPlrStatsHeadShotsEnd = "\t"
const StrfmtPlrStatsAttackerBlinds = "AttackerBlinds='"
const StrfmtPlrStatsAttackerBlindsEnd = "\t"
const StrfmtPlrStatsAssistedFlashs = "AssistedFlashs='"
const StrfmtPlrStatsAssistedFlashsEnd = "\t"
const StrfmtPlrStatsNoScopes = "NoScopes="
const StrfmtPlrStatsNoScopesEnd = "\t"
const StrfmtPlrStatsThroughSmokes = "ThroughSmokes="
const StrfmtPlrStatsThroughSmokesEnd = "\t"

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

func (p *PlrStats) fromString(str string) {
	var err error

	p.Name = strings.Split(strings.Split(str, StrfmtPlrStatsName)[1], StrfmtPlrStatsNameEnd)[0]
	p.SteamID64, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsSteamID64)[1], StrfmtPlrStatsSteamID64End)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse SteamID64: ", err)
	}
	p.statsScore, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsScore)[1], StrfmtPlrStatsScoreEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsScore: ", err)
	}
	p.statsDamage, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsDamage)[1], StrfmtPlrStatsDamageEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsDamage: ", err)
	}
	p.statsKills, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsKills)[1], StrfmtPlrStatsKillsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsKills: ", err)
	}
	p.statsAssists, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsAssists)[1], StrfmtPlrStatsAssistsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsAssists: ", err)
	}
	p.statsDeaths, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsDeaths)[1], StrfmtPlrStatsDeathsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsDeaths: ", err)
	}
	p.statsMVPs, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsMVPs)[1], StrfmtPlrStatsMVPsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsMVPs: ", err)
	}
	p.statsPing, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsPing)[1], StrfmtPlrStatsPingEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsPing: ", err)
	}
	p.statsPenetratedObjects, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsPenetratedObjects)[1], StrfmtPlrStatsPenetratedObjectsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsPenetratedObjects: ", err)
	}
	p.statsHeadShots, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsHeadShots)[1], StrfmtPlrStatsHeadShotsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsHeadShots: ", err)
	}
	p.statsAssistedFlashs, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsAssistedFlashs)[1], StrfmtPlrStatsAssistedFlashsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsAssistedFlashs: ", err)
	}
	p.statsAttackerBlinds, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsAttackerBlinds)[1], StrfmtPlrStatsAttackerBlindsEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsAttackerBlinds: ", err)
	}
	p.statsNoScopes, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsNoScopes)[1], StrfmtPlrStatsNoScopesEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsNoScopes: ", err)
	}
	p.statsThroughSmokes, err = strconv.ParseUint(strings.Split(strings.Split(str, StrfmtPlrStatsThroughSmokes)[1], StrfmtPlrStatsThroughSmokesEnd)[0], 10, 64)
	if err != nil {
		log.Println("failed to parse statsThroughSmokes: ", err)
	}
}