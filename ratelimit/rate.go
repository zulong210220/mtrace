package ratelimit

import (
	"fmt"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"go.uber.org/ratelimit"

	"context"
)

type clientWrapper struct {
	r ratelimit.Limiter
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	c.r.Take()
	return c.Client.Call(ctx, req, rsp, opts...)
}

// NewClientWrapper creates a blocking side rate limiter
func NewClientWrapper(rate int, opts ...ratelimit.Option) client.Wrapper {
	r := ratelimit.New(rate, opts...)

	return func(c client.Client) client.Client {
		return &clientWrapper{r, c}
	}
}

var (
	count int64
)

// NewHandlerWrapper creates a blocking server side rate limiter
func NewHandlerWrapper(rate int, opts ...ratelimit.Option) server.HandlerWrapper {
	r := ratelimit.New(rate, opts...)

	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			method := req.Method()
			if method == "StudentService.GetStudent" {
				r.Take()
			}
			//			r.Take()
			count++
			fmt.Println("------------|||||count:", count)
			return h(ctx, req, rsp)
		}
	}
}

/* vim: set tabstop=4 set shiftwidth=4 */
