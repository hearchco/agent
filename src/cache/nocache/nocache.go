package nocache

import (
	"time"
)

type DRV struct{}

func New() (DRV, error) { return DRV{}, nil }

func (drv DRV) Close() {}

func (drv DRV) Set(k string, v any, ttl ...time.Duration) error { return nil }

func (drv DRV) Get(k string, o any) error { return nil }

func (drv DRV) GetTTL(k string) (time.Duration, error) { return 0, nil }
