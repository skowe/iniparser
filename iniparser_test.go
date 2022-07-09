package iniparser_test

import (
	"bytes"
	"testing"

	"github.com/skowe/iniparser"
)

var expectedRaw = []byte(`[Block1]gi
;A comment line
key1=val1
key2=val2 ;A comment at the end of line
key3=valwith; a comment
key3=val with more words

[Block2]
;this = should't load
key1b2=123`)

var expectedParsed = `[Block1]
key1=val1
key2=val2
key3=valwith; a comment
key3=val with more words

[Block2]
key1b2=123`

var testini = iniparser.NewINI("./default.ini")

// Conversion to string is done in case of testing on windows

func TestRaw(t *testing.T) {

	if !bytes.Equal(testini.Raw, expectedRaw) {
		var badChPos int
		var badCh, expectedCh byte
		for i, val := range testini.Raw {
			if val != expectedRaw[i] {
				badCh = val
				badChPos = i
				expectedCh = expectedRaw[i]
				break
			}
		}
		t.Errorf("Raw data from file does not match expected values, got %d expected %d on position %d", badCh, expectedCh, badChPos)

	}
}

func TestParser(t *testing.T) {
	testini.Parse()
	s1 := string(testini.RawTrimmed)

	if !(s1 == expectedParsed) {
		if len(s1) == len(expectedParsed) {
			var badChPos int
			var badCh, expectedCh rune
			for i, val := range s1 {
				if val != rune(expectedParsed[i]) {
					badCh = val
					badChPos = i
					expectedCh = rune(expectedParsed[i])
				}
			}
			t.Errorf("Raw data from file does not match expected values, got %d expected %d on position %d", badCh, expectedCh, badChPos)
		} else {
			t.Errorf("Expected data and raw data are not the same length got %d, expected %d", len(testini.RawTrimmed), len(expectedParsed))
		}
	}
}
