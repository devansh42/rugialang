package util

type Trie struct {
	children map[byte]*Trie
	terminal bool
}

func NewTrie() *Trie {
	t := new(Trie)
	t.children = make(map[byte]*Trie)
	return t
}

func (tr *Trie) Add(word []byte) {
	if len(word) > 0 {
		child, ok := tr.children[word[0]]
		if ok {
			child.Add(word[1:])
			return
		}
		t := NewTrie()
		tr.children[word[0]] = t
		t.Add(word[1:])

	} else {
		tr.terminal = true
	}

}

func Match(tr *Trie, prefix []byte) (perfectMatch bool, prefixMatch bool) {
	if len(prefix) > 0 {
		child, ok := tr.children[prefix[0]]
		if ok {
			return Match(child, prefix[1:])
		}
		return false, false
	}
	if len(tr.children) == 0 {
		// we came upto the terminal node
		// we got an perfect match or the trie is empty only
		return true, false
	} else {
		return tr.terminal, true // we have the word as a prefix in the tree
	}
}

// func (tr Trie) String() {
// 	var str strings.Builder
// 	for k, v := range tr.children {
// 		kb := []byte{k}
// 		str.WriteString("\n(")
// 		str.WriteString(string(kb))
// 		str.WriteString(":")
// 		str.WriteString(fmt.Sprintf("%s", *v))
// 	}
// }
