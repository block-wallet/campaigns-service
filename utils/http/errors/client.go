package errors

import "fmt"

type Client struct {
	code    int
	message string
}

func NewClient(code int, message string) *Client {
	return &Client{
		code:    code,
		message: message,
	}
}

func (c *Client) Error() string {
	return fmt.Sprintf("HTTP client error with code: %d and message: %s", c.code, c.message)
}
