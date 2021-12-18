package hashtree

import "crypto/sha256"

type Node struct {
	Left *Node
	Right *Node
	Parent *Node
	Hash   [sha256.Size]byte
}

func newParentNode(left, right *Node) *Node {
	node := &Node{
		Left:  left,
		Right: right,
		Hash:  CalculateHash(append(left.getSHA(), right.getSHA()...)),
	}

	left.SetParent(node)
	right.SetParent(node)

	return node
}

func nodeByIndex(leaves []*Node, index int) *Node {
	if len(leaves) <= index {
		return nil
	}
	return leaves[index]
}

func (n *Node) getSHA() []byte {
	if n == nil {
		return nil
	}
	return n.Hash[:]
}

func (n *Node) SetParent(parent *Node) {
	if n == nil {
		return
	}
	n.Parent = parent
}