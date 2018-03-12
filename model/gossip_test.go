package model_test

import (
	"github.com/jan-g/girl/model"
	"testing"
)

func TestSynchronousGossip(t *testing.T) {
	testGossip(t, 1, 1)
}

func TestEarlierGossip(t *testing.T) {
	testGossip(t, 1, 2)
}

func TestLaterGossip(t *testing.T) {
	testGossip(t, 2, 1)
}

func testGossip(t *testing.T, e1 int64, e2 int64) {
	l1, c1 := model.NewLimiter("x", e1)
	c1.AddLimit("foo", 10, 1, 1)

	s1, ok := c1.(model.LimiterSPI)
	if !ok {
		t.Fatal("Control doesn't expose the SPI")
	}

	l2, c2 := model.NewLimiter("y", e2)
	c2.AddLimit("foo", 10, 1, 1)

	s2, ok := c2.(model.LimiterSPI)
	if !ok {
		t.Fatal("Control doesn't expose the SPI")
	}

	// Ask both for tokens
	if got := l1.Ask("foo", 10); got != 10 {
		t.Errorf("Asked first for ten units of foo, should have 10, got %d", got)
	}

	if got := l2.Ask("foo", 10); got != 10 {
		t.Errorf("Asked second ten units of foo, should have 10, got %d", got)
	}

	if got := l1.Level("foo"); got != 0 {
		t.Errorf("Asked first for foo, should have 0, got %d", got)
	}

	if got := l2.Level("foo"); got != 0 {
		t.Errorf("Asked second for foo, should have 0, got %d", got)
	}

	// Exchange gossip between the two
	gossip(s1, s2)
	gossip(s2, s1)

	s1.Tick(e1 + 1)
	s2.Tick(e2 + 1)

	gossip(s1, s2)
	gossip(s2, s1)

	s1.Tick(e1 + 1)
	s2.Tick(e2 + 1)

	gossip(s1, s2)
	gossip(s2, s1)

	if got := l1.Level("foo"); got != -8 {
		t.Errorf("Asked first for foo, should have -8, got %d", got)
	}

	if got := l2.Level("foo"); got != -8 {
		t.Errorf("Asked second for foo, should have -8, got %d", got)
	}
}

func gossip(from model.LimiterSPI, to model.LimiterSPI) {
	iHave := from.GossipOut()

	response := to.GossipIn(iHave)

	from.ReceivePush(iHave.Epoch, response.Push.Traffic)

	push := from.OriginatePush(iHave.Epoch, response.IWant)

	to.ReceivePush(iHave.Epoch, push)
}
