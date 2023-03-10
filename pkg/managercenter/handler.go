package managercenter

import (
	"encoding/json"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/exception"
	"serverless-dbapi/pkg/store"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store store.Store
}

func NewHandler(store store.Store) Handler {
	return Handler{
		store: store,
	}
}

func (h *Handler) SaveDataBase(params *valueobject.Params) tool.Result[any] {
	database := &entity.DatabaseConfig{}
	err := json.Unmarshal(params.Body, &database)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	valid := validator.New()
	if err := valid.Struct(database); err != nil {
		return tool.ErrorResult[any](exception.REQUIRE_PARAM, err.Error())
	}
	id, err := h.store.SaveDataBase(*database)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](id)
}

func (h *Handler) SaveApiGroup(params *valueobject.Params) tool.Result[any] {
	apiGroup := &entity.ApiGroupConfig{}
	err := json.Unmarshal(params.Body, &apiGroup)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	valid := validator.New()
	if err := valid.Struct(apiGroup); err != nil {
		return tool.ErrorResult[any](exception.REQUIRE_PARAM, err.Error())
	}
	id, err := h.store.SaveApiGroup(*apiGroup)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](id)
}

func (h *Handler) SaveApi(params *valueobject.Params) tool.Result[any] {
	apiInfo := &entity.ApiConfig{}
	err := json.Unmarshal(params.Body, &apiInfo)
	if err != nil {
		return tool.ErrorResult[any](exception.PARSE_REQUEST_ERROR)
	}
	valid := validator.New()
	if err := valid.Struct(apiInfo); err != nil {
		return tool.ErrorResult[any](exception.REQUIRE_PARAM, err.Error())
	}
	id, err := h.store.SaveApi(*apiInfo)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](id)
}

func (h *Handler) GetDataBases(params *valueobject.Params) tool.Result[any] {
	result := getPageInfo(params)
	if result.IsError() {
		return tool.ErrorResult[any](result.Err)
	}
	databases, err := h.store.GetDataBases(result.Data)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](databases)
}

func (h *Handler) GetDataBase(params *valueobject.Params) tool.Result[any] {
	database, err := h.store.GetDataBase(params.QueryParams["id"][0])
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](database)
}

func (h *Handler) GetApiGroups(params *valueobject.Params) tool.Result[any] {
	result := getPageInfo(params)
	if result.IsError() {
		return tool.ErrorResult[any](result.Err)
	}
	apiGroups, err := h.store.GetApiGroups(result.Data)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](apiGroups)
}

func (h *Handler) GetApis(params *valueobject.Params) tool.Result[any] {
	result := getPageInfo(params)
	if result.IsError() {
		return tool.ErrorResult[any](result.Err)
	}
	apis, err := h.store.GetApis(params.QueryParams["apiGroupId"][0], result.Data)
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](apis)
}

func (h *Handler) GetApi(params *valueobject.Params) tool.Result[any] {
	api, err := h.store.GetApi(params.QueryParams["apiId"][0])
	if err != nil {
		return tool.SimpleErrorResult[any](500, err.Error())
	}
	return tool.SuccessResult[any](api)
}

func getPageInfo(params *valueobject.Params) tool.Result[valueobject.Cursor] {
	if len(params.QueryParams["continue"]) != 1 || len(params.QueryParams["limit"]) != 1 {
		return tool.ErrorResult[valueobject.Cursor](exception.REQUIRE_PARAM, "continue, limit")
	}

	continueKey := params.QueryParams["continue"][0]
	limit, err := strconv.Atoi(params.QueryParams["limit"][0])
	if err != nil {
		return tool.ErrorResult[valueobject.Cursor](exception.REQUIRE_PARAM, "limit")
	}

	pageInfo := valueobject.Cursor{
		Continue: continueKey,
		Limit:    limit,
	}

	return tool.SuccessResult(pageInfo)
}
