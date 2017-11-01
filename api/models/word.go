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
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

// TokenString is the JSON Webtoken provided by Microsoft
// var TokenString string

type EntryListQuery struct {
	XMLName xml.Name `xml:"entry_list"`
	Entry   []Entry  `xml:"entry"`
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

type MicrosoftTranslation struct {
	Translation string `xml:",innerxml"`
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
	word.saveToMemcache(c)
	return nil
}

func (word *Word) saveToMemcache(c context.Context) error {

	wordJSON, err := json.Marshal(word)
	if err != nil {
		return err
	}
	item := &memcache.Item{
		Key:   word.Text,
		Value: wordJSON,
	}

	// Add the item to the memcache, if the key does not already exist
	if err := memcache.Add(c, item); err == memcache.ErrNotStored {
		log.Printf("item with key %q already exists", item.Key)
		if err := memcache.Set(c, item); err != nil {
			log.Printf("error setting item: %v", err)
		}
	} else if err != nil {
		log.Printf("error adding item: %v", err)
	}
	return nil
}

func (word *Word) search(c context.Context) error {
	err := word.searchInMemcache(c)
	if err != nil {
		e := datastore.Get(c, word.key(c), word)
		if e != nil {
			return e
		}
		word.saveToMemcache(c)
	}
	return nil
}

func (word *Word) searchInMemcache(c context.Context) error {

	// Get the item from the memcache
	if item, err := memcache.Get(c, word.Text); err == memcache.ErrCacheMiss {
		return err
	} else if err != nil {
		return err
	} else {
		log.Printf("the item is is %q", item.Value)
		err = json.Unmarshal(item.Value, word)
		if err != nil {
			return nil
		}
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
			u.Settings.ShowTime = true
			u.Settings.ShowTranslation = true
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
		r = bytes.NewBuffer(b)
		err = json.NewDecoder(r).Decode(&p)
		if err != nil {
			log.Println("error: " + err.Error())
			return nil, err
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

	return searchAudio(c, word)
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
	if len(eq.Entry) == 0 {
		return word, nil
	}
	log.Println(eq.Entry[0].Sound.Wav.Content)
	var fileName = eq.Entry[0].Sound.Wav.Content
	var firstLetter = string(eq.Entry[0].Sound.Wav.Content[0])
	word.Audio = baseAudioURL + firstLetter + "/" + fileName
	word.save(c)
	log.Println(word)
	return word, nil
}

func searchTranslation(c context.Context, word *Word) (*Word, error) {
	var translationURL = "https://api.microsofttranslator.com/v2/http.svc/Translate?appid=Bearer%20"

	// token, err := jwt.Parse()
	// token, err := jwt(TokenString, func(token *jwt.Token) (interface{}, error) {
	// 	return token, nil
	// })

	// // if !token.Valid {
	// // 	TokenString, err = getMicrosoftToken(c)
	// // 	if err != nil {
	// // 		return nil, err
	// // 	}
	// // }
	TokenString, err := getMicrosoftToken(c)
	if err != nil {
		return nil, err
	}

	client := urlfetch.Client(c)
	req, err := http.NewRequest("GET", translationURL+TokenString+"&text="+word.Text+"&to=ja", nil)
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
	var microsoftTranslation MicrosoftTranslation
	err = xml.NewDecoder(resp.Body).Decode(&microsoftTranslation)
	if err != nil {
		return nil, err
	}

	word.Translation = microsoftTranslation.Translation
	log.Println("\n\n\nTranslation:", microsoftTranslation)
	return word, nil
}

func getMicrosoftToken(c context.Context) (string, error) {
	var tokenURL = "https://api.cognitive.microsoft.com/sts/v1.0/issueToken?Subscription-Key="
	var subscriptionKey = "07657cb89d4c4136b5165509a16c469a"

	client := urlfetch.Client(c)
	req, err := http.NewRequest("POST", tokenURL+subscriptionKey, nil)
	if err != nil {
		return "", err
	}
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Access-Control-Allow-Origin", "*")
	// req.Header.Set("Accept", "application/jwt")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	log.Println("request URL: ", resp.Request.URL)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	token := buf.String()
	// log.Println("response Body:", token)
	return token, nil
}
