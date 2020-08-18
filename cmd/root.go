package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags
	cfgFile string

	requiredFlags = []string{"collection.feeds", "collection.words"}
	rootCmd       = &cobra.Command{
		Use:   "s7rss",
		Short: "S7rss is a service for saving content of rss feeds",
		Long: `Service for storing and 
				delivering the content of rss feeds that contain certain words.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var skippedFlags []string
			for _, requiredFlag := range requiredFlags {
				if len(viper.GetStringSlice(requiredFlag)) == 0 {
					skippedFlags = append(skippedFlags, requiredFlag)
				}
			}
			if len(skippedFlags) > 0 {
				return fmt.Errorf("required flags \"%s\" has not been set", strings.Join(skippedFlags, ", "))
			}
			return nil
		},
	}
)

func init() {
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				log.WithError(err).
					Fatal("Read config")
			}
		}
		viper.SetEnvPrefix("S7RSS")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
	})
	pf := rootCmd.PersistentFlags()

	pf.StringVarP(&cfgFile, "config", "c", "../../../config.yaml", "Config file path")

	pf.StringSlice("collection.feeds", []string{}, "List of rss feeds")
	pf.StringSlice("collection.words", []string{}, "Search words in news headlines and announcements")
	pf.Int("collection.threads", 5, "Number of collection threads")
	pf.Int("collection.timeout", 1, "Delay between collection")

	pf.String("grpc.host", "", "GRPC listen host")
	pf.Int("grpc.port", 9000, "GRPC listen port")

	pf.String("db.host", "127.0.0.1", "Database host")
	pf.Int("db.port", 27017, "Database port")
	pf.String("db.database", "s7rss", "Database name")
	pf.String("db.username", "", "Database username")
	pf.String("db.password", "", "Database password")

	pf.String("log.output", "STDOUT", "Logging output")
	pf.String("log.level", "INFO", "Logging level")
	pf.String("log.format", "CLI", "Logging format: CLI or TEXT or JSON")

	if err := viper.BindPFlags(pf); err != nil {
		log.WithError(err).
			Fatal("Bind flags")
	}

	var file *os.File
	if viper.GetString("log.output") != "STDOUT" {
		var err error
		file, err = os.OpenFile(viper.GetString("log.output"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.WithError(err).
				Fatal("Open log file")
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.WithError(err).
					Error("Close log file")
			}
		}()
	} else {
		file = os.Stdout
	}

	switch viper.GetString("log.format") {
	case "TEXT":
		log.SetHandler(text.New(file))
	case "JSON":
		log.SetHandler(json.New(file))
	default:
		log.SetHandler(cli.New(file))
	}

	if logLevel, err := log.ParseLevel(viper.GetString("log.level")); err != nil {
		log.WithError(err).
			Fatal("Parse log level")
	} else {
		log.SetLevel(logLevel)
	}
}

func Execute(run func(cmd *cobra.Command, args []string)) error {
	rootCmd.Run = run
	return rootCmd.Execute()
}
