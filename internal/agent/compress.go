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
