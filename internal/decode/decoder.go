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

	decodeIn <-chan Detection
	decodeStop <-chan struct{}
}

type Detection struct {
	state bool
	duration time.Duration
}

func NewMorseDecoder(in <-chan Detection, done <-chan struct{}, wpm int) *MorseDecoder {
	root := buildMorseDecodeTree()

	decoder := &MorseDecoder{
		decodeRoot: root,
		currentNode: root,
		decodeWpm: wpm,
		decodeIn: in,
		decodeStop: done,
	}

	return decoder
}

// StartDecode Fire up morse decoding. This will listen to the channels provided upon creation until the done channel is
//   closed.
func (md *MorseDecoder) StartDecode() {

	go func() {
		for {
			select {
			case in := <-md.decodeIn:
				md.decode(in)
				fmt.Println(in.state)
			case <-md.decodeStop:
				return
			}
		}
	}()
}

func (md *MorseDecoder) decode(d Detection) {

	wpm := md.decodeWpm

	// If on, determine if dit or dah
	// If off, determine if space between dit and dah, or end character or end word.

	if d.


}

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
		for i := range (len(code)) {
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

