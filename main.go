package main

import (
	"adc/cmd"
	_ "adc/internal/config"
	_ "adc/internal/logger"
)

func main() {
	cmd.Execute()
}
