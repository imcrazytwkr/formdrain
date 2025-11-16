package utf8util_test

import (
	"bytes"
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/utf8util"
	"golang.org/x/text/encoding/unicode"
)

var utf8test = []byte("TEST")
var utf8bomtest = append([]byte{0xEF, 0xBB, 0xBF}, utf8test...)

var utf16betest []byte
var utf16letest []byte

func init() {
	var err error
	utf16betest, err = unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM).NewEncoder().Bytes(utf8test)
	if err != nil {
		panic(err)
	}

	utf16letest, err = unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewEncoder().Bytes(utf8test)
	if err != nil {
		panic(err)
	}
}

func TestBOM(t *testing.T) {
	result, err := utf8util.FixBytes(utf8test)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, utf8test) {
		t.Fail()
	}

	result, err = utf8util.FixBytes(utf8bomtest)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, utf8test) {
		t.Logf("UTF8 BOM: %v from %v", result, utf8bomtest)
		t.Fail()
	}

	result, err = utf8util.FixBytes(utf16betest)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, utf8test) {
		t.Logf("UTF16BE BOM: %#v from %#v", result, utf16betest)
		t.Fail()
	}

	result, err = utf8util.FixBytes(utf16letest)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(result, utf8test) {
		t.Logf("UTF16LE BOM: %#v from %#v", result, utf16letest)
		t.Fail()
	}
}
