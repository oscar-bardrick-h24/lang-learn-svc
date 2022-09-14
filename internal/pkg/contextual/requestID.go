package contextual

import "context"

type requestID struct{}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestID{}, id)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(requestID{}).(string)
}
