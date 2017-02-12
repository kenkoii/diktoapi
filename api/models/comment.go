package models

import (
	"encoding/json"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/asaskevich/govalidator"
	"google.golang.org/appengine/datastore"
)

// Comment is the model for a comment struct/object
type Comment struct {
	IntIDModel
	Text     string    `json:"text"`
	PostedBy string    `json:"posted_by"`
	Created  time.Time `json:"created"`
}

func defaultCommentKeyParam(c context.Context) *datastore.Key {
	return datastore.NewIncompleteKey(c, "Comment", nil)
}

func (comment *Comment) save(c context.Context) error {
	k, err := datastore.Put(c, comment.key(c, "Comment"), comment)
	if err != nil {
		return err
	}

	comment.ID = k.IntID()
	return nil
}

func NewComment(c context.Context, r io.ReadCloser) (*Comment, error) {

	var comment Comment
	err := json.NewDecoder(r).Decode(&comment)
	if err != nil {
		return nil, err
	}

	comment.ID = 0
	comment.Created = time.Now()
	_, err = govalidator.ValidateStruct(comment)
	if err != nil {
		return nil, err
	}

	err = comment.save(c)
	// err = Save(c, comment.key(c, "Comment"), comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}
