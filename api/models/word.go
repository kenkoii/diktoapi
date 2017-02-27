package models

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"net/http"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

// Word is the model for a comment struct/object
type Word struct {
	Text    string    `json:"text"`
	Lemma   string    `json:"lemma_text"`
	Created time.Time `json:"created"`
}

func (word *Word) key(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Word", word.Text, 0, nil)
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
func GetWord(c context.Context, id string) (interface{}, error) {
	var apiURL = "https://www.enclout.com/api/v1/term/show.json?auth_token=kdXnSYS9jhyJULBXC4Bx&text="
	var word Word
	word.Text = id

	k := word.key(c)
	err := datastore.Get(c, k, &word)
	if err != nil {
		client := urlfetch.Client(c)
		req, err := http.NewRequest("GET", apiURL+word.Text, nil)
		if err != nil {
			return nil, err
		}
		//req.Header.Set("Content-Type", "application/json")
		// req.Header.Set("X-Mashape-Key", "d2frRt9uQ8mshK76dIeYMr7TC8lip1XaEwcjsnjFx4ybfoajGl")

		resp, err := client.Do(req)
		// resp, err := client.Post(apiURL+word.Word, "application/json", nil)
		if err != nil {
			return "Search Failed: Word not found", nil
		}

		defer resp.Body.Close()
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		// body, _ := ioutil.ReadAll(resp.Body)
		// var obj interface{}
		word, err := NewWord(c, resp.Body)
		if err != nil {
			return nil, err
		}
		return &word, nil
		// err = json.NewDecoder(resp.Body).Decode(&obj)
		// log.Println("response Body:", obj)
		// return &obj, nil
	}
	return &word, nil
}
