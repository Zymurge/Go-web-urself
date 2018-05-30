package main

import (
	"strings"
	"fmt"
	"net/http"
	per "webstuff/persistence"
	"webstuff/types"

	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoURL      string = "localhost:27017"
	dbName        string = "testDB"
	locCollection string = "testCollection"
)

func main() {
	mdb := per.NewMongoSession(mongoURL, dbName, 10)
	h, err := NewHandler(mdb)
	if err != nil {
		panic("Couldn't establish a Handler for some reason")
	}
	e := echo.New()
	// set routes here I guess
	e.GET("/", h.getDefault)
	e.GET("loc/:xyz", h.getLocXYZ)
	e.POST("loc/:xyz", h.postLocXYZ)

	defer e.Logger.Fatal(e.Start(":3210"))
}

// Handler encapsulates web handling with persistence
type Handler struct {
	mongoDB per.MongoAbstraction
}

// NewHandler returns a route handler instance with the injected mongo layer
func NewHandler(mdb per.MongoAbstraction) (result Handler, err error) {
	h := Handler{
		mongoDB: mdb,
	}
	return h, nil
}

func (h Handler) getDefault(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h2>This is the default page - now with format!</h2>")
}

func (h Handler) getLocXYZ(c echo.Context) (err error) {
	locID := c.Param("xyz")
	var loc types.Loc
	if err = h.mongoDB.ConnectToMongo(); err != nil {
		// TODO: do something with the err info from mongo. Log it?
		err =  c.HTML(http.StatusFailedDependency, "MongoDB not available")
		return
	}
	if loc, err = h.mongoDB.FetchFromCollection(locCollection, locID); err != nil {
		// TODO: do something with the err info from mongo. Log it?
		err = c.HTML(http.StatusNotFound, fmt.Sprintf("%s doesn't exist in DB", locID))
		return
	}
	err = c.HTML(http.StatusOK, string(loc.JSONForm()))
	return
}

func (h Handler) deleteLocXYZ(c echo.Context) (err error) {
	locID := c.Param("xyz")
	if err = h.mongoDB.ConnectToMongo(); err != nil {
		// TODO: do something with the err info from mongo. Log it?
		err = c.HTML(http.StatusFailedDependency, "MongoDB not available")
		return
	}
	if err = h.mongoDB.DeleteFromCollection(locCollection, locID); err != nil {
		// TODO: do something with the err info from mongo. Log it?
		// TODO: branch on not found (404) vs other mongo error (424). 
		if strings.Contains(err.Error(), "not found") {
			err = c.HTML(http.StatusNotFound, fmt.Sprintf("%s doesn't exist in DB", locID))
		} else {
			err = c.HTML(http.StatusFailedDependency, fmt.Sprintf("Unknown error on Mongo delete: %v", err))
		}
		return
	}
	err = c.HTML(http.StatusOK, fmt.Sprintf("%s deleted from DB", locID))
	return
}

func (h Handler) postLocXYZ(c echo.Context) (err error) {
	locString := c.Param("xyz")
	var loc types.Loc
	if loc, err = types.LocFromString(locString); err != nil {
		// TODO: do something with the err info from Loc ctor. Log it?
		err = c.HTML(http.StatusBadRequest, "Bad string for param xyz")
		return
	}
	if err = h.mongoDB.ConnectToMongo(); err != nil {
		// TODO: do something with the err info from mongo. Log it?
		err = c.HTML(http.StatusFailedDependency, "MongoDB not available")
		return
	}
	if err = h.mongoDB.WriteCollection(locCollection, loc); err != nil {
		// TODO: do something with the err info from mongo. Log it?
		if mgo.IsDup(err) {
			err = c.HTML(http.StatusAlreadyReported, fmt.Sprintf("Duplicate insert for xyz: %s", loc.GetID()))
			return
		}
		err = c.HTML(http.StatusFailedDependency, fmt.Sprintf("Unknown error on Mongo insert: %v", err))
		return
	}
	err = c.HTML(http.StatusOK, fmt.Sprintf("Inserted: %s", loc.GetID()))
	return
}
