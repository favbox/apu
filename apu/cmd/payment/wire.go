//go:build wireinject
// +build wireinject

package main

import (
	"apu/kitex_gen/payment"
	paymentService "apu/payment"
	"apu/payment/mysql"
	"github.com/google/wire"
)

// 初始化支付服务处理程序，注入依赖项。
func initHandler() payment.PaymentSvc {
	panic(wire.Build(
		mysql.ProviderSet,
		paymentService.ProviderSet,
	))
}
