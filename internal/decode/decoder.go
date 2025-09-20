package decode

import (
	"fmt"
	"time"
)

type TreeNode struct {
	// TODO: Make this with Generics?
	char byte
	left *TreeNode
	right *TreeNode
}

type MorseDecoder struct {
	decodeRoot *TreeNode
	currentNode *TreeNode
	decodeWpm int

	decodeIn chan<-Detection
	decodeStop chan<-struct{}
}

type Detection struct {
	state bool
	duration time.Duration
}

func NewMorseDecoder(in chan<-Detection, cancel chan<-struct{}, wpm int) *MorseDecoder {
	root := buildMorseDecodeTree()

	decoder := &MorseDecoder{
		decodeRoot: root,
		currentNode: root,
		decodeWpm: wpm,
		decodeIn: in,
		decodeStop: cancel,
	}

	return decoder
}

var charsE = []byte{'I', 'A', 'S', 'U', 'H', 'V', 'F', ' '}

func buildMorseDecodeTree() *TreeNode {
	var morseTable = map[byte]string{
		'A': ".-",
		'B': "-...",
		'C': "-.-.",
		'D': "-..",
		'E': ".",
		'F': "..-.",
		'G': "--.",
		'H': "....",
		'I': "..",
		'J': ".---",
		'K': "-.-",
		'L': ".-..",
		'M': "--",
		'N': "-.",
		'O': "---",
		'P': ".--.",
		'Q': "--.-",
		'R': ".-.",
		'S': "...",
		'T': "-",
		'U': "..-",
		'V': "...-",
		'W': ".--",
		'X': "-..-",
		'Y': "-.--",
		'Z': "--..",
	}

	// Build the tree based on the morse table above.
	root := &TreeNode{}
	for letter, code := range morseTable {
		index := root
		for i := 0; i < len(code); i++ {
			c := code[i]
			switch c {
			case '.':
				if index.left == nil {
					index.left = &TreeNode{}
				}
				index = index.left
			case '-':
				if index.right == nil {
					index.right = &TreeNode{}
				}
				index = index.right
			}
		}
		index.char = letter
	}

	return root
}

func (md *MorseDecoder) Decode() {

	for {
		

	}
}
