package zpm_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xakepp35/zpm"
)

func TestSummary(t *testing.T) {
	const name = "test_summary"
	// Добавляем значения
	values := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	expectedSum := 0.0
	var wg sync.WaitGroup
	wg.Add(len(values))
	for i := range values {
		value := values[i]
		expectedSum += value
		go func() {
			defer wg.Done()
			zpm.Summary(name).
				Quantiles(0.5, 0.9, 0.99).
				Observe(value)
		}()
	}
	wg.Wait()
	summary := zpm.Summary(name).State().Dto.Summary
	// Проверяем, что SampleCount и SampleSum обновились правильно
	assert.Equal(t, uint64(len(values)), *summary.SampleCount)
	assert.Equal(t, expectedSum, *summary.SampleSum)
	// Проверяем квантильные значения
	quantiles := map[float64]float64{
		0.5:  50.0,  // Медиана (середина отсортированного массива)
		0.9:  90.0,  // 90-й перцентиль
		0.99: 100.0, // 99-й перцентиль
	}
	for _, q := range summary.Quantile {
		expected, exists := quantiles[*q.Quantile]
		assert.True(t, exists, "Неизвестный квантиль %v", *q.Quantile)
		assert.InDelta(t, expected, *q.Value, 5.0, "Квантиль %v не в допустимом диапазоне", *q.Quantile)
	}
}
