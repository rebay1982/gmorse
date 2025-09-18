package decode

import "time"

type TreeNode struct {
	character string
	leftNode *TreeNode
	rightNode *TreeNode
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
	root := initTree()

	decoder := &MorseDecoder{
		decodeRoot: root,
		currentNode: root,
		decodeWpm: wpm,
		decodeIn: in,
		decodeStop: cancel,
	}

	return decoder
}

func initTree() *TreeNode {

	return nil
}

func (md *MorseDecoder) Decode() {

	for {
		

	}
}
