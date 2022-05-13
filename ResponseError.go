package common

type ResponseError struct {
	StatusCode int
	Body       string
	Info       string
}

func (r ResponseError) Error() string {
	return r.Info
}
