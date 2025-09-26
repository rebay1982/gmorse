package decode

import (
	"testing"
	"time"
)

func Test_Decode(t *testing.T) {
	testCases := []struct {
		name       string
		detections []Detection
		exp        string
	}{
		{
			name: "single_character_E",
			detections: []Detection{
				{state: true, duration: 48 * time.Millisecond},
				{state: false, duration: 3 * 48 * time.Millisecond},
			},
			exp: "E",
		},
	}

	decodeIn := make(chan Detection)
	decodeOut := make(chan string)
	done := make(chan struct{})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoder := NewMorseDecoder(decodeIn, decodeOut, done, 25)
			decoder.StartDecode()

			output := ""
			go func() {
				for {
					select {
					case msg := <-decodeOut:
						output = output + msg
					case <-done:
						return
					}
				}
			}()

			for _, detection := range tc.detections {
				decodeIn <- detection
			}
			time.Sleep(50 * time.Millisecond)
			close(done)

			if output != tc.exp {
				t.Errorf("expecting %s, got %s", tc.exp, output)
			}
		})
	}
}
