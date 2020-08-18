package mongo

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	bsonTagName = "bson"

	readTimeout  = time.Second * 1
	writeTimeout = time.Second * 2
)

type Config struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Client struct {
	ctx context.Context
	*mongo.Database
}

func New(config Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var (
		client *mongo.Client
		err    error
	)
	if len(config.Username) > 0 && len(config.Password) > 0 {
		client, err = mongo.NewClient(options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database),
		))
	} else {
		client, err = mongo.NewClient(options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s:%d/%s", config.Host, config.Port, config.Database),
		))
	}
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return &Client{
		ctx: context.Background(),
		Database: client.Database(config.Database),
	}, nil
}

func bsonFilter(filter interface{}) bson.M {
	var (
		structValue = reflect.Indirect(reflect.ValueOf(filter))
		bsonMap     = make(bson.M)
	)
	for i := 0; i < structValue.NumField(); i++ {
		fieldType, fieldValue := structValue.Type().Field(i), structValue.Field(i)
		if tag, ok := fieldType.Tag.Lookup(bsonTagName); ok && !fieldValue.IsNil() {
			var (
				k = strings.Split(tag, ",")[0]
				v = reflect.Indirect(fieldValue).Interface()
			)
			switch fieldValue.Kind() {
			case reflect.Slice:
				bsonMap[k] = bson.M{"$in": v}
			default:
				bsonMap[k] = v
			}
		}
	}
	return bsonMap
}
