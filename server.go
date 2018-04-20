package webstuff

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo"
)

func main() {
	_echo := echo.New()

	_echo.GET("/", getDefault)
	_echo.GET("/yo/:name", getYo)
	_echo.GET("/yo", func(c echo.Context) error { return c.HTML(http.StatusOK, "<p style='color:red'>Shy, eh?</p>") })
	_echo.GET("loc/:xyz", getLocXYZ)
	_echo.PUT("loc/:xyz", putLocXYZ)
	_echo.Logger.Fatal(_echo.Start(":3210"))
}

func getDefault(c echo.Context) error {
	return c.HTML(http.StatusOK, "<h2>This is the default page - now with format!</h2>")
}

func getYo(c echo.Context) error {
	return c.HTML(http.StatusOK, fmt.Sprintf("Yo <i>%s</i>! What up?", c.Param("name")))
}

func getLocXYZ(c echo.Context) error {
	loc := c.Param("xyz")
	x,y,z,err := LocConvert(loc)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, fmt.Sprintf("XYZ after split: %d, %d, %d", x, y, z))
}

func putLocXYZ(c echo.Context ) error {
	loc := c.Param("xyz")
	x,y,z,err := LocConvert(loc)
	if err != nil {
		return err
	}
	fmt.Printf("This is where I'd persist loc %d.%d.%d\n", x,y,z)
	return nil
}
