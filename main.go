package main

import (
	"github.com/gobtronic/steam-purchase-notifier/cmd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
