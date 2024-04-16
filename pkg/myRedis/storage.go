package pkg

import (
	"time"
)

type Storage struct {
	data map[string]string
}

func (s Storage) SetValue(key string, value string, duration int) bool {
	s.data[key] = value
	go s.KillValue(key, duration)
	return true
}

func (s Storage) GetValue(key string) string {
	return s.data[key]
}

func (s Storage) KillValue(key string, duration int) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
	delete(s.data, key)
}

func newStorage() Storage {
	return Storage{
		data: make(map[string]string),
	}
}
