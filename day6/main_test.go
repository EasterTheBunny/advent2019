package main

import (
	"testing"
)

func TestCreateOrbitList(t *testing.T) {
	input := []string{"COM)B", "B)C", "C)D"}
	expected := []Orbit{
		Orbit{Center: "COM", Satellite: "B"},
		Orbit{Center: "B", Satellite: "C"},
		Orbit{Center: "C", Satellite: "D"}}

	result := CreateOrbitList(input)
	if len(result) != len(expected) {
		t.Errorf("incorrect result length %v; expected %v", len(result), len(expected))
		return
	}

	for i, orb := range result {
		if orb.Center != expected[i].Center {
			t.Errorf("incorrect center %s; expected %s", orb.Center, expected[i].Center)
		}

		if orb.Satellite != expected[i].Satellite {
			t.Errorf("incorrect center %s; expected %s", orb.Satellite, expected[i].Satellite)
		}
	}
}

func TestLinkOrbits(t *testing.T) {
	input := []string{"COM)B", "COM)C", "B)D"}

	d := Satellite{Name: "D", InOrbit: []*Satellite{}}
	c := Satellite{Name: "C", InOrbit: []*Satellite{}}
	b := Satellite{Name: "B", InOrbit: []*Satellite{&d}}
	com := Satellite{Name: "COM", InOrbit: []*Satellite{&b, &c}}

	expected := &com

	list := CreateOrbitList(input)
	_, result := LinkOrbits(list)

	if result == nil {
		t.Error("incorrect nil result")
		return
	}

	if result.Name != expected.Name {
		t.Errorf("incorrect name %s; expected %s", result.Name, expected.Name)
	}

	if len(result.InOrbit) != len(expected.InOrbit) {
		t.Errorf("incorrect number of satellites %v; expected %v", len(result.InOrbit), len(expected.InOrbit))
	}

	if len(result.OrbitAt(0).InOrbit) != len(expected.OrbitAt(0).InOrbit) {
		t.Errorf("incorrect second layer orbit %v; expected %v", len(result.OrbitAt(0).InOrbit), len(expected.OrbitAt(0).InOrbit))
	}
}

func TestDepth(t *testing.T) {
	input := []string{"COM)B", "B)C", "C)D", "D)E", "E)F", "B)G", "G)H", "D)I", "E)J", "J)K", "K)L"}

	expected := 42

	list := CreateOrbitList(input)
	_, center := LinkOrbits(list)

	if center == nil {
		t.Error("incorrect nil result")
		return
	}

	result := center.Depth(0)

	if result != expected {
		t.Errorf("incorrect depth %v; expected %v", result, expected)
	}
}

func TestPathToCenter(t *testing.T) {
	input := []string{"COM)B", "B)C", "C)D", "D)E", "E)F", "B)G", "G)H", "D)I", "E)J", "J)K", "K)L", "K)YOU", "I)SAN"}
	expected := []string{"K", "J", "E", "D", "C", "B", "COM"}

	list := CreateOrbitList(input)
	orbitMap, _ := LinkOrbits(list)

	you, ok := orbitMap["YOU"]
	if !ok {
		t.Errorf("missing satellite %s in map", "YOU")
		return
	}

	path := you.PathToCenter()

	if len(path) != len(expected) {
		t.Errorf("incorrect path length %v; expected %v", len(path), len(expected))
		return
	}

	for i, point := range path {
		if point.Name != expected[i] {
			t.Errorf("incorrect point %s at index %v; expected %s", point.Name, i, expected[i])
		}
	}
}

func TestExitOrbit(t *testing.T) {
	input := []string{"COM)B", "COM)C", "B)D"}

	d := Satellite{Name: "D", InOrbit: []*Satellite{}}
	c := Satellite{Name: "C", InOrbit: []*Satellite{}}
	refD := &d
	refC := &c
	b := Satellite{Name: "B", InOrbit: []*Satellite{refD}}
	refB := &b
	com := Satellite{Name: "COM", InOrbit: []*Satellite{refB, refC}}
	refCom := &com

	refD.Orbits = refB
	refC.Orbits = refCom
	refB.Orbits = refCom

	list := CreateOrbitList(input)
	LinkOrbits(list)

	refC.ExitOrbit()

	if refC.Orbits != nil {
		t.Errorf("%s still orbiting %s; expected nil", refC.Name, refC.Orbits.Name)
	}

	found := false
	for _, or := range refCom.InOrbit {
		if or.Name == refC.Name {
			found = true
		}
	}

	if found {
		t.Errorf("satellite %s found orbiting %s; expected nil", refC.Name, refCom.Name)
	}
}

func TestEnterOrbit(t *testing.T) {
	input := []string{"COM)B", "COM)C", "B)D"}

	d := Satellite{Name: "D", InOrbit: []*Satellite{}}
	c := Satellite{Name: "C", InOrbit: []*Satellite{}}
	refD := &d
	refC := &c
	b := Satellite{Name: "B", InOrbit: []*Satellite{refD}}
	refB := &b
	com := Satellite{Name: "COM", InOrbit: []*Satellite{refB, refC}}
	refCom := &com

	refD.Orbits = refB
	refC.Orbits = refCom
	refB.Orbits = refCom

	list := CreateOrbitList(input)
	LinkOrbits(list)

	refC.EnterOrbit(refD)

	if refC.Orbits == nil {
		t.Errorf("%s not orbiting anything; expected %s", refC.Name, refD.Name)
	}

	found := false
	for _, or := range refD.InOrbit {
		if or.Name == refC.Name {
			found = true
		}
	}

	if !found {
		t.Errorf("satellite %s not found orbiting %s; expected to be orbiting", refC.Name, refD.Name)
	}

}

func TestTransferTo(t *testing.T) {
	input := []string{"COM)B", "B)C", "C)D", "D)E", "E)F", "B)G", "G)H", "D)I", "E)J", "J)K", "K)L", "K)YOU", "I)SAN"}

	list := CreateOrbitList(input)
	orbitMap, center := LinkOrbits(list)

	if center == nil {
		t.Error("incorrect nil result")
		return
	}

	s, ok := orbitMap["I"]
	if !ok {
		t.Errorf("missing satellite %s in map", "I")
		return
	}

	you, ok := orbitMap["YOU"]
	if !ok {
		t.Errorf("missing satellite %s in map", "YOU")
		return
	}

	san, ok := orbitMap["SAN"]
	if !ok {
		t.Errorf("missing satellite %s in map", "SAN")
		return
	}

	jumps := you.TransferTo(san)

	if len(s.InOrbit) != 2 {
		t.Errorf("incorrect number of satellites %v; expected %v", len(s.InOrbit), 2)
		return
	}

	if you.Orbits.Name != "I" {
		t.Errorf("incorrect orbit %s; expected %s", you.Orbits.Name, "I")
	}

	if san.Orbits.Name != "I" {
		t.Errorf("incorrect orbit %s; expected %s", san.Orbits.Name, "I")
	}

	if jumps != 4 {
		t.Errorf("incorrect number of orbital transfers %v; expected %v", jumps, 4)
	}
}
