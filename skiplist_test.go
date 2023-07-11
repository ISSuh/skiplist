package skipList

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)
}

func TestMaxLevel(t *testing.T) {
	list := New(5)
	if assert.NotNil(t, list) {
		assert.Equal(t, list.MaxLevel(), 5)
	}
}

func TestSetAndGet(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		list.Set(key, value)
	}

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)

		item := list.Get(key)

		if assert.NotNil(t, item) {
			assert.Equal(t, item.value, value)
		}
	}
}

func TestGet(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		list.Set(key, value)
	}

	item_1 := list.Get("1")
	assert.Equal(t, item_1.key, "1")
	assert.Equal(t, item_1.value, "1")

	item_empty := list.Get("222")
	assert.Equal(t, item_empty, (*Item)(nil))
}

func TestRemove(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		list.Set(key, value)
	}

	item_1 := list.Get("1")
	assert.Equal(t, item_1.key, "1")
	assert.Equal(t, item_1.value, "1")

	list.Remove("1")
	item_temp := list.Get("1")
	assert.Equal(t, item_temp, (*Item)(nil))
}

func TestConcurrency(t *testing.T) {
	list := New(10)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < 100000; i++ {
			key := strconv.Itoa(i)
			value := strconv.Itoa(i)
			list.Set(key, value)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			key := strconv.Itoa(i)
			list.Get(key)
		}
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(t, list.Length(), 100000)
}

var benchList *SkipList

func BenchmarkSet(b *testing.B) {
	b.ReportAllocs()

	benchList = New(15)

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		benchList.Set(key, value)
	}

	b.SetBytes(int64(b.N))
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		benchList.Get(key)
	}

	b.SetBytes(int64(b.N))
}
