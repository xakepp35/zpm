package algo

import (
	"math"
	"sync/atomic"
	"unsafe"
)

const bufferSize = 1024         // Фиксированный размер буфера
const defaultEpsolon = 1.0 / 64 // Точность по умолчанию

type CKMSLockless struct {
	epsilon  float64 // точность
	quantile []float64
	buffer   [bufferSize]pair
	head     atomic.Uint64 // Индекс в буфере
}

func NewCKMSLockless(quantiles ...float64) *CKMSLockless {
	return &CKMSLockless{
		epsilon:  defaultEpsolon,
		quantile: quantiles,
	}
}

func (c *CKMSLockless) Insert(value float64) {
	// Получаем текущий индекс, продвигаем его
	index := c.head.Add(1) - 1
	if index >= bufferSize {
		c.head.Store(0)
		index = 0
	}
	index %= bufferSize

	// Записываем значение в буфер
	ptr := &c.buffer[index]
	AtomicFloatStore(&ptr.value, value)
	atomic.StoreInt32((*int32)(unsafe.Pointer(&ptr.g)), 1)
}

func (c *CKMSLockless) Query(q float64) float64 {
	count := int(c.head.Load())
	if count == 0 {
		return math.NaN()
	}

	// Находим квантиль "на лету" без копирования
	// rank := int(q * float64(count))
	rank := int(math.Ceil(q * float64(count))) - 1
	if rank >= count {
		rank = count - 1
	}

	// Используем медианный алгоритм (QuickSelect)
	return quickSelect(c.buffer[:count], rank)
}

func quickSelect(arr []pair, k int) float64 {
	left, right := 0, len(arr)-1
	for left < right {
		pivot := partition(arr, left, right)
		if pivot == k {
			break
		} else if pivot < k {
			left = pivot + 1
		} else {
			right = pivot - 1
		}
	}
	return AtomicFloatLoad(&arr[k].value)
}

func partition(arr []pair, left, right int) int {
	pivot := AtomicFloatLoad(&arr[right].value)
	i := left
	for j := left; j < right; j++ {
		valJ := AtomicFloatLoad(&arr[j].value)
		if valJ < pivot {
			arr[i], arr[j] = arr[j], arr[i] // Переставляем элементы
			i++
		}
	}
	arr[i], arr[right] = arr[right], arr[i]
	return i
}

// Простая сортировка вставками (без аллокации)
func insertionSort(arr []float64) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}
