package gmongo

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
	"time"
)

var (
	database *Database
	c        *Collection
)

func TestMain(m *testing.M) {
	var err error
	client := NewClient()
	database, err = client.NewDatabase("", "", SetMaxPoolSize(10), SetMaxConnIdleTime(5*time.Second), SetMinPoolSize(1))
	if err != nil {
		log.Fatal(err)
	}
	c, err = database.GetCollection("")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

type Example struct {
	UserId string  `bson:"user_id"`
	Items  []Items `bson:"items"`
}
type Items struct {
	Id   string `bson:"id"`
	Type uint8  `bson:"type"`
}

func TestCollection_UpdateOne(t *testing.T) {
	filter := bson.M{"user_id": "1", "items.id": "Dsd0daX1"}
	update := Example{
		Items: []Items{
			{
				Type: 0,
			},
		},
	}
	err := c.UpdateOne(filter, update)
	assert.Nil(t, err)
}

func TestCollection_InsertOne(t *testing.T) {
	dom := Example{
		UserId: "1",
		Items: []Items{
			{
				Id:   "Dsd0daX1",
				Type: 1,
			},
			{
				Id:   "GpdH8J3R",
				Type: 1,
			},
		},
	}
	objectID, _ := c.InsertOne(dom)
	assert.NotNil(t, objectID)

}
