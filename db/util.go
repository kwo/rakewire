package db

import (
	"bytes"
	"compress/gzip"
	"io"
)

func zip(data []byte) ([]byte, error) {

	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func unzip(data []byte) ([]byte, error) {

	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var uncompressedData bytes.Buffer
	_, err = io.Copy(&uncompressedData, r)
	if err != nil {
		return nil, err
	}

	return uncompressedData.Bytes(), nil

}
