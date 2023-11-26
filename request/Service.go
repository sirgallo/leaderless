package request

import "net/http"
import "sync"

import "github.com/sirgallo/logger"
import "github.com/sirgallo/athn/common/utils"


func NewRequestService(opts *RequestServiceOpts) *RequestService {
	mux := http.NewServeMux()
	reqService := &RequestService{
		mux: mux,
		Port: utils.NormalizePort(opts.Port),
		system: opts.System,
		zLog: *logger.NewCustomLog(NAME),
		clientMappedResponseChannels: sync.Map{},
		RequestBuffer: make(chan *ClientRequest, REQ_BUFFER_SIZE),
		ResponseBuffer: make(chan *ClientResponse, RESP_BUFFER_SIZE),
	}

	reqService.RegisterCommandRoute()
	return reqService
}

func (reqService *RequestService) StartRequestService() {
	go func() {
		reqService.zLog.Info("http service starting up on port:", reqService.Port)

		srvErr := http.ListenAndServe(reqService.Port, reqService.mux)
		if srvErr != nil { reqService.zLog.Fatal("unable to start http service") }
	}()

	go func() {
		for response := range reqService.ResponseBuffer {
			go func(response *ClientResponse) {
				c, ok := reqService.clientMappedResponseChannels.Load(response.RequestID)

				if ok {
					clientChannel := c.(chan *ClientResponse)
					clientChannel <- response
				} else { reqService.zLog.Warn("no channel for resp associated with req id:", response.RequestID) }
			}(response)
		}
	}()
}