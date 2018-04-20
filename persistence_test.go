package webstuff

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoURL string = "localhost:27017"
	testCollection string = "testCollection"
)



// MongoClearCollection drops the specified collection. Depends on constants for mongoURL and DbName
func MongoClearCollection( collName string ) error {
	session, err := ConnectToMongo( mongoURL )
	if err != nil {
		return err
	}
	myCollection := session.DB(DbName).C(collName)
	_, err = myCollection.RemoveAll(nil)
	return err
}

// GetMongoClearedCollection clears the specified collection and returns an active session pointing to it
func GetMongoClearedCollection( t *testing.T, collName string ) (mySession *mgo.Session) {
	var err error
	if err = MongoClearCollection( collName ); err != nil {
		t.Errorf( "Failure to clear collection pre-test: %s", err )
	}
	if mySession, err = ConnectToMongo( mongoURL ); err != nil {
		t.Errorf( "Session connect threw: %s", err )
	}
	return mySession
}

func TestConnectToMongo(t *testing.T) {
	actual, err := ConnectToMongo( mongoURL )
	require.NoError( t, err, "Sucessful connect throws no error. Instead we got %s", err )
	require.IsType( t, &mgo.Session{}, actual, "Wrong type on connect: %T", actual )
}

func TestWriteCollection(t *testing.T) {
	t.Run( "Positive", func(t *testing.T) {
		mySession := GetMongoClearedCollection(t, testCollection )
		var err error
		testLoc, _ := LocFromCoords( 1, 2, 3 )
		err = WriteCollection( mySession ,testCollection, testLoc )
		assert.True( t, err == nil, "err is: %s", err )
		require.NoError(t, err, "Successful write throws no error. Instead we got %s", err )
	} )

	t.Run( "DuplicateInsertShouldError", func(t *testing.T) {
		mySession := GetMongoClearedCollection(t, testCollection )
		var err error
		testLoc, _ := LocFromCoords( 1, 2, 3 )
		if err = WriteCollection( mySession ,testCollection, testLoc ); err != nil {
			t.Errorf( "Failure to populate with required test data: %s,", err )
		}
		// write the same loc again
		err = WriteCollection( mySession ,testCollection, testLoc )
		require.Error( t, err, "Attempt to insert duplicate ID should throw")
		require.Contains( t, err.Error(), "duplicate", "Expect error text to mention this" )
	} )
}

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