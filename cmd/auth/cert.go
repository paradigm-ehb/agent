package auth

import (
	"crypto/rand"
	"crypto/rsa"

	/**
	What is x509
	X509 is a Public Key Infrastructure(PKI) standard defined by the ITU-T
	It specifies the format of public key certificates and plays a vital role in securing communications
	*/
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

const (
	caKeyFile     = "ca.key"
	caCertFile    = "ca.crt"
	agentKeyFile  = "agent.key"
	agentCertFile = "agent.crt"
)

func initTls() {

}

func GenerateCerijiijeeiijdlfjlst() (*rsa.PrivateKey, error) {

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key")
	}
	return privKey, nil
}

func LoadCa() {

}
