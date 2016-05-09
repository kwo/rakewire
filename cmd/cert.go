// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/auth"
)

// CertGen generates certificates
func CertGen(c *cli.Context) error {

	host := c.String("host")
	rsaBits := c.Int("bits")
	ecdsaCurve := c.String("curve")
	tlsCertFile := c.String("tlscert")
	tlsKeyFile := c.String("tlskey")

	return auth.GenerateCertificates(host, rsaBits, ecdsaCurve, tlsCertFile, tlsKeyFile, c.App.Version)

}
