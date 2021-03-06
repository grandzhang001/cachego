package cachego

import (
	"sync"
	"time"
)

//ErrSyncMapKeyNotFound returns an error when the key is not found.
const ErrSyncMapKeyNotFound = err("key not found")

type (
	syncMapItem struct {
		data     string
		duration int64
	}

	syncMap struct {
		storage *sync.Map
	}
)

// NewSyncMap creates an instance of SyncMap cache driver
func NewSyncMap() Cache {
	return &syncMap{&sync.Map{}}
}

func (sm *syncMap) read(key string) (*syncMapItem, error) {
	v, ok := sm.storage.Load(key)

	if !ok {
		return nil, ErrSyncMapKeyNotFound
	}

	item := v.(*syncMapItem)

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		_ = sm.Delete(key)
		return nil, ErrCacheExpired
	}

	return item, nil
}

// Contains checks if cached key exists in SyncMap storage
func (sm *syncMap) Contains(key string) bool {
	if _, err := sm.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from SyncMap storage
func (sm *syncMap) Delete(key string) error {
	sm.storage.Delete(key)
	return nil
}

// Fetch retrieves the cached value from key of the SyncMap storage
func (sm *syncMap) Fetch(key string) (string, error) {
	item, err := sm.read(key)

	if err != nil {
		return "", err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the SyncMap storage
func (sm *syncMap) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := sm.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the SyncMap storage
func (sm *syncMap) Flush() error {
	sm.storage = &sync.Map{}
	return nil
}

// Save a value in SyncMap storage by key
func (sm *syncMap) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	item := &syncMapItem{value, duration}

	sm.storage.Store(key, item)

	return nil
}
