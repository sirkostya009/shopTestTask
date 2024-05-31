package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (db *DB) PersistOrder(ctx context.Context, tx pgx.Tx, order *Order) error {
	err := tx.QueryRow(ctx, `
		insert into orders (vendee_id, total)
		values ($1, (select sum(price * quantity)
		             from json_to_recordset($2) as goods(product_id int, quantity int)
		             inner join products on products.id = product_id
		             where quantity > 0))
		returning id, order_date, total, status
	`, order.VendeeID, order.Products).Scan(&order.ID, &order.OrderDate, &order.Total, &order.Status)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		insert into order_products (order_id, product_id, quantity)
		select $1, product_id, quantity
		from json_to_recordset($2) as goods(product_id int, quantity int)
		where quantity > 0
	`, order.ID, order.Products)
	return err
}

func (db *DB) OrderExists(ctx context.Context, id int) (exists bool) {
	err := db.QueryRow(ctx, `select exists(select 1 from orders where id = $1)`, id).Scan(&exists)
	return err == nil && exists
}

func (db *DB) GetOrder(ctx context.Context, id int) (*Order, error) {
	var order Order
	err := db.QueryRow(ctx, `select id, vendee_id, order_date, total, status from orders where id = $1`, id).
		Scan(&order.ID, &order.VendeeID, &order.OrderDate, &order.Total, &order.Status)
	if err != nil {
		return nil, err
	}

	order.Products, err = db.GetOrderProducts(ctx, id)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (db *DB) GetOrderProducts(ctx context.Context, id int) ([]OrderedProduct, error) {
	rows, err := db.Query(ctx, `select product_id, quantity from order_products where order_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []OrderedProduct
	for rows.Next() {
		var product OrderedProduct
		err = rows.Scan(&product.ProductID, &product.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (db *DB) ChangeOrderStatus(ctx context.Context, tx pgx.Tx, id int, status string) error {
	_, err := tx.Exec(ctx, `update orders set status = $1 where id = $2`, status, id)
	return err
}

func (db *DB) CheckOrderStatus(ctx context.Context, id int, status string) bool {
	var currentStatus string
	err := db.QueryRow(ctx, `select status from orders where id = $1`, id).Scan(&currentStatus)
	return err == nil && status == currentStatus
}

func (db *DB) UpdateOrderProducts(ctx context.Context, tx pgx.Tx, id int, order *Order) error {
	_, err := tx.Exec(ctx, `delete from order_products where order_id = $1`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		insert into order_products (order_id, product_id, quantity)
		select $1, product_id, quantity
		from json_to_recordset($2) as goods(product_id int, quantity int)
		where quantity > 0
	`, order.ID, order.Products)
	if err != nil {
		return err
	}

	return nil
}
