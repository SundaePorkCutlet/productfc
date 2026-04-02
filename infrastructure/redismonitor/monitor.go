package redismonitor

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisStats struct {
	Hits     int64     `json:"hits"`
	Misses   int64     `json:"misses"`
	HitRate  float64   `json:"hit_rate_pct"`
	TotalOps int64     `json:"total_ops"`
	Errors   int64     `json:"errors"`
	Keys     []KeyInfo `json:"keys,omitempty"`
}

type KeyInfo struct {
	Pattern string `json:"pattern"`
	Count   int64  `json:"count"`
}

type Monitor struct {
	mu       sync.RWMutex
	hits     int64
	misses   int64
	totalOps int64
	errors   int64
	redis    *redis.Client
}

func NewMonitor(client *redis.Client) *Monitor {
	return &Monitor{
		redis: client,
	}
}

func (m *Monitor) RecordHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits++
	m.totalOps++
}

func (m *Monitor) RecordMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.misses++
	m.totalOps++
}

func (m *Monitor) RecordOp(operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalOps++
}

func (m *Monitor) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors++
	m.totalOps++
}

func (m *Monitor) GetDebugInfo(ctx context.Context) RedisStats {
	m.mu.RLock()
	stats := RedisStats{
		Hits:     m.hits,
		Misses:   m.misses,
		TotalOps: m.totalOps,
		Errors:   m.errors,
	}
	m.mu.RUnlock()

	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total) * 100
	}

	dbSize, err := m.redis.DBSize(ctx).Result()
	if err == nil {
		stats.Keys = []KeyInfo{
			{Pattern: "total_keys", Count: dbSize},
		}
	}

	return stats
}
