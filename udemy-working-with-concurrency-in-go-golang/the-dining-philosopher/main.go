package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	nPhilosophers        = 5
	DURATION_THINKING    = 1000
	DURATION_EATING      = 1000
	PHILOSOPHER_THINKING = "thinking"
	PHILOSOPHER_HUNGRY   = "hungry"
	PHILOSOPHER_EATING   = "eating"
)

type Fork struct {
	id     int
	isUsed bool
}

func (f *Fork) pickUp() {
	f.isUsed = true
}

func (f *Fork) putDown() {
	f.isUsed = false
}

type Philosopher struct {
	id        int
	state     string
	leftFork  *Fork
	rightFork *Fork
}

func (p *Philosopher) think() {
	fmt.Printf("Philosopher %d is thinking.\n", p.id)
	p.state = PHILOSOPHER_THINKING
	time.Sleep(DURATION_THINKING * time.Millisecond)
	p.state = PHILOSOPHER_HUNGRY
	fmt.Printf("Philosopher %d has finished thinking.\n", p.id)
}

func (p *Philosopher) eat() {
	if p.leftFork.isUsed || p.rightFork.isUsed {
		return
	}
	p.leftFork.pickUp()
	p.rightFork.pickUp()
	fmt.Printf("Philosopher %d is eating.\n", p.id)
	p.state = PHILOSOPHER_EATING
	time.Sleep(DURATION_EATING * time.Millisecond)
	fmt.Printf("Philosopher %d has finished eating.\n", p.id)
	p.leftFork.putDown()
	p.rightFork.putDown()
}

func (p *Philosopher) isHungry() bool {
	return p.state == PHILOSOPHER_HUNGRY
}

func run(philosophers []*Philosopher) {
	for {
		wg := sync.WaitGroup{}
		wg.Add(len(philosophers))
		for _, philosopher := range philosophers {
			go func(philosopher *Philosopher) {
				defer wg.Done()
				if philosopher.isHungry() {
					philosopher.eat()
				} else {
					philosopher.think()
				}
			}(philosopher)
		}
		wg.Wait()
	}
}

func main() {
	fmt.Println("Hello World")

	philosophers := make([]*Philosopher, nPhilosophers)
	forks := make([]*Fork, nPhilosophers)

	for i := range forks {
		forks[i] = &Fork{id: i}
	}

	for i := range philosophers {
		philosophers[i] = &Philosopher{
			id:        i,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%nPhilosophers],
			state:     PHILOSOPHER_HUNGRY,
		}
	}

	run(philosophers)
}
