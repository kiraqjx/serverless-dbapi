package actuator

import (
	"database/sql"
	"fmt"
	"serverless-dbapi/pkg/entity"
	"serverless-dbapi/pkg/exception"
	"serverless-dbapi/pkg/managercenter"
	"serverless-dbapi/pkg/tool"
	"serverless-dbapi/pkg/valueobject"
)

const API_ID_FROM_QUERY = "apiId"

type Handler struct {
	dbConn        *sql.DB
	managerCenter managercenter.ManagerCenterServer
}

func NewHandle(dbConn *sql.DB, managerCenter managercenter.ManagerCenterServer) Handler {
	return Handler{
		dbConn:        dbConn,
		managerCenter: managerCenter,
	}
}

// common handler
func (h *Handler) Handler(params *valueobject.Params) tool.Result[any] {
	apiIds := params.QueryParams[API_ID_FROM_QUERY]
	if len(apiIds) != 1 {
		return tool.ErrorResult(exception.API_ID_IS_REQUIRE)
	}
	apiId := apiIds[0]
	apiConfig := h.managerCenter.GetApiConfigByApiId(apiId)
	return h.exec(apiConfig, params)
}

func (h *Handler) exec(apiConfig entity.ApiConfig, params *valueobject.Params) tool.Result[any] {
	// create args by list order
	args := make([]any, len(apiConfig.ParamKey))
	for index, value := range apiConfig.ParamKey {
		if value, ok := params.Body[value]; ok {
			args[index] = value
		} else {
			return tool.ErrorResult(exception.REQUIRE_PARAM)
		}
	}

	// exec sql
	rows, err := h.dbConn.Query(apiConfig.Sql, args...)
	if err != nil {
		fmt.Println(err)
		return tool.SimpleErrorResult(500, err.Error())
	}
	defer rows.Close()

	// data -> map[string]any
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]any, columnLength)
	for index := range cache {
		var a any
		cache[index] = &a
	}
	var list []map[string]any
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]any)
		for i, data := range cache {
			item[columns[i]] = data
		}
		list = append(list, item)
	}

	return tool.SuccessResult(list)
}