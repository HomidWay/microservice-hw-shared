package clients

import "context"

func WithAuth(ctx context.Context, authToken string) context.Context {
	return context.WithValue(ctx, "Authorization", authToken)
}
