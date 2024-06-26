package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `jsont:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("Did not find seperator")
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	// Atoi - ascii character to integer
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	// until we find the seperator we want to keep reading the data from stdin
	// So if we haven't found it yet, we reutrn nil, fo the token so the splitfn knows to keep reading
	if !found {
		return 0, nil, nil
	}
	// Content-Length: <number>
	contentLengthBytes := header[len("Content-Length: "):]
	// Atoi - ascii character to integer
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err
	}
	if len(content) < contentLength {
		return 0, nil, nil
	}
	// 4 accounts for our \r\n\r\n
	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
