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
package llrb

import (
	"cmp"

	"github.com/alexandremahdhaoui/llrb/internal"
)

// ------------------------------------------------------------------------------
// -- TREE
//
// This is an implementation of the left-leaning Red-black Tree.
// https://sedgewick.io/wp-content/themes/sedgewick/papers/2008LLRB.pdf
// ------------------------------------------------------------------------------

type Tree[K cmp.Ordered, V any] struct {
	root *internal.Node[K, V]
}

func (t *Tree[K, V]) Search(key K) (V, bool) {
	return internal.Search(t.root, key)
}

func (t *Tree[K, V]) Insert(key K, value V) {
	t.root = internal.Insert(t.root, key, value)
	internal.SetColor(t.root, internal.ColorBlack)
}

func (t *Tree[K, V]) Delete(key K) {
	t.root = internal.Delete(t.root, key)
	internal.SetColor(t.root, internal.ColorBlack)
}
