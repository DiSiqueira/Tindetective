package main

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

// e.GET("/user/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	cUser := session.DB("tindetective").C("user")

	result := User{}

	err = cUser.Find(bson.M{"results.id": id}).Sort("-_id").Limit(1).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return c.Render(http.StatusOK, "user", result.Results)

}
