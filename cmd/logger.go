package cmd

import (
	"fmt"
	"io"
	_log "log"
	"net/http"
	"os"
	"path"
	"time"
)

type logger struct {
	file           *os.File
	error          *_log.Logger
	info           *_log.Logger
	requestLogger  *_log.Logger
	connLogger     *_log.Logger
	packetLogger   *_log.Logger
	requestEnabled bool
	connEnabled    bool
	packetEnabled  bool
}

func (l *logger) init() error {
	var output io.Writer

	if config.log.enabled {
		if !folderExists(config.log.dir) {
			if err := os.MkdirAll(config.log.dir, 0750); err != nil && !os.IsExist(err) {
				return err
			}
		}

		year, month, day := time.Now().Date()
		logFileName := fmt.Sprintf("%v_%v_%v.log", year, int(month), day)

		logFilePath := path.Join(config.log.dir, logFileName)
		if !fileExists(logFilePath) {
			if _, err := os.Create(logFilePath); err != nil {
				return err
			}
		}

		var err error
		l.file, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		if config.quiet {
			output = l.file
		} else {
			output = io.MultiWriter(os.Stdout, l.file)
		}
	} else if config.quiet {
		output = io.Discard
	} else {
		output = os.Stdout
	}

	l.requestEnabled = config.log.requests
	l.connEnabled = config.log.connections
	l.packetEnabled = config.log.packets

	flags := _log.Ldate | _log.Ltime | _log.Lmsgprefix
	l.error = _log.New(output, "[error] ", flags)
	l.info = _log.New(output, "[info] ", flags)
	l.requestLogger = _log.New(output, "[request] ", flags)
	l.connLogger = _log.New(output, "[connection] ", flags)
	l.packetLogger = _log.New(output, "[packet] ", flags)

	return nil
}

func (l *logger) close() {
	if err := l.file.Close(); err != nil {
		fmt.Printf("error while closing log file: %v\n", err)
	}
}

func (l *logger) request(next http.Handler) http.Handler {
	if l.requestEnabled {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			l.requestLogger.Printf("%s - [%s] %s %s", req.RemoteAddr, req.Method, req.Proto, req.URL)
			next.ServeHTTP(w, req)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
	})
}

func (l *logger) connection(open bool, addr string) {
	if !l.connEnabled {
		return
	}

	if open {
		l.connLogger.Printf("%s - new TCP connection", addr)
	} else {
		l.connLogger.Printf("%s - TCP connection closed", addr)
	}
}

func (l *logger) packet(op string, network string, bytes int, data []byte, addr string) {
	if !l.packetEnabled {
		return
	}

	if op == "read" {
		l.packetLogger.Printf(`[%s]  %s %s - L:%d | D:%v | T:"%s"`, op, network, addr, bytes, data, data)
	} else {
		l.packetLogger.Printf("[%s] %s %s - L:%d", op, network, addr, bytes)
	}
}
