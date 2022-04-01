package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func preorderTraversalRec(root *TreeNode) []int {
	res := []int{}
	var helper func(*TreeNode)
	helper = func(root *TreeNode) {
		if root == nil {
			return
		}
		res = append(res, root.Val)
		helper(root.Left)
		helper(root.Right)
	}
	helper(root)
	return res
}

func preorderTraversalIterative(root *TreeNode) []int {
	s := []*TreeNode{}
	res := []int{}
	for root != nil || len(s) > 0 {
		for root != nil {
			res = append(res, root.Val)
			s = append(s, root)
			root = root.Left
		}
		root = s[len(s)-1]
		s = s[:len(s)-1]
		root = root.Right
	}
	return res
}

func inorderTraversal(root *TreeNode) []int {

	ans := []int{}

	var inorder func(*TreeNode)

	inorder = func(root *TreeNode) {
		if root == nil {
			return
		}

		inorder(root.Left)
		ans = append(ans, root.Val)
		inorder(root.Right)

	}
	inorder(root)

	return ans
}

func iTraversal(root *TreeNode) []int {

	ans := []int{}

	var inorder func(*TreeNode)

	inorder = func(root *TreeNode) {
		if root == nil {
			return
		}

		inorder(root.Left)
		ans = append(ans, root.Val)
		inorder(root.Right)

	}
	inorder(root)

	return ans
}

//Bottum up approach for calculating Depth of tree. if root node is leaf node then depth =1
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	leftDepth := maxDepth(root.Left)
	rightDepth := maxDepth(root.Right)
	if leftDepth > rightDepth {
		return leftDepth + 1
	} else {
		return rightDepth + 1
	}
}

// find out the symmetric tree like mirror

func isSymmetric(root *TreeNode) bool {
	return isMiror(root, root)
}

func isMiror(t1 *TreeNode, t2 *TreeNode) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t1.Val == t2.Val && isMiror(t1.Right, t2.Left) && isMiror(t1.Left, t2.Right)
}

// calculate the sum of path from root to leaf node
func hasPathSum(root *TreeNode, sum int) bool {
	if root == nil {
		return false
	}
	if root.Left == nil && root.Right == nil && root.Val == sum {
		return true
	}
	return hasPathSum(root.Left, sum-root.Val) || hasPathSum(root.Right, sum-root.Val)
}
