package models

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// IntIDModel is the base model which id is int64
type IntIDModel struct {
	// Key *datastore.Key `datastore:"-"` //`datastore:"__key__"`
	ID int64 `json:"id" datastore:"-"`
}

func (x *IntIDModel) key(c context.Context, kind string) *datastore.Key {
	if x.ID == 0 {
		return datastore.NewIncompleteKey(c, kind, nil)
	}
	return datastore.NewKey(c, kind, "", x.ID, nil)
}

// Save supposedly is a helper function for saving to datastore
func Save(c context.Context, key *datastore.Key, src interface{}) error {
	_, err := datastore.Put(c, key, src)
	if err != nil {
		return err
	}
	return nil
}
