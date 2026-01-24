//go:build wasip1

package configurator

import (
	"context"
	"encoding/json"
	"fmt"
)

func Resolve[T any](ctx context.Context) (T, error) {
	var cfg T

	rawConfig, err := NewService().GetConfig(ctx, nil)
	if err != nil {
		return cfg, fmt.Errorf("failed to get config: %w", err)
	}

	err = json.Unmarshal([]byte(rawConfig.Json), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
