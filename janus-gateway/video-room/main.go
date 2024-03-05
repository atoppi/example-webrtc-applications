// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

//go:build !js
// +build !js

// example of how to connect Pion and Janus
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	janus "github.com/notedit/janus-go"
	gst "github.com/pion/example-webrtc-applications/v3/internal/gstreamer-sink"
	"github.com/pion/interceptor"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v4"
)

func watchHandle(handle *janus.Handle) {
	// wait for event
	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			log.Println("Event: SlowLink")
		case *janus.MediaMsg:
			log.Printf("Event: Media %v receiving %v\n", msg.Type, msg.Receiving)
		case *janus.WebRTCUpMsg:
			log.Println("Event: WebRTC Up")
		case *janus.HangupMsg:
			log.Println("Event: Hangup")
		case *janus.EventMsg:
			log.Printf("Event Msg %+v", msg.Plugindata.Data)
		}
	}
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	flag.Usage = func() {
		fmt.Println("A pion-based Janus videoroom client (subscriber-only)")
		fmt.Print("\nUsage:\n")
		fmt.Println("video-room [--ws=ws://localhost:8188/janus] --room=1234 --feed=1000")
		flag.PrintDefaults()
	}
}

func main() {

	janusWs := flag.String("ws", "ws://localhost:8188/janus", "janus websocket endpoint")
	roomId := flag.Uint64("room", 0, "the room number client will join to")
	feedId := flag.Uint64("feed", 0, "the feed number client will subscribe to")
	enableStun := flag.Bool("enable-stun", false, "true to use Google STUN servers to discover srflx candidates")
	enableRfc8888 := flag.Bool("enable-rfc8888", false, "true to enable RFC8888 support")

	flag.Parse()

	if *roomId == 0 || *feedId == 0 {
		log.Fatalf("Missing room or feed identifier\n")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go gst.StartMainLoop()

	wg := &sync.WaitGroup{}

	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.
	log.Printf("Connecting to Janus WebSocket %v\n", *janusWs)
	gateway, err := janus.Connect(*janusWs)
	if err != nil {
		log.Fatalf("Error connecting to Janus (%v)\n", err)
	}
	log.Println("WebSocket connected")

	log.Println("Creating new session")
	session, err := gateway.Create()
	if err != nil {
		log.Fatalf("Error creating session (%v)\n", err)
	}
	log.Printf("Session created (%v)\n", session.ID)

	log.Println("Attaching new handle")
	handle, err := session.Attach("janus.plugin.videoroom")
	if err != nil {
		log.Fatalf("Error attaching to videoroom plugin (%v)\n", err)
	}
	log.Printf("Handle attached (%v)\n", handle.ID)

	go func() {
		for {
			if _, keepAliveErr := session.KeepAlive(); keepAliveErr != nil {
				log.Println("Error on keepalive", keepAliveErr)
				return
			}

			time.Sleep(5 * time.Second)
		}
	}()

	go watchHandle(handle)

	log.Printf("Subscribing to feed %v in room %v\n", *feedId, *roomId)
	msg, err := handle.Message(map[string]interface{}{
		"request": "join",
		"ptype":   "subscriber",
		"room":    *roomId,
		"feed":    *feedId,
	}, nil)
	if err != nil {
		log.Fatalf("Error subscribing to feed (%v)\n", err)
	}

	if msg.Jsep == nil {
		log.Fatalf("Missing offer from response\n")
	}

	sdpVal, ok := msg.Jsep["sdp"].(string)
	if !ok {
		log.Fatalf("Failed to get SDP offer(%v)\n", err)
	}
	log.Printf("Got Offer\n\n%v\n", strings.ReplaceAll(sdpVal, "\\r\\n", "\n"))

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  sdpVal,
	}

	mediaEngine := &webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		log.Fatalf("Error registering codecs (%v)\n", err)
	}

	interceptorRegistry := &interceptor.Registry{}
	webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry)
	// Register MID extensions for audio
	if strings.Contains(offer.SDP, "urn:ietf:params:rtp-hdrext:sdes:mid") {
		if err := mediaEngine.RegisterHeaderExtension(webrtc.RTPHeaderExtensionCapability{URI: sdp.SDESMidURI}, webrtc.RTPCodecTypeAudio); err != nil {
			log.Fatalf("Error registering rtp hdr extension (%v)\n", err)
		}
	}
	if *enableRfc8888 {
		log.Println("Using RFC8888")
		err := webrtc.ConfigureCongestionControlFeedback(mediaEngine, interceptorRegistry)
		if err != nil {
			log.Fatalf("Error registering RFC8888 interceptor (%v)\n", err)
		}
	}

	// Create a new RTCPeerConnection
	settingEngine := webrtc.SettingEngine{}
	settingEngine.SetDTLSRetransmissionInterval(100 * time.Millisecond)
	settingEngine.EnableEcnParsing(*enableRfc8888)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine), webrtc.WithInterceptorRegistry(interceptorRegistry), webrtc.WithSettingEngine(settingEngine))
	config := webrtc.Configuration{
		ICETransportPolicy: webrtc.ICETransportPolicyAll,
		ICEServers:         []webrtc.ICEServer{},
		SDPSemantics:       webrtc.SDPSemanticsUnifiedPlan,
	}
	if *enableStun {
		config.ICEServers = append(config.ICEServers, webrtc.ICEServer{URLs: []string{"stun:stun.l.google.com:19302"}})
	}
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		log.Fatalf("Error creating Peer Connection (%v)\n", err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("Connection ICE State has changed -> %s \n", connectionState.String())
	})

	peerConnection.OnConnectionStateChange(func(connectionState webrtc.PeerConnectionState) {
		log.Printf("Connection State has changed -> %s \n", connectionState.String())
	})

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		mime := track.Codec().RTPCodecCapability.MimeType
		log.Printf("Got track -> %v\n", mime)

		codecName := strings.Split(track.Codec().RTPCodecCapability.MimeType, "/")[1]
		pipeline := gst.CreatePipeline(track.PayloadType(), strings.ToLower(codecName))
		log.Printf("Starting gst pipeline")
		pipeline.Start()
		wg.Add(1)
		go func() {
			defer func() {
				log.Printf("Stopping gst pipeline")
				pipeline.Stop()
				wg.Done()
			}()
			buf := make([]byte, 1400)
			for {
				i, _, readErr := track.Read(buf)
				if readErr != nil {
					return
				}

				pipeline.Push(buf[:i])
			}
		}()

	})

	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	}); err != nil {
		log.Fatalf("Error adding audio transceiver (%v)\n", err)
	} else if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	}); err != nil {
		log.Fatalf("Error adding video transceiver (%v)\n", err)
	}

	log.Println("Setting remote description")
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		log.Fatalf("Error setting remote description (%v)\n", err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	log.Println("Creating answer")
	answer, answerErr := peerConnection.CreateAnswer(nil)
	if answerErr != nil {
		log.Fatalf("Error creating answer (%v)\n", err)
	}

	log.Println("Setting local description")
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		log.Fatalf("Errorsetting local description (%v)\n", err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	log.Println("Waiting for candidate gathering")
	<-gatherComplete

	log.Printf("Sending answer\n\n%v\n", strings.ReplaceAll(peerConnection.LocalDescription().SDP, "\\r\\n", "\n"))
	// now we start
	_, err = handle.Message(map[string]interface{}{
		"request": "start",
		"room":    *roomId,
	}, map[string]interface{}{
		"type":    "answer",
		"sdp":     peerConnection.LocalDescription().SDP,
		"trickle": false,
	})
	if err != nil {
		log.Fatalf("Error received for start request (%v)\n", err)
	}

	select {
	case <-gateway.GetErrChan():
		log.Println("Connection error")
		if peerConnection != nil {
			peerConnection.Close()
		}
		wg.Wait()
		return
	case <-interrupt:
		log.Println("Intercepted interrupt signal")
		if peerConnection != nil {
			peerConnection.Close()
		}
		wg.Wait()
		return
	}
}
