package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

// MealInstance represents a meal instance with day, date, meal, and items
type MealInstance struct {
	Day   string   `json:"day"`
	Date  string   `json:"date"`
	Meal  string   `json:"meal"`
	Items []string `json:"items"`
}

// getMenuData reads the provided Excel file and returns the data as a map
func getMenuData() map[string]map[string][]string {
	xlFile, err := xlsx.OpenFile("weekly_menu.xlsx")
	if err != nil {
		fmt.Println("Error opening Excel file:", err)
		os.Exit(1)
	}

	fmt.Println(xlFile)

	menuData := make(map[string]map[string][]string)

	for _, sheet := range xlFile.Sheets {
		day := sheet.Name
		meals := make(map[string][]string)

		for _, row := range sheet.Rows {
			meal := strings.TrimSpace(row.Cells[0].String())
			items := strings.Split(strings.TrimSpace(row.Cells[1].String()), ",")

			// Filter out empty strings in the items slice
			var cleanedItems []string
			for _, item := range items {
				if item != "" {
					cleanedItems = append(cleanedItems, item)
				}
			}

			meals[meal] = cleanedItems
		}

		menuData[day] = meals
	}


	return menuData
}


// getItemsForMeal returns the items for a specific day and meal
func getItemsForMeal(day, meal string, menuData map[string]map[string][]string) []string {
	if menuData[day] == nil || menuData[day][meal] == nil {
		return nil
	}
	return menuData[day][meal]
}

// getNumItemsForMeal returns the number of items for a specific day and meal
func getNumItemsForMeal(day, meal string, menuData map[string]map[string][]string) int {
	return len(getItemsForMeal(day, meal, menuData))
}

// isItemInMeal checks if a given item is in a particular meal on a specific day
func isItemInMeal(day, meal, item string, menuData map[string]map[string][]string) bool {
	items := getItemsForMeal(day, meal, menuData)
	for _, i := range items {
		if strings.EqualFold(i, item) {
			return true
		}
	}
	return false
}

// saveMenuAsJSON converts the entire menu into JSON and saves it to a file
func saveMenuAsJSON(menuData map[string]map[string][]string) {
	menuJSON, err := json.MarshalIndent(menuData, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling menu data:", err)
		os.Exit(1)
	}

	err = os.WriteFile("menu.json", menuJSON, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		os.Exit(1)
	}
}

// printDetails prints the details of a MealInstance
func (m *MealInstance) printDetails() {
	fmt.Printf("%s (%s), %s meal:\n", m.Day, m.Date, m.Meal)
	for _, item := range m.Items {
		fmt.Printf("- %s\n", item)
	}
}

func main() {
	// Load menu data from Excel file
	menuData := getMenuData()
	fmt.Println(menuData)

	// Get user input for day and meal
	var inputDay, inputMeal string
	fmt.Print("Enter day: ")
	fmt.Scan(&inputDay)
	fmt.Print("Enter meal: ")
	fmt.Scan(&inputMeal)

	// Get items for the specified day and meal
	itemsForMeal := getItemsForMeal(inputDay, inputMeal, menuData)

	// Check if items are found and print details
	if itemsForMeal != nil {
		numItemsForMeal := getNumItemsForMeal(inputDay, inputMeal, menuData)
		fmt.Printf("Items for %s meal on %s:\n", inputMeal, inputDay)
		for _, item := range itemsForMeal {
			fmt.Printf("- %s\n", item)
		}
		fmt.Printf("Number of items: %d\n", numItemsForMeal)
	} else {
		fmt.Println("No items found for the specified day and meal.")
	}

	// Save the entire menu as JSON
	saveMenuAsJSON(menuData)

	// Create a sample MealInstance
	sampleMealInstance := MealInstance{
		Day:   "Monday",
		Date:  "2024-02-10",
		Meal:  "Breakfast",
		Items: []string{"Eggs", "Toast"},
	}

	// Print details of the sample MealInstance
	sampleMealInstance.printDetails()
}
