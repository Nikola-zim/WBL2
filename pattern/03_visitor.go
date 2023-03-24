package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Visitor_pattern
*/

// Интерфейс посетителя
type RobotVisitor interface {
	VisitDrone(*Drone)
	VisitRover(*Rover)
	VisitAndroid(*Android)
}

// Сам посетитель
type RepairVisitor struct{}

func (rv *RepairVisitor) VisitDrone(d *Drone) {
	fmt.Printf("Repairing drone %s with range %d\n", d.Model, d.Range)
}

func (rv *RepairVisitor) VisitRover(r *Rover) {
	fmt.Printf("Repairing rover %s with speed %d\n", r.Model, r.Speed)
}

func (rv *RepairVisitor) VisitAndroid(a *Android) {
	fmt.Printf("Repairing android %s with strength %d\n", a.Model, a.Strength)
}

// Конкретный элемент
type Drone struct {
	Model string
	Range int
}

func (d *Drone) Accept(visitor RobotVisitor) {
	visitor.VisitDrone(d)
}

// Конкретный элемент
type Rover struct {
	Model string
	Speed int
}

func (r *Rover) Accept(visitor RobotVisitor) {
	visitor.VisitRover(r)
}

// Конкретный элемент
type Android struct {
	Model    string
	Strength int
}

func (a *Android) Accept(visitor RobotVisitor) {
	visitor.VisitAndroid(a)
}

func VisitorWork() {
	d := &Drone{Model: "Phantom 4", Range: 10}
	r := &Rover{Model: "Curiosity", Speed: 5}
	a := &Android{Model: "Atlas", Strength: 100}

	rv := &RepairVisitor{}
	d.Accept(rv)
	r.Accept(rv)
	a.Accept(rv)
}
