package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

func updateUserData(userId string) {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cUser := session.DB("tindetective").C("user")

	var friend_result = tinderRequest("https://api.gotinder.com/user/" + userId)

	var friend User

	err2 := json.Unmarshal([]byte(friend_result), &friend)

	if err2 != nil {
		fmt.Println(friend_result)
	}

	err = cUser.Insert(&friend)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(friend.Results.Name)

	time.Sleep(300 * time.Millisecond)
}

func updateAllByCollection(col, field string) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	collection := session.DB("tindetective").C(col)

	var result []string

	err = collection.Find(nil).Distinct(field, &result)

	if err != nil {
		log.Fatal(err)
	}

	for _, element := range result {
		updateUserData(element)
	}
}

func updateFBUsers(c echo.Context) error {

	go updateAllByCollection("groupFriendsResults", "userid")

	return c.JSON(http.StatusOK, nil)
}

func updateRecCoreUsers(c echo.Context) error {

	go updateAllByCollection("recCoreResultUser", "id")

	return c.JSON(http.StatusOK, nil)
}

func updateRecSocial(c echo.Context) error {

	go updateAllByCollection("recSocialResultMember", "id")
	go updateAllByCollection("recSocialResultOwner", "id")

	return c.JSON(http.StatusOK, nil)
}

func updateAllUsersAux() {

	updateAllByCollection("groupFriendsResults", "userid")
	updateAllByCollection("recCoreResultUser", "id")
	updateAllByCollection("recSocialResultMember", "id")
	updateAllByCollection("recSocialResultOwner", "id")
}

func updateAllUsers(c echo.Context) error {

	go updateAllUsersAux()

	return c.JSON(http.StatusOK, nil)
}

func updateUser(c echo.Context) error {

	id := c.Param("id")

	updateUserData(id)

	return c.JSON(http.StatusOK, nil)
}
