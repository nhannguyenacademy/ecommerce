// Package mid is a package for common middleware functions.
package mid

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sqldb"
)

type ctxKey int

const (
	claimKey       ctxKey = 1
	transactionKey ctxKey = 2
	userIDKey      ctxKey = 3
	userKey        ctxKey = 4
	orderKey       ctxKey = 5
)

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}

func setUser(ctx context.Context, usr userbus.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (userbus.User, error) {
	v, ok := ctx.Value(userKey).(userbus.User)
	if !ok {
		return userbus.User{}, errors.New("user not found in context")
	}

	return v, nil
}

func setOrder(ctx context.Context, ord orderbus.Order) context.Context {
	return context.WithValue(ctx, orderKey, ord)
}

func GetOrder(ctx context.Context) (orderbus.Order, error) {
	v, ok := ctx.Value(orderKey).(orderbus.Order)
	if !ok {
		return orderbus.Order{}, errors.New("order not found in context")
	}

	return v, nil
}

func setTran(ctx context.Context, tx sqldb.CommitRollbacker) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

// GetTran retrieves the value that can manage a transaction.
func GetTran(ctx context.Context) (sqldb.CommitRollbacker, error) {
	v, ok := ctx.Value(transactionKey).(sqldb.CommitRollbacker)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return v, nil
}
