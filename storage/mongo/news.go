package mongo

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/qystishere/s7rss/storage"
)

var newsCollectionName = "news"

type NewsStorage struct {
	News *mongo.Collection
	*Client
}

func NewNewsStorage(c *Client) (*NewsStorage, error) {
	return &NewsStorage{
		News:   c.Collection(newsCollectionName),
		Client: c,
	}, nil
}

func (ns *NewsStorage) Store(input storage.StoreNews) (*storage.News, error) {
	ctx, cancel := context.WithTimeout(ns.ctx, writeTimeout)
	defer cancel()

	_, err := ns.News.InsertOne(ctx, input.News)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, storage.ErrNewsAlreadyExists
		}
		return nil, err
	}

	return input.News, nil
}

func (ns *NewsStorage) Get(input storage.GetNews) ([]*storage.News, error) {
	ctx, cancel := context.WithTimeout(ns.ctx, readTimeout)
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
