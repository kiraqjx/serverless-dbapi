package tool

type Response[T any] struct {
	Code int
	Msg  string
	Data any
}

const SUCCESS_CODE = 0

func ResultToResponse(result Result[any]) Response[any] {
	if result.IsError() {
		return Response[any]{
			Code: result.Err.Code,
			Msg:  result.Err.Msg,
			Data: nil,
		}
	} else {
		return Response[any]{
			Code: SUCCESS_CODE,
			Msg:  "success",
			Data: result.Data,
		}
	}
}
