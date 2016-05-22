package model

type ApiResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewErrorResp(code int, msg string) *ApiResponse {
	resp := &ApiResponse{}
	resp.Code = code
	resp.Msg = msg
	resp.Data = make(map[string]interface{})

	return resp
}

func NewSuccResp(data interface{}) *ApiResponse {
	resp := &ApiResponse{}
	resp.Code = 0
	resp.Msg = ""
	if data == nil {
		resp.Data = make(map[string]interface{})
	} else {
		resp.Data = data
	}

	return resp
}
