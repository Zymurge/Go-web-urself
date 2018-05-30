package persistence

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	mgo "gopkg.in/mgo.v2"
	"webstuff/types"
)

const (
	testMongoURL       string = "localhost:27017"
	testDbName         string = "testDB"
	testCollection 	   string = "testCollection"
)

type MongoSessionSuite struct {
	suite.Suite
	session *mgo.Session
}
// Runner for the test suite. Ensures that mongo can be reached at the default location or aborts the suite. The suite provides a 
// pre-connected session for its tests to use for setting the DB state via the SetupTest() call.
func TestMongoSessionSuite(t *testing.T) {
	// precondition is that Mongo must be connectable at the default URL for the suite to run
	session, err := mgo.Dial(testMongoURL)
	if session != nil {
		defer session.Close()
	}
	require.NoErrorf(t, err, "Mongo must be available at %s for this suite to function", testMongoURL)
	suite.Run(t, new(MongoSessionSuite))
}

func (m *MongoSessionSuite) SetupSuite() {
	m.session = GetMongoClearedCollection(m.T(), testCollection)
}

func (m *MongoSessionSuite) SetupTest() {
	err := ClearMongoCollection(m.T(), m.session, testCollection)
	m.NoError(err, "Test failed in setup clearing collection. Err: %s", err )
}

func (m *MongoSessionSuite) TestCtorDefaults() {
	result := NewMongoSession("testURL","",5)
	m.EqualValues(DefaultDbName, result.dbName, "DB name should be the default")
}

func (m *MongoSessionSuite) TestConnectToMongo() {
	ms := MongoSession{
		mongoURL:       testMongoURL,
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
	var err error
	m.T().Run( "Positive", func(t *testing.T) {
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		testLoc, _ := types.LocFromCoords(1, 2, 3)
		err = testMS.WriteCollection(testCollection, testLoc)
		require.NoError(t, err, "Successful write throws no error. Instead we got %s", err )
	} )
	m.T().Run( "DuplicateInsertShouldError", func(t *testing.T) {
		testLoc, _ := types.LocFromCoords(-1, -2, -3)
		err = AddToMongoCollection(t, m.session, testCollection, testLoc )
		require.NoError(t, err, "Test failed in setup adding to collection. Err: %s", err )

		// write the same loc again
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.WriteCollection( testCollection, testLoc )
		require.Error( t, err, "Attempt to insert duplicate ID should throw")
		require.Contains( t, err.Error(), "duplicate", "Expect error text to mention this" )
	} )
	m.T().Run( "CollectionNotExistShouldStillWrite", func(t *testing.T) {
		testBadCollection := "garbage"
		ClearMongoCollection(t, m.session, testBadCollection)

		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.WriteCollection(testBadCollection, types.Loc{})
		require.NoErrorf(t, err, "Writes should create collection on the fly. Got err: %s", err)
		writeCount,_ := m.session.DB(testDbName).C(testBadCollection).Count()
		require.True(t, writeCount == 1, "Record should have been written as only entry")
	} )
	m.T().Run("Dropped connection", func(t *testing.T){
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		testMS.mongoURL = "yo"
		testLoc, _ := types.LocFromCoords(22, 22, 33)
		err = testMS.WriteCollection(testCollection, testLoc)
		require.Error(t, err, "Should get an error if changed to unreachable URL")
		require.Contains(t, err.Error(), "no reachable servers", "Should complain about lack of connectivity")
	} )
}

func (m *MongoSessionSuite) TestDeleteFromCollection() {
	var err error
	m.T().Run("Positive", func(t *testing.T) {
		testID := "1.2.3"
		testLoc, _ := types.LocFromString(testID)
		err = AddToMongoCollection(t, m.session, testCollection, testLoc )
		require.NoError(t, err, "Test failed in setup adding to collection. Err: %s", err )

		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.DeleteFromCollection(testCollection, testID)
		require.NoError(t, err, "Successful deletions throw no errors. But this threw: %s", err )
	} )
	m.T().Run("Missing ID", func(t *testing.T) {
		testID := "1.2.3"

		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.DeleteFromCollection(testCollection, testID)
		require.Error(t, err, "Delete on missing ID should throw error")
		require.Containsf(t, err.Error(), "not found", "mgo should specify why it threw on missing ID")
	} )
	m.T().Run( "CollectionNotExist", func(t *testing.T) {
		testBadCollection := "garbage"
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.DeleteFromCollection(testBadCollection, "matters not")
		require.Error(t, err, "Should get error message when attempt to access non-existent collection")
		require.Contains(t, err.Error(), "not found", "Looking for the not found phrase, but got: %s", err)
	} )
	m.T().Run("Dropped connection", func(t *testing.T){
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		testMS.mongoURL = "yo"
		err = testMS.DeleteFromCollection(testCollection, "matters not")
		require.Error(t, err, "Should get an error if changed to unreachable URL")
		require.Contains(t, err.Error(), "no reachable servers", "Should complain about lack of connectivity")
	} )
}

func (m *MongoSessionSuite) TestUpdateCollection() {
	var err error
	m.T().Run( "Positive", func(t *testing.T) {
		testLoc, _ := types.LocFromCoords(11, 2, 13)
		err = AddToMongoCollection(t, m.session, testCollection, testLoc)
		require.NoError(t, err, "Test failed in setup adding to collection. Err: %s", err)

		testLoc.Status = "changed"
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.UpdateCollection(testCollection, testLoc)
		require.NoError(t, err, "Successful update throws no error. Instead we got %s", err)
		// TODO: validate changed element in collection
	} )
	m.T().Run( "MissingID", func(t *testing.T) {
		testLoc, _ := types.LocFromCoords(1, 12, 3)

		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.UpdateCollection(testCollection, testLoc)
		require.Error(t, err, "Missing ID should error on update")
		require.Contains(t, err.Error(), "not found", "Looking for message about ID missing, but got: %s", err)
	} )
	m.T().Run( "CollectionNotExist", func(t *testing.T) {
		testBadCollection := "garbage"
		err = m.session.DB(testDbName).C(testBadCollection).DropCollection()
		require.NoError(t, err, "Test failed in setup dropping test collection. Err: %s", err)

		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		err = testMS.UpdateCollection(testBadCollection, types.Loc{})
		require.Error(t, err, "Should get error message when attempt to access non-existent collection")
		require.Contains(t, err.Error(), "Non-existent collection for update", "Looking for missing collection, but got: %s", err)
	} )
	m.T().Run("Dropped connection", func(t *testing.T){
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		testMS.mongoURL = "yo"
		testLoc, _ := types.LocFromCoords(22, 22, 33)
		err = testMS.UpdateCollection(testCollection, testLoc)
		require.Error(t, err, "Should get an error if changed to unreachable URL")
		require.Contains(t, err.Error(), "no reachable servers", "Should complain about lack of connectivity")
	} )
}

func (m *MongoSessionSuite) TestFetchFromCollection() {
	var err error
	testLoc, _ := types.LocFromCoords(1, 2, 3)
	err = AddToMongoCollection(m.T(), m.session, testCollection, testLoc)
	m.NoError(err, "Test failed in setup adding to collection. Err: %s", err)
	testMS := NewMongoSession(testMongoURL, testDbName, 3)

	m.T().Run("Positive", func(t *testing.T) {
		var result types.Loc
		result, err = testMS.FetchFromCollection(testCollection, testLoc.GetID())
		require.NoError(t, err, "Successful lookup throws no error. Instead we got %s", err )
		require.NotNil(t, result, "Successful lookup has to actually return something")
		require.Equal(t, testLoc.ID, result.GetID() )
		require.Equal(t, testLoc.X, result.X)
	} )
	m.T().Run("Missing ID", func(t *testing.T) {
		unexpectedID := "11.12.-13"
		_, err = testMS.FetchFromCollection(testCollection, unexpectedID)
		require.Error(t, err, "Missing id should throw an error")
		require.Contains(t, err.Error(), "not found", "Message should give a clue. Instead it is %s", err)
	} )
	m.T().Run("Dropped connection", func(t *testing.T){
		testMS := NewMongoSession(testMongoURL, testDbName, 3)
		testMS.mongoURL = "yo"
		testLoc, _ := types.LocFromCoords(22, 22, 33)
		_, err = testMS.FetchFromCollection(testCollection, testLoc.GetID())
		require.Error(t, err, "Should get an error if changed to unreachable URL")
		require.Contains(t, err.Error(), "no reachable servers", "Should complain about lack of connectivity")
	} )
}

/*** Helper functions ***/


// MongoClearCollection drops the specified collection. Depends on constants for testMongoURL and DbName.
// Hardcodes 2 second timeout on connect, since it expects local mongo to work
func MongoClearCollection(collName string) error {
	to := 2 * time.Second
	session, err := mgo.DialWithTimeout(testMongoURL, to)
	defer session.Close()
	if err != nil {
		return err
	}
	myCollection := session.DB(testDbName).C(collName)
	_, err = myCollection.RemoveAll(nil)
	return err
}

// GetMongoClearedCollection clears the specified collection and returns an active session pointing to it
func GetMongoClearedCollection(t *testing.T, collName string) (session *mgo.Session) {
	var err error
	session, err = mgo.Dial(testMongoURL)
	if err != nil {
		t.Errorf("GetMongoClearedCollection failed to connect to Mongo")
	}
	myCollection := session.DB(testDbName).C(collName)
	_, err = myCollection.RemoveAll(nil)
	if err != nil {
		t.Errorf("GetMongoClearedCollection failed to clear collection %s: %s", collName, err)
	}
	return session
}

func ClearMongoCollection(t *testing.T, session *mgo.Session, collName string) error {
	var err error
	clearMe := session.DB(testDbName).C(collName)
	_, err = clearMe.RemoveAll(nil)
	if err != nil {
		t.Errorf("ClearMongoCollection failed to clear collection %s: %s", collName, err)
	}
	return err
}

func AddToMongoCollection(t *testing.T, session *mgo.Session, collName string, obj interface{}) error {
	myCollection := session.DB(testDbName).C(collName)
	return myCollection.Insert(obj)
}
