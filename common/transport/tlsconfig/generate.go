package tlsconfig

import (
	"archive/tar"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultKeyBits = 4096
	caValidity     = 10 * 365 * 24 * time.Hour // ~10 years
	leafValidity   = 5 * 365 * 24 * time.Hour  // ~5 years
)

// GenerateCA creates a self-signed CA certificate and private key in dir if they
// do not already exist. The resulting files are ca.pem and ca.key.
func GenerateCA(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	caCertPath := filepath.Join(dir, CACertFile)
	caKeyPath := filepath.Join(dir, CAKeyFile)

	if fileExists(caCertPath) && fileExists(caKeyPath) {
		return nil
	}

	key, err := generateKey()
	if err != nil {
		return fmt.Errorf("generate ca key: %w", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(0).SetBytes(mustRandomSerial()),
		Subject: pkix.Name{
			CommonName:         "beacon-ca",
			Organization:       []string{"Amuluze"},
			OrganizationalUnit: []string{"Beacon"},
			Country:            []string{"CN"},
		},
		NotBefore:             time.Now().Add(-1 * time.Minute),
		NotAfter:              time.Now().Add(caValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("create ca cert: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("parse ca cert: %w", err)
	}

	if err := writeFile(caCertPath, encodeCertToPEM(cert), 0o600); err != nil {
		return err
	}
	if err := writeFile(caKeyPath, encodePrivateKeyToPEM(key), 0o600); err != nil {
		return err
	}

	return nil
}

// GenerateLeafCert creates a leaf certificate signed by the CA in caDir and writes
// tls.crt and tls.key into dir. If dir is the same as caDir the leaf files are
// written alongside the CA files.
func GenerateLeafCert(dir, caDir, commonName string, dnsNames, ipAddresses []string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	caCert, caKey, err := loadCA(caDir)
	if err != nil {
		return err
	}

	key, err := generateKey()
	if err != nil {
		return fmt.Errorf("generate leaf key: %w", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(0).SetBytes(mustRandomSerial()),
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       []string{"Amuluze"},
			OrganizationalUnit: []string{"Beacon"},
			Country:            []string{"CN"},
		},
		NotBefore:    time.Now().Add(-1 * time.Minute),
		NotAfter:     time.Now().Add(leafValidity),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
		IPAddresses:  parseIPs(ipAddresses),
	}
	if len(tmpl.DNSNames) == 0 && len(tmpl.IPAddresses) == 0 {
		tmpl.DNSNames = []string{commonName}
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &key.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("create leaf cert: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return fmt.Errorf("parse leaf cert: %w", err)
	}

	certPath := filepath.Join(dir, TLSCertFile)
	keyPath := filepath.Join(dir, TLSKeyFile)

	if err := writeFile(certPath, encodeCertToPEM(cert), 0o600); err != nil {
		return err
	}
	if err := writeFile(keyPath, encodePrivateKeyToPEM(key), 0o600); err != nil {
		return err
	}

	return nil
}

// CreateCertPackage creates a tar.gz at outPath containing ca.pem, tls.crt and
// tls.key from the provided directories. caDir is the directory holding the CA,
// leafDir is the directory holding the leaf certificate.
func CreateCertPackage(caDir, leafDir, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	files := map[string]string{
		CACertFile:  filepath.Join(caDir, CACertFile),
		TLSCertFile: filepath.Join(leafDir, TLSCertFile),
		TLSKeyFile:  filepath.Join(leafDir, TLSKeyFile),
	}

	out, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create package: %w", err)
	}
	defer out.Close()

	zw := gzip.NewWriter(out)
	defer zw.Close()

	tw := tar.NewWriter(zw)
	defer tw.Close()

	for name, src := range files {
		data, err := os.ReadFile(src) //#nosec G304 -- src paths are constructed from admin-supplied directories
		if err != nil {
			return fmt.Errorf("read %s: %w", src, err)
		}
		hdr := &tar.Header{
			Name: name,
			Mode: 0o600,
			Size: int64(len(data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return fmt.Errorf("write tar header for %s: %w", name, err)
		}
		if _, err := tw.Write(data); err != nil {
			return fmt.Errorf("write tar body for %s: %w", name, err)
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}

	return nil
}

// CAKeyFile is the conventional private key filename for the CA.
const CAKeyFile = "ca.key"

func generateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, defaultKeyBits)
}

func encodePrivateKeyToPEM(key *rsa.PrivateKey) []byte {
	return pemEncode("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key))
}

func encodeCertToPEM(cert *x509.Certificate) []byte {
	return pemEncode("CERTIFICATE", cert.Raw)
}

func pemEncode(blockType string, data []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: data})
}

func loadCA(caDir string) (*x509.Certificate, *rsa.PrivateKey, error) {
	caCertPath := filepath.Join(caDir, CACertFile)
	caKeyPath := filepath.Join(caDir, CAKeyFile)

	caPEM, err := os.ReadFile(caCertPath) //#nosec G304 -- caDir is an admin-supplied configuration path
	if err != nil {
		return nil, nil, fmt.Errorf("read ca cert: %w", err)
	}
	block, _ := pem.Decode(caPEM)
	if block == nil {
		return nil, nil, fmt.Errorf("decode ca cert: no pem block")
	}
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse ca cert: %w", err)
	}

	keyPEM, err := os.ReadFile(caKeyPath) //#nosec G304 -- caDir is an admin-supplied configuration path
	if err != nil {
		return nil, nil, fmt.Errorf("read ca key: %w", err)
	}
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("decode ca key: no pem block")
	}
	caKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse ca key: %w", err)
	}

	return caCert, caKey, nil
}

func parseIPs(raw []string) []net.IP {
	out := make([]net.IP, 0, len(raw))
	for _, r := range raw {
		if ip := net.ParseIP(r); ip != nil {
			out = append(out, ip)
		}
	}
	return out
}

func mustRandomSerial() []byte {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("write %s: %w", filepath.Base(path), err)
	}
	return nil
}
