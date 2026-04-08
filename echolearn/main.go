package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	http.ListenAndServe(":1323", nil)
	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c echo.Context) error {
		var test string
		c.Param("some id ")
		c.QueryParam("query param")
		c.Bind(test)
		return c.JSON(200, "some-struct")
		// return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
