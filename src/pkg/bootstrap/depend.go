package bootstrap

import (
	"context"
	"errors"
	"fmt"
)

type Component interface {
	Init(ctx context.Context) error
}

type EmptyComponent struct {
	error bool
}

func (d *EmptyComponent) Init(ctx context.Context) error {
	fmt.Println("EmptyComponent Init")
	if d.error {
		return errors.New("init failed")
	}
	return nil
}

// 检查接口是否实现
var _ Component = (*EmptyComponent)(nil)
