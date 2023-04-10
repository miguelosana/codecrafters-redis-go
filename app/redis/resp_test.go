package redis

import (
	"strings"
	"testing"
)

func checkType(got, want RespValue, t *testing.T) {
	if got.respType != want.respType {
		t.Errorf("Got respType: %s, Want respType: %s", string(got.respType), string(want.respType))
	}

}
func TestDecodeSimpleString(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("+OK\r\n"))
	result, _, err := decoder.Decode()
	if err != nil {
		t.Error(err)
	}

	expectedResult := RespValue{respType: '+', str: []byte("OK")}
	if result.String() != expectedResult.String() {
		t.Errorf("Got str: %s, Want str: %s", result.String(), expectedResult.String())
	}

	checkType(result, expectedResult, t)
}

func TestDecodeInteger(t *testing.T) {
	decoder := NewDecoder(strings.NewReader(":1000\r\n"))
	result, _, err := decoder.Decode()
	if err != nil {
		t.Error(err)
	}

	expectedResult := RespValue{respType: ':', _int: 1000}
	if result._int != expectedResult._int {
		t.Errorf("Got _int: %d, Want _int: %d", result._int, expectedResult._int)
	}

	checkType(result, expectedResult, t)
}

func TestDecodeBulkString(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("$6\r\nhel\nlo\r\n"))
	result, _, err := decoder.Decode()
	if err != nil {
		t.Error(err)
	}

	expectedResult := RespValue{respType: '$', str: []byte("hel\nlo")}
	if result.String() != expectedResult.String() {
		t.Errorf("Got str: %s, Want str: %s", result.String(), expectedResult.String())
	}
	checkType(result, expectedResult, t)
}

func TestDecodeArray(t *testing.T) {
	decoder := NewDecoder(strings.NewReader("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
	result, _, err := decoder.Decode()
	if err != nil {
		t.Error(err)
	}
	expectedResult := RespValue{
		respType: '*',
		array: []RespValue{
			{
				respType: '$',
				str:      []byte("hello"),
			},
			{
				respType: '$',
				str:      []byte("world"),
			},
		},
	}
	if len(result.array) != 2 {
		t.Errorf("Got %d items, Want 2", len(result.array))
	}
	if result.array[0].String() != expectedResult.array[0].String() {
		t.Errorf("Got %s, Want %s", result.array[0].String(), expectedResult.array[0].String())
	}
	if result.array[1].String() != expectedResult.array[1].String() {
		t.Errorf("Got %s, Want %s", result.array[1].String(), expectedResult.array[1].String())
	}

}
