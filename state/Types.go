package state

import bolt "go.etcd.io/bbolt"

import "github.com/sirgallo/logger"


type State struct {
	file string
	stateDb *bolt.DB
	zLog logger.CustomLog
}

type KeyValuePair struct {
	Key []byte
	Value []byte
}

type StatePayload struct {
	RequestID [32]byte
	Operation Action
	KVPair KeyValuePair
}

type StateResponse struct {
	RequestID [32]byte
	KVPair KeyValuePair
}

type Action = string

const NAME = "State DB"