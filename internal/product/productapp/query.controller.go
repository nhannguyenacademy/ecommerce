package productapp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/query"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/response"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
)

func (a *app) queryController(c *gin.Context) {
	results, err := a.query(c.Request.Context(), parseQueryParams(c.Request))
	response.Send(c, a.log, results, err)
}

func (a *app) query(ctx context.Context, qp queryParams) (query.Result[user], error) {
	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return query.Result[user]{}, errs.NewFieldsError("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return query.Result[user]{}, err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return query.Result[user]{}, errs.NewFieldsError("order", err)
	}

	usrs, err := a.userBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return query.Result[user]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.userBus.Count(ctx, filter)
	if err != nil {
		return query.Result[user]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppUsers(usrs), total, page), nil
}
