package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func mw1(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("before mw1")
		err := next(c)
		// fmt.Println("after mw1")
		return err
	}
}

func mw2(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("before mw2")
		err := next(c)
		fmt.Println("after mw2")

		return err
	}
}

func fA() {
	defer func() {
		fmt.Println("A: before recover")
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
		fmt.Println("after recover")
	}()
	fmt.Println("A")
	fB()
}

func fB() {
	fmt.Println("B")
	fC()
	fmt.Println("b after")
}

func fC() {
	fmt.Println("C")
	panic("panic in C")
	fmt.Println("after panic")
}

func notmain() {
	fA()
}

func main() {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	// http.ListenAndServe(":8080", nil)
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Logger())
	e.Use(mw1, mw2)
	group := e.Group("groupname")
	group.GET("insidegroup", func(c echo.Context) error {
		return nil
	})
	e.GET("a/:id/c", func(c echo.Context) error {
		// var test string
		c.Param("some id ")
		c.QueryParam("query param")
		// c.Bind(test)
		fmt.Println("with param")
		// _ = c.JSON(http.StatusOK, "Hello, World!")
		// return c.JSON(400, errors.New("not error").Error())
		return c.String(http.StatusOK, "with param"+c.Param("id"))
		// return echo.NewHTTPError(400, "not found")
	})
	e.GET("a/:data/c", func(c echo.Context) error {
		// var test string
		c.Param("some id ")
		c.QueryParam("query param")
		// c.Bind(test)
		fmt.Println("with param")
		// _ = c.JSON(http.StatusOK, "Hello, World!")
		// return c.JSON(400, errors.New("not error").Error())
		return c.String(http.StatusOK, "with data param: "+c.Param("id"))
		// return echo.NewHTTPError(400, "not found")
	})
	e.GET("a/b/c", func(c echo.Context) error {
		// var test string
		c.Param("some id ")
		c.QueryParam("query param")
		// c.Bind(test)
		fmt.Println("here")
		// _ = c.JSON(http.StatusOK, "Hello, World!")
		// return c.JSON(400, errors.New("not error").Error())
		return c.String(http.StatusOK, "without")
		// return echo.NewHTTPError(400, "not found")
	})

	if err := e.Start(":8081"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
