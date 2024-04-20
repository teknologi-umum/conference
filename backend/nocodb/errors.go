package nocodb

type BadRequestError struct {
	Message string `json:"msg"`
}

func (b BadRequestError) Error() string {
	return b.Message
}
