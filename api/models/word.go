package models

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"strconv"
	"time"

	"net/http"

	"strings"

	"io/ioutil"

	"bytes"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

type EntryListQuery struct {
	XMLName xml.Name `xml:"entry_list"`
	Entry   Entry    `xml:"entry"`
}

type Entry struct {
	Sound Sound `xml:"sound"`
}

type Sound struct {
	Wav Wav `xml:"wav"`
}

type Wav struct {
	Content string `xml:",innerxml"`
}

// Query is the model for a params object
type Query struct {
	ID       string `json:"id"`
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

type WordsAPIQuery struct {
	Text          string            `json:"word"`
	Definitions   []DefinitionQuery `json:"results"`
	Pronunciation map[string]string `json:"pronunciation"`
}

type Pronunciation struct {
	PartOfSpeech string `json:"partOfSpeech"`
	IPA          string `json:"IPA"`
}

type Pron struct {
	Pronunciation string `json:"pronunciation"`
}

type Definition struct {
	Text         string `json:"definition"`
	PartOfSpeech string `json:"partOfSpeech"`
	// Synonym      string `json:"synonym"`
}

type DefinitionQuery struct {
	Text         string   `json:"definition"`
	PartOfSpeech string   `json:"partOfSpeech"`
	Synonyms     []string `json:"synonyms"`
	Examples     []string `json:"examples"`
}

// Word is the model for a comment struct/object
type Word struct {
	Text          string          `json:"text"`
	Lemma         string          `json:"lemma_text"`
	Translation   string          `json:"translation"`
	Audio         string          `json:"audio"`
	Meanings      []Definition    `json:"definition" datastore:",flatten"`
	Pronunciation []Pronunciation `json:"pronunciation"`
	Examples      []string        `json:"examples"`
	Synonyms      []string        `json:"synonyms"`
	Created       time.Time       `json:"created"`
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

func (word *Word) search(c context.Context) error {
	err := datastore.Get(c, word.key(c), word)
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
	word.Lemma = strings.TrimSpace(word.Lemma)
	err = word.save(c)
	if err != nil {
		return nil, err
	}
	return &word, nil
}

// UpdateWord decodes request body into a word struct and saves it to the database
func UpdateWord(c context.Context, r io.ReadCloser) (*Word, error) {

	var word Word
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
	word.Text = id

	if word.search(c) != nil {
		return searchWord(c, &word)
	}
	return &word, nil
}

// GetLemma takes a word and checks it in the database
func GetLemma(c context.Context, id string) (*Word, error) {
	var word Word
	word.Text = id

	if word.search(c) != nil {
		return searchLemma(c, &word)
	}
	return &word, nil
}

// FavoriteWord takes a word and user info in the body, adds word to users favorites
func FavoriteWord(c context.Context, r io.ReadCloser) (int64, error) {
	var q Query
	err := json.NewDecoder(r).Decode(&q)
	log.Println(q)
	if err != nil {
		return 0, err
	}
	var word Word
	word.Text = q.ID
	if err = word.search(c); err != nil {
		log.Println(err)
		w, err := GetWord(c, q.ID)
		word.Created = w.Created
		word.Lemma = w.Lemma
		if err != nil {
			return 0, err
		}
	}

	var u User
	u.Password, _ = strconv.ParseInt(q.Password, 10, 64)
	u.ID, err = strconv.ParseInt(q.UserID, 10, 64)

	if err = u.search(c); err != nil {
		// return 0, err
		log.Println("No account!!!")
		if verifyPassword(c, u.ID, u.Password) == 1 {
			u.Created = time.Now()
			err = u.save(c)
			if err != nil {
				return 0, err
			}
		}
	}

	favorite := Favorite{Word: word.Text, Created: time.Now(), Status: "studying"}
	for _, value := range u.Favorites {
		if value.Word == word.Text {
			return 2, nil
		}
	}

	if contains(u.Favorites, word) {
		return 2, nil
	}
	u.Favorites = append(u.Favorites, favorite)
	err = u.save(c)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// FrontendFavoriteWord takes a word and user info in the body, adds word to users favorites
func FrontendFavoriteWord(c context.Context, r io.ReadCloser) (interface{}, error) {
	var q Query
	err := json.NewDecoder(r).Decode(&q)
	log.Println(q)
	if err != nil {
		return 0, err
	}
	var word Word
	word.Text = q.ID
	if err = word.search(c); err != nil {
		log.Println(err)
		w, err := GetWord(c, q.ID)
		word.Created = w.Created
		word.Lemma = w.Lemma
		if err != nil {
			return 0, err
		}
	}

	var u User
	u.Password, _ = strconv.ParseInt(q.Password, 10, 64)
	u.ID, err = strconv.ParseInt(q.UserID, 10, 64)

	if err = u.search(c); err != nil {
		// return 0, err
		log.Println("No account!!!")
		if verifyPassword(c, u.ID, u.Password) == 1 {
			u.Created = time.Now()
			err = u.save(c)
			if err != nil {
				return 0, err
			}
		}
	}

	favorite := Favorite{Word: word.Text, Created: time.Now(), Status: "studying"}
	for _, value := range u.Favorites {
		if value.Word == word.Text {
			return 2, nil
		}
	}

	if contains(u.Favorites, word) {
		return 2, nil
	}
	u.Favorites = append(u.Favorites, favorite)
	err = u.save(c)
	if err != nil {
		return 0, err
	}
	return favorite, nil
}

// RemoveFavoriteWord takes a word and user info in the body, removes word from users favorites
func RemoveFavoriteWord(c context.Context, r io.ReadCloser) (*Word, error) {
	var q Query
	err := json.NewDecoder(r).Decode(&q)
	if err != nil {
		return nil, err
	}
	var word Word
	word.Text = q.ID
	if word.search(c) != nil {
		return nil, err
	}

	var u User
	u.ID, err = strconv.ParseInt(q.UserID, 10, 64)
	if u.search(c) != nil {
		return nil, err
	}

	for key, value := range u.Favorites {
		if value.Word == word.Text {
			u.Favorites[key] = u.Favorites[len(u.Favorites)-1]
			u.Favorites = u.Favorites[:len(u.Favorites)-1]
			break
		}
	}
	err = u.save(c)
	if err != nil {
		return nil, err
	}

	return &word, nil
}

func searchLemma(c context.Context, word *Word) (*Word, error) {
	var apiURL = "https://www.enclout.com/api/v1/term/show.json?auth_token=kdXnSYS9jhyJULBXC4Bx&text="
	client := urlfetch.Client(c)
	req, err := http.NewRequest("GET", apiURL+word.Text, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	word, err = NewWord(c, resp.Body)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func contains(f []Favorite, w Word) bool {
	for _, a := range f {
		if a.Word == w.Text {
			return true
		}
	}
	return false
}

func searchWord(c context.Context, word *Word) (*Word, error) {
	// go searchAudio(c, word)
	var apiURL = "https://wordsapiv1.p.mashape.com/words/"
	client := urlfetch.Client(c)
	req, err := http.NewRequest("GET", apiURL+word.Text, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Mashape-Key", "te6AX6SnBfmshawA0zj6VToSZO3up1MQySvjsnFmGv0qYDjUV3")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	b, err := ioutil.ReadAll(resp.Body)
	r := bytes.NewBuffer(b)

	var w WordsAPIQuery
	var p Pron
	err = json.NewDecoder(r).Decode(&w)
	if err != nil {
		// log.Println(w)
		// return nil, err
		r = bytes.NewBuffer(b)
		err = json.NewDecoder(r).Decode(&p)
		if err != nil {
			log.Println("error: " + err.Error())
		}
		word.Pronunciation = append(word.Pronunciation, Pronunciation{PartOfSpeech: "all", IPA: p.Pronunciation})
	}

	word.Created = time.Now()
	for _, v := range w.Definitions {
		for _, b := range v.Examples {
			word.Examples = append(word.Examples, b)
		}
		for _, b := range v.Synonyms {
			word.Synonyms = append(word.Synonyms, b)
		}
		word.Meanings = append(word.Meanings, Definition{PartOfSpeech: v.PartOfSpeech, Text: v.Text})
	}

	for k, v := range w.Pronunciation {
		// log.Printf("key[%s] value[%s]\n", k, v)
		word.Pronunciation = append(word.Pronunciation, Pronunciation{PartOfSpeech: k, IPA: v})
	}

	// DECORATOR PATTERN
	// err = word.save(c)
	// if err != nil {
	// 	return nil, err
	// }

	return searchAudio(c, word)
	// return word, nil
	// return nil, nil
}

func searchAudio(c context.Context, word *Word) (*Word, error) {
	var apiURL = "https://www.dictionaryapi.com/api/v1/references/collegiate/xml/"
	var apiKey = "?key=720750f6-2da7-4612-bb3e-2914b923052e"
	var baseAudioURL = "http://media.merriam-webster.com/soundc11/"
	client := urlfetch.Client(c)
	req, err := http.NewRequest("GET", apiURL+word.Text+apiKey, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	log.Println("request URL: ", resp.Request.URL)

	var eq EntryListQuery
	err = xml.NewDecoder(resp.Body).Decode(&eq)
	// data, err := ioutil.ReadAll(resp.Body)
	// err = xml.Unmarshal(data, &eq)
	if err != nil {
		return nil, err
	}
	log.Println(eq.Entry.Sound.Wav.Content)
	var fileName = eq.Entry.Sound.Wav.Content
	var firstLetter = string(eq.Entry.Sound.Wav.Content[0])
	word.Audio = baseAudioURL + firstLetter + "/" + fileName
	word.save(c)
	log.Println(word)
	return word, nil
}
