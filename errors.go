package nordigen

import "fmt"

type APIError struct {
	StatusCode int
	Body       string
	Err        error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError %v %v", e.StatusCode, e.Body)
}

func (e *APIError) Unwrap() error {
	return e.Err
}
