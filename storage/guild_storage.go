package storage

import (
	"encoding/json"
	"fmt"

	"github.com/yanzay/picasso/models"
)

// GuildStorage has methods for accessing guilds
type GuildStorage interface {
	GetGuild(int) (*models.Guild, error)
	SetGuild(*models.Guild) error
	GetAllGuilds() ([]*models.Guild, error)
}

var guildsBucket = []byte("guilds")

func (bs *boltStorage) GetGuild(id int) (*models.Guild, error) {
	bytes, err := bs.get(guildsBucket, id)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, fmt.Errorf("Guild %d not found", id)
	}
	g := &models.Guild{}
	err = json.Unmarshal(bytes, g)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (bs *boltStorage) SetGuild(g *models.Guild) error {
	bytes, err := json.Marshal(g)
	if err != nil {
		return err
	}
	return bs.set(guildsBucket, g.ID, bytes)
}

func (bs *boltStorage) GetAllGuilds() ([]*models.Guild, error) {
	guildsBytes, err := bs.getAll(guildsBucket)
	if err != nil {
		return nil, err
	}
	guilds := make([]*models.Guild, 0)
	for _, v := range guildsBytes {
		g := &models.Guild{}
		err = json.Unmarshal(v, g)
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, g)
	}
	return guilds, nil
}
