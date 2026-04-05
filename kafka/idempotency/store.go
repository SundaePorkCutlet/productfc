package idempotency

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const doneTTL = 7 * 24 * time.Hour

// Store — 동일 order_id + 토픽에 대한 중복 처리 방지 (at-least-once 소비 시).
type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{rdb: rdb}
}

func doneKey(topic string, orderID int64) string {
	return fmt.Sprintf("kafka:done:%s:%d", topic, orderID)
}

// AlreadyProcessed — 성공 처리 완료 키가 있으면 true (중복 메시지).
func (s *Store) AlreadyProcessed(ctx context.Context, topic string, orderID int64) (bool, error) {
	n, err := s.rdb.Exists(ctx, doneKey(topic, orderID)).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// MarkProcessed — 모든 비즈니스 처리 성공 후 호출.
func (s *Store) MarkProcessed(ctx context.Context, topic string, orderID int64) error {
	return s.rdb.Set(ctx, doneKey(topic, orderID), "1", doneTTL).Err()
}
