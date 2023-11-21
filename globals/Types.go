package globals

import bolt "go.etcd.io/bbolt"

import "github.com/sirgallo/logger"


type Globals struct {
	file string
	globalDb *bolt.DB
	zLog logger.CustomLog
}


const NAME = "Globals DB"