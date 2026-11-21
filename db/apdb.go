package db

import (
	"context"
	"log"

	"github.com/ffandroid/autoprime/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PResponse struct {
	Info    Info      `json:"info"`
	Servers []Servers `json:"servers"`
}
type Episode struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Episode     int    `json:"episode"`
	Season      int    `json:"season"`
	ReleaseDate string `json:"release_date"`
}
type Info struct {
	Status       string  `json:"status"`
	Type         string  `json:"type"`
	Description  string  `json:"description"`
	Title        string  `json:"title"`
	Episode      Episode `json:"episode"`
	ReleaseDate  string  `json:"release_date"`
	TvmazeID     int     `json:"tvmaze_id"`
	ImdbID       string  `json:"imdb_id"`
	TmdbID       string  `json:"tmdb_id"`
	TmdbBackdrop string  `json:"tmdb_backdrop"`
	TmdbImage    string  `json:"tmdb_image"`
}
type Servers struct {
	Name          string `json:"name"`
	FileName      any    `json:"file_name"`
	FileSize      any    `json:"file_size"`
	Key           string `json:"key"`
	Quality       any    `json:"quality"`
	AudioLanguage string `json:"audio_language"`
	AudioType     string `json:"audio_type"`
}

type PVideo struct {
	Name          string `json:"name"`
	FileName      any    `json:"file_name"`
	FileSize      any    `json:"file_size"`
	Key           string `json:"key"`
	Quality       any    `json:"quality"`
	AudioLanguage string `json:"audio_language"`
	AudioType     string `json:"audio_type"`
}

type PVideoInfo struct {
	ID     string `bson:"_id" json:"id"`
	Link   string `json:"link"`
	Host   string `json:"host"`
	HostID int    `json:"host_id"`
}

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
