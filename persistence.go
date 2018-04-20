package webstuff

import (
	mgo "gopkg.in/mgo.v2"
)

// DbName designates the default DB name in mongo
const (
	DbName string = "testDB"
)

// ConnectToMongo creates a connection to the specified mongodb instance
func ConnectToMongo( connectString string ) (*mgo.Session, error) {
	var session *mgo.Session
	session, err := mgo.Dial( connectString )
	return session, err
}

// WriteCollection writes the specified loc object to a given collection
func WriteCollection( session *mgo.Session, coll string, obj Loc ) error {
	myCollection := session.DB(DbName).C(coll)
	return myCollection.Insert(obj)
}

// UpdateCollection updates the loc object in the specified collection with a matching _id element to the passed in object
func UpdateCollection( session *mgo.Session, coll string, obj Loc ) error {
	myCollection := session.DB(DbName).C(coll)
	id := obj.ID
	return myCollection.UpdateId( id, obj )
}
