package main

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

const (
	GUIHeight     = 616
	GUIWidth      = 616
	CellSize      = 22
	StartWaitTime = 4000
)

var (
	imgBackground *ebiten.Image
	imgLight      *ebiten.Image
	carSprites    map[string]*ebiten.Image
	cars          []Car
	trafficLights []TrafficLight
	grid          [28][28]int
)

func init() {
	err := loadImages()
	if err != nil {
		log.Fatal(err)
	}
}

func loadImages() error {
	var err error
	imgBackground, _, err = ebitenutil.NewImageFromFile("res/background.png")
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
		}

		// Try to occupy the needed cell
		if grid[car.pos.y][car.pos.x] == 0 {
			// If car was in a crossing, it cannot override the traffic light cell
			if !car.inCrossing {
				grid[car.prevPos.y][car.prevPos.x] = 0
			}
			car.inCrossing = false
			grid[car.pos.y][car.pos.x] = car.id

		} else if grid[car.pos.y][car.pos.x] == -1 && !car.inCrossing {
			car.inCrossing = true
			grid[car.prevPos.y][car.prevPos.x] = 0
		}

		// Do not move car if the next cell is not yours
		if grid[car.pos.y][car.pos.x] == car.id || car.inCrossing {
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
			// [TODO] Fix collisions inside the intersection...
		}

		// Update the cars array to reflect changes
		cars[car.id-1] = car

		// Controls the car's speed
		time.Sleep(time.Duration(car.sleep) * time.Millisecond)
	}

	grid[car.pos.y][car.pos.x] = 0
	car.active = false
	cars[car.id-1] = car
}

func changeTrafficLight(tLight TrafficLight) {
	time.Sleep(time.Duration(StartWaitTime) * time.Millisecond)
	for {
		if tLight.activeCell >= len(tLight.cells) {
			tLight.activeCell = 0
		}

		for i := 0; i < len(tLight.cells); i++ {
			if i != tLight.activeCell {
				grid[tLight.cells[i].y][tLight.cells[i].x] = -2
			} else {
				grid[tLight.cells[i].y][tLight.cells[i].x] = -1
			}
		}

		trafficLights[tLight.id-1] = tLight

		tLight.activeCell += 1
		time.Sleep(time.Duration(tLight.sleep) * time.Millisecond)
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(imgBackground, nil)
	drawCars(screen)
	drawTrafficLights(screen)
}

func drawCars(screen *ebiten.Image) {
	var options *ebiten.DrawImageOptions
	for i := 0; i < len(cars); i++ {
		if cars[i].active {
			options = new(ebiten.DrawImageOptions)
			options.GeoM.Translate(
				float64(cars[i].gridPos.x),
				float64(cars[i].gridPos.y))
			sprite := cars[i].routeIndex
			if sprite < 0 {
				sprite = 0
			}
			screen.DrawImage(carSprites[cars[i].route[sprite]], options)
		}
	}
}

func drawTrafficLights(screen *ebiten.Image) {
	var options *ebiten.DrawImageOptions
	for i := 0; i < len(trafficLights); i++ {
		for j := 0; j < len(trafficLights[i].cells); j++ {
			if grid[trafficLights[i].cells[j].y][trafficLights[i].cells[j].x] == -2 {
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GUIWidth, GUIHeight
}

func initCity(newCars []Car, tLights []TrafficLight) {
	cars = newCars
	trafficLights = tLights
	ebiten.SetWindowSize(GUIWidth, GUIHeight)
	ebiten.SetWindowTitle("City Traffic")
	fmt.Println("City")

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
			grid[trafficLights[i].cells[j].y][trafficLights[i].cells[j].x] = -2
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
			grid[yCoord][xCoord] = -1
		}

		go changeTrafficLight(trafficLights[i])
	}

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

	for i := range grid {
		fmt.Println(grid[i])
	}

}
