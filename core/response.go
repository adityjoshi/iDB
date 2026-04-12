package core

func readSimpleString(data []byte) (string, int, error) {

	pos := 1

	for ; data[pos] != '\r'; pos++ {

	}

	return string(data[1:pos]), pos + 2, nil
}
