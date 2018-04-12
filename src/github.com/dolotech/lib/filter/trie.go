// trie.go: use dictionary words to build FSM
//          with the help of Aho-Corasick algorithm
//          which is proficient in searching multiple string pattern in text
//          with small usage of memory and high speed.

package filter

import (
	"container/list"
	"fmt"
	"sort"
	"unicode/utf8"
)

// In a trie, each node has many child nodes
type ChildNodeType []*Node

// Implement the interface which is needed by STL sort function
func (c ChildNodeType) Len() int {
	return len(c)
}

// Implement the interface which is needed by STL sort function
func (c ChildNodeType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Implement the interface which is needed by STL sort function
func (c ChildNodeType) Less(i, j int) bool {
	return c[i].Val < c[j].Val
}

// Node used to build trie structure
type Node struct {
	Val        rune          // val from the parent to the node,or edge, utf-8 encode
	Depth      int           // depth of the node from root,root's depth is 0
	ParentNode *Node         // parent node used to trace back to the root
	ChildNodes ChildNodeType // child node slice
	SuffixNode *Node         // suffix node of the node's longest postfix, represented by root->node1->node2....->suffix
	EOW        bool          // end-of-word tag
}

// get child node by given val
func (N *Node) GetChildNodeByVal(val rune) *Node {
	childNodeNum := len(N.ChildNodes)
	for left := 0; left <= childNodeNum-1; left++ {
		if N.ChildNodes[left].Val == val {
			return N.ChildNodes[left]
		}
	}
	return nil
}

// binary search childnodes with given val
func (N *Node) BinGetChildNodeByVal(val rune) *Node {
	right := len(N.ChildNodes) - 1
	left := 0
	mid := 0
	var midnode *Node
	for left <= right {
		mid = (left + right) / 2
		midnode = N.ChildNodes[mid]
		if midnode.Val == val {
			return midnode
		} else if midnode.Val < val {
			left = mid + 1
		} else if midnode.Val > val {
			right = mid - 1
		}
	}
	return nil
}

// simplly insert a child node
func (N *Node) InsertChildNodeByVal(val rune) *Node {
	node := new(Node)
	node.Val = val
	node.Depth = N.Depth + 1
	node.ParentNode = N
	node.ChildNodes = nil
	node.SuffixNode = nil
	node.EOW = false
	N.ChildNodes = append(N.ChildNodes, node)
	//fmt.Printf("build %c---->%c\n", N.Val, node.Val)
	return node
}

// Trie
type Trie struct {
	RootNode *Node // root
}

func (T *Trie) InitRootNode() {
	node := new(Node)
	node.Val = 0
	node.Depth = 0
	node.ParentNode = nil
	node.ChildNodes = nil
	node.SuffixNode = nil
	node.EOW = false
	T.RootNode = node
}

// dump the whole trie which is rooted in node
func (T *Trie) DumpTrie(node *Node) {

	lst := new(list.List)
	lst.PushBack(node)

	for lst.Len() > 0 {

		node := lst.Remove(lst.Front()).(*Node)
		pnode := node.ParentNode
		snode := node.SuffixNode

		var adr *Node = nil
		var padr *Node = nil
		var sadr *Node = nil
		var cadr *Node = nil

		var val rune = 0
		var pval rune = 0
		var sval rune = 0
		var cval rune = 0

		val = node.Val
		adr = node

		if pnode != nil {
			pval = pnode.Val
			padr = pnode
		}

		if snode != nil {
			sval = snode.Val
			sadr = snode
		}

		fmt.Printf("adr:%p  val:%c  depth:%d  padr:%p  pval:%c  sadr:%p  sval:%c  eow:%v\n", adr, val, node.Depth, padr, pval, sadr, sval, node.EOW)

		for _, child := range node.ChildNodes {
			cadr = child
			cval = child.Val
			fmt.Printf("-------------->cadr:%p  cval:%c\n", cadr, cval)
			lst.PushBack(child)
		}
	}
}

// trace from a node to root, root and the node are both contained in return
func (T *Trie) TraceBackToRoot(node *Node) []*Node {
	depth := node.Depth
	nodes := make([]*Node, depth+1)
	nodes[0] = T.RootNode
	for tmpnode := node; tmpnode != nil; tmpnode = tmpnode.ParentNode {
		nodes[depth] = tmpnode
		depth--
	}
	return nodes
}

// find a node the path to which can be represented by value of nodes arr
func (T *Trie) FindNodeByPath(nodes []*Node) *Node { //nodes not contain root node
	tmpnode := T.RootNode
	for i, node := range nodes {
		tmpnode = tmpnode.GetChildNodeByVal(node.Val)
		if tmpnode == nil || tmpnode.Val != node.Val {
			return nil
		} else if i == len(nodes)-1 {
			return tmpnode
		}
	}
	return nil
}

// build trie by given dictionary,each line in dictionary is a word
func (T *Trie) BuildTrie(dictionary [][]byte) {
	for _, line := range dictionary {
		if len(line) <= 0 {
			continue
		}
		parent := T.RootNode //when we handle a word, we start from rootnode
		for len(line) > 0 {
			charactor, length := utf8.DecodeRune(line) //each time get a rune and its size in bytes
			if length <= 0 {
				break //maybe not handle whole line
			}
			child := parent.GetChildNodeByVal(charactor)
			if child == nil {
				child = parent.InsertChildNodeByVal(charactor)
			}
			parent = child
			line = line[length:]
		}
		if parent != T.RootNode {
			parent.EOW = true // if len(line)>0 and rightly handle at least one charactor, we tag the node as EOW
		}
	} // now a trie is built

	lst := new(list.List)
	lst.PushBack(T.RootNode) // start from root node , we find suffix node of each node

	for lst.Len() > 0 {
		node := lst.Remove(lst.Front()).(*Node)
		sort.Sort(node.ChildNodes)              //sort the child nodes to use binary search
		for _, child := range node.ChildNodes { //each time we get a child
			lst.PushBack(child)
			startDepth := 2            //only nodes whose depth>=2 have suffix node
			searchDepth := child.Depth //search from child's depth
			var suffixNode *Node
			if searchDepth >= startDepth { // only nodes whose depth is bigger or equal to 2 are considered
				pathToRoot := T.TraceBackToRoot(child) //pathToRoot is root->level1node->level2node->level3node......
				for startDepth <= searchDepth {
					suffixNode = T.FindNodeByPath(pathToRoot[startDepth : searchDepth+1]) //just start from level2node
					if suffixNode != nil {                                                //if find,we break,because we just care the longest postfix
						break
					}
					startDepth++ //each time we add startDepth the postfix is shortened
				}
			}
			if suffixNode != nil {
				child.SuffixNode = suffixNode // found
			} else {
				child.SuffixNode = T.RootNode // if not found or level1 nodes ,root is set, so we can ensure that every node has a suffix node
			}
		}
	}
	//T.DumpTrie(T.RootNode)
}
