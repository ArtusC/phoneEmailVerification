package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	t "github.com/ArtusC/phoneEmailVerification/types"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	logger  zerolog.Logger
	session mongo.Session
}

func NewMongoRepository(logger zerolog.Logger, session mongo.Session) MongoRepository {
	return mongoRepository{
		logger:  logger,
		session: session,
	}
}

func (m mongoRepository) StoragePhoneRecord(log zerolog.Logger, data t.PhoneNumber, dbName string, collectionName string) error {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	m.logger.Info().Msgf("[storage-StoragePhoneRecord] started session, db: %s and collection: %s", dbName, collectionName)

	phoneNumber := data.PhoneInput
	data.ID = GetMD5Hash(phoneNumber)
	m.logger.Info().Msgf("[storage-StoragePhoneRecord] data.ID: %s", data.ID)

	ctx := context.Background()
	_, err := collection.InsertOne(ctx, data)

	if err != nil {
		m.logger.Error().Msgf("[storage-StoragePhoneRecord] Problem to insert data on MongoDB: %s", err.Error())
		return err
	}

	m.logger.Info().Msg("[storage-StoragePhoneRecord] data inserted in Mongo!")
	return nil
}

func (m mongoRepository) GetAllPhoneRecords(log zerolog.Logger, dbName string, collectionName string) (results t.PhoneNumberResults, err error) {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	m.logger.Info().Msgf("[storage-GetAllPhoneRecords] started session, db: %s and collection: %s", dbName, collectionName)

	ctx := context.Background()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		m.logger.Error().Msgf("[storage-GetAllPhoneRecords] Problem to generate the cursor: %s", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M

		if err := cur.Decode(&result); err != nil {
			m.logger.Error().Msgf("[storage-GetAllPhoneRecords] Problem to decode the result: %s", err.Error())
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

	m.logger.Info().Msg("[storage-GetAllPhoneRecords] got all data from Mongo!")
	return results, nil
}

func (m mongoRepository) GetPhone(log zerolog.Logger, dbName, collectionName, phoneNumber string) (t.PhoneNumber, error) {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	m.logger.Info().Msgf("[storage-GetPhone] started session, db: %s and collection: %s", dbName, collectionName)

	ctx := context.Background()
	filter := bson.M{"phoneInput": phoneNumber}
	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Info().Msg("[storage-GetPhone] document not found")
			return t.PhoneNumber{}, nil
		} else {
			m.logger.Error().Msgf("[storage-GetPhone] error fetching phone record: %s", err.Error())
			return t.PhoneNumber{}, errors.New(fmt.Sprint("error fetching phone record: ", err.Error()))
		}
	}

	phoneResult := t.PhoneNumber{
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

	m.logger.Info().Msgf("[storage-GetPhone] got the data in Mongo from phone numer %s!", phoneNumber)
	return phoneResult, nil
}

func (m mongoRepository) UpdatePhoneRecord(log zerolog.Logger, data map[string]interface{}, dbName string, collectionName string) error {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	m.logger.Info().Msgf("[storage-UpdatePhoneRecord] started session, db: %s and collection: %s", dbName, collectionName)

	if len(data) == 0 || data["_id"] == nil {
		m.logger.Error().Msgf("[storage-UpdatePhoneRecord] invalid data: %v", data)
		return errors.New("[UpdatePhoneRecord] invalid data")
	}

	id := data["_id"].(string)
	m.logger.Info().Msgf("[storage-UpdatePhoneRecord] updating this ID: %s", id)

	ctx := context.Background()

	// Fetch the existing document
	filter := bson.M{"_id": id}
	var existingDocument bson.M
	err := collection.FindOne(ctx, filter).Decode(&existingDocument)
	if err != nil {
		m.logger.Error().Msgf("[storage-UpdatePhoneRecord] Error fetching existing document: %s", err.Error())
		return err
	}

	// Check if all keys in data exist in the existing document
	for key := range data {
		if _, ok := existingDocument[key]; !ok {
			m.logger.Error().Msgf("[storage-UpdatePhoneRecord]key '%s' does not exist", key)
			return fmt.Errorf("key '%s' does not exist", key)
		}
	}

	update := bson.M{"$set": data}

	_, errUpd := collection.UpdateOne(ctx, filter, update)
	if errUpd != nil {
		m.logger.Error().Msgf("[storage-UpdatePhoneRecord] Problem to update data on MongoDB: %s", errUpd.Error())
		return err
	}

	m.logger.Info().Msg("[storage-UpdatePhoneRecord] data updated in Mongo")
	return nil
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
