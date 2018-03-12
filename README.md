# Gossiping Interval Rate Limiting

This is a little experiment in distributed rate-limiting. It uses a standard
token bucket, which is adjusted to account for peers' recent traffic. Those
adjustments can take a count negative (this doesn't happen in the standard
single-node case).

## Requirements

There are some requirements for the current implementation to work. The
major one is that nodes' clocks should be synchronised to within around one
gossip tick (this is well within typical NTP performance) - there need to be
a sufficient number of rounds of gossip to propagate traffic information across
the node set within each interval.

Although the protocol shouldn't go massively adrift as long as servers are
within an epoch of each other, it will at the least take longer to converge.

This lack of convergence opens a cluster up to a greater burst of traffic.

## The gossip protocol

This is a straightforward push-pull implementation; each node shares a map of
_node_ -> _traffic_ history for its most recently-committed epoch. On receipt
of a gossip request, a node will respond with traffic for the corresponding
epoch, accounting for this as it does; however, if that traffic is for a
different epoch then it will take no proactive effort to gossip that.

Unlike typical gossip protocols that use Lamport clocks, this protocol uses
a timestamp (with epoch-level resolution) to coordinate. (It does not seem
unreasonable to require that nodes in a cluster have closely-synchronised time.)

## Epochs and traffic bursts

Each epoch lasts for a second. Within that second, a node collects traffic
history for each bucket that it allocates tokens to. It will not allocate more
tokens than are in that bucket; however, it will not correct the bucket's
capacity until its current epoch ticks over (at which point it will add a
summary of that traffic to the gossiped content of the last-committed epoch).

Consequently, it is possible for a greedy client to request _bucket max_ tokens
from _each node_, thereby potentially exceeding the single-node maximum burst
by that factor. Under such circumstances, the nodes will gossip on the following
epoch, correcting their bucket values - which will therefore become negative.
Under this situation, it will take longer for the buckets to regain a positive
capacity.

The upshot of all that is that bursts can be steeper than configured, but the
average throughput of tokens remains as configured, in the long run.

## In case of net splits

Each partition will deliver tokens at the same, total rate. If splits last more
than an epoch (this is very likely) then the protocol has no chance at all to
recover those totals, although a rejoined network will reconverge and begin
delivering tokens at the global rate again.


#### TODO

- Performance measurements
- Smarter server framework (eg, pool connections for a small peer set)
- Anti-entropy

  Adding new limits may build on this.
- At the moment there's a BGL; profile and push this down
- Integrate with, say, libserf for membership management
- Extend to cope with desynchronised nodes.

  It's in priciple possible to cope with nodes whose clocks are further skewed.
  The accounting trick remains the same; the gossip protocol needs to offer
  more node data for exchange to do this.

#### Credits

The idea for this is inspired by a blog post describing Yahoo's "cloudbouncer";
there aren't many details in that post - this is just a quick take on one way
that might work. Any misunderstandings are entirely my own.

https://yahooeng.tumblr.com/post/111288877956/cloud-bouncer-distributed-rate-limiting-at-yahoo
