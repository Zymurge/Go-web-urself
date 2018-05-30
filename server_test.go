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
	//   - get context with target and param value
	//   - set mock mode flags
	expectedID := "5.6.7"
	expectedBody := "set me"
	mock, handler := NewHandlerWithMockMongo(t)

	t.Run("Positive", func(t *testing.T){
		expectedID = "5.6.7"
		expectedLoc, _ := types.LocFromString(expectedID)
		expectedBody = string(expectedLoc.JSONForm())
		ctx, rec := GetNewEchoContext(echo.GET, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.getLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on positive test. Got: %s", err)
		require.Equalf(t, http.StatusOK, rec.Code, "HTTP response should be success")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Missing ID", func(t *testing.T){
		expectedID = "15.16.17"
		expectedBody = fmt.Sprintf("%s doesn't exist in DB", expectedID)
		mock.queryMode = "fail"
		ctx, rec := GetNewEchoContext(echo.GET, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.getLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusNotFound, rec.Code, "HTTP response should be not found")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("No Mongo", func(t *testing.T){
		expectedBody := fmt.Sprintf("MongoDB not available")
		mock.connectMode = "no connect"
		ctx, rec := GetNewEchoContext(echo.GET, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.getLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on no mongo test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
}

func TestPostLocXYZ(t *testing.T) {
	// Base setup is handler w/mock context and param for xyz. Each case must:
	//   - get context with target and param value
	//   - set mock mode flags
	expectedID := "5.6.7"
	expectedBody := "set me"
	mock, handler := NewHandlerWithMockMongo(t)

	t.Run("Positive", func(t *testing.T){
		mock.connectMode = "positive"
		mock.writeMode = "positive"
		ctx, rec := GetNewEchoContext(echo.POST, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.postLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on positive test. Got: %s", err)
		require.Equalf(t, http.StatusOK, rec.Code, "HTTP response should be success")
	})
	t.Run("Duplicate ID", func(t *testing.T){
		expectedBody = fmt.Sprintf("Duplicate insert for xyz: %s", expectedID)
		mock.connectMode = "positive"
		mock.writeMode = "duplicate"
		ctx, rec := GetNewEchoContext(echo.POST, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.postLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusAlreadyReported, rec.Code, "HTTP response should be already reported to represent duplicate insert")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Bad Loc string", func(t *testing.T){
		expectedBody = fmt.Sprintf("Bad string for param xyz")
		badID := "a.7.tty"
		mock.connectMode = "positive"
		mock.writeMode = "positive"
		ctx, rec := GetNewEchoContext(echo.POST, "/loc/" + badID, "xyz", badID )

		err := handler.postLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on bad loc test. Got: %s", err)
		require.Equalf(t, http.StatusBadRequest, rec.Code, "HTTP response should be bad request for malformed Loc string")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Other Mongo error", func(t *testing.T){
		mockErrorMsg := "Mock error on write"
		expectedBody = fmt.Sprintf("Unknown error on Mongo insert: %s", mockErrorMsg)
		mock.connectMode = "positive"
		mock.writeMode = "fail"
		ctx, rec := GetNewEchoContext(echo.POST, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.postLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on other mongo error test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("No Mongo", func(t *testing.T){
		expectedBody := fmt.Sprintf("MongoDB not available")
		mock.connectMode = "no connect"
		ctx, rec := GetNewEchoContext(echo.POST, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.postLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on no mongo test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
}

func TestDeleteLocXYZ(t *testing.T) {
	// Base setup is handler w/mock context and param for xyz. Each case must:
	//   - get context with target and param value
	//   - set mock mode flags
	expectedID := "5.6.7"
	expectedBody := "set me"
	mock, handler := NewHandlerWithMockMongo(t)

	t.Run("Positive", func(t *testing.T){
		expectedID := "5.6.7"
		expectedBody = fmt.Sprintf("%s deleted from DB", expectedID)
		mock.connectMode = "positive"
		mock.writeMode = "positive"
		ctx, rec := GetNewEchoContext(echo.DELETE, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.deleteLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on positive test. Got: %s", err)
		require.Equalf(t, http.StatusOK, rec.Code, "HTTP response should be success")
		require.Equalf(t, expectedBody, rec.Body.String(), "Wanted the loc confirmation on delete. Got %s", rec.Body)
	})
	t.Run("Missing ID", func(t *testing.T){
		expectedID = "15.16.17"
		expectedBody = fmt.Sprintf("%s doesn't exist in DB", expectedID)
		mock.writeMode = "missing"
		ctx, rec := GetNewEchoContext(echo.DELETE, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.deleteLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on not found test. Got: %s", err)
		require.Equalf(t, http.StatusNotFound, rec.Code, "HTTP response should be not found")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("No Mongo", func(t *testing.T){
		expectedBody = "MongoDB not available"
		mock.connectMode = "no connect"
		ctx, rec := GetNewEchoContext(echo.DELETE, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.deleteLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on no mongo test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})
	t.Run("Other Mongo error", func(t *testing.T){
		mockErrorMsg := "Mock error on delete"
		expectedBody = fmt.Sprintf("Unknown error on Mongo delete: %s", mockErrorMsg)
		mock.connectMode = "positive"
		mock.writeMode = "fail"
		ctx, rec := GetNewEchoContext(echo.DELETE, "/loc/" + expectedID, "xyz", expectedID )

		err := handler.deleteLocXYZ(ctx)
		require.NoErrorf(t, err, "Didn't want an error on other mongo error test. Got: %s", err)
		require.Equalf(t, http.StatusFailedDependency, rec.Code, "HTTP response should be failed dependency")
		require.Equal(t, expectedBody, rec.Body.String())
	})

}


/*** Helper functions ***/

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

// GetNewEchoContext is a helper method to aggregate common things into a single EchoContext for use in testing
// web requests.
// The method param must be one of the known echo constants (ie - GET, PUT, UPDATE, etc) 
// It currently only supports a single param and value.
func GetNewEchoContext(method string, target string, pname string, pvalue string) (ctx echo.Context, rec *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, target, nil)
	rec = httptest.NewRecorder()
	ctx = echo.New().NewContext(req, rec)
	ctx.SetParamNames(pname)
	ctx.SetParamValues(pvalue)
	return
}

