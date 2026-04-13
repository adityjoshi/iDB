package core

import "errors"

/*
*here we are returning the string, delta and the error
delta will tell the size so pos + 2 is pos -> length from 0 till the pos + \r\n
*
*/
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

	var val int64 = 0

	for ; data[pos] != '\r'; pos++ {
		val = val*10 + int64(data[pos]-'0')
	}

	return val, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {

	pos := 1

	len, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {

	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)

	case '-':
		return readError(data)

	case ':':
		return readInt64(data)
	}

	return nil, 0, nil

}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	val, _, err := DecodeOne(data)
	return val, err
}
