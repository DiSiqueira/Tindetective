package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/mgo.v2"
	"html/template"
	"io"
	"strings"
	"time"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type TinderAuth struct {
	Token string `json:"token"`
	User  struct {
		ID              string    `json:"_id"`
		ActiveTime      time.Time `json:"active_time"`
		CanCreateSquad  bool      `json:"can_create_squad"`
		CreateDate      time.Time `json:"create_date"`
		AgeFilterMax    int       `json:"age_filter_max"`
		AgeFilterMin    int       `json:"age_filter_min"`
		APIToken        string    `json:"api_token"`
		Bio             string    `json:"bio"`
		BirthDate       time.Time `json:"birth_date"`
		ConnectionCount int       `json:"connection_count"`
		DistanceFilter  int       `json:"distance_filter"`
		FullName        string    `json:"full_name"`
		Groups          []string  `json:"groups"`
		Gender          int       `json:"gender"`
		GenderFilter    int       `json:"gender_filter"`
		Interests       []struct {
			Name        string `json:"name"`
			ID          string `json:"id"`
			CreatedTime string `json:"created_time"`
		} `json:"interests"`
		Name         string    `json:"name"`
		PingTime     time.Time `json:"ping_time"`
		Discoverable bool      `json:"discoverable"`
		Photos       []struct {
			ID             string `json:"id"`
			Main           bool   `json:"main,omitempty"`
			Shape          string `json:"shape,omitempty"`
			FileName       string `json:"fileName"`
			FbID           string `json:"fbId"`
			Extension      string `json:"extension"`
			ProcessedFiles []struct {
				Width  int    `json:"width"`
				Height int    `json:"height"`
				URL    string `json:"url"`
			} `json:"processedFiles"`
			URL              string  `json:"url"`
			YoffsetPercent   int     `json:"yoffset_percent,omitempty"`
			XdistancePercent float64 `json:"xdistance_percent,omitempty"`
			YdistancePercent float64 `json:"ydistance_percent,omitempty"`
			XoffsetPercent   float64 `json:"xoffset_percent,omitempty"`
		} `json:"photos"`
		PhotosProcessing bool `json:"photos_processing"`
		Jobs             []struct {
			Company struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Displayed bool   `json:"displayed"`
			} `json:"company"`
			Title struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Displayed bool   `json:"displayed"`
			} `json:"title"`
		} `json:"jobs"`
		Schools            []interface{} `json:"schools"`
		SquadsDiscoverable bool          `json:"squads_discoverable"`
		SquadsOnly         bool          `json:"squads_only"`
		Purchases          []interface{} `json:"purchases"`
		IsNewUser          bool          `json:"is_new_user"`
	} `json:"user"`
	Versions struct {
		ActiveText         string `json:"active_text"`
		AgeFilter          string `json:"age_filter"`
		Matchmaker         string `json:"matchmaker"`
		Trending           string `json:"trending"`
		TrendingActiveText string `json:"trending_active_text"`
	} `json:"versions"`
	Globals struct {
		Friends                  bool   `json:"friends"`
		InviteType               string `json:"invite_type"`
		RecsInterval             int    `json:"recs_interval"`
		UpdatesInterval          int    `json:"updates_interval"`
		RecsSize                 int    `json:"recs_size"`
		MatchmakerDefaultMessage string `json:"matchmaker_default_message"`
		ShareDefaultText         string `json:"share_default_text"`
		BoostDecay               int    `json:"boost_decay"`
		BoostUp                  int    `json:"boost_up"`
		BoostDown                int    `json:"boost_down"`
		Sparks                   bool   `json:"sparks"`
		Kontagent                bool   `json:"kontagent"`
		SparksEnabled            bool   `json:"sparks_enabled"`
		KontagentEnabled         bool   `json:"kontagent_enabled"`
		Mqtt                     bool   `json:"mqtt"`
		TinderSparks             bool   `json:"tinder_sparks"`
		MomentsInterval          int    `json:"moments_interval"`
		FetchConnections         bool   `json:"fetch_connections"`
		Plus                     bool   `json:"plus"`
	} `json:"globals"`
}
type Token struct {
	authToken string
}

var token *Token

func getTinderToken() *Token {

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	fbAuth := ""

	url := "https://api.gotinder.com/auth"

	request := gorequest.New()
	_, body, _ := request.Post(url).
		Send(`{"facebook_token":"`+fbAuth+`"}`).
		Set("Host", "api.gotinder.com").
		Set("User-Agent", "Tinder/6.9.1 (iPhone; iOS 10.2; Scale/2.00)").
		End()

	cAuth := session.DB("tindetective").C("auth")

	var authResponse TinderAuth

	body = strings.Replace(body, "\"main\",", "true,", -1)

	fmt.Println(body)

	err2 := json.Unmarshal([]byte(body), &authResponse)

	if err2 != nil {
		fmt.Println(body)
		panic(err2)
	}

	err = cAuth.Insert(&authResponse)

	return &Token{authToken: authResponse.Token}
}

func main() {

	token = getTinderToken()

	fmt.Println(token.authToken)

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.File("/", "public/index.html")
	e.File("/likedme", "public/liked.html")
	e.File("/all", "public/all.html")

	e.GET("/user/:id", getUser)

	e.GET("/api/refresh", refreshUsers)

	e.GET("/api/search/likedme", getLikedMe)
	e.GET("/api/search/all", getAll)
	e.GET("/api/search/recs/social", getRecsSocial)
	e.GET("/api/search/recs/core", getRecsCore)
	e.GET("/api/search/fbfriends", getFbFriends)

	e.GET("/api/update/all", updateAllUsers)
	e.GET("/api/update/user/:id", updateUser)
	e.GET("/api/update/fbusers", updateFBUsers)
	e.GET("/api/update/recs/core", updateRecCoreUsers)
	e.GET("/api/update/recs/social", updateRecSocial)

	e.GET("/api/get/likedme", getAllLikedMe)
	e.GET("/api/get/users", getAllUsers)

	e.Logger.Fatal(e.Start(":8888"))
}
