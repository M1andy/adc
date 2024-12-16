package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var AdcConfig *Config

func SetupConfig() {
	if err := viper.Unmarshal(&AdcConfig); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	srcDir := AdcConfig.SourceDirectory
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Source directory does not exist")
		os.Exit(1)
	}

	outputDir := AdcConfig.SuccessOutputDirectory
	outputDir = filepath.ToSlash(outputDir)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
