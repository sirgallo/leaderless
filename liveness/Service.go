package liveness

import "net"
import "time"
import "google.golang.org/grpc"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/utils"
import "github.com/sirgallo/athn/proto/liveness"


func NewLivenessService(opts *LivenessServiceOpts) *LivenessService {
	return &LivenessService{
		Port: utils.NormalizePort(opts.Port),
		connPool: opts.ConnectionPool,
		livenessMessages: make(chan *liveness.LivenessMessage, LIVENESS_MESSAGE_BUFFER),
		system: opts.System,
		zLog: *logger.NewCustomLog(NAME),
	}
}

func (liveSrv *LivenessService) StartLivenessService(listener *net.Listener) {
	srv := grpc.NewServer()
	liveSrv.zLog.Info("liveness gRPC server is listening on port:", liveSrv.Port)
	liveness.RegisterLivenessServer(srv, liveSrv)

	go func() {
		err := srv.Serve(*listener)
		if err != nil { liveSrv.zLog.Error("Failed to serve:", err.Error()) }
	}()

	liveSrv.startHeartbeatNeighbors()
}

func (liveSrv *LivenessService) startHeartbeatNeighbors() {
	liveSrv.heartBeatTimer = time.NewTimer(HEARTBEAT_INTERVAL)
	timeoutChan := make(chan bool)

	go func() {
		resetTimer := func() {
			if ! liveSrv.heartBeatTimer.Stop() {
				select {
					case <-liveSrv.heartBeatTimer.C:
					default:
				}
			}
		
			liveSrv.heartBeatTimer.Reset(HEARTBEAT_INTERVAL)
		}

		for range liveSrv.heartBeatTimer.C {
			timeoutChan <- true
			resetTimer()
		}
	}()
}