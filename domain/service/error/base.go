package serviceError

type ErrorApp interface {
	Error() string
	SetMessage(message string) ErrorApp
	SetData(data any) ErrorApp
	Code(code int) ErrorApp
	GetCode() int
	BadReq() ErrorApp
	UnprocessableEntity() ErrorApp
	NotFound() ErrorApp
	Unauthorized() ErrorApp
	Forbidden() ErrorApp
	Conflict() ErrorApp
	InternalServerError() ErrorApp
}

type errorApp struct {
	message string
}

func NewErrorApp(message string) *errorApp {
	return &errorApp{
		message: message,
	}
}

func (e *errorApp) Error() string {
	return e.message
}
