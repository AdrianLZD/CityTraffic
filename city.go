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
	id    int
	sleep int
	pos   Cell
	start Cell
	end   Cell
	route []string
}

var (
	cityMap           [30]byte
	intersectionPaths = map[string][][]string{
		"R": {{"D"}, {"R", "R"}, {"R", "U", "U"}},
		"D": {{"L"}, {"D", "D"}, {"D", "R", "R"}},
		"L": {{"U"}, {"L", "L"}, {"L", "D", "D"}},
		"U": {{"R"}, {"U", "U"}, {"U", "L", "L"}},
	}
)

// func moveCar(car Car) {

// 	for {
// 		if cityMap[car.pos+1] == 0 {
// 			cityMap[car.pos] = 0
// 			car.pos += 1
// 			cityMap[car.pos] = car.id
// 			fmt.Printf("%d se movio a %d: ", car.id, car.pos)
// 			fmt.Println(cityMap)
// 			time.Sleep(time.Duration(car.sleep) * time.Millisecond)
// 		}
// 	}

// }

// func test() {

// 	var car1 = Car{
// 		id:    1,
// 		pos:   3,
// 		sleep: 1000,
// 	}

// 	var car2 = Car{
// 		id:    2,
// 		pos:   5,
// 		sleep: 1200,
// 	}

// 	cityMap[car1.pos] = 1
// 	cityMap[car2.pos] = 0

// 	go moveCar(car1)
// 	go moveCar(car2)

// 	for {

// 	}

// }

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
	for grid[start.y][start.x] == "B" || grid[start.y][start.x] == "S" {
		start = Cell{rand.Intn(len(grid)), rand.Intn(len(grid[0]))}
	}

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

	// test()

	// Read grid file
	generationGrid, _ := csvToArray("grid.txt")

	// [TODO]: Initialize simulation grid

	// Initialize cars
	numCars := 5
	cars := make([]Car, numCars)
	for i := 0; i < numCars; i++ {
		route, start, end := generateRoute(generationGrid)
		cars[i] = Car{
			id:    i,
			sleep: rand.Intn(3000-1000) + 1000,
			pos:   start,
			start: start,
			end:   end,
			route: route,
		}
		//fmt.Println(cars[i])
	}

}
