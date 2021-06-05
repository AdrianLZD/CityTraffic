package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cell struct {
	x int
	y int
}

type Car struct {
	id            int
	sleep         int
	originalSleep int
	prevPos       Cell
	pos           Cell
	start         Cell
	end           Cell
	gridPos       Cell
	route         []string
	routeIndex    int
	active        bool
	inCrossing    bool
	opacity       float64
}

type TrafficLight struct {
	id          int
	cells       []Cell
	activeCell  int
	missingCell int
	sleep       int
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
		{{10, 1}, {9, 2}, {10, 2}},
		{{18, 1}, {17, 2}, {18, 2}},

		// Four in the top middle
		{{1, 9}, {2, 9}, {2, 10}},
		{{9, 9}, {10, 9}, {9, 10}, {10, 10}},
		{{17, 9}, {18, 9}, {17, 10}, {18, 10}},
		{{25, 9}, {25, 10}, {26, 10}},

		// Four in the bottom middle
		{{1, 17}, {2, 17}, {2, 18}},
		{{9, 17}, {10, 17}, {9, 18}, {10, 18}},
		{{17, 17}, {18, 17}, {17, 18}, {18, 18}},
		{{25, 17}, {25, 18}, {26, 18}},

		// Two at the bottom
		{{9, 25}, {10, 25}, {9, 26}},
		{{17, 25}, {18, 25}, {17, 26}},
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
	length := rand.Intn(50-45) + 45

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

	argsLen := len(os.Args[1:])
	numCars := 30
	numLights := 12

	if argsLen == 2 {
		if v, err := strconv.Atoi(os.Args[1]); err == nil {
			if v > 0 && v < 31 {
				numCars = v
			} else {
				fmt.Println("[ERROR]: Number of cars must be in range [1,30]")
				os.Exit(0)
			}
		} else {
			fmt.Println("[ERROR]: Number of cars <arg> must be of type int")
			os.Exit(0)
		}
		if v, err := strconv.Atoi(os.Args[2]); err == nil {
			if v > -1 && v < 13 {
				numLights = v
			} else {
				fmt.Println("[ERROR]: Number of traffic lights must be in range [0,12]")
				os.Exit(0)
			}
		} else {
			fmt.Println("[ERROR]: Number of traffic lights <arg> must be of type int")
			os.Exit(0)
		}
	} else if argsLen != 0 {
		fmt.Println("[ERROR]: Usage ./main.out <number of cars> <number of traffic lights>")
		os.Exit(0)
	}

	// Read grid file
	generationGrid, _ := csvToArray("grid.txt")

	// Initialize cars
	cars := make([]Car, numCars)
	for i := 0; i < numCars; i++ {
		route, start, end := generateRoute(generationGrid)
		speed := rand.Intn(30-5) + 5
		cars[i] = Car{
			id:            i + 1,
			sleep:         speed,
			originalSleep: speed,
			prevPos:       start,
			pos:           start,
			start:         start,
			end:           end,
			gridPos:       Cell{start.x * CellSize, start.y * CellSize},
			route:         route,
			routeIndex:    -1,
			active:        true,
			inCrossing:    false,
			opacity:       1,
		}
		//fmt.Printf("[Car %v]: Original speed %v\n", i, math.Abs(float64(speed-30)))
	}

	trafficLights := make([]TrafficLight, numLights)
	for i := 0; i < numLights; i++ {
		trafficLights[i] = TrafficLight{
			id:         i + 1,
			cells:      trafficLightsPositions[i],
			activeCell: 0,
			/* No traffic light has a cell missing by its own,
			a missing cell it is determined by the simulation space.
			Therefore, it needs to be set up while initializing a city.
			*/
			missingCell: -1,
			sleep:       3000,
		}
	}

	// Method from simulation.go
	initCity(cars, trafficLights)
}
