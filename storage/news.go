package storage

type NewsStorager interface {
	Store(input StoreNews) (*News, error)
	Get(input GetNews) ([]*News, error)
}

type News struct {
	ID          string `json:"id" bson:"_id"`
	Channel     string `json:"channel" bson:"channel"`
	Title       string `json:"title" bson:"title"`
	Link        string `json:"link" bson:"link"`
	Description string `json:"description" bson:"description"`
}

type StoreNews struct {
	*News `bson:",inline"`
}

type GetNews struct {
	ID      []string `bson:"_id"`
	Channel []string `bson:"channel"`
	Link    []string `bson:"link"`
}
