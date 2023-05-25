package onlinesim

import "fmt"

type ErrorUnexpectedResponse struct {
	Response any
}

func (e ErrorUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected response: %v", e.Response)
}
