// Package orderapp maintains the app layer.
package orderapp

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhannguyenacademy/ecommerce/internal/order/orderbus"
	"github.com/nhannguyenacademy/ecommerce/internal/product/productbus"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/auth"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/errs"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/mid"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/query"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkapp/respond"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/page"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sort"
	"github.com/nhannguyenacademy/ecommerce/internal/sdkbus/sqldb"
	"github.com/nhannguyenacademy/ecommerce/internal/user/userbus"
	"github.com/nhannguyenacademy/ecommerce/pkg/logger"
)

type app struct {
	log        *logger.Logger
	auth       *auth.Auth
	dbBeginner sqldb.Beginner
	orderBus   *orderbus.Business
	productBus *productbus.Business
	userBus    *userbus.Business
}

func New(
	log *logger.Logger,
	auth *auth.Auth,
	dbBeginner sqldb.Beginner,
	orderBus *orderbus.Business,
	productBus *productbus.Business,
	userBus *userbus.Business,
) *app {
	return &app{
		log:        log,
		auth:       auth,
		dbBeginner: dbBeginner,
		orderBus:   orderBus,
		productBus: productBus,
		userBus:    userBus,
	}
}

// newWithTx constructs a new app value using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	orderBusTx, err := a.orderBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	productBusTx, err := a.productBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		log:        a.log,
		auth:       a.auth,
		orderBus:   orderBusTx,
		productBus: productBusTx,
		userBus:    a.userBus,
	}

	return &app, nil
}

func (a *app) createHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// construct a new app value using a store transaction
	a, err := a.newWithTx(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.Internal, err))
		return
	}

	var req newOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	itemQuantityMap := make(map[uuid.UUID]int32)
	productIDs := make([]uuid.UUID, 0, len(req.Items))
	for _, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid product id: %s", err))
			return
		}
		productIDs = append(productIDs, productID)
		itemQuantityMap[productID] = item.Quantity
	}

	products, err := a.productBus.QueryByIDs(ctx, productIDs)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query products by ids: %s", err))
		return
	}
	if len(products) != len(productIDs) {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "missing products: %v", productIDs))
		return
	}

	// for mapping product data to order items
	productsMap := make(map[uuid.UUID]productbus.Product)

	// validate and update product quantity
	for _, product := range products {
		if product.Quantity < itemQuantityMap[product.ID] {
			respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "insufficient quantity: %s", product.ID))
			return
		}
		productsMap[product.ID] = product

		newQuantity := product.Quantity - itemQuantityMap[product.ID]
		_, err = a.productBus.Update(ctx, product, productbus.UpdateProduct{
			Quantity: &newQuantity,
		})
		if err != nil {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "update product: id[%s]: %s", product.ID, err))
			return
		}
	}

	// create new order
	newOrder, err := toBusNewOrder(req, productsMap)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	order, err := a.orderBus.Create(ctx, newOrder)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "create: req[%+v]: %s", req, err))
		return
	}

	respond.Success(c, a.log, toAppOrder(order))
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

	orders, err := a.orderBus.Query(ctx, filter, sortBy, page)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query: %s", err))
		return
	}

	total, err := a.orderBus.Count(ctx, filter)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "count: %s", err))
		return
	}

	respond.Success(c, a.log, query.NewResult(toAppOrders(orders), total, page))
}

func (a *app) queryUserOrdersHandler(c *gin.Context) {
	ctx := c.Request.Context()
	qp := parseQueryParams(c.Request)

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid user id: %s", err))
		return
	}

	authenUserID, err := mid.GetUserID(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.Internal, err))
		return
	}

	if userID != authenUserID {
		respond.Error(c, a.log, errs.New(errs.PermissionDenied, errors.New("user id mismatch")))
		return
	}

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
	filter.UserID = &authenUserID

	sortBy, err := sort.Parse(sortByFields, qp.SortBy, defaultSortBy)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.InvalidArgument, err))
		return
	}

	orders, err := a.orderBus.Query(ctx, filter, sortBy, page)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query: %s", err))
		return
	}

	total, err := a.orderBus.Count(ctx, filter)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "count: %s", err))
		return
	}

	respond.Success(c, a.log, query.NewResult(toAppOrders(orders), total, page))
}

func (a *app) queryByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid id: %s", err))
		return
	}

	orderWithItems, err := a.orderBus.QueryByIDWithItems(ctx, id)
	if err != nil {
		if errors.Is(err, orderbus.ErrNotFound) {
			respond.Error(c, a.log, errs.Newf(errs.NotFound, "order id[%s] not found", id))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "query order: id[%s]: %s", id, err))
		}
		return
	}

	user, err := a.userBus.QueryByID(ctx, orderWithItems.UserID)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query user: id[%s]: %s", orderWithItems.UserID, err))
		return
	}

	respond.Success(c, a.log, toAppOrderDetail(orderWithItems, user))
}

func (a *app) updateStatusHandler(c *gin.Context) {
	ctx := c.Request.Context()

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid orderID: %s", err))
		return
	}

	var req updateOrderStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, a.log, err)
		return
	}

	status, err := orderbus.ParseStatus(req.Status)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid status: %s", err))
		return
	}

	ord, err := a.orderBus.QueryByIDWithItems(ctx, orderID)
	if err != nil {
		if errors.Is(err, orderbus.ErrNotFound) {
			respond.Error(c, a.log, errs.Newf(errs.NotFound, "order orderID[%s] not found", orderID))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "query order: orderID[%s]: %s", orderID, err))
		}
		return
	}

	updatedOrder, err := a.orderBus.UpdateStatus(ctx, ord.Order, status)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "update order status: orderID[%s]: %s", orderID, err))
		return
	}

	if updatedOrder.Status.Equal(orderbus.Statuses.Cancelled) {
		// todo: revert product quantity (here or move to bus?)
	}

	respond.Success(c, a.log, toAppOrder(updatedOrder))
}

func (a *app) cancelHandler(c *gin.Context) {
	ctx := c.Request.Context()

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid orderID: %s", err))
		return
	}

	ord, err := mid.GetOrder(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.Internal, err))
		return
	}

	updatedOrder, err := a.orderBus.UpdateStatus(ctx, ord, orderbus.Statuses.Cancelled)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "cancel order: orderID[%s]: %s", orderID, err))
		return
	}

	// todo: revert product quantity (here or move to bus?)
	items, err := a.orderBus.QueryOrderItems(ctx, updatedOrder)
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "query order items: orderID[%s]: %s", orderID, err))
		return
	}
	_ = items

	respond.Success(c, a.log, toAppOrder(updatedOrder))
}

func (a *app) deleteHandler(c *gin.Context) {
	ctx := c.Request.Context()

	a, err := a.newWithTx(ctx)
	if err != nil {
		respond.Error(c, a.log, errs.New(errs.Internal, err))
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		respond.Error(c, a.log, errs.Newf(errs.InvalidArgument, "invalid orderID: %s", err))
		return
	}
	// todo: check if order has any success payments, but orderbus cannot import paymentbus, use delegate instead

	ord, err := a.orderBus.QueryByIDWithItems(ctx, orderID)
	if err != nil {
		if errors.Is(err, orderbus.ErrNotFound) {
			respond.Error(c, a.log, errs.Newf(errs.NotFound, "order orderID[%s] not found", orderID))
		} else {
			respond.Error(c, a.log, errs.Newf(errs.Internal, "query order: orderID[%s]: %s", orderID, err))
		}
		return
	}

	if err := a.orderBus.Delete(ctx, ord.Order); err != nil {
		respond.Error(c, a.log, errs.Newf(errs.Internal, "delete: orderID[%s]: %s", orderID, err))
		return
	}

	// todo: revert product quantity (here or move to bus?)

	respond.Success(c, a.log, nil)
}
