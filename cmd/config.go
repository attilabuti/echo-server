package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

type configuration struct {
	server struct {
		host string
	}

	udp struct {
		enabled bool        // UDP echo enabled
		port    int         // UDP echo port
		address net.UDPAddr // UDP echo address
	}

	tcp struct {
		enabled bool        // TCP echo enabled
		port    int         // TCP echo port
		address net.TCPAddr // TCP echo address
	}

	http struct {
		enabled bool   // HTTP server enabled
		port    int    // HTTP server port
		address string // HTTP server address
	}

	https struct {
		enabled  bool   // HTTPS server enabled
		port     int    // HTTPS server port
		address  string // HTTPS server address
		cert     string // SSL certificate file
		key      string // RSA private key file
		autoCert bool   // Automatically generate SSL certificate
	}

	content struct {
		content        string // Response body
		file           string // Path to file which contains response body
		contentType    string // Content-Type header
		addContentType bool   // Add Content-Type header
	}

	log struct {
		enabled     bool   // Logging enabled
		dir         string // Log files directory
		requests    bool   // Log HTTP(S) requests
		connections bool   // Log TCP connections
		packets     bool   // Log incoming/outgoing packets
	}

	quiet bool // Quiet mode enabled
}

func (c *configuration) init() error {
	if !c.http.enabled && !c.https.enabled && !c.udp.enabled && !c.tcp.enabled {
		return errors.New("one of the following options must be enabled: http, https, tcp echo, udp echo")
	}

	if c.http.enabled {
		if !isValidPort(c.http.port) {
			return fmt.Errorf("invalid HTTP port number: %v", c.http.port)
		}

		c.http.address = net.JoinHostPort(c.server.host, strconv.Itoa(c.http.port))
	}

	if c.https.enabled {
		if !isValidPort(c.https.port) {
			return fmt.Errorf("invalid HTTPS port number: %v", c.https.port)
		}

		c.https.address = net.JoinHostPort(c.server.host, strconv.Itoa(c.https.port))

		if len(c.https.cert) == 0 && len(c.https.key) == 0 {
			if cert, key, err := generateCert(); err != nil {
				return err
			} else {
				c.https.autoCert = true
				c.https.cert = cert
				c.https.key = key
			}
		} else {
			if len(c.https.cert) == 0 {
				return errors.New("SSL certificate file must be specified")
			} else if !fileExists(c.https.cert) {
				return fmt.Errorf("SSL certificate file specified but not found: %s", c.https.cert)
			}

			if len(c.https.key) == 0 {
				return errors.New("RSA private key file must be specified")
			} else if !fileExists(c.https.key) {
				return fmt.Errorf("RSA private key file specified but not found: %s", c.https.key)
			}
		}
	}

	if c.tcp.enabled {
		if !isValidPort(c.tcp.port) {
			return fmt.Errorf("invalid TCP echo port number: %v", c.tcp.port)
		}

		c.tcp.address = net.TCPAddr{Port: c.tcp.port, IP: net.ParseIP(c.server.host)}
	}

	if c.udp.enabled {
		if !isValidPort(c.udp.port) {
			return fmt.Errorf("invalid UDP echo port number: %v", c.udp.port)
		}

		c.udp.address = net.UDPAddr{Port: c.udp.port, IP: net.ParseIP(c.server.host)}
	}

	if len(c.content.file) > 0 {
		if !fileExists(c.content.file) {
			return fmt.Errorf("content file specified but not found: %s", c.content.file)
		} else {
			content, err := os.ReadFile(c.content.file)
			if err != nil {
				return err
			}

			c.content.content = string(content)
		}
	}

	if len(c.content.contentType) > 0 {
		c.content.addContentType = true
	}

	return nil
}
