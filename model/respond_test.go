package model_test

import (
	"github.com/jan-g/girl/model"
	"sort"
	"testing"
)

func TestRespond(t *testing.T) {
	table := map[model.Node]*model.HostTraffic{
		model.Node("hostA"): {
			Name: "hostA",
			Traffic: []*model.Traffic{
				{
					Facet: "foo",
					Usage: 1,
				},
				{
					Facet: "bar",
					Usage: 10,
				},
				{
					Facet: "baz",
					Usage: 0,
				},
			},
		},
		model.Node("hostC"): {
			Name: "hostC",
			Traffic: []*model.Traffic{
				{
					Facet: "foo",
					Usage: 3,
				},
				{
					Facet: "bar",
					Usage: 0,
				},
				{
					Facet: "baz",
					Usage: 2,
				},
			},
		},
	}
	nodes := []string{"hostA", "hostB"}
	iWant, push := model.Respond(table, nodes)

	sort.Strings(iWant)
	if len(iWant) != 1 || iWant[0] != "hostB" {
		t.Fatalf("Should only ask for hostB")
	}

	hereAre := []string{}
	for _, ht := range push.Traffic {
		hereAre = append(hereAre, ht.Name)
	}
	sort.Strings(hereAre)

	if len(hereAre) != 1 || hereAre[0] != "hostC" {
		t.Fatalf("Should only offer hostC")
	}
}

func TestRespondTheyHaveNothing(t *testing.T) {
	table := map[model.Node]*model.HostTraffic{
		model.Node("hostA"): {
			Name: "hostA",
			Traffic: []*model.Traffic{
				{
					Facet: "foo",
					Usage: 1,
				},
				{
					Facet: "bar",
					Usage: 10,
				},
				{
					Facet: "baz",
					Usage: 0,
				},
			},
		},
		model.Node("hostC"): {
			Name: "hostC",
			Traffic: []*model.Traffic{
				{
					Facet: "foo",
					Usage: 3,
				},
				{
					Facet: "bar",
					Usage: 0,
				},
				{
					Facet: "baz",
					Usage: 2,
				},
			},
		},
	}
	nodes := []string{}
	iWant, push := model.Respond(table, nodes)

	sort.Strings(iWant)
	if len(iWant) != 0 {
		t.Fatalf("Should not ask for data")
	}

	hereAre := []string{}
	for _, ht := range push.Traffic {
		hereAre = append(hereAre, ht.Name)
	}
	sort.Strings(hereAre)

	if len(hereAre) != 2 || hereAre[0] != "hostA" || hereAre[1] != "hostC" {
		t.Fatalf("Should offer hostA and hostC")
	}
}

func TestRespondTheyHaveEverything(t *testing.T) {
	table := map[model.Node]*model.HostTraffic{}
	nodes := []string{"hostA", "hostB", "hostC"}
	iWant, push := model.Respond(table, nodes)

	sort.Strings(iWant)
	if len(iWant) != 3 || iWant[0] != "hostA" || iWant[1] != "hostB" || iWant[2] != "hostC" {
		t.Fatalf("Should ask for everything")
	}

	hereAre := []string{}
	for _, ht := range push.Traffic {
		hereAre = append(hereAre, ht.Name)
	}
	sort.Strings(hereAre)

	if len(hereAre) != 0 {
		t.Fatalf("Should offer nothing")
	}
}
