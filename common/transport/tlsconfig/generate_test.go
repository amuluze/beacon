package tlsconfig

import (
	"archive/tar"
	"compress/gzip"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCACreatesPEMFiles(t *testing.T) {
	dir := t.TempDir()
	if err := GenerateCA(dir); err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, CACertFile)); err != nil {
		t.Fatalf("ca.pem not created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, CAKeyFile)); err != nil {
		t.Fatalf("ca.key not created: %v", err)
	}
}

func TestGenerateCASkipsExisting(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, CACertFile), []byte("cert"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, CAKeyFile), []byte("key"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := GenerateCA(dir); err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, CACertFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "cert" {
		t.Fatal("existing ca.pem was overwritten")
	}
}

func TestGenerateCACreatesParsableCertificate(t *testing.T) {
	dir := t.TempDir()
	if err := GenerateCA(dir); err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, CACertFile))
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		t.Fatal("ca.pem is not valid PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse ca cert: %v", err)
	}
	if !cert.IsCA {
		t.Fatal("CA certificate IsCA is false")
	}
}

func TestGenerateLeafCert(t *testing.T) {
	caDir := t.TempDir()
	if err := GenerateCA(caDir); err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	leafDir := t.TempDir()
	if err := GenerateLeafCert(leafDir, caDir, "beacon/collia", []string{"localhost"}, []string{"127.0.0.1"}); err != nil {
		t.Fatalf("GenerateLeafCert failed: %v", err)
	}

	certPath := filepath.Join(leafDir, TLSCertFile)
	keyPath := filepath.Join(leafDir, TLSKeyFile)
	if _, err := os.Stat(certPath); err != nil {
		t.Fatalf("tls.crt not created: %v", err)
	}
	if _, err := os.Stat(keyPath); err != nil {
		t.Fatalf("tls.key not created: %v", err)
	}

	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatal("tls.crt is not valid PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse leaf cert: %v", err)
	}
	if cert.Subject.CommonName != "beacon/collia" {
		t.Fatalf("unexpected CN: %s", cert.Subject.CommonName)
	}
	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "localhost" {
		t.Fatalf("unexpected DNSNames: %v", cert.DNSNames)
	}
	if len(cert.IPAddresses) != 1 || cert.IPAddresses[0].String() != "127.0.0.1" {
		t.Fatalf("unexpected IPAddresses: %v", cert.IPAddresses)
	}
}

func TestGenerateLeafCertDefaultSAN(t *testing.T) {
	caDir := t.TempDir()
	if err := GenerateCA(caDir); err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	leafDir := t.TempDir()
	if err := GenerateLeafCert(leafDir, caDir, "beacon/collia", nil, nil); err != nil {
		t.Fatalf("GenerateLeafCert failed: %v", err)
	}

	certPEM, err := os.ReadFile(filepath.Join(leafDir, TLSCertFile))
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(certPEM)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse leaf cert: %v", err)
	}
	if len(cert.DNSNames) != 1 || cert.DNSNames[0] != "beacon/collia" {
		t.Fatalf("expected CN as default DNSName, got %v", cert.DNSNames)
	}
}

func TestCreateCertPackage(t *testing.T) {
	caDir := t.TempDir()
	if err := GenerateCA(caDir); err != nil {
		t.Fatalf("GenerateCA failed: %v", err)
	}

	leafDir := t.TempDir()
	if err := GenerateLeafCert(leafDir, caDir, "node-a", nil, nil); err != nil {
		t.Fatalf("GenerateLeafCert failed: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "node-a.tar.gz")
	if err := CreateCertPackage(caDir, leafDir, outPath); err != nil {
		t.Fatalf("CreateCertPackage failed: %v", err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("open gzip: %v", err)
	}
	defer zr.Close()

	tr := tar.NewReader(zr)
	found := make(map[string]bool)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read tar: %v", err)
		}
		found[hdr.Name] = true
	}

	for _, name := range []string{CACertFile, TLSCertFile, TLSKeyFile} {
		if !found[name] {
			t.Fatalf("missing %s in package", name)
		}
	}
}
