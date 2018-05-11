package main

import (
	"net/http"
	"github.com/labstack/echo"
	"fmt"
//  "strings"
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/require"
	per "webstuff/persistence"
)

func TestWebServesSomething(t* testing.T) {
	t.Skip()
}

func TestPutLocXYZ(t *testing.T) {
	mock := &MockMongoSession{}
	handler, err := NewHandler(mock)
	require.NoErrorf(t, err, "Issue with handler construction: %s", err)
	require.IsTypef(t, Handler{}, handler, "Not sure what we got going here" )
	// fire some web traffic here and validate with the mock Fetch method
}

func TestGetLocXYZ(t *testing.T) {
	expectedID := "1.9.-10"
	expectedBody := fmt.Sprintf("XYZ is: %v", expectedID)
	_, handler := NewHandlerWithMockMongo(t)
	
	req := httptest.NewRequest(echo.GET, "/loc/" + expectedID, nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.SetParamNames("xyz")
	ctx.SetParamValues(expectedID)

	err := handler.getLocXYZ(ctx)
	require.NoErrorf(t, err, "Didn't want an error on positive test. Got: %s", err)
	require.Equalf(t, http.StatusOK, rec.Code, "HTTP response should be success")
	require.Equal(t, expectedBody, rec.Body.String())
}

func TestGetLocXYZNotFound(t *testing.T) {
	expectedID := "1.1.1"
	expectedBody := fmt.Sprintf("%s doesn't exist in DB", expectedID)
	mock, handler := NewHandlerWithMockMongo(t)
	mock.queryMode = "negative"

	req := httptest.NewRequest(echo.GET, "/loc/" + expectedID, nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.SetParamNames("xyz")
	ctx.SetParamValues(expectedID)

	err := handler.getLocXYZ(ctx)
	require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
	require.Equalf(t, http.StatusNotFound, rec.Code, "HTTP response should be not found")
	require.Equal(t, expectedBody, rec.Body.String())
}

func TestGetLocXYZNoMongo(t *testing.T) {
	expectedID := "1.1.1"
	expectedBody := fmt.Sprintf("MongoDB not available")
	mock, handler := NewHandlerWithMockMongo(t)
	mock.connectMode = "no connect"

	req := httptest.NewRequest(echo.GET, "/loc/" + expectedID, nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.SetParamNames("xyz")
	ctx.SetParamValues(expectedID)

	err := handler.getLocXYZ(ctx)
	require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
	require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
	require.Equal(t, expectedBody, rec.Body.String())
}

func NewHandlerWithMockMongo(t *testing.T) (*MockMongoSession, *Handler) {
	mock := &MockMongoSession{
		connectMode: "positive",
		queryMode: "positive",
	}
	handler, err := NewHandler(mock)
	require.NoErrorf(t, err, "Issue with handler construction: %s", err)
	require.NotNil(t, handler)
	return mock, &handler
}

type MockMongoSession struct {
	connectMode string
	queryMode string
}

func (mm *MockMongoSession) ConnectToMongo() error {
	switch {
	case mm.connectMode == "positive":
		return nil
	case mm.connectMode == "no connect":
		return fmt.Errorf( "mocked connection failure")
	}
	return fmt.Errorf("Unknown mode for ConnectToMongo: %s", mm.connectMode)
}

func (mm *MockMongoSession) WriteCollection( collectionName string, object per.Loc ) error {
	return fmt.Errorf("Mock function not implemented")
}
func (mm *MockMongoSession) UpdateCollection( collectionName string, object per.Loc ) error {
	return fmt.Errorf("Mock function not implemented")
}
func (mm *MockMongoSession) FetchFromCollection( collectionName string, id string ) (per.Loc, error) {
	var result per.Loc
	switch {
	case mm.queryMode == "positive":
		result, err := per.LocFromString(id)
		if err != nil { 
			return result, fmt.Errorf("Mock error creating loc")
		}
		return result, nil
	case mm.queryMode == "fail":
		return result, fmt.Errorf("Mock error by request")
	}
	return result, fmt.Errorf("Mock function not implemented")
}