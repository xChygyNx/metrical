package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func compress(data []byte) (res []byte, err error) {
	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error in gzip compress data: %w", err)
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("error in close gzipWriter: %w", err)
	}

	res = buf.Bytes()
	return
}

func decompress(data []byte) (res []byte, err error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error in create gzipReader in decompress: %w\n", err)
	}
	defer func() {
		err = reader.Close()
	}()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("error in read data from gzipReader in decompress: %w\n", err)
	}
	return buf.Bytes(), nil
}
