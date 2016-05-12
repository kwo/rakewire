package httpd

import (
	"github.com/kwo/rakewire/auth"
	"github.com/kwo/rakewire/model"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const (
	testUsername = "mrrobot"
	testPassword = "averybadpassword"
	testHostPort = "localhost:60606"
)

func TestSilk(t *testing.T) {

	t.SkipNow()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	certFile, keyFile := generateCerts(t)
	defer removeCerts(t, certFile, keyFile)

	cfg := &Configuration{
		DebugMode:      false,
		ListenHostPort: testHostPort,
		PublicHostPort: testHostPort,
		TLSCertFile:    certFile,
		TLSKeyFile:     keyFile,
	}

	server := NewService(cfg, db, "Rakewire", time.Now().Unix())
	if err := server.Start(); err != nil {
		t.Fatalf("Cannot start httpd: %s", err.Error())
	}
	defer server.Stop()

	// frisby.Global.Req = request.NewRequest(&http.Client{
	// 	Transport: &http.Transport{
	// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	},
	// })

	//url := fmt.Sprintf("https://%s", server.publicHostPort)

}

func generateCerts(t *testing.T) (certFile, keyFile string) {

	certFile = "certfile.crt"
	keyFile = "keyfile.key"

	if err := auth.GenerateCertificates("localhost", 2048, "", certFile, keyFile, "Rakewire"); err != nil {
		t.Fatalf("Cannot generate certificates: %s", err.Error())
	}

	return
}

func removeCerts(t *testing.T, certFile, keyFile string) {
	os.Remove(certFile)
	os.Remove(keyFile)
}

func openTestDatabase(t *testing.T) model.Database {

	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		t.Fatalf("Cannot acquire temp file: %s", err.Error())
	}
	f.Close()
	location := f.Name()

	store, err := model.Instance.Open(location)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err.Error())
	}

	err = store.Update(func(tx model.Transaction) error {
		return populateDatabase(tx)
	})
	if err != nil {
		t.Fatalf("Cannot populate database: %s", err.Error())
	}

	return store

}

func closeTestDatabase(t *testing.T, db model.Database) {

	location := db.Location()

	if err := model.Instance.Close(db); err != nil {
		t.Errorf("Cannot close database: %s", err.Error())
	}

	if err := os.Remove(location); err != nil {
		t.Errorf("Cannot remove temp file: %s", err.Error())
	}

}

func populateDatabase(tx model.Transaction) error {

	// add test user
	user := model.U.New(testUsername, testPassword)
	if err := model.U.Save(tx, user); err != nil {
		return err
	}

	return nil

}
