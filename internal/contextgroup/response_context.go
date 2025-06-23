package contextgroup

import "context"

type ResponseContext struct {
	MediaSent bool
}

type responseKeyType string

var responseKey responseKeyType = "response_key"

func SetResponseContext(ctx context.Context, rc *ResponseContext) context.Context {
	return context.WithValue(ctx, responseKey, rc)
}

func GetResponseContext(ctx context.Context) *ResponseContext {
	rc, ok := ctx.Value(responseKey).(*ResponseContext)
	if !ok || rc == nil {
		return nil
	}
	return rc
}

func MarkMediaSent(ctx context.Context) {
	if rc, ok := ctx.Value(responseKey).(*ResponseContext); ok && rc != nil {
		rc.MediaSent = true
	}
}
