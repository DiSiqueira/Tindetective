package main

import (
	"fmt"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"
)

type userSimple struct {
	Count    int
	Id       string
	Name     string
	Photo    string
	Total    float64
	Pingtime time.Time
}

func getLikedMe(c echo.Context) error {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	cRecCore := session.DB("tindetective").C("recCore")
	result := []RecsCore{}

	var all interface{}

	err2 := cRecCore.Find(all).Sort("-_id").Limit(4).All(&result)

	if err2 != nil {
		log.Fatal(err2)
	}

	resultados := []userSimple{}

	for _, elem := range result {
		for _, usuario := range elem.Results {

			found := false

			for index, m := range resultados {

				if m.Id == usuario.User.ID {
					resultados[index].Count++
					found = true
				}

			}

			if !found {
				resultados = append(resultados, userSimple{
					Count: 1,
					Id:    usuario.User.ID,
					Name:  usuario.User.Name,
					Photo: usuario.User.Photos[0].URL,
				})
			}

		}
	}

	cMatches := session.DB("tindetective").C("matches")

	for _, usu := range resultados {
		if usu.Count > 1 {

			usu.Total = float64((usu.Count * 100) / 4)

			_, err = cMatches.Upsert(&usu, &usu)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(usu.Id, " - ", usu.Name, " - ", usu.Total, "%")
		}
	}

	final := []userSimple{}

	err3 := cMatches.Find(all).Sort("-_id").All(&final)

	if err3 != nil {
		log.Fatal(err3)
	}

	return c.JSON(http.StatusOK, final)
}
