package cfdb

import (
	"context"
	"fmt"

	"github.com/ffandroid/autoprime/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CFClearance struct {
	ID         int64  `bson:"_id" json:"id"`
	Hostname   string `bson:"hostname" json:"hostname"`
	Clearance  string `bson:"clearance" json:"clearance"`
	UserAgent  string `bson:"userAgent" json:"userAgent"`
	LastUpdate int64  `bson:"lastUpdate" json:"lastUpdate"`
	IP         string `json:"ip" bson:"ip"`
}

var PrimeWireB2026 *mongo.Database

func init() {
	mongoDBHost := utils.GetEnv("MONGO_DB_HOST1", "mongodb+srv://mongo.zplay.online/?replicaSet=rs0&tls=false&readPreference=nearest&compressors=zstd&localThresholdMS=1")

	//mongoDBHost := "mongodb://100.64.243.103:27017"
	clientOptions := options.
		Client().
		SetMaxPoolSize(5).
		SetMinPoolSize(0).
		ApplyURI(mongoDBHost)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("error", err)
	}
	PrimeWireB2026 = client.Database("primeWireB2026")

}
