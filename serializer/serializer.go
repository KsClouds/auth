package serializer

// Response 基础序列化器
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

// Err 通用错误处理
func Err(errCode int, err error) Response {
	res := Response{
		Code: errCode,
		Msg:  err.Error(),
	}
	return res
}

func Success() Response {
	res := Response{
		Code: 0,
		Msg:  "success",
	}
	return res
}

func Suc(data interface{}) Response {
	return Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}
