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
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "E",
		},
		{
			name: "single_character_T",
			detections: []Detection{
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "T",
		},
		{
			name: "single_character_A", // Sequence, start short.
			detections: []Detection{
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "A",
		},
		{
			name: "single_character_N", // Sequence, start long.
			detections: []Detection{
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "N",
		},
		{
			name: "single_character_empty_node",
			detections: []Detection{
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "|?|",
		},
		{
			name: "single_character_invalid_sequence",
			detections: []Detection{
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "|?|",
		},
		{
			name: "single_character_?",
			detections: []Detection{
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "?",
		},
		{
			name: "full_word_SOS",
			detections: []Detection{
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: 3 * ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "SOS",
		},
		{
			name: "two_full_words_SOS SOS",
			detections: []Detection{
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: 3 * ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 7 * ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: 3 * ditLength},
				{State: false, Duration: 3 * ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: ditLength},
				{State: true, Duration: ditLength},
				{State: false, Duration: 3 * ditLength},
			},
			exp: "SOS SOS",
		},
		{
			name: "single_character_off_timing_positive_U",
			detections: []Detection{
				{State: true, Duration: time.Duration(float64(ditLength) * 1.2)},
				{State: false, Duration: time.Duration(float64(ditLength) * 1.2)},
				{State: true, Duration: time.Duration(float64(ditLength) * 1.2)},
				{State: false, Duration: time.Duration(float64(ditLength) * 1.2)},
				{State: true, Duration: time.Duration(float64(3*ditLength) * 1.2)},
				{State: false, Duration: time.Duration(float64(3*ditLength) * 1.2)},
			},
			exp: "U",
		},
		{
			name: "single_character_off_timing_negative_U",
			detections: []Detection{
				{State: true, Duration: time.Duration(float64(ditLength) * 0.8)},
				{State: false, Duration: time.Duration(float64(ditLength) * 0.8)},
				{State: true, Duration: time.Duration(float64(ditLength) * 0.8)},
				{State: false, Duration: time.Duration(float64(ditLength) * 0.8)},
				{State: true, Duration: time.Duration(float64(3*ditLength) * 0.8)},
				{State: false, Duration: time.Duration(float64(3*ditLength) * 0.8)},
			},
			exp: "U",
		},
	}

	config := DecoderConfig{
		Wpm:      25,
		Tolerace: 0.4,
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
