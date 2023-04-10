package redis

import (
	"bufio"
	"io"
	"strconv"
)

type Decoder struct {
	reader *bufio.Reader
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{reader: bufio.NewReader(reader)}
}

func (dec *Decoder) Decode() (RespValue, int, error) {
	value, n, err := dec.readValue()
	return value, n, err

}

func (dec *Decoder) readValue() (value RespValue, n int, err error) {
	var respType byte
	respType, err = dec.reader.ReadByte()
	if err != nil {
		return RespValue{}, 0, err
	}
	n++
	switch respType {
	case '-', '+':
		value, n, err = dec.readSimple(respType)
	case ':':
		value, n, err = dec.readInteger(respType)
	case '$':
		value, n, err = dec.readBulkString(respType)
	case '*':
		value, n, err = dec.readArray(respType)
	}
	return value, n, err
}

type RespValue struct {
	respType byte
	_int     int
	str      []byte
	array    []RespValue
}

func (value *RespValue) String() string {
	return string(value.str)
}

func (value *RespValue) Bytes() []byte {
	var bytes []byte
	bytes = append(bytes, value.respType)
	bytes = append(bytes, []byte("\r\n")...)
	switch value.respType {
	case '-', '+', '$':
		bytes = append(bytes, value.str...)
	case ':':
		bytes = append(bytes, []byte(strconv.Itoa(value._int))...)
	case '*':
		for _, item := range value.array {
			bytes = append(bytes, item.Bytes()...)
		}

	}
	return bytes

}

func (value *RespValue) Array() []RespValue {
	return value.array
}

func (dec *Decoder) readSimple(respType byte) (RespValue, int, error) {
	var line []byte
	line, n, err := dec.readLine()
	if err != nil {
		return RespValue{}, n, err
	}

	return RespValue{respType: respType, str: line}, n, nil
}

func (dec *Decoder) readInteger(respType byte) (RespValue, int, error) {
	var line []byte
	line, n, err := dec.readLine()
	if err != nil {
		return RespValue{}, n, err
	}

	_int, err := strconv.Atoi(string(line))
	if err != nil {
		return RespValue{}, n, err
	}
	return RespValue{respType: respType, _int: _int}, n, nil
}

func (dec *Decoder) readBulkString(respType byte) (RespValue, int, error) {
	size, err := dec.readSize()
	if err != nil {
		return RespValue{}, 0, err
	}
	var lines []byte
	bytesRead := 0
	for bytesRead < size {
		line, n, err := dec.readLine()
		if err != nil {
			return RespValue{}, 0, err
		}
		lines = append(lines, line...)
		bytesRead += n
	}
	return RespValue{respType: respType, str: lines}, size, nil
}

func (dec *Decoder) readSize() (int, error) {
	bytesToRead, err := dec.reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	size, err := strconv.Atoi(string(bytesToRead[:len(bytesToRead)-2]))
	if err != nil {
		return 0, err
	}
	return size, nil

}

func (dec *Decoder) readArray(respType byte) (RespValue, int, error) {
	size, err := dec.readSize()
	if err != nil {
		return RespValue{}, 0, err
	}
	respValue := RespValue{respType: respType}
	var totalSize int
	for i := 0; i < size; i++ {
		val, n, err := dec.readValue()
		if err != nil {
			return RespValue{}, 0, err
		}
		respValue.array = append(respValue.array, val)
		totalSize += n

	}
	return respValue, totalSize, nil

}

func (dec *Decoder) readLine() (line []byte, n int, err error) {
	for {
		_byte, err := dec.reader.ReadBytes('\n')
		if err != nil {
			return nil, 0, err
		}
		n += len(_byte)
		line = append(line, _byte...)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}
