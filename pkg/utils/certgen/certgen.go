package certgen

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
	"nenvoy.com/pkg/utils/osystem"
)

// CreateCA - Create Kubernetes Certificate Authority
func CreateCA(CACertPath, CAKeyPath, commonName string) (err error) {

	// Create Certificate Signing Request for CA
	caCSR := GenCACSR(commonName)

	// Gerenate the CA private key
	caPrivKey, err := GenPKey()
	if err != nil {
		return errors.Wrap(err, "could not create private key:")
	}

	// Gerenate the CA PEM encoding
	caCertPEM, caPrivKeyPEM, err := GenCACert(caCSR, caPrivKey)
	if err != nil {
		return errors.Wrap(err, "could not create PEM encoded certs:")
	}

	// Write PEM encoding to file
	err = WriteCertToFile(*caCertPEM, *caPrivKeyPEM, CACertPath, CAKeyPath)
	if err != nil {
		return errors.Wrap(err, "could not write certificates and(or) key to file:")
	}

	return nil
}

// CreateCertKeyPair - Generate certifcate and key pair
func CreateCertKeyPair(certPath, keyPath, CACertPath, CAKeyPath, commonName, organisation string, ipAddrs []net.IP, dnsNames []string) (err error) {

	// Generate certificate signing request
	adminCSR := GenCSR(commonName, organisation, ipAddrs, dnsNames)

	// Gerenate the private key
	adminPrivKey, err := GenPKey()
	if err != nil {
		return errors.Wrap(err, "failed to create private key: ")
	}

	// Import CA certificate and key
	caCSR, caPrivKey, err := ReadCertFromFile(CACertPath, CAKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to import CA key or certificate: ")
	}

	// Sign certificate signing request with CA private key
	adminCertPEM, adminPrivKeyPEM, err := GenCert(adminCSR, caCSR, adminPrivKey, caPrivKey)
	if err != nil {
		return errors.Wrap(err, "failedto generate certificate: ")
	}

	// Write certificate and key to application cert directory
	err = WriteCertToFile(*adminCertPEM, *adminPrivKeyPEM, certPath, keyPath)
	if err != nil {
		return errors.Wrap(err, "failed to write certificate or key to file path: ")
	}

	return nil
}

// GenCACSR - Create Certificate Authority certificate signing request (csr)
func GenCACSR(commonName string) (caCSR *x509.Certificate) {

	caCSR = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	return caCSR
}

// GenPKey - Generate a RSA private key, default 4096
func GenPKey() (caPrivKey *rsa.PrivateKey, err error) {

	caPrivKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Wrap(err, "failed on GenPKey:")
	}

	return caPrivKey, nil
}

// GenCACert - Generate PEM encoded CA Certificate and Private Key from csr and private key
func GenCACert(caCSR *x509.Certificate, caPrivKey *rsa.PrivateKey) (caCertPEM *bytes.Buffer, caPrivKeyPEM *bytes.Buffer, err error) {

	caBytes, err := x509.CreateCertificate(rand.Reader, caCSR, caCSR, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed on GenCACert:")
	}

	// pem encode
	caCertPEM = new(bytes.Buffer)
	pem.Encode(caCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM = new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	return caCertPEM, caPrivKeyPEM, nil

}

// GenCSR - Create certificate signing request (csr)
func GenCSR(commonName string, organisation string, ipAddresses []net.IP, dnsNames []string) (csr *x509.Certificate) {

	// Add loopback addresses for local host
	ipAddresses = append(ipAddresses, net.IPv4(127, 0, 0, 1))
	ipAddresses = append(ipAddresses, net.IPv6loopback)

	keyUsage := x509.KeyUsageDigitalSignature
	keyUsage |= x509.KeyUsageKeyEncipherment

	var subject pkix.Name

	if organisation == "" {
		subject = pkix.Name{
			CommonName: commonName,
		}
	} else {
		subject = pkix.Name{
			CommonName:   commonName,
			Organization: []string{organisation},
		}
	}

	csr = &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject:      subject,
		IPAddresses:  ipAddresses,
		DNSNames:     dnsNames,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     keyUsage,
	}
	return csr
}

// GenCert - Generate PEM encoded Certificate and Private Key from csr and CA signature
func GenCert(csr *x509.Certificate, caCSR *x509.Certificate, privKey *rsa.PrivateKey, caPrivKey *rsa.PrivateKey) (certPEM *bytes.Buffer, privKeyPEM *bytes.Buffer, err error) {

	certBytes, err := x509.CreateCertificate(rand.Reader, csr, caCSR, &privKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed on GenCert:")
	}

	certPEM = new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	privKeyPEM = new(bytes.Buffer)
	pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	return certPEM, privKeyPEM, nil

}

// WriteCertToFile - Write PEM encoded certificate and key to a files at given path
func WriteCertToFile(certPEM bytes.Buffer, privKeyPEM bytes.Buffer, certPath string, keyPath string) (err error) {

	err = ioutil.WriteFile(certPath, certPEM.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write certificate to path %s:", certPath))
	}

	err = ioutil.WriteFile(keyPath, privKeyPEM.Bytes(), 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write private key to path %s:", keyPath))
	}

	return nil

}

// ReadCertFromFile - Retrieve Certificate and Key from file paths
func ReadCertFromFile(certPath string, keyPath string) (cert *x509.Certificate, privKey *rsa.PrivateKey, err error) {
	certPEM, err := ioutil.ReadFile(certPath)
	certBlock, rest := pem.Decode(certPEM)
	if rest != nil && len(rest) > 0 {
		err := errors.New("failed to read in full certificate file, excluded data: " + string(rest))
		return nil, nil, errors.Wrap(err, "failed on ReadCertFromFile (1):")
	}
	certBytes := certBlock.Bytes
	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed on ReadCertFromFile (2):")
	}
	privKeyPEM, err := ioutil.ReadFile(keyPath)
	marshalledPrivKeyBlock, rest := pem.Decode(privKeyPEM)
	privKeyBytes := marshalledPrivKeyBlock.Bytes
	privKey, err = x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed on ReadCertFromFile (3):")
	}

	return cert, privKey, nil
}

// VerifyAll - Verify certificate, key and CA
func VerifyAll(certPath, keyPath, CACertPath, CAKeyPath string, subjectNames, dnsNames []string, IPAddrs []net.IP) (err error) {

	err = VerifyCertKeyMatch(certPath, keyPath)
	if err != nil {
		return err
	}

	err = VerifySignatureFrom(certPath, keyPath, CACertPath, CAKeyPath)
	if err != nil {
		return err
	}

	err = VerifySubjectNames(certPath, keyPath, subjectNames)
	if err != nil {
		return err
	}

	err = VerifyDNSNames(certPath, keyPath, dnsNames)
	if err != nil {
		return err
	}

	err = VerifyIPAddresses(certPath, keyPath, IPAddrs)
	if err != nil {
		return err
	}

	return nil
}

// VerifyCertKeyMatch - Verify certificate's IP Address
func VerifyCertKeyMatch(certPath, keyPath string) (err error) {

	// Check files exist
	if c, err := osystem.PathExists(certPath); c && err != nil {
		return errors.New(fmt.Sprintf("Cert not found in path: %s", certPath))
	} else if c, err := osystem.PathExists(keyPath); c && err != nil {
		return errors.New(fmt.Sprintf("Key not found in path: %s", keyPath))
	}

	// Check certificate and key are pair
	_, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return errors.New("failed to match key and certificate")
	}

	return nil
}

// VerifySignatureFrom - Verify ability to authenticate a certificate with root
func VerifySignatureFrom(certPath, keyPath, CACertPath, CAKeyPath string) (err error) {

	// Check files exist
	if c, err := osystem.PathExists(certPath); c && err != nil {
		return errors.New(fmt.Sprintf("Cert not found in path: %s", certPath))
	} else if c, err := osystem.PathExists(keyPath); c && err != nil {
		return errors.New(fmt.Sprintf("Key not found in path: %s", keyPath))
	} else if c, err := osystem.PathExists(CACertPath); c && err != nil {
		return errors.New(fmt.Sprintf("CA Cert not found in path: %s", CACertPath))
	} else if c, err := osystem.PathExists(CAKeyPath); c && err != nil {
		return errors.New(fmt.Sprintf("CA Key not found in path: %s", CAKeyPath))
	}

	// Read in certificates and keys
	cert, _, err := ReadCertFromFile(certPath, keyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read cert from file: ")
	}
	CACert, _, err := ReadCertFromFile(CACertPath, CAKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read cert from file: ")
	}

	err = cert.CheckSignatureFrom(CACert)
	if err != nil {
		return err
	}

	return nil
}

// VerifySubjectNames - Verify certificate's subject
func VerifySubjectNames(certPath, keyPath string, subjectNames []string) (err error) {

	// Check files exist
	if c, err := osystem.PathExists(certPath); c && err != nil {
		return errors.New(fmt.Sprintf("Cert not found in path: %s", certPath))
	} else if c, err := osystem.PathExists(keyPath); c && err != nil {
		return errors.New(fmt.Sprintf("Key not found in path: %s", keyPath))
	}

	// Read in certificates and keys
	cert, _, err := ReadCertFromFile(certPath, keyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read cert from file: ")
	}

	// Compare Subject names
	for _, subjectName := range subjectNames {
		matched := false
		for _, n := range cert.Subject.Names {
			if n.Value.(string) == subjectName {
				matched = true
			}
		}
		if !matched {
			return errors.New(fmt.Sprintf("failed to match subject name %s in certificate", subjectName))
		}
	}
	return nil
}

// VerifyDNSNames - Verify certificate's DNS names
func VerifyDNSNames(certPath, keyPath string, dnsNames []string) (err error) {

	// Check files exist
	if c, err := osystem.PathExists(certPath); c && err != nil {
		return errors.New(fmt.Sprintf("Cert not found in path: %s", certPath))
	} else if c, err := osystem.PathExists(keyPath); c && err != nil {
		return errors.New(fmt.Sprintf("Key not found in path: %s", keyPath))
	}

	// Read in certificates and keys
	cert, _, err := ReadCertFromFile(certPath, keyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read cert from file: ")
	}

	// Compare dns names
	for _, dnsName := range dnsNames {
		matched := false
		for _, n := range cert.DNSNames {
			if n == dnsName {
				matched = true
			}
		}
		if !matched {
			return errors.New(fmt.Sprintf("failed to match dns name %s in certificate", dnsName))
		}
	}
	return nil
}

// VerifyIPAddresses - Verify certificate's IP Address
func VerifyIPAddresses(certPath, keyPath string, IPAddrs []net.IP) (err error) {

	// Check files exist
	if c, err := osystem.PathExists(certPath); c && err != nil {
		return errors.New(fmt.Sprintf("Cert not found in path: %s", certPath))
	} else if c, err := osystem.PathExists(keyPath); c && err != nil {
		return errors.New(fmt.Sprintf("Key not found in path: %s", keyPath))
	}

	// Read in certificates and keys
	cert, _, err := ReadCertFromFile(certPath, keyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read cert from file: ")
	}

	// Compare ip addresses names
	for _, ipAddr := range IPAddrs {
		matched := false
		for _, addr := range cert.IPAddresses {
			if addr.Equal(ipAddr) {
				matched = true
			}
		}
		if !matched {
			return errors.New(fmt.Sprintf("failed to match dns name %s in certificate", ipAddr))
		}
	}
	return nil
}

// RESOURCES
// https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
