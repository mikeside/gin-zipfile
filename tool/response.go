package tool

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	Ok  = response(1, "ok")
	Err = response(0, "err")
)

func response(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func (res *Response) WithMsg(msg string) Response {
	return Response{
		Code: res.Code,
		Msg:  msg,
		Data: res.Data,
	}
}

func (res *Response) WithData(data interface{}) Response {
	return Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
}
