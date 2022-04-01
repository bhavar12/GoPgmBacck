package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var (
	flippedArray = make([]int, 0)
	index        = 0
)

func addOneRow(root *TreeNode, v int, d int) *TreeNode {
	if d == 1 {
		node := &TreeNode{
			Val: v,
		}
		node.Left = root
		return node
	}
	queue := []*TreeNode{}
	queue = append(queue, root)
	depth := 1

	for depth < d-1 {
		temp := []*TreeNode{}
		for len(queue) != 0 {
			newNode := queue[0]
			if newNode.Left != nil {
				temp = append(temp, newNode.Left)
			}
			if newNode.Right != nil {
				temp = append(temp, newNode.Right)
			}
			queue = temp[1:]
		}
		depth++
	}
	for len(queue) != 0 {
		popNode := queue[0]
		temp := popNode.Left
		popNode.Left = &TreeNode{Val: v}
		popNode.Left.Left = temp
		temp = popNode.Right
		popNode.Right = &TreeNode{Val: v}
		popNode.Right.Right = temp
		queue = queue[1:]
	}
	return root
}

func BFS(node *TreeNode) []int {
	queue := []*TreeNode{}
	queue = append(queue, node)
	result := []int{}
	return BFSUtil(queue, result)
}

//BFSUtil ...
func BFSUtil(queue []*TreeNode, res []int) []int {
	if len(queue) == 0 {
		return res
	}
	res = append(res, queue[0].Val)
	if queue[0].Right != nil {
		queue = append(queue, queue[0].Right)
	}

	if queue[0].Left != nil {
		queue = append(queue, queue[0].Left)
	}

	return BFSUtil(queue[1:], res)
}

func flipMatchVoyage(root *TreeNode, voyage []int) []int {

	dFSAlgo(root, voyage)
	return flippedArray
}

func dFSAlgo(root *TreeNode, voyage []int) {
	if root == nil || (len(flippedArray) != 0 && flippedArray[0] == -1) {
		return
	}
	v := voyage[index]
	index++
	if root.Val != v {
		flippedArray = append(flippedArray, -1)
		return
	} else if root.Left != nil && root.Left.Val != voyage[index] {
		flippedArray = append(flippedArray, root.Left.Val)
		dFSAlgo(root.Right, voyage)
		dFSAlgo(root.Left, voyage)
	} else {
		dFSAlgo(root.Left, voyage)
		dFSAlgo(root.Right, voyage)
	}

}
