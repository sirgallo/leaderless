package liveness

import "time"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/connpool"
import "github.com/sirgallo/athn/proto/liveness"
import "github.com/sirgallo/athn/system"


type LivenessServiceOpts struct {
	Port int
	ConnectionPool *connpool.ConnectionPool
	System *system.System
}

type LivenessService struct {
	liveness.UnimplementedLivenessServer
	
	Port string
	connPool *connpool.ConnectionPool
	system *system.System
	zLog logger.CustomLog

	heartBeatTimer *time.Timer
	livenessMessages chan *liveness.LivenessMessage
}

type LivenessResponseChannels struct {
	BroadcastClose chan struct{}
	Messages chan *liveness.LivenessMessage
}

const NAME = "Liveness Service"

const (
	HEARTBEAT_INTERVAL = 150 * time.Millisecond
	LIVENESS_MESSAGE_BUFFER = 10000
	LIVENESS_RPC_TIMEOUT = 500 * time.Millisecond
)