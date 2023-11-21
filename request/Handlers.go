package request

import "encoding/json"
import "net/http"

import "github.com/sirgallo/athn/common/utils"


//=========================================== Snapshot Service Handlers

/*
	Register Command Route
		path: /command
		method: POST

		request body:
			{
				action: "string",
				payload: {
					collection: "string",
					value: "string"
				}
			}

		response body:
			{
				collection: "string",
				key: "string" | nil,
				value: "string" | nil
			}

	ingest requests and pass from the HTTP Service to the replicated log service if leader,
	or the relay service if a follower.
		1.) append a both a unique identifier for the request as well as the current node that the request was sent to.
		2.) a channel for the request to be returned is created and mapped to the request id in the mapping of response channels
		2.) A context with timeout is initialized and the route either receives the response back and returns to the client,
			or the timeout is exceeded and failure is retuned to the client
*/

func (reqService *RequestService[T, U]) RegisterCommandRoute() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost { 
			var requestData *ClientRequest[T]

			decodeErr := json.NewDecoder(r.Body).Decode(&requestData)
			if decodeErr != nil {
				http.Error(w, "failed to parse JSON request body", http.StatusBadRequest)
				return
			}

			hash, hashErr := utils.GenerateSHA256HashRandom()
			if hashErr != nil {
				http.Error(w, "error producing hash for request id", http.StatusBadRequest)
				return
			}

			clientResponseChannel := make(chan *ClientResponse[U])
			reqService.clientMappedResponseChannels.Store(hash, clientResponseChannel)

			requestData.RequestID = hash

			reqService.RequestBuffer <- requestData
			responseData :=<- clientResponseChannel

			reqService.clientMappedResponseChannels.Delete(hash)

			response := &ClientResponse[U]{ Result: responseData.Result, Success: responseData.Success }
			if responseData.ErrorMsg != nil { response.ErrorMsg = responseData.ErrorMsg }

			responseJSON, encErr := json.Marshal(response)
			if encErr != nil {
				http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(responseJSON)
		
		} else { http.Error(w, "method not allowed", http.StatusMethodNotAllowed) }
	}

	reqService.mux.HandleFunc(CommandRoute, handler)
}