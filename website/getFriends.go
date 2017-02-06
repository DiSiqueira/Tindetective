package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"time"
)

type FriendList struct {
	Status  int `json:"status"`
	Results []struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Photo  []struct {
			ProcessedFiles []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"processedFiles"`
		} `json:"photo"`
		InSquad bool `json:"in_squad"`
	} `json:"results"`
}

type User struct {
	Status  int `json:"status"`
	Results struct {
		ConnectionCount int           `json:"connection_count"`
		CommonLikes     []interface{} `json:"common_likes"`
		CommonInterests []interface{} `json:"common_interests"`
		CommonFriends   []interface{} `json:"common_friends"`
		ID              string        `json:"_id"`
		Badges          []interface{} `json:"badges"`
		Bio             string        `json:"bio"`
		BirthDate       time.Time     `json:"birth_date"`
		Name            string        `json:"name"`
		PingTime        time.Time     `json:"ping_time"`
		Photos          []struct {
			URL            string `json:"url"`
			ProcessedFiles []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"processedFiles"`
			Extension string `json:"extension"`
			ID        string `json:"id"`
			FileName  string `json:"fileName"`
		} `json:"photos"`
		Jobs              []interface{} `json:"jobs"`
		Schools           []interface{} `json:"schools"`
		Teasers           []interface{} `json:"teasers"`
		SpotifyThemeTrack struct {
			Artists []struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"artists"`
			Album struct {
				Images []struct {
					Width  int    `json:"width"`
					URL    string `json:"url"`
					Height int    `json:"height"`
				} `json:"images"`
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"album"`
			PreviewURL string `json:"preview_url"`
			Name       string `json:"name"`
			ID         string `json:"id"`
			URI        string `json:"uri"`
		} `json:"spotify_theme_track"`
		Gender            int    `json:"gender"`
		BirthDateInfo     string `json:"birth_date_info"`
		DistanceMi        int    `json:"distance_mi"`
		CommonConnections []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Photo struct {
				Small  string `json:"small"`
				Medium string `json:"medium"`
				Large  string `json:"large"`
			} `json:"photo"`
			Degree int `json:"degree"`
		} `json:"common_connections"`
	} `json:"results"`
}

type RecsCore struct {
	Status  int `json:"status"`
	Results []struct {
		Type         string `json:"type"`
		GroupMatched bool   `json:"group_matched"`
		User         struct {
			DistanceMi        int           `json:"distance_mi"`
			CommonConnections []interface{} `json:"common_connections"`
			ConnectionCount   int           `json:"connection_count"`
			CommonLikes       []interface{} `json:"common_likes"`
			CommonInterests   []interface{} `json:"common_interests"`
			CommonFriends     []interface{} `json:"common_friends"`
			ContentHash       string        `json:"content_hash"`
			ID                string        `json:"_id"`
			Badges            []interface{} `json:"badges"`
			Bio               string        `json:"bio"`
			BirthDate         time.Time     `json:"birth_date"`
			Name              string        `json:"name"`
			PingTime          time.Time     `json:"ping_time"`
			Photos            []struct {
				ID             string `json:"id"`
				URL            string `json:"url"`
				ProcessedFiles []struct {
					Width  int    `json:"width"`
					Height int    `json:"height"`
					URL    string `json:"url"`
				} `json:"processedFiles"`
			} `json:"photos"`
			Jobs    []interface{} `json:"jobs"`
			Schools []interface{} `json:"schools"`
			Teaser  struct {
				String string `json:"string"`
			} `json:"teaser"`
			Teasers       []interface{} `json:"teasers"`
			SNumber       int           `json:"s_number"`
			Gender        int           `json:"gender"`
			BirthDateInfo string        `json:"birth_date_info"`
			GroupMatched  bool          `json:"group_matched"`
		} `json:"user"`
	} `json:"results"`
}

type RecsSocial struct {
	Status  int `json:"status"`
	Results []struct {
		Type  string `json:"type"`
		Group struct {
			ID    string `json:"id"`
			Owner struct {
				DistanceMi        int `json:"distance_mi"`
				CommonConnections []struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Photo struct {
						Small  string `json:"small"`
						Medium string `json:"medium"`
						Large  string `json:"large"`
					} `json:"photo"`
					Degree int `json:"degree"`
				} `json:"common_connections"`
				ConnectionCount   int           `json:"connection_count"`
				CommonLikes       []interface{} `json:"common_likes"`
				CommonInterests   []interface{} `json:"common_interests"`
				UncommonInterests []interface{} `json:"uncommon_interests"`
				CommonFriends     []interface{} `json:"common_friends"`
				ID                string        `json:"_id"`
				Badges            []interface{} `json:"badges"`
				Bio               string        `json:"bio"`
				BirthDate         time.Time     `json:"birth_date"`
				Name              string        `json:"name"`
				PingTime          time.Time     `json:"ping_time"`
				Photos            []struct {
					ID             string `json:"id"`
					URL            string `json:"url"`
					ProcessedFiles []struct {
						URL    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"processedFiles"`
				} `json:"photos"`
				Jobs    []interface{} `json:"jobs"`
				Schools []interface{} `json:"schools"`
				Teaser  struct {
					String string `json:"string"`
				} `json:"teaser"`
				Teasers       []interface{} `json:"teasers"`
				Gender        int           `json:"gender"`
				BirthDateInfo string        `json:"birth_date_info"`
			} `json:"owner"`
			CreatedDate int64 `json:"created_date"`
			Members     []struct {
				DistanceMi        int           `json:"distance_mi"`
				CommonConnections []interface{} `json:"common_connections"`
				ConnectionCount   int           `json:"connection_count"`
				CommonLikes       []interface{} `json:"common_likes"`
				CommonInterests   []interface{} `json:"common_interests"`
				CommonFriends     []interface{} `json:"common_friends"`
				ID                string        `json:"_id"`
				Badges            []interface{} `json:"badges"`
				Bio               string        `json:"bio"`
				BirthDate         time.Time     `json:"birth_date"`
				Name              string        `json:"name"`
				PingTime          time.Time     `json:"ping_time"`
				Photos            []struct {
					ID             string `json:"id"`
					URL            string `json:"url"`
					ProcessedFiles []struct {
						URL    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"processedFiles"`
				} `json:"photos"`
				Jobs    []interface{} `json:"jobs"`
				Schools []interface{} `json:"schools"`
				Teaser  struct {
					String string `json:"string"`
				} `json:"teaser"`
				Teasers       []interface{} `json:"teasers"`
				Gender        int           `json:"gender"`
				BirthDateInfo string        `json:"birth_date_info"`
			} `json:"members"`
			Expired bool   `json:"expired"`
			Status  string `json:"status"`
		} `json:"group"`
		GroupMatched bool `json:"group_matched"`
	} `json:"results"`
}

func tinderRequest(url string) string {

	request := gorequest.New()
	_, body, _ := request.Get(url).
		Set("Host", "api.gotinder.com").
		Set("Authorization", "Token token=\""+token.authToken+"\"").
		Set("x-client-version", "69105").
		Set("app-version", "1847").
		Set("Accept-Language", "en-BR;q=1, pt-BR;q=0.9").
		Set("platform", "ios").
		Set("Accept", "*/*").
		Set("User-Agent", "Tinder/6.9.1 (iPhone; iOS 10.2; Scale/2.00)").
		Set("Connection", "keep-alive").
		Set("X-Auth-Token", token.authToken).
		Set("os_version", "1000002").
		End()

	return body
}

func getFbFriends(c echo.Context) error {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cFBFriends := session.DB("tindetective").C("groupFriends")
	cFBFriendsResults := session.DB("tindetective").C("groupFriendsResults")

	var body = tinderRequest("https://api.gotinder.com/group/friends")

	var fl FriendList

	err2 := json.Unmarshal([]byte(body), &fl)

	if err2 != nil {
		panic(err2)
	}

	err = cFBFriends.Insert(&fl)

	if err != nil {
		log.Fatal(err)
	}

	for _, element := range fl.Results {

		err = cFBFriendsResults.Insert(&element)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(element.Name)
	}

	if c != nil {
		return c.JSON(http.StatusOK, fl.Results)
	}

	return nil
}

func getRecsCore(c echo.Context) error {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cRecCore := session.DB("tindetective").C("recCore")
	cRecCoreResultUser := session.DB("tindetective").C("recCoreResultUser")

	body := tinderRequest("https://api.gotinder.com/recs/core")

	var recList RecsCore

	err2 := json.Unmarshal([]byte(body), &recList)

	if err2 != nil {
		panic(err2)
	}

	err = cRecCore.Insert(&recList)

	if err != nil {
		log.Fatal(err)
	}

	for _, element := range recList.Results {

		if element.Type == "user" {

			err = cRecCoreResultUser.Insert(&element.User)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(element.User.Name)
		} else {
			fmt.Println("RecsCore retornou um grupo")
			os.Exit(2)
		}
	}

	if c != nil {
		return c.JSON(http.StatusOK, recList.Results)
	}

	return nil
}

func getRecsSocial(c echo.Context) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cRecSocial := session.DB("tindetective").C("recSocial")
	cRecSocialResult := session.DB("tindetective").C("recSocialResult")
	cRecSocialResultMember := session.DB("tindetective").C("recSocialResultMember")
	cRecSocialResultOwner := session.DB("tindetective").C("recSocialResultOwner")

	body := tinderRequest("https://api.gotinder.com/recs/social")

	var recList RecsSocial

	err2 := json.Unmarshal([]byte(body), &recList)

	if err2 != nil {
		panic(err2)
	}

	err = cRecSocial.Insert(&recList)

	if err != nil {
		log.Fatal(err)
	}

	for _, element := range recList.Results {

		if element.Type == "group" {

			err = cRecSocialResult.Insert(&element)

			if err != nil {
				log.Fatal(err)
			}

			err = cRecSocialResultOwner.Insert(&element.Group.Owner)

			if err != nil {
				log.Fatal(err)
			}

			for _, member := range element.Group.Members {

				err = cRecSocialResultMember.Insert(&member)

				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			fmt.Println("RecsSocial retornou um usuario")
			os.Exit(2)
		}
	}

	if c != nil {
		return c.JSON(http.StatusOK, recList.Results)
	}

	return nil
}

func getAllAux() {

	getRecsSocial(nil)
	getRecsCore(nil)
	getFbFriends(nil)

}

func getAll(c echo.Context) error {

	getAllAux()

	return c.JSON(http.StatusOK, nil)
}

func refreshUsersAux() {
	getAllAux()
	updateAllUsersAux()

}

func refreshUsers(c echo.Context) error {

	go refreshUsersAux()

	return c.JSON(http.StatusOK, nil)
}
