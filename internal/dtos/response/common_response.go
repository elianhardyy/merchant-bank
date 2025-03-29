package response

import (
	"encoding/json"
	"net/http"
)

func CommonResponse(w http.ResponseWriter, apiResponse ApiResponse) {
	response, _ := json.Marshal(apiResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiResponse.Status)
	w.Write(response)
}
