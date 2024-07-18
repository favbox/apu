//go:build wireinject

package main

import (
	"apu/kitex_gen/payment"
	service "apu/payment"
	"apu/payment/mysql"
	"github.com/google/wire"
)

func initHandler() payment.PaymentSvc {
	panic(wire.Build(
		mysql.ProviderSet,
		service.ProviderSet,
	))
}
