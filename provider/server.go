package provider

import (
	"context"

	"github.com/qystishere/s7rss/storage"
)

func (p *Provider) GetNews(ctx context.Context, req *GetNewsRequest) (*GetNewsResponse, error) {
	news, err := p.newsStorage.Get(storage.GetNews{
		ID:      req.Id,
		Channel: req.Channel,
		Link:    req.Link,
	})
	if err != nil {
		return nil, err
	}
	getNewsResponseNews := make([]*GetNewsResponse_News, 0, len(news))
	for _, item := range news {
		getNewsResponseNews = append(getNewsResponseNews, &GetNewsResponse_News{
			Id:          item.ID,
			Channel:     item.Channel,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		})
	}
	return &GetNewsResponse{
		News: getNewsResponseNews,
	}, nil
}
