package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/qystishere/s7rss/cmd"
	"github.com/qystishere/s7rss/parser"
	"github.com/qystishere/s7rss/provider"
	"github.com/qystishere/s7rss/storage/mongo"
)

//go:generate protoc --plugin=protoc-gen-gofast=$GOPATH/bin/protoc-gen-gofast -I=$GOPATH/src/github.com/qystishere/s7rss/resources --gofast_out=plugins=grpc:../../provider provider.proto

func main() {
	if err := cmd.Execute(run); err != nil {
		log.WithError(err).
			Fatal("Execute")
	}
}

func run(cmd *cobra.Command, args []string) {
	client, err := mongo.New(mongo.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		log.WithError(err).
			Fatal("Database connect")
	}
	newsStorage := mongo.NewNewsStorage(client)

	feedParser := parser.New(&parser.Config{
		FeedsURLs: viper.GetStringSlice("collection.feeds"),
		Words:     viper.GetStringSlice("collection.words"),
		Threads:   viper.GetInt("collection.threads"),
		Timeout:   viper.GetInt("collection.timeout"),

		NewsStorage: newsStorage,
	})
	feedProvider := provider.New(newsStorage)

	go func() {
		go feedParser.Start()

		if err = feedProvider.Listen(viper.GetString("grpc.host"), viper.GetInt("grpc.port")); err != nil {
			log.WithError(err).
				Fatal("GRPC service listen")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	log.Debug("Stopping...")
	feedParser.Stop()
	log.Debug("Parser stopped")
	feedProvider.Stop()
	log.Debug("Provider stopped")
}
