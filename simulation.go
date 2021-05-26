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
				grid[car.pos.y][car.pos.x] = 0
				break
			}

			prevCell := Cell{car.pos.x, car.pos.y}

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

			//Occupy the new cell if possible, and free the previous one
			if grid[car.pos.y][car.pos.x] == 0 {
				grid[prevCell.y][prevCell.x] = 0
				grid[car.pos.y][car.pos.x] = car.id
			}
		}

		// Do not move car if the next cell is not yours
		if grid[car.pos.y][car.pos.x] == car.id {
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

			// Try to occupy the needed cell
		} else if grid[car.pos.y][car.pos.x] == 0 {
			grid[car.pos.y][car.pos.x] = car.id
		}

		//Update the cars array to reflect changes
		cars[car.id-1] = car

		//Do not run every tick
		time.Sleep(time.Duration(car.sleep) * time.Millisecond)
	}

	car.active = false
	cars[car.id-1] = car
}

func changeTrafficLight(tLight TrafficLight) {

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
			options = new(ebiten.DrawImageOptions)
			options.GeoM.Translate(
				float64(trafficLights[i].cells[j].x*CellSize),
				float64(trafficLights[i].cells[j].y*CellSize),
			)
			screen.DrawImage(imgLight, options)
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
		for j := 0; j < len(trafficLights[i].cells); j++ {
			grid[trafficLights[i].cells[j].x][trafficLights[i].cells[j].y] = -1
		}
		go changeTrafficLight(trafficLights[i])
	}
	/*
		for i := range grid{
			fmt.Println(grid[i])
		}
	*/

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
