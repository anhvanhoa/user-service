package serviceRes

type Response interface {
	SetMessage(message string) Response
	SetData(data any) Response
	Code(code int) Response
	GetCode() int
}
