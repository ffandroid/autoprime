package db

import (
	"context"
	"log"

	"github.com/ffandroid/autoprime/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BShow struct {
	Title  string `json:"title" bson:"title"`
	ID     string `json:"id" bson:"_id"`
	ImdbID string `json:"imdbId" bson:"imdbId"`
	TmdbID int64  `json:"tmdbId" bson:"tmdbId"`
	Year   string `json:"year" bson:"year"`
}

var APrimeDB *mongo.Database

func init() {
	mongoDBHost := utils.GetEnv("MONGO_DB_APRIME", "mongodb://100.125.175.51:27025")

	clientOptions := options.
		Client().
		SetMaxPoolSize(5).
		SetMinPoolSize(0).
		ApplyURI(mongoDBHost)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error", err)
	}
	APrimeDB = client.Database("aPrimeDB")

}
