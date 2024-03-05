module github.com/pion/example-webrtc-applications/v3

go 1.19

require (
	github.com/notedit/janus-go v0.0.0-20210115013133-fdce1b146d0e
	github.com/pion/interceptor v0.1.25
	github.com/pion/rtcp v1.2.13
	github.com/pion/sdp/v3 v3.0.6
	github.com/pion/webrtc/v4 v4.0.0-beta.9
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/pion/datachannel v1.5.5 // indirect
	github.com/pion/dtls/v2 v2.2.10 // indirect
	github.com/pion/ice/v3 v3.0.3 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/mdns v0.0.12 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtp v1.8.3 // indirect
	github.com/pion/sctp v1.8.12 // indirect
	github.com/pion/srtp/v3 v3.0.1 // indirect
	github.com/pion/stun/v2 v2.0.0 // indirect
	github.com/pion/transport/v2 v2.2.4 // indirect
	github.com/pion/transport/v3 v3.0.1 // indirect
	github.com/pion/turn/v3 v3.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/pion/webrtc/v4 => ../webrtc

replace github.com/pion/ice/v3 => ../ice

replace github.com/pion/transport/v3 => ../transport

replace github.com/pion/srtp/v3 => ../srtp

replace github.com/pion/interceptor => ../interceptor
