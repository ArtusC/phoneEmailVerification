package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	// log     log.Logger
	session mongo.Session
}

func NewMongoRepository(session mongo.Session) MongoRepository {
	return mongoRepository{
		session: session,
	}
}

func (m mongoRepository) StoragePhoneRecord(data map[string]interface{}, dbName string, collectionName string) error {

	if len(data) == 0 {
		fmt.Println("The data is empty")
		return nil
	}

	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	ctx := context.Background()
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		fmt.Println("Problem to insert data on MongoDB: ", err.Error())
		return err
	}

	fmt.Println("RESULT: ", result)

	return nil
}

func (m mongoRepository) GetPhoneRecords(dbName string, collectionName string) (phoneResults map[string]interface{}, err error) {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	ctx := context.Background()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Problem to generate the cursor: ", err.Error())
		return nil, err
	}

	var results []bson.M

	for cur.Next(ctx) {
		var result bson.M

		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, result)

	}

	phoneResults = map[string]interface{}{
		"phone": results,
	}

	return phoneResults, nil

}
