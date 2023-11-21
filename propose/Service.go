package propose

import "net"
import "google.golang.org/grpc"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/utils"
import "github.com/sirgallo/athn/proto/proposal"
import "github.com/sirgallo/athn/request"


//=========================================== Athn Propose Service



func NewProposeService[T request.Payload, U request.Result](opts *ProposeServiceOpts) *ProposeService[T, U] {
	return &ProposeService[T, U]{
		Port: utils.NormalizePort(opts.Port),
		ClientReqBuffer: make(chan *request.ClientRequest[T], CLIENT_REQ_BUFFER),
		ClientRespBuffer: make(chan *request.ClientResponse[U], CLIENT_RESP_BUFFER),
		connPool: opts.ConnectionPool,
		system: opts.System,
		zLog: *logger.NewCustomLog(NAME),
	}
}

func (propSrv *ProposeService[T, U]) StartProposeService(listener *net.Listener) {
	srv := grpc.NewServer()
	propSrv.zLog.Info("liveness gRPC server is listening on port:", propSrv.Port)
	proposal.RegisterProposalServer(srv, propSrv)

	go func() {
		err := srv.Serve(*listener)
		if err != nil { propSrv.zLog.Error("Failed to serve:", err.Error()) }
	}()

	propSrv.startClientRequestListener()
}

func (propSrv *ProposeService[T, U]) startClientRequestListener() {
	for clientReq := range propSrv.ClientReqBuffer {
		
	}
}