package main

import "fmt"

type Node struct {
	next *Node
	data int
}

func NewNode(d int, n *Node) *Node {
	var v Node
	v.data = d
	v.next = n
	return &v
}

func PrintNode(node *Node) {
	for node != nil {
		fmt.Println(node.data)
		node = node.next
	}
}

func addNodeAtEnd(d int, head *Node) *Node {
	if head.next == nil {
		head.next = NewNode(d, nil)
		return head
	}

	temp := head
	for temp.next != nil {
		temp = temp.next
	}

	temp.next = NewNode(d, nil)
	return head

}

func removeLastNode(head *Node) *Node {
	if head.next == nil {
		return nil
	}
	temp := head
	for temp.next.next != nil {
		temp = temp.next
	}
	temp.next = nil
	return head

}
func main() {
	node := NewNode(10, nil)
	addNodeAtEnd(20, node)
	addNodeAtEnd(30, node)
	addNodeAtEnd(40, node)
	removeLastNode(node)
	addNodeAtEnd(50, node)
	PrintNode(node)
}
