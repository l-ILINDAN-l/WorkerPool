package cmd

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
// It is the main entry point for the application's command-line interface.
var rootCmd = &cobra.Command{
	Use:   "worker-pool",
	Short: "An application for demonstrating the work of a pool of workers",
	Long: `This application runs an interactive console for managing a pool of workers.
You can dynamically add and remove workers, as well as send them tasks.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

// init is called by Go during package initialization.
// It sets up the command-line flags.
func init() {
	cobra.OnInitialize(Config)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.worker-pool.yaml)")
}

// Config reads in config file and ENV variables if set.
// It also configures the application's logger.
func Config() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	if cfgFile != "" {

	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".worker-pool")
		viper.SetConfigType("yaml")
	}

	viper.SetDefault("workers.initial", 10)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logrus.WithError(err).Warn("Error reading config file. Using defaults.")
		} else {
			logrus.Warnf("Error config file: %v", err)
		}
	}
}
