module github.com/pion/example-webrtc-applications/v3

go 1.19

require (
	github.com/notedit/janus-go v0.0.0-20210115013133-fdce1b146d0e
	github.com/pion/interceptor v0.1.29
	github.com/pion/rtcp v1.2.14
	github.com/pion/sdp/v3 v3.0.9
	github.com/pion/webrtc/v4 v4.0.0-beta.19
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/pion/datachannel v1.5.6 // indirect
	github.com/pion/dtls/v2 v2.2.11 // indirect
	github.com/pion/ice/v3 v3.0.7 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/mdns/v2 v2.0.7 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtp v1.8.6 // indirect
	github.com/pion/sctp v1.8.16 // indirect
	github.com/pion/srtp/v3 v3.0.1 // indirect
	github.com/pion/stun/v2 v2.0.0 // indirect
	github.com/pion/transport/v2 v2.2.5 // indirect
	github.com/pion/transport/v3 v3.0.2 // indirect
	github.com/pion/turn/v3 v3.0.3 // indirect
	github.com/rs/xid v1.5.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/pion/transport/v3 => github.com/atoppi/transport/v3 v3.0.2-0.20240509121744-ab643294fc11

replace github.com/pion/interceptor => github.com/atoppi/interceptor v0.1.26-0.20240509124256-36894dc6fe23

replace github.com/pion/srtp/v3 => github.com/atoppi/srtp/v3 v3.0.2-0.20240509124413-3886575d7481

replace github.com/pion/ice/v3 => github.com/atoppi/ice/v3 v3.0.4-0.20240509123924-adf3d410d62d

replace github.com/pion/webrtc/v4 => github.com/atoppi/webrtc/v4 v4.0.0-beta.19.0.20240509124938-41d7be65e417
