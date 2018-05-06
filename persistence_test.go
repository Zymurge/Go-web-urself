package webstuff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoURL       string = "localhost:27017"
	dbName         string = "testDB"
	testCollection string = "testCollection"
)

type MongoSessionSuite struct {
	suite.Suite
	session *mgo.Session
}
// Runner for the test suite. Ensures that mongo can be reached at the default location or aborts the suite. The suite provides a 
// pre-connected session for its tests to use for setting the DB state via the SetupTest() call.
func TestMongoSessionSuite(t *testing.T) {
	// precondition is that Mongo must be connectable at the default URL for the suite to run
	session, err := mgo.Dial(mongoURL)
	if session != nil {
		defer session.Close()
	}
	require.NoErrorf(t, err, "Mongo must be available at %s for this suite to function", mongoURL)
	suite.Run(t, new(MongoSessionSuite))
}

func (m *MongoSessionSuite) SetupTest() {
	m.session = GetMongoClearedCollection(m.T(), testCollection)
}

func (m *MongoSessionSuite) xTestDummy() {
	m.Fail("Dummy test")
}

// MongoClearCollection drops the specified collection. Depends on constants for mongoURL and DbName
func MongoClearCollection(collName string) error {
	session, err := mgo.Dial(mongoURL)
	defer session.Close()
	if err != nil {
		return err
	}
	myCollection := session.DB(dbName).C(collName)
	_, err = myCollection.RemoveAll(nil)
	return err
}

// GetMongoClearedCollection clears the specified collection and returns an active session pointing to it
func GetMongoClearedCollection(t *testing.T, collName string) (session *mgo.Session) {
	var err error
	session, err = mgo.Dial(mongoURL)
	if err != nil {
		t.Errorf("GetMongoClearedCollection failed to connect to Mongo")
	}
	myCollection := session.DB(dbName).C(collName)
	_, err = myCollection.RemoveAll(nil)
	if err != nil {
		t.Errorf("GetMongoClearedCollection failed to clear collection %s: %s", collName, err)
	}
	return session
}

func ClearMongoCollection(t *testing.T, session *mgo.Session, collName string) error {
	var err error
	clearMe := session.DB(dbName).C(collName)
	_, err = clearMe.RemoveAll(nil)
	if err != nil {
		t.Errorf("ClearMongoCollection failed to clear collection %s: %s", collName, err)
	}
	return err
}

func (m *MongoSessionSuite) TestConnectToMongo() {
	ms := MongoSession{
		mongoURL:       mongoURL,
		timeoutSeconds: 3 * time.Second,
	}
	err := ms.ConnectToMongo()
	m.NoError(err, "Sucessful connect throws no error. Instead we got %s", err)
	m.IsType(MongoSession{}, ms, "Wrong type on connect: %T", ms)
}

func (m *MongoSessionSuite) TestConnectToMongoNoConnectionThrowsError() {
	ms := MongoSession{
		mongoURL:       "i.am.abad.url:12345",
		timeoutSeconds: 100 * time.Millisecond,
	}
	err := ms.ConnectToMongo()
	m.Error(err, "Should return an error when the mongo server can't be found")
	m.Containsf(err.Error(), "no reachable", "Looking for err message saying it can't find the server. Instead got %s", err)
}

func (m *MongoSessionSuite) TestWriteCollection() {
	m.T().Run( "Positive", func(t *testing.T) {
		err := ClearMongoCollection(t, m.session, testCollection)
		require.NoError(t, err, "Test failed in setup clearing collection. Err: %s", err )
		testMS := NewMongoSession(mongoURL, dbName, 3)
		testLoc, _ := LocFromCoords( 1, 2, 3 )
		err = testMS.WriteCollection(testCollection, testLoc)
		require.NoError(t, err, "Successful write throws no error. Instead we got %s", err )
	} )
/*
	t.Run( "DuplicateInsertShouldError", func(t *testing.T) {
		mySession := GetMongoClearedCollection(t, testCollection )
		testLoc, _ := LocFromCoords( 1, 2, 3 )
		if err = WriteCollection( mySession ,testCollection, testLoc ); err != nil {
			t.Errorf( "Failure to populate with required test data: %s,", err )
		}
		// write the same loc again
		err = WriteCollection( mySession ,testCollection, testLoc )
		require.Error( t, err, "Attempt to insert duplicate ID should throw")
		require.Contains( t, err.Error(), "duplicate", "Expect error text to mention this" )
	} )
*/
}
/*
func TestUpdateCollection(t *testing.T) {
	t.Run( "Positive", func(t *testing.T ) {
		mySession := GetMongoClearedCollection(t, testCollection )
		var err error
		testLoc, _ := LocFromCoords( 1, 2, 3 )
		if err = WriteCollection( mySession ,testCollection, testLoc ); err != nil {
			t.Errorf( "Failure to populate with required test data: %s,", err )
		}
		testLoc.Status = "changed"
		err = UpdateCollection( mySession, testCollection, testLoc )
		require.NoError(t, err, "Successful update throws no error. Instead we got %s", err )
		// TODO: validate changed element in collection
	} )

	t.Run( "CollectionNotExist", func(t *testing.T) {
		mySession := GetMongoClearedCollection(t, testCollection )
		var err error
		testLoc, _ := LocFromCoords( 21, -21, 0 )
		err = UpdateCollection( mySession, testCollection, testLoc )
		require.Error(t, err, "Should get error message when attempt to update non-existent ID" )
		require.Contains( t, err.Error(), "not found", "Looking for the not found phrase, but got: %s", err )
	} )
}

func TestFetchFromCollection(t *testing.T) {
	mySession := GetMongoClearedCollection(t, testCollection )
	var err error
	testLoc, _ := LocFromCoords( 1, 2, 3 )
	if err = WriteCollection( mySession ,testCollection, testLoc ); err != nil {
		t.Errorf( "Failure to populate with required test data: %s,", err )
	}

	t.Run( "IDExists", func(t *testing.T) {
		var result Location
		result, err = FetchFromCollection( mySession, testCollection, testLoc.GetID() )
		require.NoError(t, err, "Successful lookup throws no error. Instead we got %s", err )
		require.NotNil(t, result, "Successful lookup has to actually return something")
		require.Equal(t, testLoc.ID, result.GetID() )
		rloc, ok := result.(Loc)
		require.True(t, ok, "Returned object type not of concrete type Loc")
		require.Equal(t, testLoc.X, rloc.X)
	})

	t.Run( "IDNotExists", func(t *testing.T) {
		unexpectedID := "11.12.-13"
		_, err = FetchFromCollection( mySession, testCollection, unexpectedID )
		require.Error(t, err, "Missing id should throw an error")
		require.Contains(t, err.Error(), "not found", "Message should give a clue. Instead it is %s", err)
	})
}
*/
