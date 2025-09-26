package decode

import (
	"testing"
	"time"
)

func Test_Decode(t *testing.T) {
	ditLength := 48 * time.Millisecond

	testCases := []struct {
		name       string
		detections []Detection
		exp        string
	}{
		{
			name: "single_character_E",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "E",
		},
		{
			name: "single_character_T",
			detections: []Detection{
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "T",
		},
		{
			name: "single_character_A",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "A",
		},
		{
			name: "single_character_N",
			detections: []Detection{
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "N",
		},
		{
			name: "single_character_SOS",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "SOS",
		},
		{
			name: "single_character_SOS SOS",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 7 * ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "SOS SOS",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wpm := 25
			decodeOut := make(chan string)
			decodeIn := make(chan Detection)
			done := make(chan struct{})
			defer close(decodeIn)
			defer close(done)

			decoder := NewMorseDecoder(decodeIn, decodeOut, done, wpm)
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
			time.Sleep(1 * time.Millisecond)

			if output != tc.exp {
				t.Errorf("expecting [%s], got [%s]", tc.exp, output)
			}
		})
	}
}
