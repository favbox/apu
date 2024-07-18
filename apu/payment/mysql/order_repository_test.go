package mysql

import (
	"context"
	"fmt"
	"testing"

	"apu/payment"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_Create(t *testing.T) {
	r := NewRepository(NewEntClient())
	o := &payment.Order{
		TotalAmount: 100,
	}
	err := r.Create(context.Background(), o)
	assert.Nil(t, err)
	fmt.Println(o.ID)
}

func TestOrderRepository_GetByOutOrderNo(t *testing.T) {
	r := NewRepository(NewEntClient())
	assert.NotNil(t, r)
}
