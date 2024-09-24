// Package productapp maintains the app layer api for the product domain.
package productapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/domain/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/query"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdk/sdkbus/sort"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

type app struct {
	log        *logger.Logger
	auth       *auth.Auth
	productBus *productbus.Business
}

func New(
	log *logger.Logger,
	auth *auth.Auth,
	productBus *productbus.Business,
) *app {
	return &app{
		log:        log,
		auth:       auth,
		productBus: productBus,
	}
}

func (a *app) createHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req newProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	newProduct, err := toBusNewProduct(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	prod, err := a.productBus.Create(ctx, newProduct)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "create: req[%+v]: %s", req, err))
		return
	}

	respond.Success(c, a.log, toAppProduct(prod))
}

func (a *app) updateHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req updateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	updateProduct, err := toBusUpdateProduct(req)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid productID: %s", err))
		return
	}

	prd, err := a.productBus.QueryByID(ctx, productID)
	if err != nil {
		if errors.Is(err, productbus.ErrNotFound) {
			respond.Error(c, a.log, errs.Newf(errs.NotFound, "update: productID[%s] req[%+v]: %s", productID, req, err))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "update: productID[%s] req[%+v]: %s", productID, req, err))
		}
		return
	}

	updatedProduct, err := a.productBus.Update(ctx, prd, updateProduct)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "update: productID[%s] req[%+v]: %s", productID, req, err))
		return
	}

	respond.Success(c, a.log, toAppProduct(updatedProduct))
}

func (a *app) queryHandler(c *gin.Context) {
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

	products, err := a.productBus.Query(ctx, filter, sortBy, page)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query: %s", err))
		return
	}

	total, err := a.productBus.Count(ctx, filter)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "count: %s", err))
		return
	}

	respond.Success(c, a.log, query.NewResult(toAppProducts(products), total, page))
}

func (a *app) deleteHandler(c *gin.Context) {
	ctx := c.Request.Context()

	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid productID: %s", err))
		return
	}

	prd, err := a.productBus.QueryByID(ctx, productID)
	if err != nil {
		if errors.Is(err, productbus.ErrNotFound) {
			respond.Error(c, a.log, errs.Newf(errs.NotFound, "delete: productID[%s]: %s", productID, err))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "delete: productID[%s]: %s", productID, err))
		}
		return
	}

	if err := a.productBus.Delete(ctx, prd); err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "delete: productID[%s]: %s", productID, err))
		return
	}

	respond.Success(c, a.log, nil)
}

func (a *app) queryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	productID, err := uuid.Parse(c.Param("product_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid productID: %s", err))
		return
	}

	prd, err := a.productBus.QueryByID(ctx, productID)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "querybyid: %s", err))
		return
	}

	respond.Success(c, a.log, toAppProduct(prd))
}
