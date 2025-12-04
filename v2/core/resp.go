package core

import (
	"errors"
	"fmt"
)

func readSimpleString(data []byte) (string, int, error) {

	pos := 1
	for ; data[pos] != '\r'; pos++ {
	}

	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	var size int = 0
	for ; data[pos] != '\r'; pos++ {
		size = size*10 + int(data[pos]-'0')
	}
	pos += 2
	return string(data[pos : pos+size]), pos + size + 2, nil
}

func readArray(data []byte) ([]any, int, error) {
	pos := 1
	var count int = 0
	for ; data[pos] != '\r'; pos++ {
		count = count*10 + int(data[pos]-'0')
	}
	pos += 2

	elems := make([]any, count)
	for i := 0; i < count; i++ {
		val, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = val
		pos += delta
	}
	return elems, pos, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}
    fmt.Println(value)
	ts := value.([]any)
	tokens := make([]string, len(ts))

	for i := range tokens {
		tokens[i] = ts[i].(string)
	}
	return tokens, nil
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("empty data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)

	case ':':
		return readInt64(data)

	case '-':
		return readError(data)
	case '*':
		return readArray(data)
	case '$':
		return readBulkString(data)
	}
	return nil, 0, nil
}
func Decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}

func Encode(value any,isSimple bool) []byte{
    switch v:=value.(type){
    case string:
        if isSimple{
            return []byte(fmt.Sprintf("+%s\r\n",v))
        }
        return []byte(fmt.Sprintf("$%d\r\n%s\r\n",len(v),v))
	case int64:
		return []byte(fmt.Sprintf(":%d\r\n",value))

	default:
		return RESP_NIL
    }
}