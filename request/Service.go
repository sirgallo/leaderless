package request

import "net/http"
import "sync"

import "github.com/sirgallo/logger"
import "github.com/sirgallo/athn/common/utils"


//=========================================== Request Service


/*
	create a new service instance with passable options
	--> initialize the mux server and register route handlers on it, in this case the command route
		for sending operations to perform on the state machine
*/

func NewRequestService[T Payload, U Result](opts *RequestServiceOpts) *RequestService[T, U] {
	mux := http.NewServeMux()
	reqService := &RequestService[T, U]{
		mux: mux,
		Port: utils.NormalizePort(opts.Port),
		system: opts.System,
		zLog: *logger.NewCustomLog(NAME),
		clientMappedResponseChannels: sync.Map{},
		RequestBuffer: make(chan *ClientRequest[T], RequestChannelSize),
		ResponseBuffer: make(chan *ClientResponse[U], ResponseChannelSize),
	}

	reqService.RegisterCommandRoute()
	return reqService
}

/*
	Start Request Service
		separate go routines:
			1.) http server
				--> start the server to begin listening for client requests
			2.) handle response channel 
				--> for incoming respones, check the request id against the mapping of client response channels
					if the channel exists for the response, pass the response back to the route so it can be 
					returned to the client
*/

func (reqService *RequestService[T, U]) StartRequestService() {
	go func() {
		reqService.zLog.Info("http service starting up on port:", reqService.Port)

		srvErr := http.ListenAndServe(reqService.Port, reqService.mux)
		if srvErr != nil { reqService.zLog.Fatal("unable to start http service") }
	}()

	go func() {
		for response := range reqService.ResponseBuffer {
			go func(response *ClientResponse[U]) {
				c, ok := reqService.clientMappedResponseChannels.Load(response.RequestID)

				if ok {
					clientChannel := c.(chan *ClientResponse[U])
					clientChannel <- response
				} else { reqService.zLog.Warn("no channel for resp associated with req id:", response.RequestID) }
			}(response)
		}
	}()
}