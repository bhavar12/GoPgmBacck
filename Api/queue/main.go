package main

import "fmt"

type Queue struct {
	data []int
}

func (q *Queue) enqueue(n int) {
	q.data = append(q.data, n)
	fmt.Println("Element added to queue..", n)
}

func (q *Queue) dequeue() int {
	if len(q.data) > 0 {
		n := q.data[0]
		q.data = q.data[1:]
		fmt.Println(n, " deleted from queue")
		return n
	} else {
		fmt.Println("Queue is empty")
	}
	return -1
}
func main() {
	q := &Queue{}
	q.enqueue(1)
	q.enqueue(2)
	q.enqueue(3)
	q.dequeue()
	q.dequeue()
	q.enqueue(4)
	fmt.Println("Final queue ele..", q.data)

}
