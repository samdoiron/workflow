package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiURL = "http://localhost:8080/mortgage/api"

type APIResponse struct {
	OK      bool
	Message string
}

// PutEmployerInfo submits MortgageEmployerInfo to the mortgage API
func PutEmployerInfo(id int, info MortgageEmployerInfo) (APIResponse, error) {
	encoded, err := json.Marshal(info)
	if err != nil {
		return APIResponse{}, err
	}

	url := fmt.Sprintf("%v/applications/%v/employer-info", apiURL, id)
	response, err := put(url, encoded)
	if err != nil {
		return APIResponse{}, err
	}

	return decodeAPIResponse(response)
}

func decodeAPIResponse(response []byte) (APIResponse, error) {
	var decoded APIResponse
	if err := json.Unmarshal(response, &decoded); err != nil {
		return decoded, err
	}
	return decoded, nil
}

func put(url string, body []byte) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}

	request.ContentLength = int64(len(body))
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}
