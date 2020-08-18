package main

import (
	"github.com/apex/log"

	"github.com/qystishere/s7rss/cmd"
)

func main() {
	if err := cmd.Execute(nil); err != nil {
		log.WithError(err).
			Fatal("Execute")
	}
}
