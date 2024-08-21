package productapp

import (
	"github.com/gin-gonic/gin"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/query"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
)

func (a *app) queryController(c *gin.Context) {
	ctx := c.Request.Context()
	qp := parseQueryParams(c.Request)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	filter, err := parseFilter(qp)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	sortBy, err := sort.Parse(sortByFields, qp.SortBy, defaultSortBy)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	output, err := a.productBus.Query(ctx, filter, sortBy, page)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query: %s", err))
		return
	}

	total, err := a.productBus.Count(ctx, filter)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "count: %s", err))
		return
	}

	respond.Success(c, a.log, query.NewResult(toAppProducts(output), total, page))
}
