package decode

import (
	"fmt"
	"time"
)

// TODO: Make this with Generics?
type TreeNode struct {
	char  byte
	left  *TreeNode
	right *TreeNode
}

type DecoderConfig struct {
	wpm      int
	tolerace float64
}

type MorseDecoder struct {
	config DecoderConfig

	root     *TreeNode
	currentNode *TreeNode

	decodeIn   <-chan Detection
	decodeOut  chan<- string
	decodeStop <-chan struct{}
}

type Detection struct {
	state    bool
	duration time.Duration
}

func NewMorseDecoder(in <-chan Detection, out chan<- string, done <-chan struct{}, cfg DecoderConfig) *MorseDecoder {
	root := buildMorseDecodeTree()
	decoder := &MorseDecoder{
		config:      cfg,
		root:     	 root,
		currentNode: root,
		decodeIn:    in,
		decodeOut:   out,
		decodeStop:  done,
	}

	return decoder
}

// StartDecode Fire up morse decoding. This will listen to the channels provided upon creation until the done channel is
//
//	closed.
func (md *MorseDecoder) StartDecode() {

	// Reset
	md.currentNode = md.root

	go func() {
		for {
			select {
			case in := <-md.decodeIn:
				md.decode(in)
			case <-md.decodeStop:
				if md.currentNode != md.root {
					md.decodeOut <- string(md.currentNode.char)
				}
				close(md.decodeOut)
				return
			}
		}
	}()
}

func (md *MorseDecoder) decode(d Detection) {
	// If on, determine if dit or dah
	// If off, determine if space between dit and dah, or end character or end word.
	if d.state {
		if md.approxDitLength(d.duration) {
			if md.currentNode.left != nil {
				// Dit, go left.
				md.currentNode = md.currentNode.left

			} else {
				// Code doesn't exist.
				md.currentNode = md.root
				md.decodeOut <- "?"
			}

		} else if md.approxDahLength(d.duration) {
			if md.currentNode.right != nil {
				// Dah, go right.
				md.currentNode = md.currentNode.right

			} else {
				// Char doesn't exist.
				md.currentNode = md.root
				md.decodeOut <- "?"
			}

		} else {
			// Scrap word and start over.
			md.currentNode = md.root
		}

	} else {
		if md.approxBetweenBeepLength(d.duration) {
			// Do nothing, wait for next

		} else if md.approxBetweenCharLength(d.duration) {
			md.decodeOut <- string(md.currentNode.char)
			md.currentNode = md.root

		} else if md.approxBetweenWordLength(d.duration) {
			md.decodeOut <- fmt.Sprintf("%s ", string(md.currentNode.char))
			md.currentNode = md.root

		} else {
			// And transmission? assume so...
			md.decodeOut <- fmt.Sprintf("%s ", string(md.currentNode.char))
			md.currentNode = md.root
		}
	}
}

func (md *MorseDecoder) approxDitLength(d time.Duration) bool {
	ditLength := float64(60000) / float64(50*md.config.wpm)
	max := int64(ditLength + ditLength*md.config.tolerace)
	min := int64(ditLength - ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxDahLength(d time.Duration) bool {
	ditLength := float64(60000) / float64(50*md.config.wpm)
	max := 3 * int64(ditLength+ditLength*md.config.tolerace)
	min := 3 * int64(ditLength-ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxBetweenBeepLength(d time.Duration) bool {
	ditLength := float64(60000) / float64(50*md.config.wpm)
	max := int64(ditLength + ditLength*md.config.tolerace)
	min := int64(ditLength - ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxBetweenCharLength(d time.Duration) bool {
	ditLength := float64(60000) / float64(50*md.config.wpm)
	max := 3 * int64(ditLength+ditLength*md.config.tolerace)
	min := 3 * int64(ditLength-ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxBetweenWordLength(d time.Duration) bool {
	ditLength := float64(60000) / float64(50*md.config.wpm)
	max := 7 * int64(ditLength+ditLength*md.config.tolerace)
	min := 7 * int64(ditLength-ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
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
		for i := range len(code) {
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
