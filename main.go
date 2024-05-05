package main

import (
	"github.com/anton-uvarenko/headway_test/nasa/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cmd.Execute()
}
