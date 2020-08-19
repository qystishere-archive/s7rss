package parser

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/mmcdole/gofeed"

	"github.com/qystishere/s7rss/storage"
)

var timeout = time.Second * 1

type Config struct {
	FeedsURLs []string
	Words     []string
	Threads   int
	Timeout   int

	NewsStorage storage.NewsStorager
}

type Parser struct {
	feeds []*Feed

	wg   *sync.WaitGroup
	stop bool

	sync.Mutex
	*Config
}

func New(config *Config) *Parser {
	log.Debugf("Feeds: %s", config.FeedsURLs)
	log.Debugf("Words: %s", config.Words)
	feeds := make([]*Feed, 0, len(config.FeedsURLs))
	for _, feedURL := range config.FeedsURLs {
		feeds = append(feeds, &Feed{
			URL: feedURL,

			UpdatedAt: time.Time{},
		})
	}
	return &Parser{
		feeds: feeds,

		Config: config,
	}
}

func (p *Parser) Start() {
	var (
		feedParser = gofeed.NewParser()

		wg = &sync.WaitGroup{}
	)
	p.wg = wg
	p.stop = false

	wg.Add(p.Threads)
	for i := 0; i < p.Threads; i++ {
		go func() {
			defer wg.Done()

			for range time.NewTicker(timeout).C {
				if p.stop {
					return
				}
				p.Lock()
				var feed *Feed
				for _, f := range p.feeds {
					if f.Processing {
						continue
					}

					if time.Since(f.UpdatedAt).Minutes() >= float64(p.Config.Timeout) {
						feed = f
					}
				}
				if feed != nil {
					feed.Processing = true
					p.Unlock()
				} else {
					p.Unlock()
					continue
				}

				fp, err := feedParser.ParseURL(feed.URL)
				if err != nil {
					log.WithError(err).
						WithField("feedUrl", feed.URL).
						Error("Parse feed from url")
					continue
				}

				news := make([]*storage.News, 0)
				for _, item := range fp.Items {
					var found bool
					for _, word := range p.Config.Words {
						if strings.Contains(item.Title, word) || strings.Contains(item.Description, word) {
							found = true
							break
						}
					}
					if !found {
						continue
					}
					news = append(news, &storage.News{
						ID:          fmt.Sprintf("%x", md5.Sum([]byte(item.GUID))),
						Channel:     fp.Link,
						Title:       item.Title,
						Link:        item.Link,
						Description: item.Description,
					})
				}

				if len(news) > 0 {
					log.Debugf("Got '%d' articles with collection words from '%s'", len(news), fp.Title)
				}

				success := true
				if err := p.NewsStorage.Batch(storage.BatchNews{
					News: news,
				}); err != nil {
					if combinedError, ok := err.(*storage.CombinedError); ok {
						for _, err := range combinedError.Errors {
							if errors.Is(err, storage.ErrNewsAlreadyExists) {
								continue
							}

							log.WithError(err).
								Error("Batching news")
							success = false
						}
					} else {
						log.WithError(err).
							Error("Batching news")
						success = false
					}
				}

				p.Lock()
				if success {
					feed.UpdatedAt = time.Now()
				}
				feed.Processing = false
				p.Unlock()
			}
		}()
	}
	wg.Wait()
}

func (p *Parser) Stop() {
	p.stop = true
	p.wg.Wait()
}
