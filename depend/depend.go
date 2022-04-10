package depend

import (
	"context"
)

type Depend interface {
	Init(ctx context.Context) error
}
