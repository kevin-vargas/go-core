package ctx

import "context"

const (
	contextUsername = contextKey("username")
)

type contextKey string

func (c contextKey) String() string {
	return "request_" + string(c)
}

func GetUsername(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(contextUsername).(string)
	return id, ok
}

func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, contextUsername, username)
}
