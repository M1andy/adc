package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "adc/internal/config"
	. "adc/internal/crawler"
	. "adc/internal/logger"
)

var CfgFilePath string

var rootCmd = &cobra.Command{
	Use:   "adc",
	Short: "adc is a fast av data organizer",
	Run: func(cmd *cobra.Command, args []string) {
		main(cmd)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return nil
}

func main(cmd *cobra.Command) {
	SetupConfig()
	SetupLogger()

	isWatchDogMode, _ := cmd.Root().PersistentFlags().GetBool("watch")
	if !isWatchDogMode {
		oneTimeMode()
		return
	}

	watchDogMode()
}

func oneTimeMode() {
	Logger.Infoln("Init as one-time mode.")
	StartTasks("one-time")
	Logger.Infoln("Finish all jobs! Exiting...")
}

func watchDogMode() {
	Logger.Infoln("Init as watch dog mode.")

	done := make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		<-sigs
		Logger.Infoln("Exiting...")
		ticker.Stop()
		done <- true
	}()

	var wg conc.WaitGroup

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			wg.Go(func() {
				StartTasks("watchdog")
			})
			wg.Wait()
		}
	}
}

func init() {
	cobra.OnInitialize(initViper)

	rootCmd.PersistentFlags().BoolP("watch", "w", false, "whether to enable watch dog mode.")
	rootCmd.PersistentFlags().StringVarP(&CfgFilePath, "config", "c", "", "config file path (default is ./adc.toml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func initViper() {
	if CfgFilePath != "" {
		viper.SetConfigFile(CfgFilePath)
		fmt.Println("Using custom config file:", CfgFilePath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")

		viper.AddConfigPath("/etc/adc")
		viper.AddConfigPath("$HOME/.config/adc")
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed to read config file"))
		os.Exit(1)
	}

	if used := viper.ConfigFileUsed(); used != "" {
		fmt.Println("Using config file:", used)
	}
}
