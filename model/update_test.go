package model_test

import (
	"github.com/jan-g/girl/model"
	"testing"
)

// Auxiliary functions

func makeCB(record map[model.Facet]int64) func(model.Facet, int64) {
	return func(kind model.Facet, used int64) {
		record[kind] = used
	}
}

func V(kind string, used int64) func(m map[model.Facet]int64) {
	return func(m map[model.Facet]int64) {
		m[model.Facet(kind)] = used
	}
}

func compareResults(t *testing.T, r map[model.Facet]int64, values ...func(m map[model.Facet]int64)) bool {
	m := make(map[model.Facet]int64)
	for _, v := range values {
		v(m)
	}
	return compareHostTrafficMaps(t, r, m)
}

func compareHostTrafficMaps(t *testing.T, r map[model.Facet]int64, m map[model.Facet]int64) bool {
	// Check keys in r
	for k, v := range r {
		v2, ok := m[k]
		if !ok {
			t.Fatalf("Extra key in update: %v = %v", k, v)
		}
		if v != v2 {
			t.Fatalf("Value for key %v doesn't match: wanted %v, got %v", k, v2, v)
		}
	}

	// Check for missing values
	for k, v := range m {
		if _, ok := r[k]; !ok {
			t.Fatalf("Missing key %v: expected %v", k, v)
		}
	}

	return true
}

func U(host string, facet string, total int64) func(map[model.Node]*model.HostTraffic) {
	node := model.Node(host)
	return func(m map[model.Node]*model.HostTraffic) {
		ht, ok := m[node]
		if !ok {
			ht = &model.HostTraffic{
				Name:    host,
				Traffic: []*model.Traffic{},
			}
			m[node] = ht
		}
		ht.Traffic = append(ht.Traffic, &model.Traffic{Facet: facet, Usage: total})
		m[node] = ht
	}
}

func compareTotals(t *testing.T, r map[model.Node]*model.HostTraffic, values ...func(m map[model.Node]*model.HostTraffic)) bool {
	m := make(map[model.Node]*model.HostTraffic)
	for _, v := range values {
		v(m)
	}

	// Check keys in r
	for h, ht := range r {
		ht2, ok := m[h]
		if !ok {
			t.Fatalf("Extra host in totals: %v = %v", h, ht)
		}
		if !compareUsageSets(t, ht, ht2) {
			t.Fatalf("Values for host %v doesn't match: wanted %v, got %v", h, ht2, ht)
		}
	}

	// Check for missing values
	for h, ht := range m {
		if _, ok := r[h]; !ok {
			t.Fatalf("Missing host %v: expected %v", h, ht)
		}
	}

	return true
}

func compareUsageSets(t *testing.T, got *model.HostTraffic, wanted *model.HostTraffic) bool {
	r := summariseHostTraffic(got)
	m := summariseHostTraffic(wanted)
	return compareHostTrafficMaps(t, r, m)
}

func summariseHostTraffic(ht *model.HostTraffic) map[model.Facet]int64 {
	m := make(map[model.Facet]int64)
	for _, l := range ht.Traffic {
		m[model.Facet(l.Facet)] = l.Usage
	}
	return m
}

// Tests follow

func TestNullUpdate(t *testing.T) {
	epochData := make(map[model.Node]*model.HostTraffic)
	traffic := []*model.HostTraffic{}

	results := make(map[model.Facet]int64)

	model.Update(epochData, traffic, makeCB(results))

	compareResults(t, results)
}

func TestNoCBForZeroTraffic(t *testing.T) {
	epochData := map[model.Node]*model.HostTraffic{
		model.Node("hostA"): {
			Name: "hostA",
			Traffic: []*model.Traffic{
				{Facet: "foo", Usage: 0},
			},
		},
	}

	traffic := []*model.HostTraffic{
		{
			Name: "hostA",
			Traffic: []*model.Traffic{
				{Facet: "foo", Usage: 0},
			},
		},
	}

	results := make(map[model.Facet]int64)

	model.Update(epochData, traffic, makeCB(results))

	compareResults(t, results)
}

func TestUpdate(t *testing.T) {
	epochData := map[model.Node]*model.HostTraffic{
		model.Node("hostA"): {
			Name: "hostA",
			Traffic: []*model.Traffic{
				{Facet: "foo", Usage: 2},
				{Facet: "bar", Usage: 4},
			},
		},
		model.Node("hostB"): {
			Name: "hostB",
			Traffic: []*model.Traffic{
				{Facet: "foo", Usage: 1},
				{Facet: "bar", Usage: 3},
			},
		},
	}

	traffic := []*model.HostTraffic{
		epochData[model.Node("hostA")],
		{
			Name: "hostC",
			Traffic: []*model.Traffic{
				{Facet: "foo", Usage: 3},
			},
		},
		{
			Name: "hostD",
			Traffic: []*model.Traffic{
				{Facet: "foo", Usage: 2},
				{Facet: "baz", Usage: 2},
			},
		},
	}

	results := make(map[model.Facet]int64)

	model.Update(epochData, traffic, makeCB(results))

	compareResults(t, results, V("foo", 5), V("baz", 2))
	compareTotals(t, epochData,
		U("hostA", "foo", 2), U("hostA", "bar", 4),
		U("hostB", "foo", 1), U("hostB", "bar", 3),
		U("hostC", "foo", 3),
		U("hostD", "foo", 2), U("hostD", "baz", 2),
	)
}
