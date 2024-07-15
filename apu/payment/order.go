package payment

import (
	"context"

	"apu/kitex_gen/payment"
)

// Order 定义订单领域类型。
type Order struct {
	ID int

	MerchantID string
	Channel    string
	PayWay     string

	OutOrderNo  string
	TotalAmount uint64
	Body        string
	OrderStatus int8

	AuthCode string

	WxAppid   string
	SubOpenid string

	JumpURL   string
	NotifyURL string

	ClientIP        string
	Attach          string
	OrderExpiration string
	ExtendParams    string
}

// Reader 定义订单的读数据接口。
type Reader interface {
	GetByOutOrderNo(ctx context.Context, outOrderNo string) (*Order, error)
}

// Writer 定义订单的写数据接口。
type Writer interface {
	Create(ctx context.Context, order *Order) error
	UpdateOrderStatus(ctx context.Context, outOrderNo string, orderStatus int8) error
}

// Repository 定义订单的数据持久化接口。
type Repository interface {
	Reader
	Writer
}

// UseCase 定义订单的用例接口。
type UseCase interface {
	payment.PaymentSvc
}
