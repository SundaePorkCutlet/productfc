package kafkamonitor

import (
	"sync"
)

// Monitor — Kafka 컨슈머 처리 건수 (프로세스 메모리, 디버그용).
type Monitor struct {
	mu               sync.RWMutex
	StockUpdatedOK   int64 `json:"stock_updated_ok"`
	StockUpdatedDup  int64 `json:"stock_updated_duplicate_skipped"`
	StockUpdatedDLQ  int64 `json:"stock_updated_dlq"`
	RollbackOK       int64 `json:"stock_rollback_ok"`
	RollbackDup      int64 `json:"stock_rollback_duplicate_skipped"`
	RollbackDLQ      int64 `json:"stock_rollback_dlq"`
	UnmarshalErrors  int64 `json:"unmarshal_errors"`
	SchemaRejected   int64 `json:"schema_version_rejected"`
}

func NewMonitor() *Monitor {
	return &Monitor{}
}

func (m *Monitor) IncStockUpdatedOK() {
	m.mu.Lock()
	m.StockUpdatedOK++
	m.mu.Unlock()
}

func (m *Monitor) IncStockUpdatedDup() {
	m.mu.Lock()
	m.StockUpdatedDup++
	m.mu.Unlock()
}

func (m *Monitor) IncStockUpdatedDLQ() {
	m.mu.Lock()
	m.StockUpdatedDLQ++
	m.mu.Unlock()
}

func (m *Monitor) IncRollbackOK() {
	m.mu.Lock()
	m.RollbackOK++
	m.mu.Unlock()
}

func (m *Monitor) IncRollbackDup() {
	m.mu.Lock()
	m.RollbackDup++
	m.mu.Unlock()
}

func (m *Monitor) IncRollbackDLQ() {
	m.mu.Lock()
	m.RollbackDLQ++
	m.mu.Unlock()
}

func (m *Monitor) IncUnmarshalErr() {
	m.mu.Lock()
	m.UnmarshalErrors++
	m.mu.Unlock()
}

func (m *Monitor) IncSchemaRejected() {
	m.mu.Lock()
	m.SchemaRejected++
	m.mu.Unlock()
}

func (m *Monitor) Snapshot() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return map[string]int64{
		"stock_updated_ok":                 m.StockUpdatedOK,
		"stock_updated_duplicate_skipped":    m.StockUpdatedDup,
		"stock_updated_dlq":                m.StockUpdatedDLQ,
		"stock_rollback_ok":                m.RollbackOK,
		"stock_rollback_duplicate_skipped": m.RollbackDup,
		"stock_rollback_dlq":               m.RollbackDLQ,
		"unmarshal_errors":                 m.UnmarshalErrors,
		"schema_version_rejected":          m.SchemaRejected,
	}
}
