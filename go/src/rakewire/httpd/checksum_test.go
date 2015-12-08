package httpd

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestChecksum(t *testing.T) {

	hash := sha256.New()

	for name, file := range _escData {
		if file.IsDir() {
			continue
		}
		t.Logf("filename: %s", name)
		if data, err := FSByte(false, name); err == nil {
			if _, err := hash.Write(data); err != nil {
				t.Fatalf("Error hashing: %s", err.Error())
			}
		} else {
			t.Fatalf("Error getting file data: %s", err.Error())
		}
	} // files

	t.Logf("hash: %s", hex.EncodeToString(hash.Sum(nil)))

}
