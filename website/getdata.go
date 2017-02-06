package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func getAllLikedMe(c echo.Context) error {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	cMatches := session.DB("tindetective").C("matches")

	final := []userSimple{}
	var all interface{}

	err3 := cMatches.Find(all).Sort("-_id").All(&final)

	if err3 != nil {
		log.Fatal(err3)
	}

	return c.JSON(http.StatusOK, final)
}

type online struct {
	_id      string
	pingtime time.Time
}

func getAllUsers(c echo.Context) error {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	cUser := session.DB("tindetective").C("user")

	result := []bson.M{}

	//db.getCollection('user').aggregate(
	//	{
	//		$group:
	//		{
	//			"_id": "$results.id",
	//			"pingtime":
	//			{
	//				"$max": "$results.pingtime"
	//			}
	//		}
	//	},
	//	{
	//		$sort:
	//		{
	//			"pingtime": -1
	//		}
	//	}
	//)

	group := bson.M{
		"$group": bson.M{
			"_id": "$results.id",
			"pingtime": bson.M{
				"$max": "$results.pingtime",
			},
		},
	}

	sort := bson.M{
		"$sort": bson.M{
			"pingtime": -1,
		},
	}

	all := []bson.M{group, sort}
	err = cUser.Pipe(all).All(&result)

	//err = cUser.Find(nil).Distinct("results.id", &result)

	if err != nil {
		log.Fatal(err)
	}

	var final []userSimple

	for _, user := range result {

		result := User{}

		err = cUser.Find(bson.M{"results.id": user["_id"]}).Sort("-_id").Limit(1).One(&result)
		if err != nil {
			log.Fatal(err)
		}

		var userFinal userSimple

		userFinal.Id = result.Results.ID
		userFinal.Name = result.Results.Name
		userFinal.Pingtime = result.Results.PingTime

		if len(result.Results.Photos) > 0 {
			userFinal.Photo = result.Results.Photos[0].URL
		} else {
			fmt.Println("Sem foto?", result.Results.ID)
			fmt.Println(result)
		}

		final = append(final, userFinal)
	}

	return c.JSON(http.StatusOK, final)
}
