package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var AdcConfig *Config

func SetupConfig() error {
	if err := viper.Unmarshal(&AdcConfig); err != nil {
		return errors.Wrap(err, "unmarshal config error")
	}

	srcDir := AdcConfig.SourceDirectory
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return errors.Errorf("source directory %s does not exist", srcDir)
	}

	outputDir := AdcConfig.SuccessOutputDirectory
	outputDir = filepath.ToSlash(outputDir)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return errors.Wrap(err, "create output directory error")
		}
	}
	return nil
}
