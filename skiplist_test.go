package skiplist

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var charSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = charSet[rand.Int63()%int64(len(charSet))]
	}
	return string(b)
}

func TestOverLevelOnNode(t *testing.T) {
	fistNode := &SkipListNode{
		levels:    5,
		prevNode:  make([]*SkipListNode, 5),
		nextNode:  make([]*SkipListNode, 5),
		item:      SkipListItem{key: "1", value: []byte("1")},
		isEndNode: false,
	}

	secondNode := &SkipListNode{
		levels:    3,
		prevNode:  make([]*SkipListNode, 3),
		nextNode:  make([]*SkipListNode, 3),
		item:      SkipListItem{key: "2", value: []byte("2")},
		isEndNode: false,
	}

	fistNode.appendOnLevel(secondNode, 0)
	fistNode.appendOnLevel(secondNode, 1)
	fistNode.appendOnLevel(secondNode, 2)

	temp := fistNode.next(0)
	assert.Equal(t, temp, secondNode)

	temp = fistNode.next(1)
	assert.Equal(t, temp, secondNode)

	temp = fistNode.next(2)
	assert.Equal(t, temp, secondNode)

	temp = fistNode.next(3)
	assert.Equal(t, temp, (*SkipListNode)(nil))

	temp = fistNode.next(6)
	assert.Equal(t, temp, (*SkipListNode)(nil))
}

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
		list.Set(key, []byte(value))
	}

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)

		item := list.Get(key)

		if assert.NotNil(t, item) {
			assert.Equal(t, item.Value(), []byte(value))
		}
	}
}

func TestGet(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		list.Set(key, []byte(value))
	}

	item_1 := list.Get("1")
	assert.Equal(t, item_1.Key(), "1")
	assert.Equal(t, item_1.Value(), []byte("1"))

	item_empty := list.Get("222")
	assert.Equal(t, item_empty, (*SkipListItem)(nil))
}

func TestSize(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	for i := 0; i < 2; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		list.Set(key, []byte(value))
	}

	assert.Equal(t, list.Size(), uint64(4))
}

func TestUpdate(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	list.Set("1", []byte("1"))

	item_1 := list.Get("1")
	assert.Equal(t, item_1.Key(), "1")
	assert.Equal(t, item_1.Value(), []byte("1"))

	list.Set("1", []byte("11"))
	item_1 = list.Get("1")
	assert.Equal(t, item_1.Key(), "1")
	assert.Equal(t, item_1.Value(), []byte("11"))
}

func TestRemove(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	for i := 0; i < 5; i++ {
		key := strconv.Itoa(i)
		value := strconv.Itoa(i)
		list.Set(key, []byte(value))
	}

	item_1 := list.Get("1")
	assert.Equal(t, item_1.key, "1")
	assert.Equal(t, item_1.value, []byte("1"))

	list.Remove("1")
	item_temp := list.Get("1")
	assert.Equal(t, item_temp, (*SkipListItem)(nil))

	list.Remove("1")
	item_temp = list.Get("1")
	assert.Equal(t, item_temp, (*SkipListItem)(nil))
}

func TestIterateNext(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	var temp []string
	for i := 0; i < 30; i++ {
		word := randomString(5)

		key := word
		value := word
		list.Set(key, []byte(value))

		temp = append(temp, word)
	}

	sort.Strings(temp)

	node := list.Front()
	for _, word := range temp {
		if assert.NotNil(t, node) {
			assert.Equal(t, node.Key(), word)
			assert.Equal(t, node.Value(), []byte(word))
		}
		node = node.Next()
	}
}

func TestIteratePrev(t *testing.T) {
	list := New(5)
	assert.NotEqual(t, list, nil)

	var temp []string
	for i := 0; i < 30; i++ {
		word := randomString(5)

		key := word
		value := word
		list.Set(key, []byte(value))

		temp = append(temp, word)
	}

	sort.Strings(temp)

	node := list.Back()
	for i := range temp {
		word := temp[len(temp)-1-i]
		if assert.NotNil(t, node) {
			assert.Equal(t, node.Key(), word)
			assert.Equal(t, node.Value(), []byte(word))
		}
		node = node.Prev()
	}
}

func TestConcurrency(t *testing.T) {
	list := New(10)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for i := 0; i < 100000; i++ {
			key := strconv.Itoa(i)
			value := strconv.Itoa(i)
			list.Set(key, []byte(value))
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
		benchList.Set(key, []byte(value))
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
