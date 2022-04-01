package main

type MyCircularQueue struct {
	arr   []int
	front int
	rear  int
	size  int
}

func Constructor(k int) MyCircularQueue {
	return MyCircularQueue{
		arr:   make([]int, k),
		front: -1,
		rear:  -1,
		size:  k,
	}
}

func (this *MyCircularQueue) EnQueue(value int) bool {
	if this.IsFull() {
		return false
	}
	if this.IsEmpty() {
		this.front = 0
	}
	this.rear = (this.rear + 1) % this.size
	this.arr[this.rear] = value
	return true
}

func (this *MyCircularQueue) DeQueue() bool {

	if this.IsEmpty() {
		return false
	}
	if this.front == this.rear {
		this.front = -1
		this.rear = -1
		return true
	}
	this.front = (this.front + 1) % this.size
	return true
}

func (this *MyCircularQueue) Front() int {
	if this.IsEmpty() {
		return -1
	}
	return this.arr[this.front]
}

func (this *MyCircularQueue) Rear() int {
	if this.IsEmpty() {
		return -1
	}
	return this.arr[this.rear]
}

func (this *MyCircularQueue) IsEmpty() bool {
	return this.front == -1
}

func (this *MyCircularQueue) IsFull() bool {
	return ((this.rear + 1) % this.size) == this.front
}

func main() {

	// Your MyCircularQueue object will be instantiated and called as such:
	obj := Constructor(k)
	param_1 := obj.EnQueue(value)
	param_2 := obj.DeQueue()
	param_3 := obj.Front()
	param_4 := obj.Rear()
	param_5 := obj.IsEmpty()
	param_6 := obj.IsFull()

}
