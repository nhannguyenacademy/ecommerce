package orderapp

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

func (a *app) query(ctx context.Context, qp queryParams) (query.Result[product], error) {
	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return query.Result[product]{}, errs.NewFieldsError("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return query.Result[product]{}, err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return query.Result[product]{}, errs.NewFieldsError("order", err)
	}

	prds, err := a.productBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return query.Result[product]{}, errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.productBus.Count(ctx, filter)
	if err != nil {
		return query.Result[product]{}, errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppProducts(prds), total, page), nil
}
