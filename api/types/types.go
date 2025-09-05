package types

import "sync"

type ContextKeyType string

type ServerConfig struct {
	Port       string
	Dir        string
	Secret     []byte
	MU         sync.RWMutex
	ContextKey ContextKeyType
}

type ContextDataType struct {
	Uid   int
	Gid   int
	Name  string
	Group string
}
