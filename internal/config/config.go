package config

import (
	"os"
	"runtime"
	"strconv"
	"time"
)

var Config = struct {
	WorkerPoolSize   int64
	AllocThreshold   uint64
	ReadMemStatsRate time.Duration
}{}

func init() {
	Config.WorkerPoolSize = int64(parseEnvInt("WORKER_POOL_SIZE", runtime.NumCPU()))
	Config.AllocThreshold = uint64(parseEnvInt("ALLOC_THRESHOLD", 256e8))
	Config.ReadMemStatsRate = parseEnvDuration("READ_MEM_STATS_RATE", 100*time.Millisecond)
}

func parseEnvInt(key string, def int) int {
	val, ok := os.LookupEnv(key)

	if !ok {
		return def
	}

	i, err := strconv.Atoi(val)

	if err != nil {
		return def
	}

	return i
}

func parseEnvDuration(key string, def time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)

	if !ok {
		return def
	}

	d, err := time.ParseDuration(val)

	if err != nil {
		return def
	}

	return d
}
