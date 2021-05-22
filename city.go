package main

import (
	"fmt"
	"time"
)

type Car struct {
	id    byte
	pos   int64
	sleep int64
}

var (
	cityMap [30]byte
)

func moveCar(car Car) {
	for {
		if cityMap[car.pos+1] == 0 {
			cityMap[car.pos] = 0
			car.pos += 1
			cityMap[car.pos] = car.id
			fmt.Printf("%d se movio a %d: ", car.id, car.pos)
			fmt.Println(cityMap)
			time.Sleep(time.Duration(car.sleep) * time.Millisecond)
		}
	}

}

func main() {

	var car1 = Car{
		id:    1,
		pos:   3,
		sleep: 1000,
	}

	var car2 = Car{
		id:    2,
		pos:   5,
		sleep: 1200,
	}

	cityMap[car1.pos] = 1
	cityMap[car2.pos] = 0

	go moveCar(car1)
	go moveCar(car2)

	for {

	}

}
