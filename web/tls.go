package web

import (
	"crypto/tls"
	"fmt"
	"sync"
)

// TLSCertReloader is a TLS certificate provider that can reload certificates
// from disk. This is used when a TLS cert is configured, with Reload called
// once an hour. This means updating the cert and key on disk will start being
// used within an hour without needing to restart the server.
type TLSCertReloader struct {
	certPath string
	keyPath  string
	cert     *tls.Certificate
	sync.Mutex
}

// NewTLSCertReloader returns a new NewTLSCertReloader after an initial call
// to Reload, which could return an error.
func NewTLSCertReloader(certPath, keyPath string) (*TLSCertReloader, error) {
	cr := &TLSCertReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	err := cr.Reload()
	if err != nil {
		return nil, err
	}
	return cr, nil
}

// Reload performs a simple tls.LoadX509KeyPair with the given paths and atomically
// makes it available via the
func (cr *TLSCertReloader) Reload() error {
	c, err := tls.LoadX509KeyPair(cr.certPath, cr.keyPath)
	if err != nil {
		return err
	}
	cr.Lock()
	defer cr.Unlock()
	cr.cert = &c
	return nil
}

// GetCertificate is used with tls.Config to provide a tls.Certificate
func (cr *TLSCertReloader) GetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cr.Lock()
	defer cr.Unlock()
	if cr.cert == nil {
		return nil, fmt.Errorf("no certificate available")
	}
	return cr.cert, nil
}
