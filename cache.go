package main

import (
	"bytes"
	"time"
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
	cache    = make(map[string]*cacheImage)
	creating = make(map[string]bool)
)

func PutCache(key string, data *bytes.Buffer) {
	cache[key] = &cacheImage{
		image:    data,
		lifeTime: time.Now().Add(cacheTimeInSeconds * time.Second).Unix(),
	}
}

func GetCache(key string) (*bytes.Buffer, bool) {
	if img, ok := cache[key]; ok {
		img.lifeTime = time.Now().Add(cacheTimeInSeconds * time.Second).Unix()
		return img.image, true
	}

	return nil, false
}