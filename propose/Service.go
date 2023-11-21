package propose

import "net"
import "google.golang.org/grpc"

import "github.com/sirgallo/logger"

import "github.com/sirgallo/athn/common/utils"
import "github.com/sirgallo/athn/proto/proposal"
import "github.com/sirgallo/athn/request"


func NewProposeService(opts *ProposeServiceOpts) *ProposeService {
	return &ProposeService{
		Port: utils.NormalizePort(opts.Port),
		ClientReqBuffer: make(chan *request.ClientRequest, CLIENT_REQ_BUFFER),
		ClientRespBuffer: make(chan *request.ClientResponse, CLIENT_RESP_BUFFER),
		connPool: opts.ConnectionPool,
		system: opts.System,
		zLog: *logger.NewCustomLog(NAME),
	}
}

func (propSrv *ProposeService) StartProposeService(listener *net.Listener) {
	srv := grpc.NewServer()
	propSrv.zLog.Info("liveness gRPC server is listening on port:", propSrv.Port)
	proposal.RegisterProposalServer(srv, propSrv)

	go func() {
		err := srv.Serve(*listener)
		if err != nil { propSrv.zLog.Error("Failed to serve:", err.Error()) }
	}()

	propSrv.startClientRequestListener()
}

func (propSrv *ProposeService) startClientRequestListener() {
	for clientReq := range propSrv.ClientReqBuffer {
		
	}
}