package main

import "fmt"

//Node ...
type Node struct {
	val  int
	next *Node
}

//LinkedList ...
type LinkedList struct {
	head *Node
	len  int
}

//Insert ...
func (l *LinkedList) Insert(val int) {
	n := Node{}
	n.val = val
	if l.len == 0 {
		l.head = &n
		l.len++
		return
	}
	ptr := l.head
	for i := 0; i < l.len; i++ {
		if ptr.next == nil {
			ptr.next = &n
			l.len++
			return
		}
		ptr = ptr.next
	}
}

func (l *LinkedList) insertAt(val, pos int) {
	if pos < 0 {
		return
	}
	newNode := Node{}
	newNode.val = val
	if pos == 0 {
		l.head = &newNode
		l.len++
	}
	if pos > l.len {
		return
	}
	n := l.GetAt(pos)
	newNode.next = n
	prevNode := l.GetAt(pos - 1)
	prevNode.next = &newNode
	l.len++
}

//Print ...
func (l *LinkedList) Print() {
	if l.len == 0 {
		fmt.Println("No nodes in list")
	}
	ptr := l.head
	for i := 0; i < l.len; i++ {
		fmt.Println("Node: ", ptr)
		ptr = ptr.next
	}
}

// GetAt returns node at given position from linked list
func (l *LinkedList) GetAt(pos int) *Node {
	ptr := l.head
	if pos < 0 {
		return ptr
	}
	if pos > (l.len - 1) {
		return nil
	}
	for i := 0; i < pos; i++ {
		ptr = ptr.next
	}
	return ptr
}

// Search returns node position with given value from linked list
func (l *LinkedList) Search(val int) int {
	ptr := l.head
	for i := 0; i < l.len; i++ {
		if ptr.val == val {
			return i
		}
		ptr = ptr.next
	}
	return -1
}

func main() {
	listA := LinkedList{}
	listB := LinkedList{}
	fmt.Println("\n************* Insert *************")
	listA.Insert(1)
	listA.Insert(9)
	listA.Insert(1)
	listA.Insert(2)
	listA.Insert(4)

	//list B
	listB.Insert(3)
	listB.Insert(2)
	listB.Insert(4)
	fmt.Println("************* Print *************")
	//interSectVal := 2
	fmt.Println(getIntersectionNode(listA.head, listB.head))

}

func getIntersectionNode(headA, headB *Node) *Node {
	for headA != nil {
		pB := headB
		for pB != nil {
			if headA == pB {
				return headA
			}
			pB = pB.next
		}
		headA = headA.next
	}
	return nil
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func swapNodes(head *ListNode, k int) *ListNode {
	n := 0
	rootNode := head
	arr := make([]int, 0)
	for rootNode != nil {
		arr = append(arr, rootNode.Val)
		n++
		rootNode = rootNode.Next
	}
	temp := arr[k-1]
	arr[k-1] = arr[n-k]
	arr[n-k] = temp
	rootNode = &ListNode{Val: arr[0]}
	tempNode := rootNode
	for i := 1; i < n; i++ {
		rootNode.Next = &ListNode{Val: arr[i]}
		rootNode = rootNode.Next
	}
	return tempNode
}

func isPalindrome(head *ListNode) bool {
	n := 0
	node1 := head
	arr := make([]int, 0)
	for node1 != nil {
		arr = append(arr, node1.Val)
		n++
		node1 = node1.Next
	}
	for i := 0; i < n; i++ {
		if arr[i] == arr[n-1] {
			n--
		} else {
			return false
		}
	}
	return true
}
