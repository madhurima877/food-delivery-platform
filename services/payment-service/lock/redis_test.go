package lock

import (
	"context"
	"testing"
)

func TestRedisLock(t *testing.T) {
	t.Log("Redis lock test started")
	ctx := context.Background()
	redisLock := NewRedisLock()
	acquired, err := redisLock.Acquire(ctx, "22")
	if err != nil {
		t.Fatal(err)
	}
	if !acquired {
		t.Error("expected first lock to be acquired")
	}
	acquiredAgain, err := redisLock.Acquire(ctx, "22")
	if err != nil {
		t.Fatal(err)
	}
	if acquiredAgain {
		t.Error("expected second lock attempt to fail")
	}
	t.Log("First acquire:", acquired)
	t.Log("Second acquire:", acquiredAgain)
	err = redisLock.Release(ctx, "22")
	if err != nil {
		t.Fatal(err)
	}
	acquiredAfterRelease, err := redisLock.Acquire(ctx, "22")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("After release", acquiredAfterRelease)
	if !acquiredAfterRelease {
		t.Error("expected lock to be acquired after release")
	}
}
