package main

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"webstuff/types"
)

// MockMongoSession provides a mock abstraction to mongo
type MockMongoSession struct {
	connectMode string
	queryMode   string
	writeMode   string
}

/*
type mgoError struct {
	Code int
	message string
}

func (m *mgoError) Error() string {
	return m.message
}
*/

// ConnectToMongo mock. Controlled by mm.connectMode values 'positive' and 'no connect'
func (mm *MockMongoSession) ConnectToMongo() error {
	switch {
	case mm.connectMode == "positive":
		return nil
	case mm.connectMode == "no connect":
		return fmt.Errorf("mocked connection failure")
	}
	return fmt.Errorf("Unknown mode for ConnectToMongo: %s", mm.connectMode)
}

// WriteCollection mock. Controlled by mm.writeMode values 'positive', 'fail' and 'duplicate'
func (mm *MockMongoSession) WriteCollection(collectionName string, object types.Loc) error {
	switch {
	case mm.writeMode == "positive":
		return nil
	case mm.writeMode == "fail":
		return fmt.Errorf("Mock error on write")
	case mm.writeMode == "duplicate":
		err := mgo.QueryError {
			Code: 11000,
			Message: "Mock duplicate on write",
		}
		return &err
	}
	return fmt.Errorf("Unknown mode for WriteCollection: %s", mm.writeMode)
}

// UpdateCollection mock. Controlled by mm.writeMode values 'positive', 'fail' and 'missing'
func (mm *MockMongoSession) UpdateCollection(collectionName string, object types.Loc) error {
	switch {
	case mm.writeMode == "positive":
		return nil
	case mm.writeMode == "fail":
		return fmt.Errorf("Mock error on update")
	case mm.writeMode == "missing":
		err := mgo.QueryError {
			Code: 11000, // TODO: find the right error Code and type
			Message: "Mock not found on update",
		}
		return &err
	}
	return fmt.Errorf("Unknown mode for UpdateCollection: %s", mm.writeMode)
}

// FetchFromCollection mock. Controlled by mm.queryMode values 'positive' and 'fail'
func (mm *MockMongoSession) FetchFromCollection(collectionName string, id string) (types.Loc, error) {
	var result types.Loc
	switch {
	case mm.queryMode == "positive":
		result, err := types.LocFromString(id)
		if err != nil {
			return result, fmt.Errorf("Mock error creating loc")
		}
		return result, nil
	case mm.queryMode == "fail":
		return result, fmt.Errorf("Mock error on get")
	}
	return result, fmt.Errorf("Unknown mode for FetchFromCollection: %s", mm.queryMode)
}

// DeleteFromCollection mock. Controlled by mm.queryMode values 'positive' and 'fail'
func (mm *MockMongoSession) DeleteFromCollection(collectionName string, id string) (error) {
	switch {
	case mm.writeMode == "positive":
		_, err := types.LocFromString(id)
		if err != nil {
			return fmt.Errorf("Mock error creating loc")
		}
		return nil
	case mm.writeMode == "fail":
		return fmt.Errorf("Mock error on delete")
	case mm.writeMode == "missing":
		err := mgo.QueryError {
			Code: 11000, // TODO: find the right error Code and type
			Message: "Mock not found on delete",
		}
		return &err
	}
	return fmt.Errorf("Unknown mode for DeleteFromCollection: %s", mm.queryMode)
}
