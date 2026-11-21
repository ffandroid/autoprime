package cfrepos

import (
	"context"
	"time"

	"github.com/ffandroid/autoprime/cfp/cfdb"
	"go.mongodb.org/mongo-driver/bson"
)

func FindCFClearanceByHostname(hostname string) cfdb.CFClearance {

	var cFClearance cfdb.CFClearance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"hostname": hostname}

	cfdb.PrimeWireB2026.
		Collection("cFClearance").
		FindOne(ctx, filter).Decode(&cFClearance)

	return cFClearance
}
