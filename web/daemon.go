package web

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

// InitializeDaemon will attempt to setup a TLSCertReloader if TLSCertPath is
// configured.
func (c *Component) InitializeDaemon() error {
	if c.TLSCertPath != "" {
		cr, err := NewTLSCertReloader(c.TLSCertPath, c.TLSKeyPath)
		if err != nil {
			return err
		}
		c.cr = cr
	}
	return nil
}

// Serve starts listening and serving HTTP and HTTPS if configured
func (c *Component) Serve() {
	// set up normal http server
	c.http = &http.Server{
		Addr:    c.ListenAddr,
		Handler: c,
	}

	// set up TLS http server if cert reloader was setup
	if c.cr != nil {
		c.https = &http.Server{
			Addr:    c.TLSAddr,
			Handler: c,
			TLSConfig: &tls.Config{
				GetCertificate: c.cr.GetCertificate,
			},
		}

		// start an hourly job to reload
		go func() {
			for {
				time.Sleep(1 * time.Hour)
				err := c.cr.Reload()
				if err != nil {
					c.Log.Debug("cert reload error:", err)
				}
			}
		}()
	}

	var wg sync.WaitGroup

	// start listening and serving for http
	wg.Add(1)
	go func() {
		c.Log.Info("http listening on", c.http.Addr)
		if err := c.http.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				c.Log.Error(err)
			}
		}
		wg.Done()
	}()

	// start listening and serving for http if setup
	if c.https != nil {
		wg.Add(1)
		go func() {
			c.Log.Info("https listening on", c.https.Addr)
			if err := c.https.ListenAndServeTLS("", ""); err != nil {
				if err != http.ErrServerClosed {
					c.Log.Error(err)
				}
			}
			wg.Done()
		}()
	}

	// block until all serving returns
	wg.Wait()
}

// Stop will shutdown any HTTP servers
func (c *Component) Stop() {
	if c.https != nil {
		c.https.Shutdown(context.Background())
	}
	c.http.Shutdown(context.Background())
}
