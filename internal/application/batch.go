package application

import (
	"context"
)

type BatchProcessor interface {
	Process(ctx context.Context, input Input)
}
