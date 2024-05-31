package db

import "time"

type OrderedProduct struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

const (
	PendingOrder   = "pending"
	CancelledOrder = "cancelled"
	CompletedOrder = "completed"
)

type Order struct {
	ID        int              `json:"id"`
	VendeeID  int              `json:"vendee_id"`
	OrderDate time.Time        `json:"order_date"`
	Total     float64          `json:"total"`
	Status    string           `json:"status"`
	Products  []OrderedProduct `json:"products"`
}

func (o *Order) Validate() bool {
	if o.Products != nil {
		validProducts := 0
		for _, product := range o.Products {
			if product.ProductID > 0 || product.Quantity > 0 {
				validProducts++
			}
		}
		if validProducts == 0 {
			return false
		}
	}

	return o.Status == "" || o.Status == CancelledOrder || o.Status == CompletedOrder || o.Status == PendingOrder
}
