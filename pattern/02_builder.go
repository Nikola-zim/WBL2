package pattern

import (
	"fmt"
	"strings"
)

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/
//Структура создаваемого объекта
type Robot struct {
	robotType      string
	driverType     string
	motorNum       int
	navigationType []string
}

// Опишем интерфейс строителя
type IBuilder interface {
	setRobotType()
	setDriverType()
	setMotorNum()
	setNavigationType()
	getRobot() Robot
}

func getBuilder(builderType string) IBuilder {
	if builderType == "industrial" {
		return newIndustrialRobotBuilder()
	}
	if builderType == "civil" {
		return newCivilRobotBuilder()
	}
	return nil
}

// Конкретный строитель для первого типа
type IndustrialRobotBuilder struct {
	robotType      string
	driverType     string
	navigationType []string
	motorNum       int
}

func newIndustrialRobotBuilder() *IndustrialRobotBuilder {
	return &IndustrialRobotBuilder{}
}

func (ir *IndustrialRobotBuilder) setRobotType() {
	ir.robotType = "Cartesian robot"
}

func (ir *IndustrialRobotBuilder) setDriverType() {
	ir.driverType = "stepper"
}

func (ir *IndustrialRobotBuilder) setMotorNum() {
	ir.motorNum = 6
}
func (ir *IndustrialRobotBuilder) setNavigationType() {
	ir.navigationType = make([]string, 1)
	ir.navigationType[0] = "encoder"
}

func (ir *IndustrialRobotBuilder) getRobot() Robot {
	return Robot{
		robotType:      ir.robotType,
		driverType:     ir.driverType,
		motorNum:       ir.motorNum,
		navigationType: ir.navigationType,
	}
}

// Конкретный строитель для второго типа
type CivilRobotBuilder struct {
	robotType      string
	driverType     string
	navigationType []string
	motorNum       int
}

func newCivilRobotBuilder() *CivilRobotBuilder {
	return &CivilRobotBuilder{}
}

func (ir *CivilRobotBuilder) setRobotType() {
	ir.robotType = "mobile robot"
}

func (ir *CivilRobotBuilder) setDriverType() {
	ir.driverType = "DC"
}

func (ir *CivilRobotBuilder) setMotorNum() {
	ir.motorNum = 3
}
func (ir *CivilRobotBuilder) setNavigationType() {
	ir.navigationType = make([]string, 2)
	ir.navigationType[0] = "encoder"
	ir.navigationType[1] = "GPS"
}

func (ir *CivilRobotBuilder) getRobot() Robot {
	return Robot{
		robotType:      ir.robotType,
		driverType:     ir.driverType,
		motorNum:       ir.motorNum,
		navigationType: ir.navigationType,
	}
}

// Часть паттерна - "Директор"
type Director struct {
	builder IBuilder
}

func newDirector(b IBuilder) *Director {
	return &Director{
		builder: b,
	}
}

func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}

func (d *Director) buildRobot() Robot {
	d.builder.setRobotType()
	d.builder.setDriverType()
	d.builder.setMotorNum()
	d.builder.setNavigationType()
	return d.builder.getRobot()
}

// Клиентский код

func BuilderWork() {
	industrialBuilder := getBuilder("industrial")
	civilBuilder := getBuilder("civil")

	director := newDirector(industrialBuilder)
	industrialRobot := director.buildRobot()

	fmt.Printf("Industrial robot type: %s\n", industrialRobot.robotType)
	fmt.Printf("Industrial robot driver type: %s\n", industrialRobot.driverType)
	fmt.Printf("Industrial robot motor number: %d\n", industrialRobot.motorNum)
	fmt.Printf("Industrial robot navigation type: %s\n", strings.Join(industrialRobot.navigationType, ", "))

	director.setBuilder(civilBuilder)
	civilRobot := director.buildRobot()

	fmt.Printf("Civil robot type: %s\n", civilRobot.robotType)
	fmt.Printf("Civil robot driver type: %s\n", civilRobot.driverType)
	fmt.Printf("Civil robot motor number: %d\n", civilRobot.motorNum)

	fmt.Printf("Civil robot navigation type: %s\n", strings.Join(civilRobot.navigationType, ", "))

}
