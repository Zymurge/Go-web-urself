package persistence

import (
	"fmt"
	"time"
	//"fmt"
	mgo "gopkg.in/mgo.v2"
	"webstuff/types"
)

// MongoAbstraction defines the set of DAL functions for accessing this Mongo collection
type MongoAbstraction interface {
	ConnectToMongo() error
	WriteCollection(collectionName string, object types.Loc) error
	UpdateCollection(collectionName string, object types.Loc) error
	FetchFromCollection(collectionName string, id string) (types.Loc, error)
	DeleteFromCollection(collectionName string, id string) error
}

// MongoSession defines an instantiation of a Mongo DAL. The session maintains a connected state to Mongodb.
type MongoSession struct {
	session        *mgo.Session
	db             *mgo.Database
	mongoURL       string
	dbName         string
	timeoutSeconds time.Duration
}

// DbName designates the default DB name in mongo
const (
	DefaultDbName  string        = "defaultDB"
	DefaultTimeout time.Duration = 10 * time.Second
)

// NewMongoSession is a factory method to create a fresh MongoSession for a given connection string and DB
func NewMongoSession(mongoURL string, dbName string, toDuration int64) *MongoSession {
	result := &MongoSession{
		mongoURL: mongoURL,
		dbName:   dbName,
	}
	if result.dbName == "" {
		result.dbName = DefaultDbName
	}
	result.timeoutSeconds = time.Duration(toDuration) * time.Second

	return result
}

// ConnectToMongo creates a connection to the specified mongodb instance
func (ms *MongoSession) ConnectToMongo() (err error) {
	//to := ms.timeoutSeconds != nil ? ms.timeoutSeconds : DefaultTimeout
	ms.session, err = mgo.DialWithTimeout(ms.mongoURL, ms.timeoutSeconds)
	if err != nil {
		return
	}
	ms.db = ms.session.DB(ms.dbName)
	return
}

// CheckAndReconnect ensures that there is an active DB connection to mongo. Attempts to reestablish connection if needed
func (ms *MongoSession) CheckAndReconnect() (err error) {
	if ms.db == nil {
		err = ms.ConnectToMongo()
	}
	return
}

// WriteCollection writes the specified loc object to a given collection
func (ms *MongoSession) WriteCollection(coll string, obj types.Loc) error {
	if err := ms.CheckAndReconnect(); err != nil {
		return err
	}
	myCollection := ms.db.C(coll)
	return myCollection.Insert(obj)
}

// UpdateCollection updates the loc object in the specified collection with a matching _id element to the passed in object
func (ms *MongoSession) UpdateCollection(collName string, obj types.Loc) error {
	if err := ms.CheckAndReconnect(); err != nil {
		return err
	}
	if !ms.collectionExists(collName) {
		return fmt.Errorf("Non-existent collection for update: %s", collName)
	}
	id := obj.GetID()
	myCollection := ms.db.C(collName)
	return myCollection.UpdateId(id, obj)
}

// FetchFromCollection fetches the Loc by ID from the specified collection
func (ms *MongoSession) FetchFromCollection(coll string, id string) (result types.Loc, err error) {
	result = types.Loc{}
	if err = ms.CheckAndReconnect(); err != nil {
		return
	}
	fmt.Printf("Calling mongo collection: %s\nDB = %v", coll, ms.db)
	myCollection := ms.db.C(coll)
	q := myCollection.FindId(id)
	err = q.One(&result)
	return
}

// DeleteFromCollection removes the Loc by ID from the specified collection
func (ms *MongoSession) DeleteFromCollection(coll string, id string) (err error) {
	if err := ms.CheckAndReconnect(); err != nil {
		return err
	}
	myCollection := ms.db.C(coll)
	return myCollection.RemoveId(id)
}

func (ms *MongoSession) collectionExists(collName string) bool {
	names, err := ms.db.CollectionNames()
	if err != nil { 
		return false 
	}
	for _, n := range names {
		if collName == n {
			return true
		}
	}
	return false
}
