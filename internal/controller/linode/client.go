package linode

import (
	"context"

	"github.com/linode/linodego/v2"
)

type clientContextKey struct{}

// NewContext returns a context with the linodego client attached.
func NewContext(ctx context.Context, client *linodego.Client) context.Context {
	return context.WithValue(ctx, clientContextKey{}, client)
}

// FromContext returns the linodego client from the context, or nil.
func FromContext(ctx context.Context) *linodego.Client {
	client, _ := ctx.Value(clientContextKey{}).(*linodego.Client)
	return client
}
