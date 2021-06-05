package main

import (
	"bufio"
	"fmt"
	_ "image/png"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

const (
	GUIHeight     = 616
	GUIWidth      = 900
	LogStart      = 620
	LogLineHeight = 14
	CellSize      = 22
	StartWaitTime = 4000
)

var (
	imgBackground  *ebiten.Image
	imgLight       *ebiten.Image
	carSprites     map[string]*ebiten.Image
	cars           []Car
	trafficLights  []TrafficLight
	grid           [28][28]int
	gridTLights    [28][28]int
	finishedCars   []Car
	carsFinished   int
	routeDisplayed int
)

func init() {
	err := loadImages()
	if err != nil {
		log.Fatal(err)
	}
}

func loadImages() error {
	var err error
	imgBackground, _, err = ebitenutil.NewImageFromFile("res/backgroundDetailed.png")
	if err != nil {
		return err
	}
	carSprites = make(map[string]*ebiten.Image)
	imgCarU, _, err := ebitenutil.NewImageFromFile("res/carU.png")
	if err != nil {
		return err
	}
	carSprites["U"] = imgCarU

	imgCarR, _, err := ebitenutil.NewImageFromFile("res/carR.png")
	if err != nil {
		return err
	}
	carSprites["R"] = imgCarR

	imgCarD, _, err := ebitenutil.NewImageFromFile("res/carD.png")
	if err != nil {
		return err
	}
	carSprites["D"] = imgCarD

	imgCarL, _, err := ebitenutil.NewImageFromFile("res/carL.png")
	if err != nil {
		return err
	}
	carSprites["L"] = imgCarL

	imgLight, _, err = ebitenutil.NewImageFromFile("res/trafficLight.png")
	if err != nil {
		return err
	}

	return nil
}

func moveCar(car Car) {
	time.Sleep(time.Duration(StartWaitTime) * time.Millisecond)
	for {
		posToGrid := Cell{car.pos.x * CellSize, car.pos.y * CellSize}

		//The car has reached the desired position, receive next instruction
		if posToGrid == car.gridPos {
			car.routeIndex += 1

			//If the route is over, exit the loop
			if car.routeIndex >= len(car.route) {
				break
			}

			car.prevPos = Cell{car.pos.x, car.pos.y}

			switch car.route[car.routeIndex] {
			case "U":
				car.pos.y -= 1
			case "R":
				car.pos.x += 1
			case "D":
				car.pos.y += 1
			case "L":
				car.pos.x -= 1
			}

			if car.routeIndex+1 < len(car.route) {
				// Id there is no car in advance and speed had been reduced return to original speed
				nextPos := getNextPos(car)
				if nextPos == 0 && car.sleep != car.originalSleep {
					//fmt.Printf("[Car %v]: Speed %v <- %v\n", car.id, math.Abs(float64(car.sleep-30)), math.Abs(float64(car.originalSleep-30)))
					car.sleep = car.originalSleep
				}
			}

		} else if car.routeIndex+1 < len(car.route) {
			// Verify if next position in route has a slower car to slow in advance
			nextPos := getNextPos(car)
			if nextPos != 0 && car.sleep < cars[nextPos-1].sleep {
				newSleep := cars[nextPos-1].sleep
				//fmt.Printf("[Car %v]: Speed %v -> %v\n", car.id, math.Abs(float64(car.sleep-30)), math.Abs(float64(newSleep-30)))
				car.sleep = newSleep
			}
		}

		// Check if the traffic light is open
		if gridTLights[car.pos.y][car.pos.x] == 1 {
			car.inCrossing = true
		} else if gridTLights[car.pos.y][car.pos.x] == 0 {
			car.inCrossing = false
		}

		/*Try to occupy the needed cell.
		If the car is already in a crossing it will try to move
		until it exit the intersection*/
		if grid[car.pos.y][car.pos.x] == 0 && (gridTLights[car.pos.y][car.pos.x] == 0 || car.inCrossing) {
			grid[car.prevPos.y][car.prevPos.x] = 0
			grid[car.pos.y][car.pos.x] = car.id
		} /*else if grid[car.pos.y][car.pos.x] != 0 && car.sleep < cars[grid[car.pos.y][car.pos.x]-1].sleep {
			newSleep := cars[grid[car.pos.y][car.pos.x]-1].sleep
			fmt.Printf("[Car %v]: Speed %v -> %v\n", car.id, math.Abs(float64(car.sleep-30)), math.Abs(float64(newSleep-30)))
			car.sleep = newSleep
		}*/

		// Do not move car if the next cell is not yours
		if grid[car.pos.y][car.pos.x] == car.id /*|| car.inCrossing*/ {
			switch car.route[car.routeIndex] {
			case "U":
				car.gridPos.y -= 1
			case "R":
				car.gridPos.x += 1
			case "D":
				car.gridPos.y += 1
			case "L":
				car.gridPos.x -= 1
			}
		}

		// Update the cars array to reflect changes
		cars[car.id-1] = car

		// Controls the car's speed
		time.Sleep(time.Duration(car.sleep) * time.Millisecond)
	}

	//Rules to stop car in simulation
	finishedCars[carsFinished] = car
	carsFinished += 1
	grid[car.pos.y][car.pos.x] = 0
	car.active = false
	cars[car.id-1] = car

	// Creates a fade out effect
	for car.opacity >= 0 {
		car.opacity -= 0.05
		cars[car.id-1] = car
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

}

func changeTrafficLight(tLight TrafficLight) {
	time.Sleep(time.Duration(StartWaitTime) * time.Millisecond)
	for {
		if tLight.activeCell >= len(tLight.cells) {
			tLight.activeCell = 0
		}

		for i := 0; i < len(tLight.cells); i++ {
			if i != tLight.activeCell {
				gridTLights[tLight.cells[i].y][tLight.cells[i].x] = 2
			} else {
				gridTLights[tLight.cells[i].y][tLight.cells[i].x] = 1
			}
		}

		trafficLights[tLight.id-1] = tLight

		tLight.activeCell += 1
		time.Sleep(time.Duration(tLight.sleep) * time.Millisecond)
	}
}

func getNextPos(car Car) int {
	var aheadPos Cell

	switch car.route[car.routeIndex+1] {
	case "U":
		aheadPos = Cell{car.pos.x, car.pos.y - 1}
	case "R":
		aheadPos = Cell{car.pos.x + 1, car.pos.y}
	case "D":
		aheadPos = Cell{car.pos.x, car.pos.y + 1}
	case "L":
		aheadPos = Cell{car.pos.x - 1, car.pos.y}
	}

	return grid[aheadPos.y][aheadPos.x]
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(imgBackground, nil)
	drawCarRoute(screen)
	drawCars(screen)
	drawTrafficLights(screen)
	drawLog(screen)
}

func drawCarRoute(screen *ebiten.Image) {
	if routeDisplayed < 1 {
		return
	}
	var options *ebiten.DrawImageOptions
	posToGrid := Cell{cars[routeDisplayed-1].start.x * CellSize, cars[routeDisplayed-1].start.y * CellSize}

	// Display the starting cell
	options = new(ebiten.DrawImageOptions)
	options.GeoM.Translate(
		float64(posToGrid.x),
		float64(posToGrid.y))
	options.ColorM.RotateHue(180)
	screen.DrawImage(imgLight, options)
	ebitenutil.DebugPrintAt(screen, "START", posToGrid.x, posToGrid.y+2)

	// Traverse the instructions and display the resultant cell
	for i := 0; i < len(cars[routeDisplayed-1].route); i++ {
		switch cars[routeDisplayed-1].route[i] {
		case "U":
			posToGrid = Cell{posToGrid.x, posToGrid.y - CellSize}
		case "R":
			posToGrid = Cell{posToGrid.x + CellSize, posToGrid.y}
		case "D":
			posToGrid = Cell{posToGrid.x, posToGrid.y + CellSize}
		case "L":
			posToGrid = Cell{posToGrid.x - CellSize, posToGrid.y}
		}

		options = new(ebiten.DrawImageOptions)
		options.GeoM.Translate(
			float64(posToGrid.x),
			float64(posToGrid.y),
		)

		// Check if this is the END cell
		if i != len(cars[routeDisplayed-1].route)-1 {
			screen.DrawImage(imgLight, options)
		} else {
			options.ColorM.RotateHue(180)
			screen.DrawImage(imgLight, options)
			ebitenutil.DebugPrintAt(screen, "END", posToGrid.x, posToGrid.y+2)
		}

	}
}

func drawCars(screen *ebiten.Image) {
	var options *ebiten.DrawImageOptions
	for i := 0; i < len(cars); i++ {
		options = new(ebiten.DrawImageOptions)
		options.GeoM.Translate(
			float64(cars[i].gridPos.x),
			float64(cars[i].gridPos.y))
		sprite := cars[i].routeIndex
		if sprite < 0 {
			sprite = 0
		}
		if cars[i].active {
			screen.DrawImage(carSprites[cars[i].route[sprite]], options)
			sleepStr := fmt.Sprintf("%.0f", math.Abs(float64(cars[i].sleep-30)))
			ebitenutil.DebugPrintAt(screen, sleepStr, cars[i].gridPos.x+CellSize/2, cars[i].gridPos.y-CellSize/2)
		} else if cars[i].opacity > 0 {
			options.ColorM.Scale(1, 1, 1, cars[i].opacity)
			screen.DrawImage(carSprites[cars[i].route[len(cars[i].route)-1]], options)
		}
	}
}

func drawTrafficLights(screen *ebiten.Image) {
	var options *ebiten.DrawImageOptions
	for i := 0; i < len(trafficLights); i++ {
		for j := 0; j < len(trafficLights[i].cells); j++ {
			if gridTLights[trafficLights[i].cells[j].y][trafficLights[i].cells[j].x] == 2 {
				options = new(ebiten.DrawImageOptions)
				options.GeoM.Translate(
					float64(trafficLights[i].cells[j].x*CellSize),
					float64(trafficLights[i].cells[j].y*CellSize),
				)
				screen.DrawImage(imgLight, options)
			}

		}
	}
}

func drawLog(screen *ebiten.Image) {
	carLen := fmt.Sprintf("%d", len(cars))

	ebitenutil.DebugPrintAt(screen, "City Traffic Simulation", LogStart, LogLineHeight*0)
	ebitenutil.DebugPrintAt(screen, "Simulating a total of "+carLen+" cars.", LogStart, LogLineHeight*1)

	ebitenutil.DebugPrintAt(screen, "Cars that have completed their route:", LogStart, LogLineHeight*3)
	line := 4
	for i := 0; i < len(finishedCars); i++ {
		if finishedCars[i].id != 0 {
			idStr := fmt.Sprintf("%d", finishedCars[i].id)
			startStr := fmt.Sprintf("[%d, %d]", finishedCars[i].start.x, finishedCars[i].start.y)
			endStr := fmt.Sprintf("[%d, %d]", finishedCars[i].end.x, finishedCars[i].end.y)

			toPrint := "Car " + idStr + " (Start: " + startStr + " Finish: " + endStr + ")"
			ebitenutil.DebugPrintAt(screen, toPrint, LogStart, LogLineHeight*line)
			line += 1
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GUIWidth, GUIHeight
}

func readInput(scanner *bufio.Scanner) {
	for scanner.Scan() {
		text := scanner.Text()
		id, err := strconv.Atoi(text)
		if err != nil || id > len(finishedCars) {
			fmt.Println("Please enter a valid ID")
			fmt.Print("Enter a car ID: ")
			continue
		}

		carFound := false
		for i := 0; i < len(finishedCars); i++ {
			if finishedCars[i].id == id {
				routeDisplayed = id
				fmt.Println(finishedCars[i].route)
				fmt.Println("Displaying route for Car " + text)
				fmt.Print("Enter a car ID: ")
				carFound = true
				break
			}
		}
		if !carFound {
			fmt.Println("Car " + text + " has not finished its route")
			fmt.Print("Enter a car ID: ")
		}

	}

}

func initCity(newCars []Car, tLights []TrafficLight) {
	cars = newCars
	finishedCars = make([]Car, len(cars))
	trafficLights = tLights
	ebiten.SetWindowSize(GUIWidth, GUIHeight)
	ebiten.SetWindowTitle("City Traffic")
	fmt.Println("City Traffic Simulation")
	fmt.Println("To see the route of a car type its ID")
	fmt.Print("Enter a car ID: ")
	go readInput(bufio.NewScanner(os.Stdin))

	// Let the cars occupy their initial cell
	for i := 0; i < len(cars); i++ {
		grid[cars[i].pos.y][cars[i].pos.x] = cars[i].id
		go moveCar(cars[i])
	}

	// Start all the traffic lights cells
	for i := 0; i < len(trafficLights); i++ {
		lenCells := len(trafficLights[i].cells)

		xMap := make(map[int]int)
		yMap := make(map[int]int)

		for j := 0; j < lenCells; j++ {
			xMap[trafficLights[i].cells[j].x] += 1
			yMap[trafficLights[i].cells[j].y] += 1
			gridTLights[trafficLights[i].cells[j].y][trafficLights[i].cells[j].x] = 2
		}

		// Fill missing intersection cell with a "free" space
		if lenCells < 4 {
			xCoord := 0
			yCoord := 0
			for k, v := range xMap {
				if v == 1 {
					xCoord = k
					break
				}
			}
			for k, v := range yMap {
				if v == 1 {
					yCoord = k
					break
				}
			}
			gridTLights[yCoord][xCoord] = 1
		}

		go changeTrafficLight(trafficLights[i])
	}

	if err := ebiten.RunGame(&Game{}); err != nil {

		log.Fatal(err)
	}
	/*
		for i := range grid {
			fmt.Println(grid[i])
		}*/

}
