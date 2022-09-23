package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type appServer struct {
	http        http.Server
	https       http.Server
	udpConn     *net.UDPConn
	tcpListener *net.TCPListener
	udpClosed   bool
	tcpClosed   bool
	idle        chan struct{}
	errors      chan error
}

func (s *appServer) start() {
	s.errors = make(chan error)
	s.idle = make(chan struct{})

	if config.http.enabled || config.https.enabled {
		s.handleFunctions()
	}

	if config.http.enabled {
		s.http = http.Server{
			Addr:     config.http.address,
			ErrorLog: log.error,
		}

		go func() {
			log.info.Printf("HTTP server listening on %v\n", config.http.address)
			s.errors <- s.http.ListenAndServe()
		}()
	}

	if config.https.enabled {
		s.https = http.Server{
			Addr:     config.https.address,
			ErrorLog: log.error,
		}

		go func() {
			log.info.Printf("HTTPS server listening on %v\n", config.https.address)
			s.errors <- s.https.ListenAndServeTLS(config.https.cert, config.https.key)
		}()
	}

	if config.tcp.enabled {
		go func() {
			s.tcpEcho()
		}()
	}

	if config.udp.enabled {
		go func() {
			s.udpEcho()
		}()
	}

	sigint := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		s.stop()
	}()

	log.error.Println(<-s.errors)

	<-s.idle
}

func (s *appServer) stop() {
	if config.http.enabled {
		if err := s.http.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.error.Printf("HTTP server shutdown: %v\n", err)
		} else {
			log.info.Println("HTTP server shutdown")
		}
	}

	if config.https.enabled {
		if err := s.https.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.error.Printf("HTTPS server shutdown: %v\n", err)
		} else {
			log.info.Println("HTTPS server shutdown")
		}

		if config.https.autoCert {
			if err := os.Remove(config.https.cert); err != nil {
				log.error.Printf("error while removing cert file: %v\n", err)
			}

			if err := os.Remove(config.https.key); err != nil {
				log.error.Printf("error while removing key file: %v\n", err)
			}
		}
	}

	if config.udp.enabled {
		log.info.Println("UDP echo server shutdown")

		s.udpClosed = true
		if err := s.udpConn.Close(); err != nil {
			log.error.Printf("udpConn.Close() error: %s\n", err)
		}
	}

	if config.tcp.enabled {
		log.info.Println("TCP echo server shutdown")

		s.tcpClosed = true
		if err := s.tcpListener.Close(); err != nil {
			log.error.Printf("TCPListener.Close() error: %s\n", err)
		}
	}

	if config.log.enabled {
		log.close()
	}

	close(s.errors)
	close(s.idle)
}

func (s *appServer) handleFunctions() {
	http.Handle("/", log.request(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			if config.content.addContentType {
				w.Header().Set("Content-Type", config.content.contentType)
			}

			w.Write([]byte(config.content.content))
		}),
	))

	http.Handle("/headers", log.request(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if reqHeadersBytes, err := json.Marshal(req.Header); err != nil {
				log.error.Println("could not marshal request headers:", err)

				w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
			} else {
				w.Write([]byte(reqHeadersBytes))
			}
		}),
	))
}

func (s *appServer) tcpEcho() {
	var err error
	s.tcpListener, err = net.ListenTCP("tcp", &config.tcp.address)
	if err != nil {
		log.error.Printf("net.ListenTCP() error: %s\n", err)
		return
	}

	log.info.Printf("TCP echo server listening on %v\n", s.tcpListener.Addr().String())

	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				if s.tcpClosed {
					return
				}

				log.error.Fatalf("net.ReadFromUDP() error: %s\n", err)
			}

			log.error.Printf("TCPListener.Accept() error: %s\n", err)
		} else {
			go s.handleTCPConnection(conn)
		}
	}
}

func (s *appServer) handleTCPConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.connection(true, remoteAddr)

	defer conn.Close()
	defer log.connection(false, remoteAddr)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.error.Printf("net.Read() error: %s\n", err)
			return
		}

		if n == 0 {
			return
		}

		log.packet("read", "TCP", n, buf[:n], remoteAddr)

		wn, werr := conn.Write(buf[:n])
		if err != nil {
			log.error.Printf("net.Write() error: %s\n", werr)
		} else {
			log.packet("write", "TCP", wn, []byte{}, remoteAddr)
		}
	}
}

func (s *appServer) udpEcho() {
	var err error
	s.udpConn, err = net.ListenUDP("udp", &config.udp.address)
	if err != nil {
		log.error.Fatalf("net.ListenUDP() error: %s\n", err)
		return
	}

	log.info.Printf("UDP echo server listening on %v\n", s.udpConn.LocalAddr())

	buf := make([]byte, 4096)
	for {
		n, addr, err := s.udpConn.ReadFromUDP(buf)
		remoteAddr := addr.String()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				if s.udpClosed {
					return
				}

				log.error.Fatalf("net.ReadFromUDP() error: %s\n", err)
			}

			log.error.Printf("net.ReadFromUDP() error: %s\n", err)

			continue
		}

		if n > 0 {
			log.packet("read", "UDP", n, buf[:n], remoteAddr)

			wn, werr := s.udpConn.WriteTo(buf[:n], addr)
			if err != nil {
				log.error.Printf("net.WriteTo() error: %s\n", werr)
			} else {
				log.packet("write", "UDP", wn, nil, remoteAddr)
			}
		}
	}
}
