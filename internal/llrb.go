/*
 * Copyright 2025 Alexandre Mahdhaoui
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package internal

import "cmp"

// ------------------------------------------------------------------------------
// -- NODE
// ------------------------------------------------------------------------------

type Direction int

const (
	Left Direction = iota
	Right
)

// rbnode is the node datastructure for a red/black tree.
type Node[K cmp.Ordered, V any] struct {
	Key   K
	Value V

	parent   *Node[K, V]
	children [2]*Node[K, V]
	isBlack  bool
}

func (n *Node[K, V]) Left() *Node[K, V] {
	return n.children[Left]
}

func (n *Node[K, V]) Right() *Node[K, V] {
	return n.children[Right]
}

func NewNode[K cmp.Ordered, V any](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		Key:      key,
		Value:    value,
		parent:   nil,
		children: [2]*Node[K, V]{},
		isBlack:  false,
	}
}

// ------------------------------------------------------------------------------
// -- SEARCH
// ------------------------------------------------------------------------------

func Search[K cmp.Ordered, V any](root *Node[K, V], key K) (V, bool) {
	for n := root; n != nil; {
		if key == n.Key {
			return n.Value, true
		}

		if key < n.Key {
			n = n.children[Right]
		} else {
			n = n.children[Left]
		}
	}

	var zeroVal V
	return zeroVal, false
}

// SearchMin implements the equivalentof the following recursive implementation.
//
//	```go
//	  if root.Left() == nil {
//		return root
//	  }
//
//	  return SearchMin(root.Left())
//	```
func SearchMin[K cmp.Ordered, V any](root *Node[K, V]) *Node[K, V] {
	n := root
	for {
		if n.Left() == nil {
			return n
		}
		n = n.Left()
	}
}

// ------------------------------------------------------------------------------
// -- INSERTION
// ------------------------------------------------------------------------------

func Insert[K cmp.Ordered, V any](root *Node[K, V], key K, value V) *Node[K, V] {
	if root == nil {
		return NewNode(key, value)
	}

	if key == root.Key {
		root.Value = value
	} else {
		var direction Direction
		if key < root.Key {
			direction = Left
		} else {
			direction = Right
		}

		root.children[direction] = Insert(root.children[direction], key, value)
	}

	return FixUp(root)
}

// ------------------------------------------------------------------------------
// -- DELETION
// ------------------------------------------------------------------------------

func Delete[K cmp.Ordered, V any](root *Node[K, V], key K) *Node[K, V] {
	if key < root.Key {
		if !IsRed(root.Left()) && !IsRed(root.Left().Left()) {
			root = MoveRedLeft(root)
		}

		root.children[Left] = Delete(root.Left(), key)
		return FixUp(root)
	}

	if IsRed(root.Left()) {
		root = Rotate(root, Right)
	}

	if key == root.Key && root.Right() == nil {
		return nil
	}

	if !IsRed(root.Left()) && !IsRed(root.Right()) {
		root = MoveRedRight(root)
	}

	if key == root.Key {
		minNode := SearchMin(root.Right())
		root.Key = minNode.Key
		root.Value = minNode.Value
		root.children[Right] = DeleteMin(root.Right())

		return FixUp(root)
	}

	root.children[Right] = Delete(root.Right(), key)

	return FixUp(root)
}

func DeleteMin[K cmp.Ordered, V any](root *Node[K, V]) *Node[K, V] {
	if root.Left() == nil {
		return nil
	}

	if !IsRed(root.Left()) && !IsRed(root.Left().Left()) {
		root = MoveRedLeft(root)
	}

	root.children[Left] = DeleteMin(root.Left())
	return FixUp(root)
}

// ------------------------------------------------------------------------------
// -- ROTATIONS
// ------------------------------------------------------------------------------

// Takes the root of a subtree, performs a rotation and returns the new root of the subtree.
// These operations are recursively performed top-down in llrb so we don't need to update the parents.
//
// We perform a rotation when the link between the root of the subtree and the pivot `x` is red, i.e
// the pivot is red.
//
// This rotation moves the red link from left to right or right to left from the perspective of the
// root of the subtree.
//
//	## Rotations
//
//	### INITIAL
//
//	      D
//	    /   \
//	  B       F
//	 / \     / \
//	A   C   E   G
//
//	### RIGHT ROTATION
//
// In the right rotation the pivot is the left child.
//
//	    B
//	  /   \
//	A       D
//	       / \
//	      C   F
//	         / \
//	        E   G
//
//	### LEFT ROTATION
//
// In the left rotation the pivot is the right child.
//
//		       F
//		     /   \
//		   D       G
//		  / \
//		 B   E
//		/ \
//	   A   C
func Rotate[K cmp.Ordered, V any](
	root *Node[K, V],
	direction Direction,
) *Node[K, V] {
	x := root.children[1-direction]
	root.children[1-direction] = x.children[direction]
	x.children[direction] = root

	// -- swap colors
	x.isBlack = root.isBlack
	root.isBlack = false

	return x
}

// ------------------------------------------------------------------------------
// -- LLRB HELPERS
// ------------------------------------------------------------------------------

// Fixes up the root of a subtree.
//
// If its right child is red, then we perform a left rotation to move the red link
// to the left.
//
// If the left subtree is made of 2 consecutive red links (i.e. left is red and
// left.left is also red): then we must fix it to ensure we have at most 1
// consecutive red link. By performing a right rotation, we move the first left
// red link to the right.
// This allows us in the next step to move both right and left red links up to the
// parent.
//
// Finally, if both children are red we flip colors, moving the red links up to
// the parent as shown in the figure entitled "Passing a red link up in a LLRB tree"
// on page 4 of the following paper:
// - https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
func FixUp[K cmp.Ordered, V any](root *Node[K, V]) *Node[K, V] {
	if IsRed(root.Right()) {
		root = Rotate(root, Left)
	}

	if IsRed(root.Left()) && IsRed(root.Left().Left()) {
		root = Rotate(root, Right)
	}

	if IsRed(root.Left()) && IsRed(root.Right()) {
		FlipColor(root)
	}

	return root
}

func FlipColor[K cmp.Ordered, V any](node *Node[K, V]) {
	node.isBlack = !node.isBlack

	if left := node.Left(); left != nil {
		left.isBlack = !left.isBlack
	}

	if right := node.Right(); right != nil {
		right.isBlack = !right.isBlack
	}
}

func IsRed[K cmp.Ordered, V any](node *Node[K, V]) bool {
	return node != nil && !node.isBlack
}

func MoveRedLeft[K cmp.Ordered, V any](root *Node[K, V]) *Node[K, V] {
	FlipColor(root)

	if IsRed(root.Right().Left()) {
		root.children[Right] = Rotate(root.Right(), Right)
		root = Rotate(root, Left)

		FlipColor(root)
	}

	return root
}

func MoveRedRight[K cmp.Ordered, V any](root *Node[K, V]) *Node[K, V] {
	FlipColor(root)

	if IsRed(root.Left().Left()) {
		root = Rotate(root, Right)
	}

	return root
}

// ------------------------------------------------------------------------------
// -- NODE HELPERS
// ------------------------------------------------------------------------------

type Color int

const (
	ColorBlack Color = iota
	ColorRed
)

func SetColor[K cmp.Ordered, V any](node *Node[K, V], color Color) {
	switch color {
	case ColorBlack:
		node.isBlack = true
	case ColorRed:
		node.isBlack = false
	default:
		panic("please")
	}
}
