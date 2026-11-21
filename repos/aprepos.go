package repos

import (
	"context"
	"log"
	"time"

	"github.com/ffandroid/autoprime/db"
	"go.mongodb.org/mongo-driver/bson"
)

func FindPVideoInfo(key string) db.PVideoInfo {

	var pVideoInfo db.PVideoInfo
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{"_id": key}
	err := db.APrimeDB.Collection("pVideoInfo").FindOne(ctx, filter).Decode(&pVideoInfo)
	if err != nil {
		log.Println("Error saving pVideoInfo:", err)
	}
	return pVideoInfo
}

func SavePVideoInfo(pVideoInfo db.PVideoInfo) db.PVideoInfo {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := db.APrimeDB.Collection("pVideoInfo").InsertOne(ctx, &pVideoInfo)
	if err != nil {
		log.Println("Error saving pVideoInfo:", err)
	}
	return pVideoInfo
}
