package main

import "github.com/quic-go/quic-go"


type conn struct{
    quicConn  quic.Connection
	// transport *transport

	localPeer      int
	localMultiaddr int

	remotePeerID    int
	remoteMultiaddr int

}


type stream struct {
	quic.Stream
}


func (s *stream) Read(b []byte) (n int, err error) {
	n, err = s.Stream.Read(b)
	if err != nil {

	}
	return n, err
}

func (s *stream) Write(b []byte) (n int, err error) {
	n, err = s.Stream.Write(b)
	if err != nil  {
	}
	return n, err
}