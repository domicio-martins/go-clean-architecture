package wallet

import (
"errors"
"fmt"
	"os/user"
	"time"

"github.com/PicPay/picpay-dev-ms-template-manager/pkg/newrelic"

"github.com/PicPay/picpay-dev-ms-template-manager/pkg/log"
. "github.com/gobeam/mongo-go-pagination"
"github.com/google/uuid"
"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/bson/primitive"
"go.mongodb.org/mongo-driver/mongo/options"

"go.mongodb.org/mongo-driver/mongo"
)

type WalletRepository struct {
	*mongo.Database
	newrelic.NewRelic
}

func NewRepository(database *mongo.Database, monitor newrelic.NewRelic) *WalletRepository {
	return &WalletRepository{
		database,
		monitor,
	}
}

func (r *WalletRepository) Update(user *user.User, amount float64) (*Wallet, error) {
	wallet, err := r.find(bson.M{"user": user.ID})
	if err != nil {
		log.Error(fmt.Sprintf("Error to find screen to update, ScreenIdentifier: %s. Error: ", identifier), err, nil)
		return nil, errors.New("screen identifier not found")
	}

	wallet.Amount += amount

	r.NewRelic.FinishTransaction()
	if err != nil {
		return nil, err
	}

	screen, err = r.findScreen(bson.M{"_id": res.InsertedID.(primitive.ObjectID)})
	return &screen, err
}

func (r *WalletRepository) find(filter interface{}) (Screen, error) {
	var screen Screen
	screenCol := r.Collection("screen")
	err := screenCol.FindOne(
		r.NewRelic.StartTransaction(),
		filter,
		options.FindOne().SetSort(bson.D{{"_id", -1}}),
	).Decode(&screen)
	r.NewRelic.FinishTransaction()
	return screen, err
}
