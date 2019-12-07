package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Orbit ...
type Orbit struct {
	Center    string
	Satellite string
}

// Satellite ...
type Satellite struct {
	Name    string
	InOrbit []*Satellite
	Orbits  *Satellite
}

// Depth ...
func (s *Satellite) Depth(i int) int {
	depth := i
	i++
	for _, or := range s.InOrbit {
		depth += or.Depth(i)
	}
	return depth
}

// OrbitAt ...
func (s *Satellite) OrbitAt(position int) *Satellite {
	if position >= 0 && position < len(s.InOrbit) {
		return s.InOrbit[position]
	}
	return nil
}

// PathToCenter ...
func (s *Satellite) PathToCenter() []*Satellite {
	path := []*Satellite{}
	nextPoint := s
	for ok := true; ok; ok = nextPoint.Orbits != nil {
		path = append(path, nextPoint.Orbits)
		if nextPoint.Orbits.Name == "COM" {
			break
		}
		nextPoint = nextPoint.Orbits
	}
	return path
}

// ExitOrbit ...
func (s *Satellite) ExitOrbit() {
	if s.Orbits == nil {
		return
	}

	orbiting := s.Orbits
	for i, o := range orbiting.InOrbit {
		if o.Name == s.Name {
			if len(orbiting.InOrbit) == 1 {
				orbiting.InOrbit = []*Satellite{}
			} else {
				orbiting.InOrbit = append(orbiting.InOrbit[:i], orbiting.InOrbit[i+1:]...)
			}
			break
		}
	}
	s.Orbits = nil
}

// EnterOrbit ...
func (s *Satellite) EnterOrbit(sat *Satellite) {
	s.Orbits = sat
	sat.InOrbit = append(sat.InOrbit, s)
}

// TransferTo ...
func (s *Satellite) TransferTo(destination *Satellite) int {
	jumps := 0

	if s.Orbits != nil && s.Orbits.Name != destination.Orbits.Name {
		pathA := s.PathToCenter()
		pathB := destination.PathToCenter()
		var common *Satellite

		lenA := 0
		lenB := 0

		for i, a := range pathA {
			lenA = i
			for j, b := range pathB {
				lenB = j
				if a.Name == b.Name {
					common = a
					break
				}
			}
			if common != nil {
				break
			}
		}

		if common != nil {
			jumps = lenA + lenB
			s.ExitOrbit()
			s.EnterOrbit(destination.Orbits)
		}
	}

	return jumps
}

// LinkOrbits ...
func LinkOrbits(list []Orbit) (map[string]*Satellite, *Satellite) {
	lookup := make(map[string]*Satellite)
	var center *Satellite

	for _, o := range list {
		_, ok := lookup[o.Center]
		if !ok {
			new := Satellite{Name: o.Center, InOrbit: []*Satellite{}}
			n := &new
			if o.Center == "COM" {
				center = n
			}
			lookup[o.Center] = n
		}

		_, ok = lookup[o.Satellite]
		if !ok {
			new := Satellite{Name: o.Satellite, InOrbit: []*Satellite{}}
			lookup[o.Satellite] = &new
		}

		cen, _ := lookup[o.Center]
		sat, _ := lookup[o.Satellite]

		cen.InOrbit = append(cen.InOrbit, sat)
		sat.Orbits = cen
	}

	return lookup, center
}

// CreateOrbitList ...
func CreateOrbitList(in []string) []Orbit {
	orbits := []Orbit{}

	for _, s := range in {
		v := strings.Split(s, ")")
		orbits = append(orbits, Orbit{Center: v[0], Satellite: v[1]})
	}

	return orbits
}

// ReadOrbits ...
func ReadOrbits() []string {
	codes := []string{}

	file, err := os.Open("./orbits.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return codes
}

func main() {
	orbitMap, center := LinkOrbits(CreateOrbitList(ReadOrbits()))
	result := center.Depth(0)

	you, ok := orbitMap["YOU"]
	if !ok {
		panic("YOU not found")
	}

	san, ok := orbitMap["SAN"]
	if !ok {
		panic("SAN not found")
	}

	jumps := you.TransferTo(san)

	println(fmt.Sprintf("orbits: %v", result))
	println(fmt.Sprintf("jumps: %v", jumps))
}
