package trie

import "log"

type Trie struct {
	Root *Node
}

type Node struct {
	Children map[rune]*Node
	End bool
}

func NewTrie() *Trie {
	trie := &Trie{}
	trie.Root = NewTrieNode()
	return trie
}

func NewTrieNode() *Node {
	node := &Node{}
	node.Children = make(map[rune]*Node)
	node.End = false
	return node
}

func (t *Trie)Add (word string) {
	chars := []rune(word)
	if len(chars) == 0 {
		log.Fatal("illegal word")
		return
	}
	node := t.Root
	for _, char := range chars {
		if _, exists := node.Children[char]; !exists {
			node.Children[char] = NewTrieNode()
		}
		node = node.Children[char]
	}
	node.End = true
}

func (t *Trie) Replace(text string) (string, []string) {
	chars := []rune(text)
	result := []rune(text)
	found := make([]string, 0, 10)
	length := len(chars)
	node := t.Root
	for i := 0; i < length; i ++ {
		if _, exists := node.Children[chars[i]]; exists {
			node = node.Children[chars[i]]
			for j := i + 1; j < length; j ++ {
				if _, exists := node.Children[chars[j]]; !exists {
					break
				}
				node = node.Children[chars[j]]
				if node.End == true {
					for k := i; k <= j; k ++ {
						result[k] = ' '
					}
					found = append(found, string(chars[i:j+1]))
					i = j
					node = t.Root
					break
				}
			}
			node = t.Root
		}
	}
	return string(result), found
}



