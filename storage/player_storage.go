package storage

import (
	"encoding/json"

	"github.com/yanzay/picasso/models"
)

// PlayerStorage has methods for accessing players
type PlayerStorage interface {
	GetPlayer(int) (*models.Player, error)
	GetAllPlayers() ([]*models.Player, error)
	SetPlayer(*models.Player) error
	SetAllPlayers([]*models.Player) error
	DeletePlayer(int) error
}

var playersBucket = []byte("players")

func (bs *boltStorage) GetAllPlayers() ([]*models.Player, error) {
	playersBytes, err := bs.getAll(playersBucket)
	if err != nil {
		return nil, err
	}
	players := make([]*models.Player, 0)
	for _, v := range playersBytes {
		p := &models.Player{}
		err = json.Unmarshal(v, p)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, nil
}

func (bs *boltStorage) GetPlayer(id int) (*models.Player, error) {
	pBytes, err := bs.get(playersBucket, id)
	if err != nil {
		return nil, err
	}
	if len(pBytes) == 0 {
		return models.NewPlayer(id), nil
	}
	p := &models.Player{}
	err = json.Unmarshal(pBytes, p)
	return p, err
}

func (bs *boltStorage) SetPlayer(p *models.Player) error {
	pBytes, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return bs.set(playersBucket, p.ID, pBytes)
}

func (bs *boltStorage) SetAllPlayers(players []*models.Player) error {
	playersBytes := make(map[int][]byte)
	for _, player := range players {
		pBytes, err := json.Marshal(player)
		if err != nil {
			return err
		}
		playersBytes[player.ID] = pBytes
	}
	return bs.setAll(playersBucket, playersBytes)
}

func (bs *boltStorage) DeletePlayer(id int) error {
	return bs.del(playersBucket, id)
}
