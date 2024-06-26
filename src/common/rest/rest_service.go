package rest

import "encoding/json"

type RestClient struct {
	BaseURL string
}

func NewRestClient(baseURL string) *RestClient {
	return &RestClient{BaseURL: baseURL}
}

func (r *RestClient) Get(path string, args ...any) (*Response, error) {
	queryParams := make(map[string]string)
	headers := make(map[string]string)

	if len(args) > 0 {
		queryParams = args[0].(map[string]string)
	}

	if len(args) > 1 {
		headers = args[1].(map[string]string)
	}

	request := Request{
		Method:      Get,
		BaseURL:     r.BaseURL + path,
		Headers:     headers,
		QueryParams: queryParams,
	}
	return Send(request)
}

func (r *RestClient) Post(path string, args ...any) (*Response, error) {
	body := []byte{}
	headers := make(map[string]string)

	if len(args) > 0 {
		bytes, err := json.Marshal(args[0])
		if err != nil {
			return nil, err
		}
		body = bytes
	}

	if len(args) > 1 {
		headers = args[1].(map[string]string)
	}

	request := Request{
		Method:  Post,
		BaseURL: r.BaseURL + path,
		Headers: headers,
		Body:    body,
	}
	return Send(request)
}

func (r *RestClient) Put(path string, args ...any) (*Response, error) {
	body := []byte{}
	headers := make(map[string]string)

	if len(args) > 0 {
		bytes, err := json.Marshal(args[0])
		if err != nil {
			return nil, err
		}
		body = bytes
	}

	if len(args) > 1 {
		headers = args[1].(map[string]string)
	}

	request := Request{
		Method:  Put,
		BaseURL: r.BaseURL + path,
		Headers: headers,
		Body:    body,
	}
	return Send(request)
}
