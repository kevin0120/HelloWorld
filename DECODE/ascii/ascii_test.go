package ascii

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TEST_STRINGS = "123" + "3.14" + "7410" + "15.3" + "-897" + "99" + "1" + "0" + "sn001" + "666"
)

type Header struct {
	TOOL string `start:"1"  end:"5"`
	Sn   int    `start:"6"  end:"8"`
}

type OpenProtocol struct {
	L      int64   `start:"1"  end:"3"`
	LEN    float32 `start:"4"  end:"7"`
	MID    string  `start:"8"  end:"11"`
	MD     float64 `start:"12" end:"15"`
	M      int     `start:"16" end:"19"`
	Faa    uint    `start:"20" end:"21"`
	B      bool    `start:"22" end:"22"`
	C      bool    `start:"23" end:"23"`
	Header `start:"24" end:"..."`
}

func Test_Ascii(t *testing.T) {
	//var he =Header{}
	var assertdata = OpenProtocol{
		L:      123,
		LEN:    3.14,
		MID:    "7410",
		MD:     15.3,
		M:      -897,
		Faa:    99,
		B:      true,
		C:      false,
		Header: Header{TOOL: "sn001", Sn: 666},
	}
	fmt.Println(fmt.Sprintf("Test Data: %# 20X", TEST_STRINGS))
	var testop OpenProtocol
	err := Unmarshal(TEST_STRINGS, &testop)
	assert.Nil(t, err)
	//fmt.Printf("%+v\n", testop)
	if assert.NotNil(t, testop) {
		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.Equal(t, assertdata, testop)
	}
	//	fmt.Println(assertdata)
}
