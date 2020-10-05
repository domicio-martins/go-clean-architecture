package mongodb

import (
	"context"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func IDFromResult(insertedID interface{}) string {
	switch v := insertedID.(type) {
	case string:
		return v
	case primitive.ObjectID:
		return v.Hex()
	default:
		return ""
	}
}

func Open(uri string, newrelic bool) (*mongo.Client, error) {
	conf := options.Client().ApplyURI(uri)
	if newrelic {
		conf.SetMonitor(nrmongo.NewCommandMonitor(nil))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, conf)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "timeout: could not connect to mongo")
	}

	return client, err
}

func CreateIndexes(col *mongo.Collection, indexes []mongo.IndexModel) error {
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	_, err := col.Indexes().CreateMany(context.TODO(), indexes, opts)
	return err
}
