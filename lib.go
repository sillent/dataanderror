package dataanderror

import (
	"fmt"
	"maps"
	"strings"
	"sync"
)

type from int

const (
	data from = iota
	err
)

type Keyable interface {
	comparable
	~string
}
type DataAndError[K Keyable, D any, E any] struct {
	rw sync.RWMutex

	data  map[K]D
	error map[K]E
}

func New[K Keyable, D, E any]() *DataAndError[K, D, E] {
	return &DataAndError[K, D, E]{
		rw:    sync.RWMutex{},
		data:  make(map[K]D),
		error: make(map[K]E),
	}
}

func (de *DataAndError[K, D, E]) Store(key K, data D) {
	de.rw.Lock()
	de.data[key] = data
	delete(de.error, key)
	de.rw.Unlock()
}

func (de *DataAndError[K, D, E]) StoreError(key K, error E) {
	de.rw.Lock()
	defer de.rw.Unlock()
	if _, exist := de.data[key]; exist {
		return
	}
	de.error[key] = error
}

func (de *DataAndError[K, D, E]) Load(key K) (D, bool) {
	de.rw.RLock()
	defer de.rw.RUnlock()
	value, exist := de.data[key]
	return value, exist
}

func (de *DataAndError[K, D, E]) LoadError(key K) (E, bool) {
	de.rw.RLock()
	defer de.rw.RUnlock()
	value, exist := de.error[key]
	return value, exist
}

func (de *DataAndError[K, D, E]) CopiedData() map[K]D {
	de.rw.RLock()
	defer de.rw.RUnlock()
	copied := make(map[K]D, len(de.data))
	maps.Copy(copied, de.data)
	return copied
}

func (de *DataAndError[K, D, E]) CopiedError() map[K]E {
	de.rw.RLock()
	defer de.rw.RUnlock()
	copied := make(map[K]E, len(de.error))
	maps.Copy(copied, de.error)
	return copied
}

func (de *DataAndError[K, D, E]) DataString() string {
	de.rw.RLock()
	defer de.rw.RUnlock()
	totalStr := ""
	for k, v := range de.data {
		kv := strings.Join([]string{string(k), fmt.Sprintf("%v", v)}, ":")
		totalStr = strings.Join([]string{totalStr, kv}, "|")
	}
	return totalStr
}

func (de *DataAndError[K, D, E]) ErrorString() string {
	de.rw.RLock()
	defer de.rw.RUnlock()
	totalStr := ""
	for k, v := range de.error {
		kv := strings.Join([]string{string(k), fmt.Sprintf("%v", v)}, ":")
		totalStr = strings.Join([]string{totalStr, kv}, "|")
	}
	return totalStr
}

func (de *DataAndError[K, D, E]) Remove(key K) {
	de.remove(key, data)
}

func (de *DataAndError[K, D, E]) RemoveError(key K) {
	de.remove(key, err)
}

func (de *DataAndError[K, D, E]) remove(key K, f from) {
	de.rw.Lock()
	defer de.rw.Unlock()
	switch f {
	case data:
		delete(de.data, key)
	case err:
		delete(de.error, key)
	}
}
