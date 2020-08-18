package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/qystishere/s7rss/storage"
)

type NewsSuite struct {
	news []*storage.News

	suite.Suite
}

func (s *NewsSuite) SetupSuite() {
	err := newsStorage.News.Drop(context.Background())
	s.Require().NoError(err)
}

func (s *NewsSuite) Test1Store() {
	s.news = []*storage.News{{
		ID:          "id",
		Channel:     "channel",
		Title:       "title",
		Link:        "link",
		Description: "description",
	}, {
		ID:          "id2",
		Channel:     "channel2",
		Title:       "title2",
		Link:        "link2",
		Description: "description2",
	}}

	for _, article := range s.news {
		_, err := newsStorage.Store(storage.StoreNews{
			News: article,
		})
		s.NoError(err)
	}

}

func (s *NewsSuite) Test2Get() {
	news, err := newsStorage.Get(storage.GetNews{})
	s.Require().NoError(err)
	s.ElementsMatch(news, s.news)

	for _, article := range s.news {
		news, err = newsStorage.Get(storage.GetNews{
			ID:      []string{article.ID},
			Channel: []string{article.Channel},
			Link:    []string{article.Link},
		})
		s.NoError(err)
		s.ElementsMatch(news, []*storage.News{article})
	}
}

func TestNewsStorage(t *testing.T) {
	suite.Run(t, new(NewsSuite))
}