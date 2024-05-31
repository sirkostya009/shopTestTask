package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"shopTestTask/db"
	"strconv"
)

// PlaceOrder POST /order
func (h *Handlers) PlaceOrder(ctx *Context) error {
	var order db.Order
	err := json.NewDecoder(ctx.Body).Decode(&order)
	if err != nil {
		return &httpErr{http.StatusBadRequest, err.Error()}
	}

	if !order.Validate() {
		return &httpErr{http.StatusBadRequest, "invalid data"}
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = h.db.PersistOrder(ctx, tx, &order)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(order)
}

// QueryOrder GET /order/{id}
func (h *Handlers) QueryOrder(ctx *Context) error {
	id, err := strconv.Atoi(ctx.PathValue("id"))
	if err != nil {
		return &httpErr{http.StatusBadRequest, err.Error()}
	}

	if !h.db.OrderExists(ctx, id) {
		return &httpErr{http.StatusNotFound, "not found"}
	}

	order, err := h.db.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	return ctx.JSON(order)
}

// QueryOrderProducts GET /order/{id}/products
func (h *Handlers) QueryOrderProducts(ctx *Context) error {
	id, err := strconv.Atoi(ctx.PathValue("id"))
	if err != nil {
		return &httpErr{http.StatusBadRequest, err.Error()}
	}

	if !h.db.OrderExists(ctx, id) {
		return &httpErr{http.StatusNotFound, "not found"}
	}

	rows, err := h.db.GetOrderProducts(ctx, id)
	if err != nil {
		return err
	}

	return ctx.JSON(rows)
}

// CancelOrder DELETE /order/{id}
func (h *Handlers) CancelOrder(ctx *Context) error {
	id, err := strconv.Atoi(ctx.PathValue("id"))
	if err != nil {
		return &httpErr{http.StatusBadRequest, err.Error()}
	}

	if !h.db.OrderExists(ctx, id) {
		return &httpErr{http.StatusNotFound, "not found"}
	}

	if h.db.CheckOrderStatus(ctx, id, db.CancelledOrder) {
		return &httpErr{http.StatusBadRequest, "order is not pending"}
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = h.db.ChangeOrderStatus(ctx, tx, id, db.CancelledOrder)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	ctx.WriteHeader(http.StatusNoContent)
	return nil
}

// ChangeOrder PATCH /order/{id}
func (h *Handlers) ChangeOrder(ctx *Context) error {
	id, err := strconv.Atoi(ctx.PathValue("id"))
	if err != nil {
		return &httpErr{http.StatusBadRequest, err.Error()}
	}

	if !h.db.OrderExists(ctx, id) {
		return &httpErr{http.StatusNotFound, "not found"}
	}

	if h.db.CheckOrderStatus(ctx, id, db.CancelledOrder) {
		return &httpErr{http.StatusBadRequest, "order is not pending"}
	}

	var order db.Order
	err = json.NewDecoder(ctx.Body).Decode(&order)
	if err != nil {
		return &httpErr{http.StatusBadRequest, err.Error()}
	}
	order.ID = id

	if !order.Validate() {
		return &httpErr{http.StatusBadRequest, "invalid data"}
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if order.Products != nil {
		err = h.db.UpdateOrderProducts(ctx, tx, id, &order)
		if err != nil {
			return err
		}
	}

	if order.Status != "" {
		err = h.db.ChangeOrderStatus(ctx, tx, id, order.Status)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	o, err := h.db.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	return ctx.JSON(o)
}
