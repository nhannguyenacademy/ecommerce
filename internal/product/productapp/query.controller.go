package productapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/query"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/order"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
)

func (a *app) queryController(c *gin.Context) {
	ctx := c.Request.Context()
	qp := parseQueryParams(c.Request)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		respond.Error(c, a.log, errs.NewFieldsError("page", err))
		return
	}

	filter, err := parseFilter(qp)
	if err != nil {
		respond.Error(c, a.log, errs.NewFieldsError("filter", err))
		return
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		respond.Error(c, a.log, errs.NewFieldsError("order", err))
		return
	}

	prds, err := a.productBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query: %s", err))
		return
	}

	total, err := a.productBus.Count(ctx, filter)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "count: %s", err))
		return
	}

	respond.Success(c, a.log, query.NewResult(toAppProducts(prds), total, page))
}
