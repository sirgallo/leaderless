package system

import "sync"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/globals"
import "github.com/sirgallo/athn/state"


type SystemOpts struct {
	Seed []byte
	Globals *globals.Globals
	State *state.State
}

type System struct {
	zLog *logger.CustomLog
	Host string
	NodeId[32]byte
	PropagationFactor uint64

	Neighbors *sync.Map
	PreviousNeighbors *sync.Map

	Globals *globals.Globals
	State *state.State

	SystemMutex sync.Mutex
}


const NAME = "System"