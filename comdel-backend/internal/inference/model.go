package inference

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const DETECTED int = 1;
const NOT_DETECTED int = 0;
const ERROR = -1;

type ModelAPI struct {
	Result 		int		`json:"result"`;
}

type ModelResponse struct {
	result		int;
	message 	string;
	err 		error;
}

func (m *ModelAPI) Detect(comment string) *ModelResponse {
	var response ModelResponse;
	var query string = url.QueryEscape(comment);
	var endpoint string = fmt.Sprintf("http://model:8000/comment/detect/?comment=%s", query);
	resp, err := http.Get(endpoint);

	if err != nil {
		response = ModelResponse{result: -1, message: "Failed to get the endpoint", err: err}
		return &response;
	}

	defer resp.Body.Close();

	bodyByte, err := io.ReadAll(resp.Body)

	if err != nil {
		response = ModelResponse{result: -1, message: "failed to parse http response", err: err}
		return &response;
	}

	err = json.Unmarshal(bodyByte, &m);

	if err != nil {
		response = ModelResponse{result: -1, message: "Invalid Response type of JSON", err: err};
		return &response;
	}

	response = ModelResponse{result: m.Result, message: "Success detecting comments", err: nil}
	return &response
}

func (mr *ModelResponse) Get() (int, string, error) {
	return mr.result, mr.message, mr.err;
}