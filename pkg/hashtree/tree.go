package hashtree

import (
	"bufio"
	"fmt"
	"io"
)

type Tree struct {
	head *Node
	leaves []*Node
}

func (t Tree) Hash() [HashSize]byte {
	return t.head.Hash
}

func (t Tree) LeafHashAt(index int) ([HashSize]byte, error) {
	if index >= len(t.leaves) {
		return emptyHash, fmt.Errorf("index is out of bounds: %d from %d", index, len(t.leaves)-1)
	}

	return t.leaves[index].Hash, nil
}

func (t Tree) UpdateLeafHashAt(index int, newHash [HashSize]byte) (Tree, error) {
	if index >= len(t.leaves) {
		return Tree{}, fmt.Errorf("index is out of bounds: %d from %d", index, len(t.leaves)-1)
	}

	leaf := t.leaves[index]
	newLeaf := &Node{Parent: leaf.Parent, Hash: newHash}

	newLeaves := append(t.leaves[:index], newLeaf)
	newLeaves = append(newLeaves, t.leaves[index+1:]...)

	newParent := newLeaf
	nodeHash := leaf.Hash
	for newParent.Parent != nil {
		parent := *newParent.Parent
		switch nodeHash {
		case parent.Left.Hash:
			newParent = newParentNode(newParent, parent.Right)
		case parent.Right.Hash:
			newParent = newParentNode(parent.Left, newParent)
		default:
			return Tree{}, fmt.Errorf("something wrong with hash tree")
		}
		newParent.Parent = parent.Parent
		nodeHash = parent.Hash
	}

	return Tree{
		head: newParent,
		leaves: newLeaves,
	}, nil
}

func ByLinesFromReader(r *bufio.Reader) (Tree, error) {
	leaves, err := readLeaves(r)
	if err != nil {
		return Tree{}, err
	}

	return buildTree(leaves)
}

func readLeaves(r *bufio.Reader) ([]*Node, error) {
	var (
		line []byte
		leaves []*Node
	)
	for {
		data, isPrefix, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				return leaves, nil
			}
			return nil, err
		}
		line = append(line, data...)
		if !isPrefix {
			leaves = append(leaves, &Node{
				Hash: CalculateHash(line),
			})
			line = nil
		}
	}
}

func buildTree(leaves []*Node) (Tree, error) {
	var (
		levelNodes = leaves
		nextLevelNodes []*Node
	)
	for len(nextLevelNodes) != 1 {
		nextLevelNodes = nil
		for i := 0; i < len(levelNodes); i += 2 {
			nextLevelNodes = append(nextLevelNodes, newParentNode(levelNodes[i], nodeByIndex(levelNodes, i+1)))
		}
		levelNodes = nextLevelNodes
	}

	return Tree{
		head: nextLevelNodes[0],
		leaves: leaves,
	}, nil
}
