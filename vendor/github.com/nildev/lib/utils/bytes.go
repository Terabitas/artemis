package utils

import "bytes"

// PopLine removes first line
func PopLine(in []byte) ([]byte, error) {
	r := bytes.NewBuffer(in)
	total := len(in)
	l, e := r.ReadString(10)
	if e != nil {
		return []byte{}, e
	}

	remaining := make([]byte, total-len([]byte(l)))

	r.Read(remaining)

	return remaining, nil
}
