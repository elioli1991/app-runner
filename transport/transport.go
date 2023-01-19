package transport

import (
	"context"
	"net/url"
)

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type EndPointer interface {
	EndPoint() (*url.URL, error)
}
