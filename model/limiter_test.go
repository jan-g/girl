package model_test

import (
	"github.com/jan-g/girl/model"
	"testing"
)

func TestNewLimiter(t *testing.T) {
	limiter, control := model.NewLimiter("x", 1)
	control.AddLimit("foo", 5, 5, 1)

	spi, ok := control.(model.LimiterSPI)
	if !ok {
		t.Fatal("Control doesn't expose the SPI")
	}

	if got := limiter.Ask("foo", 5); got != 5 {
		t.Fatalf("Asked for five units of foo, should have 5, got %d", got)
	}

	if got := limiter.Ask("foo", 10); got != 0 {
		t.Fatalf("Asked for 10 units of foo, should have 0, got %d", got)
	}

	spi.Tick(2)

	if got := limiter.Ask("foo", 10); got != 5 {
		t.Fatalf("Asked for 10 units of foo, should have 5, got %d", got)
	}

	spi.Tick(3)

	if got := limiter.Ask("foo", 1); got != 1 {
		t.Fatalf("Asked for 1 unit of foo, should have 1, got %d", got)
	}

	spi.Tick(4)

	if got := limiter.Ask("foo", 6); got != 5 {
		t.Fatalf("Asked for 6 unit of foo, should have 5, got %d", got)
	}
}
