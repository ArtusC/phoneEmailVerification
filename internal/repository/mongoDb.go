package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	t "github.com/ArtusC/phoneEmailVerification/types"
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

func (m mongoRepository) StoragePhoneRecord(data t.PhoneNumber, dbName string, collectionName string) error {

	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	// dataInt := data.(map[string]interface{})

	phoneNumber := data.PhoneInput
	data.ID = GetMD5Hash(phoneNumber)

	ctx := context.Background()
	_, err := collection.InsertOne(ctx, data)

	if err != nil {
		fmt.Println("Problem to insert data on MongoDB: ", err.Error())
		return err
	}

	return nil
}

func (m mongoRepository) GetPhoneRecords(dbName string, collectionName string) (results t.PhoneNumberResults, err error) {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	ctx := context.Background()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("Problem to generate the cursor: ", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M

		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		phoneNumber := t.PhoneNumber{
			ID:                  result["_id"].(string),
			PhoneInput:          result["phoneInput"].(string),
			IsValid:             result["isValid"].(bool),
			E164Format:          result["e164Format"].(string),
			InternationalFormat: result["internationalFormat"].(string),
			NationalFormat:      result["nationalFormat"].(string),
			Location:            result["location"].(string),
			LineType:            result["lineType"].(string),
			Country: t.Country{
				IsoAlpha2:        result["country"].(bson.M)["isoAlpha2"].(string),
				IsoAlpha3:        result["country"].(bson.M)["isoAlpha3"].(string),
				Name:             result["country"].(bson.M)["name"].(string),
				IsoName:          result["country"].(bson.M)["isoName"].(string),
				IsoNameFull:      result["country"].(bson.M)["isoNameFull"].(string),
				UnRegion:         result["country"].(bson.M)["unRegion"].(string),
				CallingCode:      int32(result["country"].(bson.M)["callingCode"].(int32)),
				CountryFlagEmoji: result["country"].(bson.M)["countryFlagEmoji"].(string),
				WikidataID:       result["country"].(bson.M)["wikidataId"].(string),
				GeonameID:        result["country"].(bson.M)["geonameId"].(string),
				IsIndependent:    result["country"].(bson.M)["isIndependent"].(bool),
				Currency: t.Currency{
					NumericCode: int32(result["country"].(bson.M)["currency"].(bson.M)["numericCode"].(int32)),
					Code:        result["country"].(bson.M)["currency"].(bson.M)["code"].(string),
					Name:        result["country"].(bson.M)["currency"].(bson.M)["name"].(string),
					MinorUnits:  int32(result["country"].(bson.M)["currency"].(bson.M)["minorUnits"].(int32)),
				},
				WbRegion: t.WbRegion{
					ID:       result["country"].(bson.M)["wbRegion"].(bson.M)["id"].(string),
					Iso2Code: result["country"].(bson.M)["wbRegion"].(bson.M)["iso2Code"].(string),
					Value:    result["country"].(bson.M)["wbRegion"].(bson.M)["value"].(string),
				},
				WbIncomeLevel: t.WbIncomeLevel{
					ID:       result["country"].(bson.M)["wbIncomeLevel"].(bson.M)["id"].(string),
					Iso2Code: result["country"].(bson.M)["wbIncomeLevel"].(bson.M)["iso2Code"].(string),
					Value:    result["country"].(bson.M)["wbIncomeLevel"].(bson.M)["value"].(string),
				},
			},
		}

		results = append(results, phoneNumber)
	}

	return results, nil
}

func (m mongoRepository) UpdatePhoneRecord(data map[string]interface{}, dbName string, collectionName string) error {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	if len(data) == 0 || data["_id"] == nil {
		return errors.New("[UpdatePhoneRecord] invalid data")
	}

	id := data["_id"].(string)
	ctx := context.Background()

	// Fetch the existing document
	filter := bson.M{"_id": id}
	var existingDocument bson.M
	err := collection.FindOne(ctx, filter).Decode(&existingDocument)
	if err != nil {
		fmt.Println("Error fetching existing document:", err.Error())
		return err
	}

	// Check if all keys in data exist in the existing document
	for key := range data {
		if _, ok := existingDocument[key]; !ok {
			return fmt.Errorf("key '%s' does not exist in the existing document", key)
		}
	}

	update := bson.M{"$set": data}

	_, errUpd := collection.UpdateOne(ctx, filter, update)
	if errUpd != nil {
		fmt.Println("Problem to update data on MongoDB: ", errUpd.Error())
		return err
	}

	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
