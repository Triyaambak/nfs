package types

import "sync"

type ServerConfig struct {
	Port   string
	Dir    string
	Secret []byte
	MU     sync.RWMutex
}
