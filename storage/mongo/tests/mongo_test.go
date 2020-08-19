package tests

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qystishere/s7rss/cmd"
	"github.com/qystishere/s7rss/storage/mongo"
)

var newsStorage *mongo.NewsStorage

func TestMain(m *testing.M) {
	if err := cmd.Execute(func(cmd *cobra.Command, args []string) {
		client, err := mongo.New(mongo.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Database: viper.GetString("db.database"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
		})
		if err != nil {
			panic(err)
		}

		newsStorage = mongo.NewNewsStorage(client)
	}); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
