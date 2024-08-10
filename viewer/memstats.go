package viewer

import (
	"runtime"
	"sync"
)

type statsEntity struct {
	stats runtime.MemStats
	ts    string
}

var (
	statsEntityMut   sync.Mutex
	innerStatsEntity = &statsEntity{}
)

func getStatsEntity() statsEntity {
	statsEntityMut.Lock()
	defer statsEntityMut.Unlock()

	newStats := *innerStatsEntity
	return newStats
}

func updateStatsEntity(stats *statsEntity) {
	statsEntityMut.Lock()
	defer statsEntityMut.Unlock()

	innerStatsEntity = stats
}

func unitMB(n uint64) float64 {
	return float64(n) / 1024 / 1024
}
