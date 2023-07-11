package skipList

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	DefaultMaxLevel int = 5
)

type Item struct {
	key   string
	value string
}

type Node struct {
	maxLevel  int
	prevNode  []*Node
	nextNode  []*Node
	item      Item
	isEndNode bool
}

func (node *Node) IsEndNode() bool {
	return node.isEndNode
}

func (node *Node) Match(key string) bool {
	return key == node.item.key
}

func (node *Node) MaxLevel() int {
	return node.maxLevel
}

func (node *Node) Key() string {
	return node.item.key
}

func (node *Node) Value() string {
	return node.item.value
}

func (node *Node) Next(targetLevel int) *Node {
	if node.maxLevel < targetLevel {
		return nil
	}
	return node.nextNode[targetLevel]
}

func (node *Node) Prev(targetLevel int) *Node {
	if node.maxLevel < targetLevel {
		return nil
	}
	return node.prevNode[targetLevel]
}

func (node *Node) AppendOnLevel(newNode *Node, targetLevel int) {
	if node.nextNode[targetLevel] != nil {
		node.nextNode[targetLevel].prevNode[targetLevel] = newNode
	}

	newNode.prevNode[targetLevel] = node
	newNode.nextNode[targetLevel] = node.nextNode[targetLevel]

	node.nextNode[targetLevel] = newNode
}

func (node *Node) RemoveOnLevel(targetLevel int) {
	if node.nextNode[targetLevel] != nil {
		node.nextNode[targetLevel].prevNode[targetLevel] = node.prevNode[targetLevel]
	}

	if node.prevNode[targetLevel] != nil {
		node.prevNode[targetLevel].nextNode[targetLevel] = node.nextNode[targetLevel]
	}
}

// type NodeList struct {
// 	level int
// }

type SkipList struct {
	maxLevel   int
	randSource rand.Source
	head       *Node
	tail       *Node
}

func New(maxLevel int) *SkipList {
	headNode := &Node{
		maxLevel:  maxLevel,
		prevNode:  make([]*Node, maxLevel),
		nextNode:  make([]*Node, maxLevel),
		item:      Item{},
		isEndNode: true,
	}

	tailNode := &Node{
		maxLevel:  maxLevel,
		prevNode:  make([]*Node, maxLevel),
		nextNode:  make([]*Node, maxLevel),
		item:      Item{},
		isEndNode: true,
	}

	list := SkipList{
		maxLevel:   maxLevel,
		randSource: rand.New(rand.NewSource(time.Now().UnixNano())),
		head:       headNode,
		tail:       tailNode,
	}

	for i := 0; i < DefaultMaxLevel; i++ {
		list.head.AppendOnLevel(list.tail, i)
	}

	return &list
}

func (list *SkipList) MaxLevel() int {
	return list.maxLevel
}

func (list *SkipList) Set(key, value string) {
	history := make([]*Node, list.maxLevel)
	node := list.findInternal(key, history)
	if node != nil {
		node.item.value = value
		return
	}

	list.insertNode(key, value, history)
}

func (list *SkipList) Remove(key string) {
}

func (list *SkipList) Get(key string) *Item {
	history := make([]*Node, list.maxLevel)
	node := list.findInternal(key, history)
	if node == nil {
		return nil
	}
	return &node.item
}

func (list *SkipList) findInternal(key string, history []*Node) *Node {
	current := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for list.tail != current.Next(i) && current.Next(i).item.key < key {
			current = current.Next(i)
		}
		history[i] = current
	}

	current = current.Next(0)
	if current.isEndNode || !current.Match(key) {
		return nil
	}
	return current
}

func (list *SkipList) insertNode(key, value string, history []*Node) {
	randomLevel := list.randomLevel()
	fmt.Println("randomLevel : ", randomLevel)

	node := &Node{
		maxLevel:  randomLevel,
		prevNode:  make([]*Node, randomLevel),
		nextNode:  make([]*Node, randomLevel),
		item:      Item{key: key, value: value},
		isEndNode: false,
	}

	for i := 0; i < randomLevel; i++ {
		history[i].AppendOnLevel(node, i)
	}
}

func (list *SkipList) randomLevel() int {
	const prob = 1 << 62
	maxLevel := list.maxLevel
	rand := list.randSource

	level := 1
	for ; (level < maxLevel) && (rand.Int63() > prob); level++ {
	}

	return level
}

func (list *SkipList) Print() {
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

		current = current.Next(level)
		if current == nil {
			current = list.tail
		}
	}
	fmt.Println()
}
