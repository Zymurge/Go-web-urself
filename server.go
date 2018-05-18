package main

import (
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
	e.PUT("loc/:xyz", h.putLocXYZ)

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

func (h Handler) getLocXYZ(c echo.Context) error {
	locID := c.Param("xyz")
	if err := h.mongoDB.ConnectToMongo(); err != nil {
		return c.HTML(http.StatusFailedDependency, "MongoDB not available")
	}
	loc, err := h.mongoDB.FetchFromCollection(locCollection, locID)
	if err != nil {
		return c.HTML(http.StatusNotFound, fmt.Sprintf("%s doesn't exist in DB", locID))
	}
	return c.HTML(http.StatusOK, string(loc.JSONForm()))
}

func (h Handler) putLocXYZ(c echo.Context) error {
	locString := c.Param("xyz")
	loc, err := types.LocFromString(locString)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "Bad string for param xyz")
	}
	if err := h.mongoDB.ConnectToMongo(); err != nil {
		return c.HTML(http.StatusFailedDependency, "MongoDB not available")
	}
	err = h.mongoDB.WriteCollection(locCollection, loc)
	if err != nil {
		if mgo.IsDup(err) {
			return c.HTML(http.StatusAlreadyReported, fmt.Sprintf("Duplicate insert for xyz: %s", loc.GetID()))
		}
		return c.HTML(http.StatusFailedDependency, fmt.Sprintf("Unknown error on Mongo insert: %v", err))
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("Inserted: %s", loc.GetID()))
}
