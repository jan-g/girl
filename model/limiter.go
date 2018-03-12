package model

import (
	"sync"

	"github.com/sirupsen/logrus"
)

// User-facing interface
type Limiter interface {
	// We want a particular amount of a limit. Tell us how much we can have.
	Ask(limit string, want int64) int64

	// Get the current level
	Level(limit string) int64
}

// Interface for managing new limits
type LimiterControl interface {
	AddLimit(limit string, upper int64, additional int64, every int64)
}

// Service-facing interface
type LimiterSPI interface {
	// Report the current epoch
	Epoch() int64

	Tick(newEpoch int64)

	// Prepare to gossip with a node
	GossipOut() *IHave

	// Respond to initial gossip request from a peer
	GossipIn(ih *IHave) *ResponderHandshake

	// Respond to a push update, either side
	ReceivePush(epoch int64, traffic []*HostTraffic)

	// Send a requested set of host details
	OriginatePush(epoch int64, wanted []string) []*HostTraffic
}

type limiter interface {
	Limiter
	LimiterControl
	LimiterSPI
}

type Node string
type Facet string

type epochLimiter struct {
	// For access to these top-level pieces. Just start with a BGL.
	mx    sync.RWMutex
	epoch int64

	name Node

	limits map[Facet]*bucket

	// This epoch's traffic, organised by limit label
	traffic map[Facet]int64

	// The previous epoch's gossiped traffic details.
	// Node -> []{limit, traffic}
	lastEpoch map[Node]*HostTraffic

	// The epoch *before that*. Nodes that are still in the previous epoch will gossip this.
	oldEpoch map[Node]*HostTraffic

	// The current epoch. Nodes that are ahead of us may gossip this one.
	// The current Node is not added to this until it Tick()s over
	currentEpoch map[Node]*HostTraffic
}

type bucket struct {
	mx sync.RWMutex

	// Traffic levels
	curr point
	last point
	old  point

	// Bucket parameters
	max   int64
	add   int64
	every int64
}

type point struct {
	// This gives the level of the bucket
	level int64

	// This records how much spare capacity could have been added to the
	// bucket at the *end* of this epoch. When rolling forward usage from the
	// past, we can offset that with anything that remains here.
	available int64
}

var _ limiter = &epochLimiter{}

func NewLimiter(name string, epoch int64) (Limiter, LimiterControl) {
	l := &epochLimiter{
		name:         Node(name),
		epoch:        epoch,
		limits:       make(map[Facet]*bucket),
		traffic:      make(map[Facet]int64),
		lastEpoch:    make(map[Node]*HostTraffic),
		oldEpoch:     make(map[Node]*HostTraffic),
		currentEpoch: make(map[Node]*HostTraffic),
	}
	return l, l
}

func (l *epochLimiter) Ask(limit string, want int64) int64 {
	l.mx.RLock()
	defer l.mx.RUnlock()

	b := l.limits[Facet(limit)]
	if b == nil {
		// Ya cen't cam in.
		return 0
	}

	b.mx.Lock()
	defer b.mx.Unlock()

	var rcv int64 = 0
	if b.curr.level >= want {
		rcv = want
	} else if b.curr.level > 0 {
		rcv = b.curr.level
	}

	b.curr.level -= rcv

	// record traffic
	l.traffic[Facet(limit)] += rcv

	return rcv
}

func (l *epochLimiter) Level(limit string) int64 {
	l.mx.RLock()
	defer l.mx.RUnlock()

	b := l.limits[Facet(limit)]
	if b == nil {
		// Ya cen't cam in.
		return 0
	}

	b.mx.RLock()
	defer b.mx.RUnlock()

	return b.curr.level
}

func (l *epochLimiter) AddLimit(limit string, upper int64, additional int64, every int64) {
	l.mx.Lock()
	defer l.mx.Unlock()

	available := int64(0)
	if l.epoch%every == 0 {
		available = additional
	}
	l.limits[Facet(limit)] = &bucket{
		curr: point{level: upper, available: available},

		max:   upper,
		add:   additional,
		every: every,
	}
}

func (l *epochLimiter) Epoch() int64 {
	return l.epoch
}

func (l *epochLimiter) Tick(newEpoch int64) {
	l.mx.Lock()
	defer l.mx.Unlock()

	for _, b := range l.limits {
		b.mx.Lock()
		b.old = b.last
		b.last = b.curr
		if newEpoch%b.every == 0 {
			b.curr.available = b.add
		} else {
			b.curr.available = 0
		}
		b.curr.level += b.last.available
		if b.curr.level > b.max {
			b.last.available = b.curr.level - b.max
			b.curr.level = b.max
		} else {
			b.last.available = 0
		}
		b.mx.Unlock()
	}

	l.epoch = newEpoch

	// Shift epoch records around
	// We may still receive gossip for this from nodes we are ahead of
	l.oldEpoch = l.lastEpoch

	// The current epoch. Add a record for ourselves.
	l.lastEpoch = l.currentEpoch

	ht := &HostTraffic{}
	ht.Name = string(l.name)
	traffic := []*Traffic{}
	for facet, use := range l.traffic {
		if use != 0 {
			traffic = append(traffic, &Traffic{
				Facet: string(facet),
				Usage: use,
			})
		}
	}
	ht.Traffic = traffic
	l.lastEpoch[l.name] = ht

	// The next epoch. We may receive gossip for this.
	l.currentEpoch = make(map[Node]*HostTraffic)

	// Zero traffic counters for this new epoch
	l.traffic = make(map[Facet]int64)
}

// Return the first part of a gossip handshake.
// We are only interested in gossiping about the previously-commited epoch,
// although we will respond to alternative epochs to help out neighbours
func (l *epochLimiter) GossipOut() *IHave {
	l.mx.RLock()
	defer l.mx.RUnlock()

	hosts := []string{}
	for host := range l.lastEpoch {
		hosts = append(hosts, string(host))
	}
	ih := &IHave{
		Epoch: l.epoch - 1,
		Hosts: hosts,
	}
	return ih
}

// Receive the first part of a handshake
func (l *epochLimiter) GossipIn(ih *IHave) *ResponderHandshake {
	l.mx.Lock()
	defer l.mx.Unlock()

	// Which epoch are we concerned with?
	if ih.Epoch == l.epoch-1 {
		// We are synchronised. Gossip about the previous period.
		iWant, push := Respond(l.lastEpoch, ih.Hosts)
		return &ResponderHandshake{
			IWant: iWant,
			Push:  push,
		}
	} else if ih.Epoch == l.epoch {
		// They are ahead. Exchange gossip about the currentEpoch, but
		// don't account for it yet with limits
		iWant, push := Respond(l.currentEpoch, ih.Hosts)
		return &ResponderHandshake{
			IWant: iWant,
			Push:  push,
		}
	} else if ih.Epoch == l.epoch-2 {
		// They are behind us. Exchange gossip about the older period, so they
		// can catch up.
		iWant, push := Respond(l.oldEpoch, ih.Hosts)
		return &ResponderHandshake{
			IWant: iWant,
			Push:  push,
		}
	} else {
		logrus.
			WithField("my.epoch", l.epoch).
			WithField("their.epoch", ih.Epoch).
			Warning("Gossip from node out of sync")
	}
	return &ResponderHandshake{}
}

func Respond(table map[Node]*HostTraffic, theyHave []string) ([]string, *Push) {
	iWant := []string{}
	traffic := []*HostTraffic{}
	theySee := map[string]bool{}

	for _, node := range theyHave {
		_, ok := table[Node(node)]
		if !ok {
			iWant = append(iWant, node)
		}
		theySee[node] = true
	}

	for node, usage := range table {
		if !theySee[string(node)] {
			traffic = append(traffic, usage)
		}
	}

	return iWant, &Push{Traffic: traffic}
}

// Receive the first part of a handshake
func (l *epochLimiter) ReceivePush(epoch int64, traffic []*HostTraffic) {
	l.mx.Lock()
	defer l.mx.Unlock()

	// Which epoch are we concerned with?
	if epoch == l.epoch-1 {
		// We are synchronised. Gossip about the previous period.
		Update(l.lastEpoch, traffic, func(kind Facet, used int64) {
			// We simply adjust buckets here, but take into account any available
			// capacity left over at the end of the last epoch
			bucket := l.limits[kind]
			bucket.last.level -= used
			// At the end of the last epoch, we may have some available bonus left over.
			if bucket.last.available > 0 && used <= bucket.last.available {
				bucket.last.available -= used
				used = 0 // This has no effect on current levels
			} else if bucket.last.available > 0 {
				used -= bucket.last.available
				bucket.last.available = 0
			}
			bucket.curr.level -= used
		})
	} else if epoch == l.epoch {
		// They are ahead. Exchange gossip about the currentEpoch, but
		// don't account for it yet with limits
		Update(l.currentEpoch, traffic, func(kind Facet, used int64) {
			// We are learning about this ahead of time, but adjust our limits
			l.limits[kind].curr.level -= used
		})
	} else if epoch == l.epoch-2 {
		// They are behind us. Exchange gossip about the older period, so they
		// can catch up. We can process it too.
		Update(l.oldEpoch, traffic, func(kind Facet, used int64) {
			// We simply adjust buckets here, but take into account any available
			// capacity left over at the end of the last epoch and the one before that
			bucket := l.limits[kind]
			bucket.old.level -= used
			// At the end of the old epoch, we may have some available bonus left over.
			if bucket.old.available > 0 && used <= bucket.old.available {
				bucket.old.available -= used
				used = 0 // This has no effect on later levels
			} else if bucket.old.available > 0 {
				used -= bucket.old.available
				bucket.old.available = 0
			}
			bucket.last.level -= used
			// At the end of the last epoch, we may have some available bonus left over.
			if bucket.last.available > 0 && used <= bucket.last.available {
				bucket.last.available -= used
				used = 0 // This has no effect on current levels
			} else if bucket.last.available > 0 {
				used -= bucket.last.available
				bucket.last.available = 0
			}
			bucket.curr.level -= used
		})
	} else {
		logrus.
			WithField("my.epoch", l.epoch).
			WithField("their.epoch", epoch).
			Warning("Push from node out of sync")
	}
}

// Run an update.
// For each node in the incoming traffic report we do not yet know about, update
// the traffic data for that node.
// Assemble the combined traffic data updates for each Facet so named, and then
// call the callback function which knows how to process those.
func Update(epochData map[Node]*HostTraffic, traffic []*HostTraffic, cb func(kind Facet, used int64)) {
	totalTraffic := make(map[Facet]int64)
	for _, usage := range traffic {
		node := Node(usage.Name)
		if _, ok := epochData[node]; !ok {
			// We have a new node we're learning the traffic data for
			epochData[node] = usage
			for _, limit := range usage.Traffic {
				totalTraffic[Facet(limit.Facet)] += limit.Usage
			}
		}
	}
	for facet, usage := range totalTraffic {
		if usage != 0 {
			cb(facet, usage)
		}
	}
}

// Return a requested set of Node traffic information
func (l *epochLimiter) OriginatePush(epoch int64, wanted []string) []*HostTraffic {
	l.mx.RLock()
	defer l.mx.RUnlock()

	// Which epoch are we concerned with?
	if epoch == l.epoch-1 {
		// We are synchronised. Gossip about the previous period.
		return Wanted(l.lastEpoch, wanted)
	} else if epoch == l.epoch {
		// They are ahead. Exchange gossip about the currentEpoch, but
		// don't account for it yet with limits
		return Wanted(l.currentEpoch, wanted)
	} else if epoch == l.epoch-2 {
		// They are behind us. Exchange gossip about the older period, so they
		// can catch up.
		return Wanted(l.oldEpoch, wanted)
	} else {
		logrus.
			WithField("my.epoch", l.epoch).
			WithField("their.epoch", epoch).
			Warning("Gossip request from node out of sync")
	}
	return []*HostTraffic{}
}

func Wanted(epochData map[Node]*HostTraffic, wanted []string) []*HostTraffic {
	reply := []*HostTraffic{}
	for _, host := range wanted {
		if ht, ok := epochData[Node(host)]; ok {
			reply = append(reply, ht)
		}
	}
	return reply
}
