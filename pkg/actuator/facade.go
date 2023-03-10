package actuator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"serverless-dbapi/pkg/exception"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"

	edclient "github.com/kiraqjx/ed-client"
)

// show the actuator server api
type ActuatorServer interface {
	ApiActuator(params *valueobject.Params) tool.Result[any]
}

// impl by http client
// TODO server LB
type HttpActuatorServer struct {
	Lb     *edclient.Lb
	client *http.Client
}

func NewHttpActuatorServer() ActuatorServer {
	return &HttpActuatorServer{
		client: &http.Client{},
	}
}

func (h *HttpActuatorServer) ApiActuator(params *valueobject.Params) tool.Result[any] {
	node := h.Lb.Lb()
	url, err := tool.BuildURL(node.Server, params.QueryParams)
	if err != nil {
		return tool.ErrorResult[any](exception.BUILD_URL_ERROR)
	}
	bodyBytes, err := json.Marshal(params.Body)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	resp, err := h.client.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return tool.ErrorResult[any](exception.RPC_ERROR)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tool.ErrorResult[any](exception.RPC_RESPONSE_PARSE_ERROR)
	}
	data := &tool.Response[map[string]any]{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return tool.ErrorResult[any](exception.RPC_RESPONSE_PARSE_ERROR)
	}
	return tool.SuccessResult[any](data)
}
