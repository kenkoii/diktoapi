package models

import (
	"encoding/json"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/asaskevich/govalidator"
	"google.golang.org/appengine/datastore"
)

// Word is the model for a comment struct/object
type Word struct {
	Word    string    `json:"word"`
	Lemma   string    `json:"lemma"`
	Created time.Time `json:"created"`
}

func (word *Word) key(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Word", word.Word, 0, nil)
}

func (word *Word) save(c context.Context) error {
	_, err := govalidator.ValidateStruct(word)
	if err != nil {
		return err
	}

	_, err = datastore.Put(c, word.key(c), word)
	if err != nil {
		return err
	}

	return nil
}

// NewWord decodes request body into a word struct and saves it to the database
func NewWord(c context.Context, r io.ReadCloser) (*Word, error) {

	var word Word
	word.Created = time.Now()
	err := json.NewDecoder(r).Decode(&word)
	if err != nil {
		return nil, err
	}

	err = word.save(c)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

// GetWord takes a word and checks it in the database
func GetWord(c context.Context, id string) (*Word, error) {
	var word Word
	word.Word = id

	k := word.key(c)
	err := datastore.Get(c, k, &word)
	if err != nil {
		return nil, err
	}

	return &word, nil
}
