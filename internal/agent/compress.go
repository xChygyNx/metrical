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
		err = fmt.Errorf("error in gzip compress data: %w", err)
	}

	res = buf.Bytes()
	return
}

//func decompress(data []byte) (res []byte, err error) {
//	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
//	if err != nil {
//		return nil, fmt.Errorf("error in create gzip reader: %w", err)
//	}
//	defer func() {
//		err = gzipReader.Close()
//	}()
//
//	var buf bytes.Buffer
//	_, err = buf.ReadFrom(gzipReader)
//	if err != nil {
//		return nil, fmt.Errorf("error in decompress data: %w", err)
//	}
//	return buf.Bytes(), nil
//}
