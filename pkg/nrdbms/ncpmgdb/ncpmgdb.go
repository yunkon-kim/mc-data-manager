package ncpmgdb

import (
	"context"

	"github.com/cloud-barista/cm-data-mold/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NCPMongoDBMS struct {
	provider utils.Provider
	dbName   string

	client *mongo.Client
	db     *mongo.Database
	ctx    context.Context
}

type NCPMongoDBOption func(*NCPMongoDBMS)

func New(client *mongo.Client, databaseName string, opts ...NCPMongoDBOption) *NCPMongoDBMS {
	dms := &NCPMongoDBMS{
		provider: utils.NCP,
		dbName:   databaseName,
		client:   client,
		ctx:      context.TODO(),
		db:       client.Database(databaseName),
	}

	for _, opt := range opts {
		opt(dms)
	}

	return dms
}

// list table
func (n *NCPMongoDBMS) ListTables() ([]string, error) {
	return n.db.ListCollectionNames(n.ctx, bson.D{})
}

// delete table
func (n *NCPMongoDBMS) DeleteTables(tableName string) error {
	return n.client.Database(n.dbName).Collection(tableName).Drop(n.ctx)
}

// create table
func (n *NCPMongoDBMS) CreateTable(tableName string) error {
	_, err := n.db.Collection(tableName).InsertOne(n.ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = n.db.Collection(tableName).DeleteOne(n.ctx, map[string]interface{}{})
	return err
}

// import table
func (n *NCPMongoDBMS) ImportTable(tableName string, srcData *[]map[string]interface{}) error {
	for _, data := range *srcData {
		_, err := n.db.Collection(tableName).InsertOne(n.ctx, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// export table
func (n *NCPMongoDBMS) ExportTable(tableName string, dstData *[]map[string]interface{}) error {
	cursor, err := n.db.Collection(tableName).Find(n.ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	defer cursor.Close(n.ctx)

	for cursor.Next(n.ctx) {
		var result map[string]interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return err
		}

		if oid, ok := result["_id"].(primitive.ObjectID); ok {
			result["_id"] = oid.Hex()
		}

		*dstData = append(*dstData, result)
	}
	return nil
}
