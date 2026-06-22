package tlsconfig

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"
)

func TestCertificateMatchesName(t *testing.T) {
	cert := &x509.Certificate{
		Subject:     pkix.Name{CommonName: "amprobe"},
		DNSNames:    []string{"collia-host-a"},
		IPAddresses: []net.IP{net.ParseIP("10.0.0.11")},
	}

	tests := []struct {
		name string
		want bool
	}{
		{name: "amprobe", want: true},
		{name: "collia-host-a", want: true},
		{name: "10.0.0.11", want: true},
		{name: "collia-host-b", want: false},
	}

	for _, tt := range tests {
		if got := certificateMatchesName(cert, tt.name); got != tt.want {
			t.Fatalf("certificateMatchesName(%q) = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCertificateMatchesAnyName(t *testing.T) {
	cert := &x509.Certificate{Subject: pkix.Name{CommonName: "amprobe"}}

	if !certificateMatchesAnyName(cert, []string{"other", "amprobe"}) {
		t.Fatal("expected certificate to match one allowed name")
	}
	if certificateMatchesAnyName(cert, []string{"other"}) {
		t.Fatal("expected certificate to reject unknown names")
	}
}
