package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, totalPizza int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {

		fmt.Printf("Received order number #%d\n", pizzaNumber)

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     "",
			success:     false,
		}

		rnd := rand.Intn(12) + 1
		if rnd < 5 {
			p.success = false
			if rnd <= 2 {
				p.message = fmt.Sprintf("*** run out of ingredients for pizza #%d!", pizzaNumber)
			} else {
				p.message = fmt.Sprintf("*** the cook quit while makeing pizza #%d!", pizzaNumber)
			}
			pizzasFailed++
		} else {
			p.success = true
			p.message = fmt.Sprintf("*** Pizza is ready #%d!", pizzaNumber)
			pizzasMade++
		}
		totalPizza++

		delay := rand.Intn(5) + 1
		fmt.Printf("Making pizza #%d. It will take %d seconds ...\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay))

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}

}

func pizzeria(pizzaMaker *Producer) {

	var i = 0

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}

}

func main() {
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(pizzaJob)

	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The consumer is really mad!")
			}
		} else {
			color.Cyan("Done making pizza!")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channels!", err)
			}
		}
	}

}
