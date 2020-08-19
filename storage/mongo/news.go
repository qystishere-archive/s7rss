package mongo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/qystishere/s7rss/storage"
)

const (
	newsCollectionName = "news"

	newsBatchTimeout  = time.Second * 5
	newsStoreTimeout  = time.Second * 2
	newsGetTimeout    = time.Second * 1
)

type NewsStorage struct {
	News *mongo.Collection
	*Client
}

func NewNewsStorage(c *Client) *NewsStorage {
	return &NewsStorage{
		News:   c.Collection(newsCollectionName),
		Client: c,
	}
}

func (ns *NewsStorage) Batch(input storage.BatchNews) error {
	ctx, cancel := context.WithTimeout(ns.ctx, newsBatchTimeout)
	defer cancel()

	var documents []interface{}
	for _, news := range input.News {
		documents = append(documents, news)
	}

	_, err := ns.News.InsertMany(ctx, documents)
	if err, ok := err.(*mongo.BulkWriteException); ok {
		var errors []error
		for _, bulkWriteError := range err.WriteErrors {
			if strings.Contains(bulkWriteError.Error(), "duplicate") {
				errors = append(errors,
					fmt.Errorf("%w: %s", storage.ErrNewsAlreadyExists, bulkWriteError.Error()),
				)
				continue
			}
			errors = append(errors, bulkWriteError)
		}
		if len(errors) > 0 {
			return storage.CombineErrors(errors...)
		}
	} else {
		if err != nil {
			return err
		}
	}
	return nil
}

func (ns *NewsStorage) Store(input storage.StoreNews) error {
	ctx, cancel := context.WithTimeout(ns.ctx, newsStoreTimeout)
	defer cancel()

	_, err := ns.News.InsertOne(ctx, input.News)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%w: %s", storage.ErrNewsAlreadyExists, err.Error())
		}
		return err
	}

	return nil
}

func (ns *NewsStorage) Get(input storage.GetNews) ([]*storage.News, error) {
	ctx, cancel := context.WithTimeout(ns.ctx, newsGetTimeout)
	defer cancel()

	cur, err := ns.News.Find(ctx, bsonFilter(input))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var news []*storage.News
	for cur.Next(nil) {
		var article storage.News
		if err := cur.Decode(&article); err != nil {
			return nil, err
		}

		news = append(news, &article)
	}

	return news, err

}