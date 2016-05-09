package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// GenerateCertificates generates certificates
func GenerateCertificates(hostlist string, rsaBits int, ecdsaCurve, tlsCertFile, tlsKeyFile, organizationName string) error {

	var priv interface{}
	var err error
	switch ecdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		fmt.Printf("Unrecognized elliptic curve: %s\n", ecdsaCurve)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("failed to generate private key: %s\n", err)
		os.Exit(1)
	}

	notBefore := time.Now().Truncate(time.Hour)
	notAfter := notBefore.Add(90 * 24 * time.Hour) // cert expires after 90 days

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		fmt.Printf("failed to generate serial number: %s\n", err)
		os.Exit(1)
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		Subject: pkix.Name{
			Organization: []string{organizationName},
		},
	}

	hosts := strings.Split(hostlist, ",")
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		fmt.Printf("Failed to create certificate: %s\n", err)
		os.Exit(1)
	}

	statCertFile, errCertFile := os.Stat(tlsCertFile)
	statKeyFile, errKeyFile := os.Stat(tlsKeyFile)
	exitFlag := false
	if statCertFile != nil && errCertFile == nil {
		fmt.Printf("File already exists, will not overwrite: %s\n", tlsCertFile)
		exitFlag = true
	}
	if statKeyFile != nil && errKeyFile == nil {
		fmt.Printf("File already exists, will not overwrite: %s\n", tlsKeyFile)
		exitFlag = true
	}
	if exitFlag {
		os.Exit(1)
	}

	certOut, err := os.OpenFile(tlsCertFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("failed to open %s for writing: %s\n", tlsCertFile, err)
		os.Exit(1)
	}

	keyOut, err := os.OpenFile(tlsKeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Failed to open %s for writing: %s\n", tlsKeyFile, err)
		os.Exit(1)
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	fmt.Printf("TLS certificate saved to %s\n", tlsCertFile)

	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	fmt.Printf("TLS key saved to %s\n", tlsKeyFile)

	return nil

}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
