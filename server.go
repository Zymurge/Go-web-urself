package main

import (
	"fmt"
	"net/http"
	per "webstuff/persistence"

	"github.com/labstack/echo"
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
	e.GET("/yo/:name", h.getYo)
	e.GET("/yo", func(c echo.Context) error { return c.HTML(http.StatusOK, "<p style='color:red'>Shy, eh?</p>") })
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

func (h Handler) getYo(c echo.Context) error {
	err := h.mongoDB.ConnectToMongo()
	if err != nil {
		return fmt.Errorf("Got me an error: %s", err)
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("Yo <i>%s</i>! What up?", c.Param("name")))
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
	return c.String(http.StatusOK, fmt.Sprintf("XYZ is: %v", loc.GetID()))
}

func (h Handler) putLocXYZ(c echo.Context) error {
	loc := c.Param("xyz")
	/* 	x,y,z,err := per.LocConvert(loc)
	   	if err != nil {
	   		return err
	   	}
	   	fmt.Printf("This is where I'd persist loc %d.%d.%d\n", x,y,z) */
	return c.String(http.StatusMethodNotAllowed, fmt.Sprintf("putLocXYZ: %s", loc))
}
