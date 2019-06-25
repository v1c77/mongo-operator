package mongo

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"github.com/globalsign/mgo"
	"github.smartx.com/mongo-operator/pkg/utils"
	"net"
	"time"
)

var logger = utils.NewLogger("mongocluster.service.mongo")

// We specify a timeout to mgo.Dial, to prevent
// mongod failures hanging the reconcile.
const mgoDialTimeout = 60 * time.Second

type Client interface {
	DialDirect() (*mgo.Session, error)
}

type client struct {
	// addr holds the address of the MongoDB server
	addrs []string

	// MgoPort holds the port of the MongoDB server.
}

// ======================== mongo client ===============================

func NewClient(addrs ...string) Client {

	return &client{
		addrs: addrs,
	}
}

// Certs holds the certificates and keys required to make a secure
// SSL connection.
type Certs struct {
	// CACert holds the CA certificate. This must certify the private key that
	// was used to sign the server certificate.
	CACert *x509.Certificate
	// ServerCert holds the certificate that certifies the server's
	// private key.
	ServerCert *x509.Certificate
	// ServerKey holds the server's private key.
	ServerKey *rsa.PrivateKey
}

// MgoDialInfo returns a DialInfo suitable
// for dialling an MgoInstance at any of the
// given addresses, optionally using TLS.
func MgoDialInfo(certs *Certs, addrs ...string) *mgo.DialInfo {
	var dial func(addr net.Addr) (net.Conn, error)
	if certs != nil {
		pool := x509.NewCertPool()
		pool.AddCert(certs.CACert)
		tlsConfig := &tls.Config{
			RootCAs:    pool,
			ServerName: "anything",
		}
		dial = func(addr net.Addr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			if err != nil {
				logger.Errorf("tls.Dial(%s) failed with %v", addr, err)
				return nil, err
			}
			return conn, nil
		}
	} else {
		dial = func(addr net.Addr) (net.Conn, error) {
			conn, err := net.Dial("tcp", addr.String())
			if err != nil {
				logger.Errorf("net.Dial(%s) failed with %v", addr, err)
				return nil, err
			}
			return conn, nil
		}
	}
	return &mgo.DialInfo{Addrs: addrs, Dial: dial, Timeout: mgoDialTimeout}
}

func (c *client) DialDirect() (*mgo.Session, error) {
	dialInfo := MgoDialInfo(nil, c.addrs...)
	dialInfo.Direct = true

	return mgo.DialWithInfo(dialInfo)
}

func (c *client) Dial() (*mgo.Session, error) {
	dialInfo := MgoDialInfo(nil, c.addrs...)
	return mgo.DialWithInfo(dialInfo)
}
