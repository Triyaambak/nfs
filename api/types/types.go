package types

import "sync"

type ServerConfig struct {
	Port string
	Dir  string
	MU   sync.RWMutex
}
