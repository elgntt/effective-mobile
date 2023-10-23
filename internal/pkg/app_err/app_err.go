package app_err

type BusinessError struct {
	code    string
	message string
}

func (b BusinessError) Error() string {
	return b.message
}

func (b BusinessError) Code() string {
	return b.code
}

func NewBusinessError(message string) error {
	return BusinessError{
		code:    "BusinessError",
		message: message,
	}
}
