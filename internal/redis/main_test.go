package redis

import (
	"net"
	"testing"

	"github.com/go-redis/redis"
	"github.com/ory/dockertest"
)

func TestRedisStore(t *testing.T) *Client {
	addr, destroyFunc := startRedis(t)

	defer destroyFunc()

	redis := NewWithAddr(addr, "")

	return redis
}

func startRedis(t *testing.T) (string, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to start Dockertest: %+v", err)
	}

	resource, err := pool.Run("redis", "5-alpine", nil)
	if err != nil {
		t.Fatalf("Failed to start redis: %+v", err)
	}

	// determine the port the container is listening on
	addr := net.JoinHostPort("localhost", resource.GetPort("6379/tcp"))

	// wait for the container to be ready
	err = pool.Retry(func() error {
		var e error
		client := redis.NewClient(&redis.Options{Addr: addr})
		defer client.Close()

		_, e = client.Ping().Result()
		return e
	})

	if err != nil {
		t.Fatalf("Failed to ping Redis: %+v", err)
	}

	destroyFunc := func() {
		pool.Purge(resource)
	}

	return addr, destroyFunc
}
