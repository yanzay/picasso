package storage

import (
	"encoding/json"
	"fmt"

	"github.com/yanzay/picasso/models"
)

// BattleStorage has methods for accessing battles
type BattleStorage interface {
	GetBattle(int) (*models.Battle, error)
	SetBattle(*models.Battle) error
	GetAllBattles() ([]*models.Battle, error)
	DeleteBattle(int) error
}

var battlesBucket = []byte("battles")

func (bs *boltStorage) GetBattle(id int) (*models.Battle, error) {
	bytes, err := bs.get(battlesBucket, id)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, fmt.Errorf("Battle %d not found", id)
	}
	b := &models.Battle{}
	err = json.Unmarshal(bytes, b)
	if err != nil {
		return nil, err
	}
	b.ID = id
	return b, nil
}

func (bs *boltStorage) SetBattle(b *models.Battle) error {
	bytes, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return bs.set(battlesBucket, b.ID, bytes)
}

func (bs *boltStorage) GetAllBattles() ([]*models.Battle, error) {
	battlesBytes, err := bs.getAll(battlesBucket)
	if err != nil {
		return nil, err
	}
	battles := make([]*models.Battle, 0)
	for _, v := range battlesBytes {
		b := &models.Battle{}
		err = json.Unmarshal(v, b)
		if err != nil {
			return nil, err
		}
		battles = append(battles, b)
	}
	return battles, nil
}

func (bs *boltStorage) DeleteBattle(id int) error {
	return bs.del(battlesBucket, id)
}
