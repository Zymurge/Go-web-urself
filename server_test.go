package main

import (
	"net/http"
	"github.com/labstack/echo"
	"fmt"
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/require"
	"webstuff/types"
)

func TestWebServesSomething(t* testing.T) {
	t.Skip()
}

func TestGetLocXYZ(t *testing.T) {
	// Base setup is handler w/mock context and param for xyz. Each case must:
	//   - set param value
	//   - set mock mode flags
	expectedID := "5.6.7"
	expectedBody := "set me"
	mock, handler := NewHandlerWithMockMongo(t)
	req := httptest.NewRequest(echo.GET, "/loc/" + expectedID, nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.SetParamNames("xyz")

	t.Run("Positive", func(t *testing.T){
		expectedID = "5.6.7"
		expectedLoc, _ := types.LocFromString(expectedID)
		expectedBody = string(expectedLoc.JSONForm())
		ctx.SetParamValues(expectedID)

		err := handler.getLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on positive test. Got: %s", err)
		require.Equalf(t, http.StatusOK, rec.Code, "HTTP response should be success")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Missing ID", func(t *testing.T){
		expectedID = "15.16.17"
		expectedBody = fmt.Sprintf("%s doesn't exist in DB", expectedID)
		mock.queryMode = "fail"
		req = httptest.NewRequest(echo.GET, "/loc/" + expectedID, nil)
		rec = httptest.NewRecorder()
		ctx = echo.New().NewContext(req, rec)
		ctx.SetParamNames("xyz")
		ctx.SetParamValues(expectedID)

		err := handler.getLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusNotFound, rec.Code, "HTTP response should be not found")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("No Mongo", func(t *testing.T){
		expectedBody := fmt.Sprintf("MongoDB not available")
		mock.connectMode = "no connect"
		req = httptest.NewRequest(echo.GET, "/loc/" + expectedID, nil)
		rec = httptest.NewRecorder()
		ctx = echo.New().NewContext(req, rec)
		ctx.SetParamNames("xyz")
		ctx.SetParamValues(expectedID)

		err := handler.getLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
}

func TestPutLocXYZ(t *testing.T) {
	// Base setup is handler w/mock context and param for xyz. Each case must:
	//   - set param value
	//   - set mock mode flags
	expectedID := "5.6.7"
	//putLoc, _ := types.LocFromString(expectedID)
	expectedBody := "set me"
	mock, handler := NewHandlerWithMockMongo(t)
	req := httptest.NewRequest(echo.PUT, "/loc/" + expectedID, nil)
	rec := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, rec)
	ctx.SetParamNames("xyz")

	t.Run("Positive", func(t *testing.T){
		mock.connectMode = "positive"
		mock.writeMode = "positive"
		ctx.SetParamValues(expectedID)

		err := handler.putLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on positive test. Got: %s", err)
		require.Equalf(t, http.StatusOK, rec.Code, "HTTP response should be success")
	})
	t.Run("Duplicate ID", func(t *testing.T){
		expectedBody = fmt.Sprintf("Duplicate insert for xyz: %s", expectedID)
		mock.connectMode = "positive"
		mock.writeMode = "duplicate"
		req = httptest.NewRequest(echo.PUT, "/loc/" + expectedID, nil)
		rec = httptest.NewRecorder()
		ctx = echo.New().NewContext(req, rec)
		ctx.SetParamNames("xyz")
		ctx.SetParamValues(expectedID)

		err := handler.putLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusAlreadyReported, rec.Code, "HTTP response should be already reported to represent duplicate insert")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Bad Loc string", func(t *testing.T){
		expectedBody = fmt.Sprintf("Bad string for param xyz")
		badID := "a.7.tty"
		mock.connectMode = "positive"
		mock.writeMode = "positive"
		req = httptest.NewRequest(echo.PUT, "/loc/" + expectedID, nil)
		rec = httptest.NewRecorder()
		ctx = echo.New().NewContext(req, rec)
		ctx.SetParamNames("xyz")
		ctx.SetParamValues(badID)

		err := handler.putLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusBadRequest, rec.Code, "HTTP response should be bad request for malformed Loc string")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Other Mongo error", func(t *testing.T){
		mockErrorMsg := "Mock error on write"
		expectedBody = fmt.Sprintf("Unknown error on Mongo insert: %s", mockErrorMsg)
		mock.connectMode = "positive"
		mock.writeMode = "fail"
		req = httptest.NewRequest(echo.PUT, "/loc/" + expectedID, nil)
		rec = httptest.NewRecorder()
		ctx = echo.New().NewContext(req, rec)
		ctx.SetParamNames("xyz")
		ctx.SetParamValues(expectedID)

		err := handler.putLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("No Mongo", func(t *testing.T){
		expectedBody := fmt.Sprintf("MongoDB not available")
		mock.connectMode = "no connect"
		req = httptest.NewRequest(echo.PUT, "/loc/" + expectedID, nil)
		rec = httptest.NewRecorder()
		ctx = echo.New().NewContext(req, rec)
		ctx.SetParamNames("xyz")
		ctx.SetParamValues(expectedID)

		err := handler.putLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
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

