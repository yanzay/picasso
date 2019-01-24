package storage

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/boltdb/bolt"
	"github.com/yanzay/log"

	"github.com/yanzay/picasso/models"
)

// Storage is generic storage, contains all game data access
type Storage interface {
	PlayerStorage
	SurveyStorage
	BattleStorage
	SessionStorage
	GuildStorage
	Backup(io.Writer) error
}

type boltStorage struct {
	db *bolt.DB
}

// New creates new boltDB storage
func New(filename string) Storage {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		log.Fatalf("can't open storage file %s: %q", filename, err)
	}
	storage := &boltStorage{db: db}
	err = storage.createBuckets()
	if err != nil {
		log.Fatalf("bucket creation error: %q", err)
	}
	err = storage.populateData()
	if err != nil {
		log.Fatalf("unable to populate data for db: %q", err)
	}
	return storage
}

func (bs *boltStorage) Backup(w io.Writer) error {
	return bs.db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(w)
		return err
	})
}

func (bs *boltStorage) createBuckets() error {
	buckets := [][]byte{playersBucket, surveysBucket, battlesBucket, sessionsBucket, guildsBucket}
	for _, bucket := range buckets {
		err := bs.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket)
			return err
		})
		if err != nil {
			return fmt.Errorf("can't create bucket %s: %q", string(bucket), err)
		}
	}
	return nil
}

func (bs *boltStorage) populateData() error {
	newGuilds := make([]models.Guild, 0)
	bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(guildsBucket)
		for _, guild := range models.Guilds {
			bytes := b.Get(bytesFromID(guild.ID))
			if len(bytes) == 0 {
				newGuilds = append(newGuilds, guild)
			}
		}
		return nil
	})
	for _, guild := range newGuilds {
		err := bs.SetGuild(&guild)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bs *boltStorage) set(bucket []byte, id int, data []byte) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put(bytesFromID(id), data)
	})
}

func (bs *boltStorage) del(bucket []byte, id int) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Delete(bytesFromID(id))
	})
}

func (bs *boltStorage) get(bucket []byte, id int) ([]byte, error) {
	data := make([]byte, 0)
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		data = b.Get(bytesFromID(id))
		return nil
	})
	return data, err
}

func (bs *boltStorage) getAll(bucket []byte) (map[int][]byte, error) {
	data := make(map[int][]byte)
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.ForEach(func(k, v []byte) error {
			id, err := idFromBytes(k)
			if err != nil {
				return err
			}
			data[id] = v
			return nil
		})
	})
	return data, err
}

func (bs *boltStorage) setAll(bucket []byte, data map[int][]byte) error {
	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		for id, bytes := range data {
			err := b.Put(bytesFromID(id), bytes)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func bytesFromID(id int) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(id))
	return bytes
}

func idFromBytes(buf []byte) (int, error) {
	id := binary.LittleEndian.Uint64(buf)
	return int(id), nil
}
