package storage

import (
	"github.com/boltdb/bolt"

	"github.com/yanzay/picasso/models"
)

// SurveyStorage is a storage for interactive multiple-questions surveys
type SurveyStorage interface {
	GetSurvey(string, int) (*models.Survey, error)
	SetSurvey(string, int, *models.Survey) error
}

var surveysBucket = []byte("surveys")

func (bs *boltStorage) GetSurvey(surveyName string, id int) (*models.Survey, error) {
	var surveyBytes []byte
	err := bs.db.View(func(tx *bolt.Tx) error {
		survb := tx.Bucket(surveysBucket)
		b := survb.Bucket([]byte(surveyName))
		if b == nil {
			return nil
		}
		surveyBytes = b.Get(bytesFromID(id))
		return nil
	})
	return &models.Survey{Asking: string(surveyBytes)}, err
}

func (bs *boltStorage) SetSurvey(surveyName string, id int, survey *models.Survey) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		survb := tx.Bucket(surveysBucket)
		b, err := survb.CreateBucketIfNotExists([]byte(surveyName))
		if err != nil {
			return err
		}
		return b.Put(bytesFromID(id), []byte(survey.Asking))
	})
}
