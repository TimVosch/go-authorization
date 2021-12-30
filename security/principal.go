package security

import "context"

var principalKey = principalKeyT{}

type principalKeyT struct{}

// Principal ...
type Principal interface {
	Type() PrincipalType
}

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalKey, principal)
}

func GetPrincipal(ctx context.Context) Principal {
	p, ok := ctx.Value(principalKey).(Principal)
	if !ok {
		return nil
	}
	return p
}
