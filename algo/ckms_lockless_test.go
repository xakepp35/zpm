package algo

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Тест на вставку и запрос квантилей
//
//	🔍 Суть: Проверяем, что после вставки 1000 элементов квантильные оценки соответствуют ожидаемым.
//	✅ Если да: Значения 50-го и 90-го квантилей около 499 и 900 (с допустимой погрешностью).
//	❌ Если нет: Значения не соответствуют ожиданиям, что указывает на ошибки в алгоритме CKMS.
func TestCKMSLockless_InsertAndQuery(t *testing.T) {
	c := NewCKMSLockless(0.5, 0.9)
	// Вставляем 1000 элементов
	for i := 0; i < 1000; i++ {
		c.Insert(float64(i))
	}
	// Проверяем квантили
	assert.InDelta(t, 499, c.Query(0.5), 10, "Медиана должна быть около 499")
	assert.InDelta(t, 900, c.Query(0.9), 20, "90-й квантиль должен быть около 900")
}

// Тест на конкурентную вставку
//
//	🔍 Суть: Проверяем, что параллельные вставки в несколько горутин работают без гонок и повреждения данных.
//	✅ Если да: Запрос квантиля возвращает осмысленное значение (500.0), без NaN и мусора.
//	❌ Если нет: Возвращается NaN, некорректные числа или происходит паника.
func TestCKMSLockless_ConcurrentInsert(t *testing.T) {
	c := NewCKMSLockless(0.5)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				c.Insert(float64(base*100 + j))
			}
		}(i)
	}
	wg.Wait()
	// Должны быть корректные данные
	result := c.Query(0.5)
	assert.False(t, math.IsNaN(result), "Квантиль не должен быть NaN")
	assert.Equal(t, 500.0, result)
}

// Тест на запрос из пустого CKMS
//
//	🔍 Суть: Проверяем, что при отсутствии вставленных значений функция Query() возвращает NaN.
//	✅ Если да: Вернётся NaN, как и ожидается, поскольку нет данных для оценки квантилей.
//	❌ Если нет: Возвращается любое число кроме NaN, что говорит о неправильном поведении при пустом CKMS.
func TestCKMSLockless_EmptyQuery(t *testing.T) {
	c := NewCKMSLockless(0.5)
	result := c.Query(0.5)
	assert.True(t, math.IsNaN(result), "На пустом CKMS результат должен быть NaN")
}

// Тест на гонку Insert vs Query при переполнении буфера
//
//	🛠 Суть: Вставляем данные в несколько потоков до полного заполнения буфера, затем читаем на лету. Если есть гонка, будут некорректные значения или паника.
//	✅ Если гонки нет, тест пройдет стабильно
//	❌ Если есть гонка, тест может зафейлиться с паникой
func TestCKMSLockless_RaceInsertQuery(t *testing.T) {
	c := NewCKMSLockless(0.5)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	// Вставляем данные в фоне, пока выполняется Query
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < bufferSize*5; i++ {
			c.Insert(float64(i))
			if i%bufferSize == 0 {
				time.Sleep(time.Microsecond) // Имитация перезаписи
			}
		}
		close(stop)
	}()
	// Читаем квантили в фоне, пока идет вставка
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_ = c.Query(0.5) // Должно быть без паники!
			}
		}
	}()
	wg.Wait()
}

// Тест на изменение буфера в quickSelect()
// 🛠 Суть: quickSelect() изменяет arr, но arr — это буфер c.buffer, который атомарно читается в Insert. Если буфер изменяется во время работы Query(), может быть повреждение данных.
// ✅ Если quickSelect() работает корректно, тест пройдет
// ❌ Если quickSelect() портит буфер, тест зафейлится
func TestCKMSLockless_BrokenQuickSelect(t *testing.T) {
	c := NewCKMSLockless(0.5)

	// Заполняем буфер
	for i := 0; i < bufferSize; i++ {
		c.Insert(float64(i))
	}

	before := make([]float64, bufferSize)
	for i := 0; i < bufferSize; i++ {
		before[i] = AtomicFloatLoad(&c.buffer[i].value)
	}

	// Запускаем Query, который внутри вызывает quickSelect
	_ = c.Query(0.5)

	// Проверяем, что данные в буфере не изменились
	for i := 0; i < bufferSize; i++ {
		after := AtomicFloatLoad(&c.buffer[i].value)
		assert.Equal(t, before[i], after, "quickSelect() изменил буфер!")
	}
}

// Тест на повреждение данных из-за unsafe.Pointer
//
//	🛠 Суть: Если unsafe.Pointer используется неправильно, может произойти разрыв данных при многопоточных операциях.
//	✅ Если Insert() работает безопасно, тест пройдет
//	❌ Если unsafe.Pointer используется некорректно, тест зафейлится
func TestCKMSLockless_UnsafePointer(t *testing.T) {
	c := NewCKMSLockless(0.5)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			c.Insert(val)
		}(float64(i))
	}
	wg.Wait()

	// Проверяем, что значения не содержат битых данных
	for i := 0; i < bufferSize; i++ {
		val := AtomicFloatLoad(&c.buffer[i].value)
		assert.False(t, math.IsNaN(val), "Найдено поврежденное значение NaN!")
		assert.Less(t, val, 1e10, "Найдено подозрительно большое значение")
	}
}
