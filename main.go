package db

import (
	"context"

	"github.com/ffandroid/proj0425/logx"
	"github.com/ffandroid/proj1024go/commons/utils"
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

var BstsrsDB2025 *mongo.Database

func init() {
	mongoDBHost := utils.GetEnv("MONGO_DB_HOST1", "mongodb+srv://mongo.zplay.online/?replicaSet=rs0&tls=false&readPreference=nearest&compressors=zstd&localThresholdMS=1")

	clientOptions := options.
		Client().
		SetMaxPoolSize(5).
		SetMinPoolSize(0).
		ApplyURI(mongoDBHost)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logx.Println("error", err)
	}
	BstsrsDB2025 = client.Database("bstsrsDB2025")

}
