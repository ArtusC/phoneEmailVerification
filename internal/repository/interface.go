package repository

type MongoRepository interface {
	StoragePhoneRecord(data map[string]interface{}, dbName string, collectionName string) error
	GetPhoneRecords(dbName string, collectionName string) (results map[string]interface{}, err error)
}
