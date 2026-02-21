package bridge

import (
	"context"

	"github.com/Miklakapi/esplogbridge/internal/config"
)

type Bridge struct {
}

func New(c config.Config) (Bridge, error) {
	return Bridge{}, nil
}

func (*Bridge) Run(ctx context.Context) error {
	return nil
}
