package main

import (
	"sync"
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

var PlrMutex = &sync.Mutex{}

func (p *PlrStats) appendStats(plr *common.Player) {
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
