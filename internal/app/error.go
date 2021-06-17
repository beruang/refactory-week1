package app

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

type errorCode int

const (
	InternalCode errorCode = iota + 1
	DuplicateCode
	UnauthenticatedCode
	UnauthorizedCode
	NotFoundCode
	InvalidCode
	BadRequestCode
)

func (e errorCode) Int() int {
	return int(e)
}

func (e errorCode) String() string {
	return [...]string{"Internal Server Error", "Data already exists",
		"Unauthenticated", "Unauthorized", "Data Not Found", "Invalid data", "Bad Request"}[e-1]
}

var (
	InternalError = Error{
		Code:    InternalCode.Int(),
		Message: InternalCode.String(),
	}
	DuplicateError = Error{
		Code:    DuplicateCode.Int(),
		Message: DuplicateCode.String(),
	}
	UnauthenticateError = Error{
		Code:    UnauthenticatedCode.Int(),
		Message: UnauthenticatedCode.String(),
	}
	UnauthorizedError = Error{
		Code:    UnauthorizedCode.Int(),
		Message: UnauthorizedCode.String(),
	}
	NotFoundError = Error{
		Code:    NotFoundCode.Int(),
		Message: NotFoundCode.String(),
	}
	InvalidCodeError = Error{
		Code:    InvalidCode.Int(),
		Message: InvalidCode.String(),
	}
	BadRequestError = Error{
		Code:    BadRequestCode.Int(),
		Message: BadRequestCode.String(),
	}
)
