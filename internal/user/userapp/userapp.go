// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"context"
	"errors"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/query"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
)

// App manages the set of app layer api functions for the user domain.
type App struct {
	userBus *userbus.Business
}

// NewApp constructs a user app API for use.
func NewApp(userBus *userbus.Business) *App {
	return &App{
		userBus: userBus,
	}
}

// query returns a list of users with paging.
func (a *App) query(ctx context.Context, qp QueryParams) (query.Result[User], error) {
	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return query.Result[User]{}, errs.NewFieldsError("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return query.Result[User]{}, err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return query.Result[User]{}, errs.NewFieldsError("order", err)
	}

	usrs, err := a.userBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return query.Result[User]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.userBus.Count(ctx, filter)
	if err != nil {
		return query.Result[User]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppUsers(usrs), total, page), nil
}

// create adds a new user to the system.
func (a *App) create(ctx context.Context, app NewUser) (User, error) {
	nc, err := toBusNewUser(app)
	if err != nil {
		return User{}, errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return User{}, errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return User{}, errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr), nil
}

// queryByID returns a user by its Ia.
func (a *App) queryByID(ctx context.Context) (User, error) {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppUser(usr), nil
}
