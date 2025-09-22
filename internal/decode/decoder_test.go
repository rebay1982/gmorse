package decode

import (
	"fmt"
	"testing"
	"time"
)

func Test_Decode(t *testing.T) {
	decodeIn := make(chan Detection)
	decodeOut := make(chan string)
	done := make(chan struct{})
	expected := "ESA"

	decoder := NewMorseDecoder(decodeIn, decodeOut, done, 25)
	decoder.StartDecode()

	output := ""
	go func() {
		for {
			select {
			case msg := <-decodeOut:
				output = output + msg
			case <-done:
				fmt.Println("Done receiving")
				return
			}
		}
	}()

	decodeIn <- Detection{state: true, duration: 58 * time.Millisecond}
	decodeIn <- Detection{state: false, duration: 3 * 48 * time.Millisecond}
	decodeIn <- Detection{state: true, duration: 37 * time.Millisecond}
	decodeIn <- Detection{state: false, duration: 48 * time.Millisecond}
	decodeIn <- Detection{state: true, duration: 48 * time.Millisecond}
	decodeIn <- Detection{state: false, duration: 48 * time.Millisecond}
	decodeIn <- Detection{state: true, duration: 48 * time.Millisecond}
	decodeIn <- Detection{state: false, duration: 3 * 48 * time.Millisecond}
	decodeIn <- Detection{state: true, duration: 48 * time.Millisecond}
	decodeIn <- Detection{state: false, duration: 48 * time.Millisecond}
	decodeIn <- Detection{state: true, duration: 3 * 48 * time.Millisecond}
	decodeIn <- Detection{state: false, duration: 3 * 48 * time.Millisecond}
	time.Sleep(50 * time.Millisecond)
	close(done)

	if output != expected {
		t.Errorf("expecting %s, got %s", expected, output)
	}
}
