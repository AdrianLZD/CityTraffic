package main

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Cell struct {
	x int
	y int
}

type Car struct {
	id         int
	sleep      int
	pos        Cell
	start      Cell
	end        Cell
	gridPos    Cell
	route      []string
	routeIndex int
	active     bool
}

type TrafficLight struct {
	id         int
	cells      []Cell
	activeCell int
	sleep      int
}

var (
	intersectionPaths = map[string][][]string{
		"R": {{"D"}, {"R", "R"}, {"R", "U", "U"}},
		"D": {{"L"}, {"D", "D"}, {"D", "R", "R"}},
		"L": {{"U"}, {"L", "L"}, {"L", "D", "D"}},
		"U": {{"R"}, {"U", "U"}, {"U", "L", "L"}},
	}
	usedStartPos           = make(map[Cell]bool)
	trafficLightsPositions = [][]Cell{
		// Two at the top
		{{9, 1}, {10, 1}, {9, 2}, {10, 2}},
		{{17, 1}, {18, 1}, {17, 2}, {18, 2}},

		// Four in the top middle
		{{1, 9}, {2, 9}, {1, 10}, {2, 10}},
		{{9, 9}, {10, 9}, {9, 10}, {10, 10}},
		{{17, 9}, {18, 9}, {17, 10}, {18, 10}},
		{{25, 9}, {26, 9}, {25, 10}, {26, 10}},

		// Four in the bottom middle
		{{1, 17}, {2, 17}, {1, 18}, {2, 18}},
		{{9, 17}, {10, 17}, {9, 18}, {10, 18}},
		{{17, 17}, {18, 17}, {17, 18}, {18, 18}},
		{{25, 17}, {26, 17}, {25, 18}, {26, 18}},

		// Two at the bottom
		{{9, 25}, {10, 25}, {9, 26}, {10, 26}},
		{{17, 25}, {18, 25}, {17, 26}, {18, 26}},
	}
)

func csvToArray(path string) ([][]string, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.Split(scanner.Text(), ", "))
	}

	return lines, scanner.Err()

}

func generateRoute(grid [][]string) ([]string, Cell, Cell) {

	// Get random route with length [25, 30]
	rand.Seed(time.Now().UnixNano())
	length := rand.Intn(30-25) + 25

	// Get random initial position
	start := Cell{rand.Intn(len(grid)), rand.Intn(len(grid[0]))}
	for grid[start.y][start.x] == "B" || grid[start.y][start.x] == "S" || usedStartPos[start] {
		start = Cell{rand.Intn(len(grid)), rand.Intn(len(grid[0]))}
	}
	usedStartPos[start] = true

	// Travel randomly following street directions
	var route []string
	var steps []string
	x := start.x
	y := start.y
	count := 0

	for count < length {
		pos := grid[y][x]

		if pos == "S" {
			steps = getIntersectionSteps(grid, x, y, route[count-1])
		} else {
			steps = []string{pos}
		}

		for _, s := range steps {
			switch {
			case s == "R":
				route = append(route, s)
				x++
			case s == "D":
				route = append(route, s)
				y++
			case s == "L":
				route = append(route, s)
				x--
			case s == "U":
				route = append(route, s)
				y--
			}
		}
		count += len(steps)
	}

	return route, start, getLastPosition(x, y, route[len(route)-1])

}

func getIntersectionSteps(grid [][]string, x int, y int, lastDir string) []string {

	var validOptions []int

	switch {
	case lastDir == "R":
		if grid[y+1][x] != "B" {
			validOptions = append(validOptions, 0)
		}
		if grid[y][x+2] != "B" {
			validOptions = append(validOptions, 1)
		}
		if grid[y-2][x+1] != "B" {
			validOptions = append(validOptions, 2)
		}
	case lastDir == "D":
		if grid[y][x-1] != "B" {
			validOptions = append(validOptions, 0)
		}
		if grid[y+2][x] != "B" {
			validOptions = append(validOptions, 1)
		}
		if grid[y+1][x+2] != "B" {
			validOptions = append(validOptions, 2)
		}
	case lastDir == "L":
		if grid[y-1][x] != "B" {
			validOptions = append(validOptions, 0)
		}
		if grid[y][x-2] != "B" {
			validOptions = append(validOptions, 1)
		}
		if grid[y+2][x-1] != "B" {
			validOptions = append(validOptions, 2)
		}
	case lastDir == "U":
		if grid[y][x+1] != "B" {
			validOptions = append(validOptions, 0)
		}
		if grid[y-2][x] != "B" {
			validOptions = append(validOptions, 1)
		}
		if grid[y-1][x-2] != "B" {
			validOptions = append(validOptions, 2)
		}
	}

	rand.Seed(time.Now().UnixNano())
	return intersectionPaths[lastDir][validOptions[rand.Intn(len(validOptions))]]

}

func getLastPosition(x int, y int, lastDir string) Cell {
	switch {
	case lastDir == "R":
		x--
	case lastDir == "D":
		y--
	case lastDir == "L":
		x++
	case lastDir == "U":
		y++
	}
	return Cell{x, y}
}

func main() {

	// Read grid file
	generationGrid, _ := csvToArray("grid.txt")

	// Initialize cars
	numCars := 30
	cars := make([]Car, numCars)
	for i := 0; i < numCars; i++ {
		route, start, end := generateRoute(generationGrid)
		cars[i] = Car{
			id:         i + 1,
			sleep:      rand.Intn(20-10) + 10,
			pos:        start,
			start:      start,
			end:        end,
			gridPos:    Cell{start.x * CellSize, start.y * CellSize},
			route:      route,
			routeIndex: -1,
			active:     true,
		}
	}

	numLights := 12
	trafficLights := make([]TrafficLight, numLights)
	for i := 0; i < numLights; i++ {
		trafficLights[i] = TrafficLight{
			id:         i + 1,
			cells:      trafficLightsPositions[i],
			activeCell: 0,
			sleep:      3000,
		}
	}

	// Method from simulation.go
	initCity(cars, trafficLights)
}
