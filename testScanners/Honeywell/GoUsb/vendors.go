package GoUsb

import (
	"errors"
	"github.com/bep/debounce"
	"time"
)

const (
	KeyShift    = 2
	IndexTarget = 2
	IndexShift  = 0
)

var keymap = map[byte][]string{
	4:  {"a", "A"},
	5:  {"b", "B"},
	6:  {"c", "C"},
	7:  {"d", "D"},
	8:  {"e", "E"},
	9:  {"f", "F"},
	10: {"g", "G"},
	11: {"h", "H"},
	12: {"i", "I"},
	13: {"j", "J"},
	14: {"k", "K"},
	15: {"l", "L"},
	16: {"m", "M"},
	17: {"n", "N"},
	18: {"o", "O"},
	19: {"p", "P"},
	20: {"q", "Q"},
	21: {"r", "R"},
	22: {"s", "S"},
	23: {"t", "T"},
	24: {"u", "U"},
	25: {"v", "V"},
	26: {"w", "W"},
	27: {"x", "X"},
	28: {"y", "Y"},
	29: {"z", "Z"},
	30: {"1", "!"},
	31: {"2", "@"},
	32: {"3", "#"},
	33: {"4", "$"},
	34: {"5", "%"},
	35: {"6", "^"},
	36: {"7", "&"},
	37: {"8", "*"},
	38: {"9", "("},
	39: {"0", ")"},
	51: {";", ":"},
	53: {"`", "~"},
	45: {"-", "_"},
	46: {"=", "+"},
	47: {"[", "{"},
	48: {"]", "}"},
	52: {"'", "\""},
	54: {",", "<"},
	55: {".", ">"},
	56: {"/", "?"},
}


func CommonParse(buf []byte) (string, error) {
	v := buf[IndexTarget]
	if v == 0 || keymap[v] == nil {
		return "", errors.New("invalid byte")
	}

	str := keymap[v][0]
	if buf[IndexShift] == KeyShift {
		str = keymap[v][1]
	}

	return str, nil
}



func (s1 *GoUsb) ResetDebounce() {
	if s1.init {
		s1.init = false
	}

	s1.debounceTrigger = false
}

func (s1 *GoUsb) TriggerDebounce() {
	if !s1.debounceTrigger {
		//debInit, debCommon := s.devInfo.Debounce()
		//deb := debCommon
		//if s.init {
		//	deb = debInit
		//}

		s1.Debounced = debounce.New(100 * time.Millisecond)
		s1.debounceTrigger = true
	}
}
