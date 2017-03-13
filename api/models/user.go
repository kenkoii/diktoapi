package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"io/ioutil"

	"github.com/asaskevich/govalidator"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

// Settings is the model for a user
type Settings struct {
	SortAZ            bool `json:"sortAZ"`
	ShowPronunciation bool `json:"showPronunciation"`
	ShowTime          bool `json:"showTime"`
	ShowTranslation   bool `json:"showTranslation"`
}

// ErrorMessage is the model for a user
type ErrorMessage struct {
	Message string `json:"error"`
}

// User is the model for a user
type User struct {
	ID        int64      `json:"id"`
	Password  int64      `json:"password"`
	Favorites []Favorite `json:"favorites"`
	Created   time.Time  `json:"created"`
	Settings  Settings   `json:"settings"`
}

func (user *User) key(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "User", "", user.ID, nil)
}

func (user *User) save(c context.Context) error {
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}

	_, err = datastore.Put(c, user.key(c), user)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) search(c context.Context) error {
	err := datastore.Get(c, user.key(c), user)
	if err != nil {
		return err
	}
	return nil
}

func NewUser(c context.Context, r io.ReadCloser) (*User, error) {

	var user User
	user.Created = time.Now()
	err := json.NewDecoder(r).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = user.save(c)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(c context.Context, id int64, r io.ReadCloser) (*User, error) {

	var user User
	user.Created = time.Now()
	err := json.NewDecoder(r).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = user.save(c)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(c context.Context, id int64, password int64) (interface{}, error) {
	var user User
	user.ID = id
	k := user.key(c)
	err := datastore.Get(c, k, &user)
	if err != nil {
		switch verifyPassword(c, id, password) {
		case 1:
			user.Created = time.Now()
			user.Password = password
			err = user.save(c)
			if err != nil {
				return nil, err
			}
		default:
			return ErrorMessage{Message: "user not found"}, nil
		}
	}
	if user.Password != password {
		return ErrorMessage{Message: "password mismatch"}, nil
	}
	return &user, nil
}

func verifyPassword(c context.Context, id int64, password int64) uint64 {
	var apiURL = "http://englishstoryserver.appspot.com/ConfirmPassword"
	var reqURL = apiURL + "?userId=" + strconv.FormatInt(id, 10) + "&password=" + strconv.FormatInt(password, 10)
	client := urlfetch.Client(c)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return 0
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	log.Println("request URL: ", resp.Request.URL)

	serverData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	x, _ := strconv.ParseUint(string(serverData), 10, 64)
	log.Println(x)
	return x
}
