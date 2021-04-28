package dataStructures

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	value     int
	leftNode  *Node
	rightNode *Node
}

func (currentNode *Node) addValue(value int) Node {

	newNode := Node{value: value}

	if currentNode.value < value {
		if currentNode.rightNode != nil {
			currentNode.rightNode.addValue(value)
		} else {
			currentNode.rightNode = &newNode
		}
	}

	if currentNode.value > value {
		if currentNode.leftNode != nil {
			currentNode.leftNode.addValue(value)
		} else {
			currentNode.leftNode = &newNode
		}
	}

	return newNode
}

func (currentNode *Node) findValue(value int, steps int) (*Node, int) {
	if currentNode.value > value {
		return currentNode.leftNode.findValue(value, steps+1)
	}

	if currentNode.value < value {

		return currentNode.rightNode.findValue(value, steps+1)
	}

	if currentNode.value == value {
		return currentNode, steps
	}

	return nil, steps
}

func (currentNode *Node) findMaxValue() *Node {

	if currentNode.rightNode != nil {
		return currentNode.rightNode.findMaxValue()
	}

	return currentNode
}

func (currentNode *Node) findMinValue() *Node {

	if currentNode.leftNode != nil {
		return currentNode.leftNode.findMinValue()
	}

	return currentNode
}

func TestBST() {
	var node Node

	for i := 0; i < 1000000; i++ {
		if i == 0 {
			node = Node{value: rand.Intn(100000)}
		}
		node.addValue(rand.Intn(100000))
	}

	searchStart := time.Now()
	result, steps := node.findValue(35345, 0)
	maxValue := node.findMaxValue()
	minValue := node.findMinValue()

	// result := Find(randomNodes, rand.Intn(999))

	fmt.Println("result ", result)
	fmt.Println("maxValue ", maxValue)
	fmt.Println("minValue ", minValue)

	fmt.Println(fmt.Printf("found result in %s seconds and %d steps", time.Since(searchStart), steps))
}
