package storage

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

type redisBreakerState struct {
    State      BreakerState `json:"state"`
    RetryAfter int64        `json:"retry_after_unix"`
    Reason     string       `json:"reason"`
}

// RedisDistributedStateStore implements DistributedStateStore using Redis.
type RedisDistributedStateStore struct {
    client *redis.Client
}

func NewRedisDistributedStateStore(client *redis.Client) *RedisDistributedStateStore {
    return &RedisDistributedStateStore{client: client}
}

func (r *RedisDistributedStateStore) GetBreakerState(
    ctx context.Context,
    key string,
) (BreakerState, time.Time, string, error) {
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return BreakerClosed, time.Time{}, "", nil
    }
    if err != nil {
        return BreakerClosed, time.Time{}, "", err
    }

    var s redisBreakerState
    // Using JSON is acceptable here as this is only called in the recovery/discovery path (OPEN state), 
    // never on the Titanium happy path.
    if err := json.Unmarshal([]byte(val), &s); err != nil {
        return BreakerClosed, time.Time{}, "", err
    }

    return s.State, time.Unix(s.RetryAfter, 0), s.Reason, nil
}

func (r *RedisDistributedStateStore) SetBreakerState(
    ctx context.Context,
    key string,
    state BreakerState,
    retryAfter time.Time,
    reason string,
	ttl time.Duration,
) error {
    s := redisBreakerState{
        State:      state,
        RetryAfter: retryAfter.Unix(),
        Reason:     reason,
    }
    data, err := json.Marshal(s)
    if err != nil {
        return err
    }
    return r.client.Set(ctx, key, data, ttl).Err()
}
