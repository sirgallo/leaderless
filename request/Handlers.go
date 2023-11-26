package request

import "encoding/json"
import "net/http"

import "github.com/sirgallo/athn/common/utils"


func (reqService *RequestService) RegisterCommandRoute() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost { 
			var requestData *ClientRequest

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

			clientResponseChannel := make(chan *ClientResponse)
			reqService.clientMappedResponseChannels.Store(hash, clientResponseChannel)

			requestData.RequestID = hash

			reqService.RequestBuffer <- requestData
			responseData :=<- clientResponseChannel

			reqService.clientMappedResponseChannels.Delete(hash)

			response := &ClientResponse{ Result: responseData.Result, Success: responseData.Success }
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

	reqService.mux.HandleFunc(COMMAND_ROUTE, handler)
}