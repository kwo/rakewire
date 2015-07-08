package db

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

// Marshal serialize FeedInfo object to bytes
func (z *FeedInfo) Marshal() ([]byte, error) {

	data, err := json.Marshal(z)
	if err != nil {
		return nil, err
	}

	return zip(data)

}

// Unmarshal serialize FeedInfo object to bytes
func (z *FeedInfo) Unmarshal(gzData []byte) error {

	data, err := unzip(gzData)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, z)

}

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
