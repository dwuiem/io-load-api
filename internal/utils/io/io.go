package io

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

const (
	minProcessingDuration         = time.Second * 5
	maxProcessingDuration         = time.Second * 30
	failingProbability    float32 = 0.3
)

// SimulateIOProcessing simulates IO load by using random time from minProcessingDuration to maxProcessingDuration
// And randomly generate failings with failingProbability
func SimulateIOProcessing(ctx context.Context) error {
	duration := minProcessingDuration + time.Duration(rand.Int63n(int64(maxProcessingDuration-minProcessingDuration)))

	select {
	case <-time.After(duration):
	case <-ctx.Done():
		return ctx.Err()
	}

	if rand.Float32() < failingProbability {
		return errors.New("task failed")
	}
	return nil
}
