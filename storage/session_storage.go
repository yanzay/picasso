package storage

type SessionStorage interface {
	Set(int64, string)
	Get(int64) string
	Reset(int64)
}

var sessionsBucket = []byte("sessions")

func (bs *boltStorage) Get(id int64) string {
	bytes, err := bs.get(sessionsBucket, int(id))
	if err != nil {
		return ""
	}
	if len(bytes) == 0 {
		return ""
	}
	return string(bytes)
}

func (bs *boltStorage) Set(id int64, path string) {
	bs.set(sessionsBucket, int(id), []byte(path))
}

func (bs *boltStorage) Reset(id int64) {
	bs.del(sessionsBucket, int(id))
}
