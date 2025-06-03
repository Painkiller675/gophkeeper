package cmd

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/Painkiller675/gophkeeper/pkg/version"
)

var (
	cfgFile  string
	defaults = map[string]interface{}{
		"grpc.address": "127.0.0.1:8081",
		"db.url":       "",
		"auth.key":     "",
		"hasher.key":   "",
	}
)

var rootCmd = &cobra.Command{
	Use:     "gophkeeper",
	Short:   "GophKeeper Server",
	Version: version.Info(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute command")
	}
}

func init() {
	// a try to get some options from the .env file
	if err := godotenv.Load(); err != nil {
		log.Debug().Msg("No .env file found")
	}
	// set some persistent flags
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile, "config", "c", "", "Server config filepath")

	rootCmd.PersistentFlags().StringP(
		"grpc-address", "g", "", "Server grpc address")

	rootCmd.PersistentFlags().StringP(
		"db-url", "d", "", "Database dns")

	rootCmd.PersistentFlags().StringP(
		"auth-key", "k", "", "Auth key")

	rootCmd.PersistentFlags().DurationP(
		"expiration-time", "e", 24*time.Hour, "Auth key expiration time")

	rootCmd.PersistentFlags().StringP(
		"hasher-key", "i", "", "Hash key")
	// it will be executed before the Execute func
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// set some defaults options
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// set some multiple paths to find cfgFile
		viper.AddConfigPath("/etc/gophkeeper")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")

		viper.SetConfigName("server-config")
	}
	// read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal().Err(err).Msg("Failed to load server config")
		}
	} else {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
	// read some .env to replace options from .yaml with them
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	// checking
	rootCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		key := strings.ReplaceAll(flag.Name, "-", ".")
		if err := viper.BindPFlag(key, flag); err != nil {
			log.Fatal().Err(err).Msg("Failed to bind flag")
		}
	})
}
