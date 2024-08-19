package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

func (m mongoRepository) GetAllPhones(log zerolog.Logger, dbName string, collectionName string) (results t.PhoneNumberResults, err error) {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	m.logger.Info().Msgf("[storage-GetAllPhones] started session, db: %s and collection: %s", dbName, collectionName)

	ctx := context.Background()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		m.logger.Error().Msgf("[storage-GetAllPhones] Problem to generate the cursor: %s", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M

		if err := cur.Decode(&result); err != nil {
			m.logger.Error().Msgf("[storage-GetAllPhones] Problem to decode the result: %s", err.Error())
			return nil, errors.New(fmt.Sprintln("problem to decode the result: ", err.Error()))
		}

		phoneResult, err := PrepareGetOperation(result)
		if err != nil {
			m.logger.Error().Msgf("[storage-GetAllPhones] Failed to prepare the phoneResult: %s", err.Error())
			return nil, errors.New(fmt.Sprintln("Failed to prepare the phoneResult: ", err.Error()))
		}

		results = append(results, phoneResult)
	}

	m.logger.Info().Msg("[storage-GetAllPhones] got all data from Mongo!")

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
			return t.PhoneNumber{}, errors.New(fmt.Sprintln("error fetching phone record: ", err.Error()))
		}
	}

	phoneResult, err := PrepareGetOperation(result)
	if err != nil {
		m.logger.Error().Msgf("[storage-GetPhone] Failed to prepare the phoneResult: %s", err.Error())
		return t.PhoneNumber{}, errors.New(fmt.Sprintln("Failed to prepare the phoneResult: ", err.Error()))
	}

	m.logger.Info().Msgf("[storage-GetPhone] got the data in Mongo from phone numer %s!", phoneNumber)
	return phoneResult, nil
}

func (m mongoRepository) UpdatePhoneRecord(log zerolog.Logger, data t.PhoneNumber, dbName string, collectionName string) error {
	session := m.session
	collection := session.Client().Database(dbName).Collection(collectionName)

	m.logger.Info().Msgf("[storage-UpdatePhoneRecord] started session, db: %s and collection: %s", dbName, collectionName)

	ctx := context.Background()

	if data.PhoneInput == "" {
		m.logger.Error().Msg("[storage-UpdatePhoneRecord] phone number is null")
		return errors.New("phone number is null")

	}

	if data.ID == "" {
		data.ID = GetMD5Hash(data.PhoneInput)
	}

	m.logger.Info().Msgf("[storage-UpdatePhoneRecord] updating this record: %+v", data)

	// Fetch the existing document
	filter := bson.M{"_id": data.ID}
	var existingDocument t.PhoneNumber

	err := collection.FindOne(ctx, filter).Decode(&existingDocument)
	if err != nil {
		m.logger.Error().Msgf("[storage-UpdatePhoneRecord] Error fetching existing document: %s", err.Error())
		return err
	}

	c := PrepareUpdateOperation(existingDocument.Country, data.Country)

	update := bson.M{
		"$set": bson.M{
			"phoneInput":          data.PhoneInput,
			"isValid":             data.IsValid,
			"e164Format":          data.E164Format,
			"internationalFormat": data.InternationalFormat,
			"nationalFormat":      data.NationalFormat,
			"location":            data.Location,
			"lineType":            data.LineType,
			"country":             c,
		},
	}
	_, errUpd := collection.UpdateOne(ctx, filter, update)

	if errUpd != nil {
		m.logger.Error().Msgf("[storage-UpdatePhoneRecord] Problem to update data on MongoDB: %s", errUpd.Error())
		return err
	}

	m.logger.Info().Msg("[storage-UpdatePhoneRecord] data updated in Mongo")
	return nil
}

func PrepareGetOperation(result bson.M) (t.PhoneNumber, error) {
	output := make(map[string]interface{})

	for k, v := range result {
		if result[k] != nil {
			output[k] = v
		}

		if k == "country" {
			for k2, v2 := range result[k].(bson.M) {
				if mapValue, ok := v2.(map[string]interface{}); ok {
					output["country"].(map[string]interface{})[k2] = mapValue
				}
			}
		}
	}

	phoneNumberBytes, err := json.Marshal(output)
	if err != nil {
		return t.PhoneNumber{}, errors.New(fmt.Sprintln("Failed to marshal output: ", err.Error()))
	}

	var phoneResult t.PhoneNumber
	err = json.Unmarshal(phoneNumberBytes, &phoneResult)
	if err != nil {
		return t.PhoneNumber{}, errors.New(fmt.Sprintln("Failed to unmarshal output: ", err.Error()))
	}

	return phoneResult, nil
}

func PrepareUpdateOperation(existing, updated t.Country) bson.M {

	setFields := unmarshalData(existing)

	if existing.M49Code != updated.M49Code {
		setFields["m49Code"] = updated.M49Code
	}
	if existing.IsoAlpha2 != updated.IsoAlpha2 {
		setFields["isoAlpha2"] = updated.IsoAlpha2
	}
	if existing.IsoAlpha3 != updated.IsoAlpha3 {
		setFields["isoAlpha3"] = updated.IsoAlpha3
	}
	if existing.Name != updated.Name {
		setFields["name"] = updated.Name
	}
	if existing.IsoName != updated.IsoName {
		setFields["isoName"] = updated.IsoName
	}
	if existing.IsoNameFull != updated.IsoNameFull {
		setFields["isoNameFull"] = updated.IsoNameFull
	}
	if existing.UnRegion != updated.UnRegion {
		setFields["unRegion"] = updated.UnRegion
	}

	if existing.Currency.NumericCode != updated.Currency.NumericCode {
		setFields["currency"].(map[string]interface{})["numericCode"] = updated.Currency.NumericCode
	}
	if existing.Currency.Code != updated.Currency.Code {
		setFields["currency"].(map[string]interface{})["code"] = updated.Currency.Code
	}
	if existing.Currency.Name != updated.Currency.Name {
		setFields["currency"].(map[string]interface{})["name"] = updated.Currency.Name
	}
	if existing.Currency.MinorUnits != updated.Currency.MinorUnits {
		setFields["currency"].(map[string]interface{})["minorUnits"] = updated.Currency.MinorUnits
	}

	if existing.WbRegion.ID != updated.WbRegion.ID {
		setFields["wbRegion"].(map[string]interface{})["id"] = updated.WbRegion.ID
	}
	if existing.WbRegion.Iso2Code != updated.WbRegion.Iso2Code {
		setFields["wbRegion"].(map[string]interface{})["iso2Code"] = updated.WbRegion.Iso2Code
	}
	if existing.WbRegion.Value != updated.WbRegion.Value {
		setFields["wbRegion"].(map[string]interface{})["value"] = updated.WbRegion.Value
	}

	if existing.WbIncomeLevel.ID != updated.WbIncomeLevel.ID {
		setFields["wbIncomeLevel"].(map[string]interface{})["id"] = updated.WbIncomeLevel.ID
	}
	if existing.WbIncomeLevel.Iso2Code != updated.WbIncomeLevel.Iso2Code {
		setFields["wbIncomeLevel"].(map[string]interface{})["iso2Code"] = updated.WbIncomeLevel.Iso2Code
	}
	if existing.WbIncomeLevel.Value != updated.WbIncomeLevel.Value {
		setFields["wbIncomeLevel"].(map[string]interface{})["value"] = updated.WbIncomeLevel.Value
	}

	if existing.CallingCode != updated.CallingCode {
		setFields["callingCode"] = updated.CallingCode
	}
	if existing.CountryFlagEmoji != updated.CountryFlagEmoji {
		setFields["countryFlagEmoji"] = updated.CountryFlagEmoji
	}
	if existing.WikidataID != updated.WikidataID {
		setFields["wikidataId"] = updated.WikidataID
	}
	if existing.GeonameID != updated.GeonameID {
		setFields["geonameId"] = updated.GeonameID
	}
	if existing.IsIndependent != updated.IsIndependent {
		setFields["isIndependent"] = updated.IsIndependent
	}

	if len(existing.IsoAdminLanguages) != len(updated.IsoAdminLanguages) {
		setFields["isoAdminLanguages"] = updated.IsoAdminLanguages
	} else {
		for i := range updated.IsoAdminLanguages {
			if existing.IsoAdminLanguages[i].IsoAlpha3 != updated.IsoAdminLanguages[i].IsoAlpha3 {
				setFields["isoAdminLanguages"].([]interface{})[i].(map[string]interface{})["isoAlpha3"] = updated.IsoAdminLanguages[i].IsoAlpha3
			}

			if existing.IsoAdminLanguages[i].IsoAlpha2 != updated.IsoAdminLanguages[i].IsoAlpha2 {
				setFields["isoAdminLanguages"].([]interface{})[i].(map[string]interface{})["isoAlpha2"] = updated.IsoAdminLanguages[i].IsoAlpha2
			}

			if existing.IsoAdminLanguages[i].IsoName != updated.IsoAdminLanguages[i].IsoName {
				setFields["isoAdminLanguages"].([]interface{})[i].(map[string]interface{})["isoName"] = updated.IsoAdminLanguages[i].IsoName
			}

			if existing.IsoAdminLanguages[i].NativeName != updated.IsoAdminLanguages[i].NativeName {
				setFields["isoAdminLanguages"].([]interface{})[i].(map[string]interface{})["nativeName"] = updated.IsoAdminLanguages[i].NativeName
			}
		}
	}

	return setFields
}

func unmarshalData(data t.Country) map[string]interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	d := bson.M{}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return nil
	}

	return d
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
