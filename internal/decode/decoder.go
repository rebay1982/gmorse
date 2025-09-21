package decode

import (
	"fmt"
	"time"

	"honnef.co/go/tools/analysis/code"
)

type TreeNode struct {
	// TODO: Make this with Generics?
	char byte
	left *TreeNode
	right *TreeNode
}

type decoderConfig struct {
	root *TreeNode
	wpm int
	tolerace float64
}
type MorseDecoder struct {
	config decoderConfig

	currentNode *TreeNode

	decodeIn <-chan Detection
	decodeStop <-chan struct{}
}

type Detection struct {
	state bool
	duration time.Duration
}

func NewMorseDecoder(in <-chan Detection, done <-chan struct{}, wpm int) *MorseDecoder {
	config := decoderConfig{
		root: buildMorseDecodeTree(),
		wpm: wpm,
		tolerace: 0.4,
	}

	decoder := &MorseDecoder{
		config: config,
		currentNode: config.root,
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
	// If on, determine if dit or dah
	// If off, determine if space between dit and dah, or end character or end word.
	if d.state {
		if md.approxDitLength(d.duration) {
			if md.currentNode.left != nil {
				// Dit, go left.
				md.currentNode = md.currentNode.left
			} else {
				// Code doesn't exist.
				md.currentNode = md.config.root
				fmt.Print("?")
			}
		}	else if md.approxDahLength(d.duration) {
			if md.currentNode.right != nil {
				// Dah, go right.
				md.currentNode = md.currentNode.right
			} else {
				// Char doesn't exist.
				md.currentNode = md.config.root
				fmt.Print("?")
			}
		} else {
			// Scrap word and start over.
			md.currentNode = md.config.root
		}

	} else {
		if md.approxBetweenBeepLength(d.duration) {
			fmt.Print(".\033[1C")
		}	else if md.approxBetweenCharLength(d.duration) {
			fmt.Print(md.currentNode.char)	
			md.currentNode = md.config.root

		}	else if md.approxBetweenWordLength(d.duration) {
			fmt.Print(md.currentNode.char)	
			md.currentNode = md.config.root

		} else {
			fmt.Printf("%s??", string(md.currentNode.char))
			md.currentNode = md.config.root
		}
	}
}

func (md *MorseDecoder) approxDitLength(d time.Duration) bool {
	ditLength := float64(60000)/float64(50 * md.config.wpm)
	max := int64(ditLength + ditLength*md.config.tolerace)
	min := int64(ditLength - ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxDahLength(d time.Duration) bool {
	ditLength := float64(60000)/float64(50 * md.config.wpm)
	max := 3*int64(ditLength + ditLength*md.config.tolerace)
	min := 3*int64(ditLength - ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxBetweenBeepLength(d time.Duration) bool {
	ditLength := float64(60000)/float64(50 * md.config.wpm)
	max := int64(ditLength + ditLength*md.config.tolerace)
	min := int64(ditLength - ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxBetweenCharLength(d time.Duration) bool {
	ditLength := float64(60000)/float64(50 * md.config.wpm)
	max := 3*int64(ditLength + ditLength*md.config.tolerace)
	min := 3*int64(ditLength - ditLength*md.config.tolerace)

	dm := d.Milliseconds()
	return dm >= min && dm <= max
}

func (md *MorseDecoder) approxBetweenWordLength(d time.Duration) bool {
	ditLength := float64(60000)/float64(50 * md.config.wpm)
	max := 7*int64(ditLength + ditLength*md.config.tolerace)
	min := 7*int64(ditLength - ditLength*md.config.tolerace)

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

