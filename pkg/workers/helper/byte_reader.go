package helper

import "bytes"

func BytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
