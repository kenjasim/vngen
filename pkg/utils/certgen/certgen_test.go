package certgen_test

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/pkg/errors"
	certgen "nenvoy.com/pkg/utils/certgen"
	printing "nenvoy.com/pkg/utils/printing"
)

var (
	TestDir           = "/tmp/nenvoy/test/certgen"
	CACommonName      = "CA-CommonName-Test"
	CACertPath        = "/tmp/nenvoy/test/certgen/ca-test.crt"
	CAKeyPath         = "/tmp/nenvoy/test/certgen/ca-test.key"
	ClientSubjects    = []string{"Client-CommonName-Test", "Client-Organisation-Test"}
	ClientIPAddresses = []net.IP{net.ParseIP("10.0.0.1"), net.ParseIP("10.0.0.2"), net.ParseIP("10.0.0.3")}
	ClientDNSNames    = []string{"DNS-Name1-Test", "DNS-Name2-Test"}
	ClientCertPath    = "/tmp/nenvoy/test/certgen/client-test.crt"
	ClientKeyPath     = "/tmp/nenvoy/test/certgen/client-test.key"
)

// TestCertificateCreate
func TestCertificateCreate(t *testing.T) {

	// Create test directory
	err := os.MkdirAll(TestDir, os.ModePerm)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, fmt.Sprintf("failed to create test directory: %s", TestDir)))
	}

	//////////////////////////////////////////////////////////////////////
	// Create Certificate Authority
	//////////////////////////////////////////////////////////////////////

	err = certgen.CreateCA(CACertPath, CAKeyPath, CACommonName)
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Log(printing.SprintSuccess(fmt.Sprintf("CA certificate and private key written to: %s", TestDir)))

	//////////////////////////////////////////////////////////////////////
	// Create Client Certificates and Private Keys
	//////////////////////////////////////////////////////////////////////

	err = certgen.CreateCertKeyPair(ClientCertPath, ClientKeyPath, CACertPath, CAKeyPath, ClientSubjects[0], ClientSubjects[1], ClientIPAddresses, ClientDNSNames)
	if err != nil {
		t.Fatalf("%s", err)
	}

	// // Generate certificate signing request
	// clientCSR := certgen.GenCSR(ClientSubjects[0], ClientSubjects[1], ClientIPAddresses, ClientDNSNames)

	// // Gerenate the CA private key
	// clientPrivKey, err := certgen.GenPKey()
	// if err != nil {
	// 	t.Fatalf("%s", errors.Wrap(err, "failed on to create client private key:"))
	// }

	// // Re-import CA certificate and key
	// caCSR, caPrivKey, err := certgen.ReadCertFromFile(CACertPath, CAKeyPath)
	// if err != nil {
	// 	t.Fatalf("%s", errors.Wrap(err, "failed to import certificate or private key from :"+CACertPath))
	// }

	// // Sign client certificate signing request with CA private key
	// clientCSRPEM, clientPrivKeyPEM, err := certgen.GenCert(clientCSR, caCSR, clientPrivKey, caPrivKey)
	// if err != nil {
	// 	t.Fatalf("%s", errors.Wrap(err, "failed to sign client's csr:"))
	// }

	// // Write certificate and key to application cert directory
	// err = certgen.WriteCertToFile(*clientCSRPEM, *clientPrivKeyPEM, ClientCertPath, ClientKeyPath)
	// if err != nil {
	// 	t.Fatalf("%s", errors.Wrap(err, "failed to write client's certificate to:"+ClientCertPath))
	// }

	t.Log(printing.SprintSuccess(fmt.Sprintf("Client certificate and private key written to: %s", TestDir)))

	//////////////////////////////////////////////////////////////////////
	// Verify Certificate and key match
	//////////////////////////////////////////////////////////////////////

	err = certgen.VerifyCertKeyMatch(ClientCertPath, ClientKeyPath)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, "failed to match client cert and key pair:"))
	}

	t.Log(printing.SprintSuccess("Verified client certificate and key are pair"))

	//////////////////////////////////////////////////////////////////////
	// Verify Verify CA Issuer
	//////////////////////////////////////////////////////////////////////

	err = certgen.VerifySignatureFrom(ClientCertPath, ClientKeyPath, CACertPath, CAKeyPath)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, "failed to verify if client certificate was signed by CA:"))
	}

	t.Log(printing.SprintSuccess("Verified client certificates signed by CA"))

	//////////////////////////////////////////////////////////////////////
	// Verify Subject
	//////////////////////////////////////////////////////////////////////

	err = certgen.VerifySubjectNames(ClientCertPath, ClientKeyPath, ClientSubjects)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, "failed to verify subject name in client certificate:"))
	}

	t.Log(printing.SprintSuccess("Verified client subject name present in certificate"))

	//////////////////////////////////////////////////////////////////////
	// Verify DNS Name
	//////////////////////////////////////////////////////////////////////

	err = certgen.VerifyDNSNames(ClientCertPath, ClientKeyPath, ClientDNSNames)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, "failed to verify dns name in client certificate:"))
	}

	t.Log(printing.SprintSuccess("Verified client DNS names present in certificate"))

	//////////////////////////////////////////////////////////////////////
	// Verify IP Addresses
	//////////////////////////////////////////////////////////////////////

	err = certgen.VerifyIPAddresses(ClientCertPath, ClientKeyPath, ClientIPAddresses)
	if err != nil {
		t.Fatalf("%s", errors.Wrap(err, "failed to verify ip address in client certificate:"))
	}

	t.Log(printing.SprintSuccess("Verified client ip addresses present in certificate"))

}

// // TestWriteToMissingDirectory
// func TestWriteToMissingDirectory(t *testing.T) {

// 	// Create Certificate Signing Request for CA
// 	CSR := certgen.GenCACSR(CACommonName)

// 	// Gerenate the CA private key
// 	privKey, err := certgen.GenPKey()
// 	if err != nil {
// 		t.Fatalf("%s", errors.Wrap(err, "could not create private key:"))
// 	}

// 	// Gerenate the CA certificate (PEM encoded)
// 	caCertPEM, caPrivKeyPEM, err := certgen.GenCACert(CSR, privKey)
// 	if err != nil {
// 		t.Fatalf("%s", errors.Wrap(err, "could not create PEM encoded certs:"))
// 	}

// 	// Write PEM encoding to file
// 	certPath := "/tmp/nenvoy/test/xxxxxx/ca-test.crt"
// 	keyPath := "/tmp/nenvoy/test/xxxxxx/ca-test.crt"
// 	err = certgen.WriteCertToFile(*caCertPEM, *caPrivKeyPEM, certPath, keyPath)
// 	if err != nil {
// 		t.Fatalf("%s", errors.Wrap(err, "could not write certificates and(or) key to file:"))
// 	}

// }
