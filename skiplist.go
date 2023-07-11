// Copyright 2011 ISSuh. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

// Package skiplist implement skip list data structure.
// Reference: https://en.wikipedia.org/wiki/Skip_list

// Example

// list := skipList.New(5)
// list.Set("key", "value")

// item := list.Get("key")
// fmt.Printf("key : %s / value : %s", item.key, item.value)

package skipList

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Item struct {
	key   string
	value string
}

type Node struct {
	levels    int
	prevNode  []*Node
	nextNode  []*Node
	item      Item
	isEndNode bool
}

func (node *Node) match(key string) bool {
	return key == node.item.key
}

func (node *Node) nodeLevel() int {
	return node.levels
}

func (node *Node) next(targetLevel int) *Node {
	if node.levels < targetLevel {
		return nil
	}
	return node.nextNode[targetLevel]
}

func (node *Node) appendOnLevel(newNode *Node, targetLevel int) {
	if node.nextNode[targetLevel] != nil {
		node.nextNode[targetLevel].prevNode[targetLevel] = newNode
	}

	newNode.prevNode[targetLevel] = node
	newNode.nextNode[targetLevel] = node.nextNode[targetLevel]

	node.nextNode[targetLevel] = newNode
}

func (node *Node) removeOnLevel(targetLevel int) {
	if node.nextNode[targetLevel] != nil {
		node.nextNode[targetLevel].prevNode[targetLevel] = node.prevNode[targetLevel]
	}

	if node.prevNode[targetLevel] != nil {
		node.prevNode[targetLevel].nextNode[targetLevel] = node.nextNode[targetLevel]
	}
}

type SkipList struct {
	maxLevel int
	length   int
	head     *Node
	tail     *Node
	rand     *rand.Rand
	mutex    sync.RWMutex
	history  []*Node
}

func New(maxLevel int) *SkipList {
	headNode := &Node{
		levels:    maxLevel,
		prevNode:  make([]*Node, maxLevel),
		nextNode:  make([]*Node, maxLevel),
		item:      Item{},
		isEndNode: true,
	}

	tailNode := &Node{
		levels:    maxLevel,
		prevNode:  make([]*Node, maxLevel),
		nextNode:  make([]*Node, maxLevel),
		item:      Item{},
		isEndNode: true,
	}

	list := SkipList{
		maxLevel: maxLevel,
		length:   0,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		head:     headNode,
		tail:     tailNode,
		history:  make([]*Node, maxLevel),
	}

	for i := 0; i < maxLevel; i++ {
		list.head.appendOnLevel(list.tail, i)
	}

	return &list
}

func (list *SkipList) MaxLevel() int {
	return list.maxLevel
}

func (list *SkipList) Length() int {
	return list.length
}

func (list *SkipList) Set(key, value string) {
	node := list.findInternal(key, list.history)
	if node != nil {
		node.item.value = value
		return
	}

	list.insertNode(key, value, list.history)
}

func (list *SkipList) Get(key string) *Item {
	node := list.findInternal(key, list.history)
	if node == nil {
		return nil
	}
	return &node.item
}

func (list *SkipList) Remove(key string) {
	node := list.findInternal(key, list.history)
	if node == nil {
		return
	}

	list.deleteNode(node)
	list.length--
}

func (list *SkipList) findInternal(key string, history []*Node) *Node {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	current := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for list.tail != current.next(i) && current.next(i).item.key < key {
			current = current.next(i)
		}
		history[i] = current
	}

	current = current.next(0)
	if current.isEndNode || !current.match(key) {
		return nil
	}
	return current
}

func (list *SkipList) insertNode(key, value string, history []*Node) {
	randomLevel := list.randomLevel()

	node := &Node{
		levels:    randomLevel,
		prevNode:  make([]*Node, randomLevel),
		nextNode:  make([]*Node, randomLevel),
		item:      Item{key: key, value: value},
		isEndNode: false,
	}

	list.mutex.Lock()
	defer list.mutex.Unlock()

	for i := 1; i <= randomLevel; i++ {
		randomLevelIndex := i - 1
		history[randomLevelIndex].appendOnLevel(node, randomLevelIndex)
	}

	list.length++
}

func (list *SkipList) deleteNode(node *Node) {
	list.mutex.Lock()
	defer list.mutex.Unlock()

	for i := 0; i < node.nodeLevel(); i++ {
		node.removeOnLevel(i)
	}
}

func (list *SkipList) randomLevel() int {
	const prob = 1 << 30
	maxLevel := list.maxLevel
	rand := list.rand

	level := 1
	for ; (level < maxLevel) && (rand.Int31() > prob); level++ {
	}

	return level
}

func (list *SkipList) PrintForDebug() {
	for i := list.maxLevel - 1; i >= 0; i-- {
		list.printLevel(i)
	}
	fmt.Println("----------------------------------")
}

func (list *SkipList) printLevel(level int) {
	fmt.Printf("[%d] : ", level)
	current := list.head
	for current != list.tail {
		if current != list.head {
			fmt.Printf("[%s, %s]\t", current.item.key, current.item.value)
		}

		current = current.next(level)
		if current == nil {
			current = list.tail
		}
	}
	fmt.Println()
}
