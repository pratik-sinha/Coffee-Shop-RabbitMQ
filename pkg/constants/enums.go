package constants

type Status int8

const (
	StatusPlaced Status = iota
	StatusInProcess
	StatusFulfilled
)

func (e Status) String() string {
	return []string{
		"STATUS_PLACED",
		"STATUS_IN_PROCESS",
		"STATUS_FULFILLED",
	}[e]
}

type ItemType int32

const (
	ItemTypeCappuccino ItemType = iota
	ItemTypeCoffeeBlack
	ItemTypeCoffeeWithRoom
	ItemTypeEspresso
	ItemTypeEspressoDouble
	ItemTypeLatte
	ItemTypeCakePop
	ItemTypeCroissant
	ItemTypeMuffin
	ItemTypeCroissantChocolate
)

func (e ItemType) String() string {
	return []string{
		"CAPPUCCINO",
		"COFFEE_BLACK",
		"COFFEE_WITH_ROOM",
		"ESPRESSO",
		"ESPRESSO_DOUBLE",
		"LATTE",
		"CAKEPOP",
		"CROISSANT",
		"MUFFIN",
		"CROISSANT_CHOCOLATE",
		"CAPPUCCINO",
	}[e]
}

type KitchenType int32

const (
	Kitchen KitchenType = iota
	Barista
)

func (k KitchenType) String() string {
	return []string{
		"KITCHEN",
		"BARISTA",
	}[k]
}
