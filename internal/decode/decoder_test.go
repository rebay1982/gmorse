package decode

import (
	"fmt"
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
			name: "single_character_A", // Sequence, start short.
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "A",
		},
		{
			name: "single_character_N", // Sequence, start long.
			detections: []Detection{
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "N",
		},
		{
			name: "single_character_empty_node",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "|?|",
		},
		{
			name: "single_character_invalid_sequence",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "|?|",
		},
		{
			name: "single_character_?",
			detections: []Detection{
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: 3 * ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: ditLength},
				{state: true, duration: ditLength},
				{state: false, duration: 3 * ditLength},
			},
			exp: "?",
		},
		{
			name: "full_word_SOS",
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
			name: "two_full_words_SOS SOS",
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
		{
			name: "single_character_off_timing_positive_U",
			detections: []Detection{
				{state: true, duration: time.Duration(float64(ditLength) * 1.2)},
				{state: false, duration: time.Duration(float64(ditLength) * 1.2)},
				{state: true, duration: time.Duration(float64(ditLength) * 1.2)},
				{state: false, duration: time.Duration(float64(ditLength) * 1.2)},
				{state: true, duration: time.Duration(float64(3*ditLength) * 1.2)},
				{state: false, duration: time.Duration(float64(3*ditLength) * 1.2)},
			},
			exp: "U",
		},
		{
			name: "single_character_off_timing_negative_U",
			detections: []Detection{
				{state: true, duration: time.Duration(float64(ditLength) * 0.8)},
				{state: false, duration: time.Duration(float64(ditLength) * 0.8)},
				{state: true, duration: time.Duration(float64(ditLength) * 0.8)},
				{state: false, duration: time.Duration(float64(ditLength) * 0.8)},
				{state: true, duration: time.Duration(float64(3*ditLength) * 0.8)},
				{state: false, duration: time.Duration(float64(3*ditLength) * 0.8)},
			},
			exp: "U",
		},
	}

	config := DecoderConfig{
		wpm:      25,
		tolerace: 0.4,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decodeOut := make(chan string)
			decodeIn := make(chan Detection)
			done := make(chan struct{})
			defer close(decodeIn)
			defer close(done)

			decoder := NewMorseDecoder(decodeIn, decodeOut, done, config)
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
			time.Sleep(5 * time.Millisecond)

			if output != tc.exp {
				fmt.Println([]byte(output))
				fmt.Println(len(output))
				t.Errorf("expecting [%s], got [%s]", tc.exp, output)
			}
		})
	}
}
