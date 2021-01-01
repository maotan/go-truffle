package routes

import (
	"errors"
	"github.com/maotan/go-truffle/truffle"
	"net/http"
	"github.com/gin-gonic/gin"
)
type Person struct {
	Name  string
	Age     int
	Email string
}

func test() (int, error) {
	return 8, errors.New("678899")
}

func addPingRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/ping")
	ping.GET("/", func(c *gin.Context) {
		//panic(truffle.NewWarnError(500,"12345"))
		test()
		var p Person
		p.Name = "123"
		p.Age = 3
		base := truffle.Success(p)
		c.JSON(http.StatusCreated, base)
	})
}