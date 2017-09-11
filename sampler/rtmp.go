package sampler

import (
	"github.com/nareix/joy4/format/rtmp"
	"github.com/nareix/joy4/av"
)

type RTMPServer struct {
	server *rtmp.Server
	Connections chan *rtmp.Conn
}

func CreateRTMPServer() *RTMPServer {
	rtmpServer := RTMPServer{Connections: make(chan *rtmp.Conn, 10)}
	rtmpServer.server = &rtmp.Server{}

	rtmpServer.server.HandlePublish = func(conn *rtmp.Conn) {
		println("Publish")

		rtmpServer.Connections <- conn
	}

	return &rtmpServer
}

func(server *RTMPServer) Demuxers() <-chan av.Demuxer {
	demuxers := make(chan av.Demuxer, 10)
	go func() {
		for con := range server.Connections {
			demuxers <- con
		}
	}()
	return demuxers
}

func(server *RTMPServer) Listen()  {
	server.server.ListenAndServe()
}
