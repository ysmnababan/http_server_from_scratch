package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	// http.ListenAndServe(":8080", nil)
	e := echo.New()
	e.Use(middleware.RequestLogger())

	// e.GET("/oke/mantap/wew", func(c echo.Context) error {
	// 	var test string
	// 	c.Param("some id ")
	// 	c.QueryParam("query param")
	// 	c.Bind(test)
	// 	fmt.Println("here")
	// 	// _ = c.JSON(http.StatusOK, "Hello, World!")
	// 	// return c.JSON(400, errors.New("not error").Error())
	// 	return c.String(http.StatusOK, "panjang")
	// 	// return echo.NewHTTPError(400, "not found")
	// })
	e.GET("/oke", func(c echo.Context) error {
		var test string
		c.Param("some id ")
		c.QueryParam("query param")
		c.Bind(test)
		fmt.Println("here")
		// _ = c.JSON(http.StatusOK, "Hello, World!")
		// return c.JSON(400, errors.New("not error").Error())
		return c.String(http.StatusOK, "pendek")
		// return echo.NewHTTPError(400, "not found")
	})

	if err := e.Start(":8081"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
