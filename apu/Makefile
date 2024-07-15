MODULE = apu

# 通过IDL生成客户端
.PHONY: gen
gen:
	kitex -module ${MODULE} idl/payment.thrift
	kitex -module ${MODULE} idl/common.thrift