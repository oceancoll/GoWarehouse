package service

import "Generics/Interface"

var AnimalService Interface.Animal[any]

func SetUp() {
	AnimalService = Interface.NewAnimalService()
}
