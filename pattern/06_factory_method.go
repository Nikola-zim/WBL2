package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern

В Go невозможно реализовать классический вариант паттерна Фабричный метод, поскольу в языке отсутствуют возможности ООП, в том числе классы и наследственность.
Несмотря на это, мы все же можем реализовать базовую версию этого паттерна — Простая фабрика.
Паттерн реализует создание различных схватов для роботов
*/

type RobotGrib interface {
	setTypeName(name string)
	setTorque(torque int)
	printGribInfo()
}

// Схват
type Grib struct {
	typeName string
	torque   int
}

func (g *Grib) setTypeName(name string) {
	g.typeName = name
}

func (g *Grib) setTorque(torque int) {
	g.torque = torque
}

func (g *Grib) printGribInfo() {
	fmt.Printf("Тип захвата: %s, Момент: %d \n\n", g.typeName, g.torque)
}

// Конкретный схват 1
type VacuumGrib struct {
	Grib
}

func newVacuumGrib() RobotGrib {
	return &VacuumGrib{
		Grib{
			typeName: "vacuum",
			torque:   60,
		},
	}
}

// Конкретный схват 2
type MechanicalGrib struct {
	Grib
}

func newMechanicalGrib() RobotGrib {
	return &MechanicalGrib{
		Grib{
			typeName: "mechanical",
			torque:   150,
		},
	}
}

// Конкретный схват 3
type ElectricalGrib struct {
	Grib
}

func newElectricalGrib() RobotGrib {
	return &ElectricalGrib{
		Grib{
			typeName: "electrical",
			torque:   150,
		},
	}
}

// Фабрика
func getGrib(gribType string) (RobotGrib, error) {
	switch gribType {
	case "vacuum":
		return newVacuumGrib(), nil
	case "mechanical":
		return newMechanicalGrib(), nil
	case "electrical":
		return newElectricalGrib(), nil
	}
	return nil, fmt.Errorf("такого типа нет")
}

func FactoryMethodWork() {
	vacuumG, _ := getGrib("vacuum")
	mechanicalG, _ := getGrib("mechanical")
	electricalG, _ := getGrib("electrical")
	vacuumG.printGribInfo()
	mechanicalG.printGribInfo()
	electricalG.printGribInfo()
}
