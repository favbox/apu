package payment

import (
	"context"

	"apu/kitex_gen/common"
	"apu/kitex_gen/payment"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewService)

var _ UseCase = (*Service)(nil)

type Service struct {
	repo Repository
}

// UnifyPay 实现统一支付服务。
func (s *Service) UnifyPay(ctx context.Context, req *payment.UnifyPayReq) (r *payment.UnifyPayResp, err error) {
	o := &Order{
		PayWay:     req.PayWay,
		SubOpenid:  req.SubOpenId,
		MerchantID: req.MerchantId,
		OutOrderNo: req.OutOrderNo,
		Channel:    "1",
	}

	err = s.repo.Create(ctx, o)
	if err != nil {
		return nil, err
	}

	return &payment.UnifyPayResp{
		MerchantId:    o.MerchantID,
		SubMerchantId: o.SubOpenid,
		OutOrderNo:    o.OutOrderNo,
		JspayInfo:     "xxxxx",
		PayWay:        o.PayWay,
	}, nil
}

func (s *Service) uniqueOutOrderNo(ctx context.Context, outOrderNo string) (bool, error) {
	return true, nil
}

// QRPay 实现扫码支付服务。
func (s *Service) QRPay(ctx context.Context, req *payment.QRPayReq) (r *payment.QRPayResp, err error) {
	onlyOne, err := s.uniqueOutOrderNo(ctx, req.OutOrderNo)
	if err != nil {
		klog.CtxErrorf(ctx, "err: %v", err)
		return nil, err
	}
	if !onlyOne {
		return nil, kerrors.NewBizStatusError(int32(common.Err_DuplicateOutOrderNo), common.Err_DuplicateOutOrderNo.String())
	}

	o := &Order{
		PayWay:      req.OutOrderNo,
		TotalAmount: uint64(req.TotalAmount),
		MerchantID:  req.MerchantId,
		OutOrderNo:  req.OutOrderNo,
		Channel:     "1",
	}
	err = s.repo.Create(ctx, o)
	if err != nil {
		return nil, err
	}
	return &payment.QRPayResp{
		MerchantId: o.MerchantID,
		SubOpenid:  o.SubOpenid,
		OutOrderNo: o.OutOrderNo,
		PayWay:     o.PayWay,
	}, nil
}

// QueryOrder 实现订单查询服务。
func (s *Service) QueryOrder(ctx context.Context, req *payment.QueryOrderReq) (r *payment.QueryOrderResp, err error) {
	order, err := s.repo.GetByOutOrderNo(ctx, req.OutOrderNo)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return &payment.QueryOrderResp{
		OrderStatus: order.OrderStatus,
	}, nil
}

// CloseOrder 实现关闭订单服务。
func (s *Service) CloseOrder(ctx context.Context, req *payment.CloseOrderReq) (r *payment.CloseOrderResp, err error) {
	if err = s.repo.UpdateOrderStatus(ctx, req.OutOrderNo, 9); err != nil {
		return nil, err
	}
	return &payment.CloseOrderResp{}, nil
}

func NewService(r Repository) payment.PaymentSvc {
	return &Service{
		repo: r,
	}
}
