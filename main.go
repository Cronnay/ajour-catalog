package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Cronnay/catalog-go/addons"
)

var (
	curseAPIKey *string
)

func init() {
	curseAPIKey = flag.String("capikey", "", "API Key to fetch data from curse")
	flag.Parse()
}

func main() {
	env := os.Getenv("APP_ENV")
	if env == "dev" {
		fmt.Println(*curseAPIKey)
	}

	if curseAPIKey == nil || *curseAPIKey == "" {
		panic("No API key was provided")
	}
	c := addons.NewCurse(*curseAPIKey)
	curseAddons, curseError := c.GetAddons()
	if curseError != nil {
		fmt.Println(curseError)
		panic(curseError)
	}

	allAddons := make([]addons.Addon, 0)
	allAddons = append(allAddons, curseAddons...)

	addonsAsBytes, jsonErr := json.Marshal(allAddons)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		panic(jsonErr)
	}

	writeErr := os.WriteFile("./catalog-0.2.0.json", addonsAsBytes, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
		panic(writeErr)
	}
}
