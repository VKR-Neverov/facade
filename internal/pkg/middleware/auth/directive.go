package auth

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/fidesy-pay/facade/internal/pkg/model"
)

func Directive(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	_, ok := ctx.Value(SessionCtxName).(Session)
	if !ok {
		return nil, errors.New(model.ErrorNoAuth.String())
	}

	return next(ctx)
}
