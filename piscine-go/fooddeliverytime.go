package piscine

type food struct {
	preptime int
}

// Define a map to hold the menu items and their cooking times.
var menu = map[string]food{
	"burger":  {preptime: 15},
	"chips":   {preptime: 10},
	"nuggets": {preptime: 12},
}

func FoodDeliveryTime(order string) int {
	// Check if the order item exists in the menu.
	if food, exists := menu[order]; exists {
		// Return the cooking time for the item.
		return food.preptime
	}
	// If the item does not exist, return 0.
	return 404
}
