package models

import (
	"encoding/json"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/asaskevich/govalidator"
	"google.golang.org/appengine/datastore"
)

// Settings is the model for a user
type Settings struct {
	SortAZ            bool `json:"sortAZ"`
	ShowPronunciation bool `json:"showPronunciation"`
	ShowTime          bool `json:"showTime"`
	ShowTranslation   bool `json:"showTranslation"`
}

// User is the model for a user
type User struct {
	ID        int64      `json:"id"`
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

func GetUser(c context.Context, id int64) (*User, error) {
	var user User
	user.ID = id
	k := user.key(c)
	err := datastore.Get(c, k, &user)
	if err != nil {
		user.Created = time.Now()
		err = user.save(c)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}
