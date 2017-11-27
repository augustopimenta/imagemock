package main

import (
	"bytes"
	"time"
	"sync"
	"fmt"
	"runtime/debug"
)

type cacheImage struct {
	image    *bytes.Buffer
	lifeTime int64
}

const (
	cacheTimeInSeconds              = 100
	cacheRemoveRoutineTimeInSeconds = 5
)

var (
	mt sync.RWMutex
	cache    = make(map[string]*cacheImage)
)

func putCache(key string, data *bytes.Buffer) {
	mt.Lock()
	defer mt.Unlock()

	cache[key] = &cacheImage{
		image:    data,
		lifeTime: time.Now().Add(cacheTimeInSeconds * time.Second).Unix(),
	}
}

func getCache(key string) (*bytes.Buffer, bool) {
	mt.RLock()
	defer mt.RUnlock()

	if img, ok := cache[key]; ok {
		img.lifeTime = time.Now().Add(cacheTimeInSeconds * time.Second).Unix()

		return img.image, true
	}

	return nil, false
}

func remCache(key string) {
	mt.Lock()
	defer mt.Unlock()

	delete(cache, key)
}

func clearOldCache() {
	for {
		mt.RLock()
		for k, v := range cache {
			if v.lifeTime < time.Now().Unix() {
				remCache(k)
				fmt.Println("Removendo cache, index: ", k)
			}
		}
		mt.RUnlock()

		time.Sleep(cacheRemoveRoutineTimeInSeconds * time.Second)
		debug.FreeOSMemory()
	}
}