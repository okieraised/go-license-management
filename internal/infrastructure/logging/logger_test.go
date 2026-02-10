package logging

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestLoggerWithCustomFields_NoRaceCondition tests that WithCustomFields is thread-safe
func TestLoggerWithCustomFields_NoRaceCondition(t *testing.T) {
	logger := NewECSLogger()

	// Run concurrent calls to WithCustomFields
	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Each goroutine creates its own derived logger
				derivedLogger := logger.WithCustomFields(
					zap.String("goroutine_id", string(rune(id))),
					zap.Int("iteration", j),
				)
				// Use the derived logger (prevents optimization)
				assert.NotNil(t, derivedLogger)
			}
		}(i)
	}

	wg.Wait()
	// If there was a race condition, this test would fail or panic
}

// TestLoggerWithCustomStringFields_NoRaceCondition tests that WithCustomStringFields is thread-safe
func TestLoggerWithCustomStringFields_NoRaceCondition(t *testing.T) {
	logger := NewECSLogger()

	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Each goroutine creates its own derived logger
				derivedLogger := logger.WithCustomStringFields("test_key", "test_value")
				assert.NotNil(t, derivedLogger)
			}
		}(i)
	}

	wg.Wait()
}

// TestLoggerWithCustomFields_DoesNotMutateOriginal tests that derived loggers don't affect original
func TestLoggerWithCustomFields_DoesNotMutateOriginal(t *testing.T) {
	logger := NewECSLogger()
	original := logger.GetLogger()

	// Create derived loggers with fields
	derived1 := logger.WithCustomFields(zap.String("field1", "value1"))
	derived2 := logger.WithCustomFields(zap.String("field2", "value2"))
	derived3 := logger.WithCustomStringFields("field3", "value3")

	// Verify they are different instances
	assert.NotEqual(t, original, derived1, "Derived logger should be different from original")
	assert.NotEqual(t, derived1, derived2, "Derived loggers should be different from each other")
	assert.NotEqual(t, derived2, derived3, "Derived loggers should be different from each other")

	// Original logger should remain unchanged
	assert.Equal(t, original, logger.GetLogger(), "Original logger should not be modified")
}

// TestLoggerWithCustomFields_IndependentDerivedLoggers tests that derived loggers are independent
func TestLoggerWithCustomFields_IndependentDerivedLoggers(t *testing.T) {
	logger := NewECSLogger()

	// Create multiple derived loggers concurrently
	results := make(chan *zap.Logger, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			derived := logger.WithCustomFields(zap.Int("id", id))
			results <- derived
		}(i)
	}

	// Collect all derived loggers
	derivedLoggers := make([]*zap.Logger, 10)
	for i := 0; i < 10; i++ {
		derivedLoggers[i] = <-results
	}

	// Verify all are unique instances
	for i := 0; i < len(derivedLoggers); i++ {
		for j := i + 1; j < len(derivedLoggers); j++ {
			assert.NotEqual(t, derivedLoggers[i], derivedLoggers[j],
				"Each derived logger should be independent")
		}
	}
}

// TestLoggerGetInstance_Singleton tests singleton pattern
func TestLoggerGetInstance_Singleton(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()

	assert.Equal(t, instance1, instance2, "GetInstance should return the same singleton")
}

// TestNewECSLogger_CreatesNewInstance tests that NewECSLogger creates new instances
func TestNewECSLogger_CreatesNewInstance(t *testing.T) {
	logger1 := NewECSLogger()
	logger2 := NewECSLogger()

	assert.NotEqual(t, logger1, logger2, "NewECSLogger should create new instances")
	assert.NotEqual(t, logger1.GetLogger(), logger2.GetLogger(),
		"Each logger instance should have its own zap.Logger")
}

// TestLoggerThreadSafety_RaceDetector tests with race detector
// Run with: go test -race
func TestLoggerThreadSafety_RaceDetector(t *testing.T) {
	logger := NewECSLogger()

	// This test is designed to trigger race detector if there's a problem
	var wg sync.WaitGroup
	const concurrency = 50

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			// Mix different operations
			_ = logger.GetLogger()
			_ = logger.GetSugarLogger()
			_ = logger.WithCustomFields(zap.String("test", "value"))
			_ = logger.WithCustomStringFields("test", "value")
		}()
	}

	wg.Wait()
}

// BenchmarkLoggerWithCustomFields benchmarks field addition performance
func BenchmarkLoggerWithCustomFields(b *testing.B) {
	logger := NewECSLogger()
	fields := []zap.Field{
		zap.String("key1", "value1"),
		zap.String("key2", "value2"),
		zap.Int("key3", 123),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.WithCustomFields(fields...)
	}
}

// BenchmarkLoggerWithCustomStringFields benchmarks string field performance
func BenchmarkLoggerWithCustomStringFields(b *testing.B) {
	logger := NewECSLogger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = logger.WithCustomStringFields("key", "value")
	}
}

// BenchmarkLoggerConcurrent benchmarks concurrent logger access
func BenchmarkLoggerConcurrent(b *testing.B) {
	logger := NewECSLogger()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = logger.WithCustomFields(zap.String("test", "value"))
		}
	})
}
