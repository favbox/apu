package mysql

import (
	"context"
	"fmt"
	"log"

	"apu/payment"
	"apu/payment/mysql/ent"
	Order "apu/payment/mysql/ent/order"
	"github.com/google/wire"

	_ "github.com/go-sql-driver/mysql"
)

var ProviderSet = wire.NewSet(NewEntClient, NewRepository)

var _ payment.Repository = (*Repository)(nil)

// Repository 定义MySQL实现的订单仓库。
type Repository struct {
	db *ent.Client
}

// GetByOutOrderNo 实现根据外部订单号获取订单。
func (r *Repository) GetByOutOrderNo(ctx context.Context, outOrderNo string) (*payment.Order, error) {
	row, err := r.db.Order.Query().
		Where(Order.OutOrderNo(outOrderNo)).
		First(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &payment.Order{
		ID:              row.ID,
		MerchantID:      row.MerchantID,
		Channel:         row.Channel,
		PayWay:          row.PayWay,
		OutOrderNo:      row.OutOrderNo,
		TotalAmount:     row.TotalAmount,
		Body:            row.Body,
		OrderStatus:     row.OrderStatus,
		AuthCode:        row.AuthCode,
		WxAppid:         row.WxAppid,
		SubOpenid:       row.SubOpenid,
		JumpURL:         row.JumpURL,
		NotifyURL:       row.NotifyURL,
		ClientIP:        row.ClientIP,
		Attach:          row.Attach,
		OrderExpiration: row.OrderExpiration,
		ExtendParams:    row.ExtendParams,
	}, nil
}

// Create 实现订单的创建。
func (r *Repository) Create(ctx context.Context, order *payment.Order) error {
	ret, err := r.db.Order.Create().
		SetMerchantID(order.MerchantID).
		SetChannel(order.Channel).
		SetPayWay(order.PayWay).
		SetOutOrderNo(order.OutOrderNo).
		SetTotalAmount(order.TotalAmount).
		SetBody(order.Body).
		SetOrderStatus(order.OrderStatus).
		SetAuthCode(order.AuthCode).
		SetWxAppid(order.WxAppid).
		SetSubOpenid(order.SubOpenid).
		SetJumpURL(order.JumpURL).
		SetNotifyURL(order.NotifyURL).
		SetClientIP(order.ClientIP).
		SetAttach(order.Attach).
		SetOrderExpiration(order.OrderExpiration).
		SetExtendParams(order.ExtendParams).
		Save(ctx)
	if err != nil {
		return err
	}
	order.ID = ret.ID
	return nil
}

// UpdateOrderStatus 实现订单状态的更新。
func (r *Repository) UpdateOrderStatus(ctx context.Context, outOrderNo string, orderStatus int8) error {
	return r.db.Order.Update().Where(Order.OutOrderNo(outOrderNo)).SetOrderStatus(orderStatus).Exec(ctx)
}

// NewEntClient 使用默认配置创建一个 ent 客户端。
func NewEntClient() *ent.Client {
	entClient, err := ent.Open(
		"mysql",
		fmt.Sprintf("root:asdfasdf@tcp(127.0.0.1:3306)/apu?charset=utf8mb4&parseTime=true"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = entClient.Schema.Create(context.TODO()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return entClient
}

// NewRepository 创建一个新的订单仓库。
// 这是 MySQL 具体实现。
func NewRepository(dbClient *ent.Client) payment.Repository {
	return &Repository{
		db: dbClient,
	}
}
