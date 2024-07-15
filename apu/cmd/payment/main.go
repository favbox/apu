package main

import (
	"apu/kitex_gen/payment/paymentsvc"
	"github.com/cloudwego/kitex/pkg/klog"
)

func main() {
	handler := initHandler()

	svr := paymentsvc.NewServer(handler)
	err := svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
