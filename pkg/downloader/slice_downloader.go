package downloader

import (
	"context"
	"log"
	"sync"
	"time"

	pubUserFile "github.com/city404/v6-public-rpc-proto/go/v6/userfile"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v4"
)

type SliceDownloader struct {
	// authService       *auth.AuthService
	candidatesMux     sync.Mutex
	pendingCandidates []string
	rtcConfiguration  webrtc.Configuration
	peerConnection    *webrtc.PeerConnection
	clientID          string
	fileClient        pubUserFile.PubUserFileClient
}

func NewSliceDownloader(fileClient pubUserFile.PubUserFileClient) *SliceDownloader {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	return &SliceDownloader{
		// authService:       authService,
		pendingCandidates: make([]string, 0),
		rtcConfiguration:  config,
		clientID:          uuid.NewString(),
		fileClient:        fileClient,
	}
}

func (s *SliceDownloader) StartDownload(fileID string) error {
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(s.rtcConfiguration)
	if err != nil {
		return err
	}

	s.peerConnection = peerConnection

	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		log.Printf("CreateDataChannel error: %v", err)
		return err
	}
	checking := true

	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		if s == webrtc.PeerConnectionStateConnected || s == webrtc.PeerConnectionStateDisconnected || s == webrtc.PeerConnectionStateFailed || s == webrtc.PeerConnectionStateClosed {
			checking = false
		}
		log.Printf("Peer Connection State has changed: %s\n", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			log.Println("Peer Connection has gone to failed exiting")
			// os.Exit(0)
		}

		if s == webrtc.PeerConnectionStateClosed {
			// PeerConnection was explicitly closed. This usually happens from a DTLS CloseNotify
			log.Println("Peer Connection has gone to closed exiting")
			// os.Exit(0)
		}
	})

	peerConnection.OnICECandidate(s.OnICECandidate)

	dataChannel.OnOpen(func() {
		log.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	dataChannel.OnClose(func() {
		log.Printf("Data channel '%s'-'%d' closed\n", dataChannel.Label(), dataChannel.ID())
	})

	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Printf("CreateOffer error: %v", err)
		return err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		log.Printf("SetLocalDescription error: %v", err)
		return err
	}

	go func() {
		for range time.After(time.Second) {
			if !checking {
				return
			}
			data := s.getICECandidate()
			//if data != nil {
			for _, c := range data {
				if len(c) == 0 {
					break
				}
				candidate := webrtc.ICECandidateInit{
					Candidate: c,
				}
				if err = peerConnection.AddICECandidate(candidate); err != nil {
					log.Printf("AddICECandidate error: %v", err)
				}
			}
			//}
		}
	}()

	operationCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	offerResult, err := s.fileClient.CreateDownloadOffer(operationCtx, &pubUserFile.RTCFileRequest{
		ContentIdentity: fileID,
		ClientIdentity:  s.clientID,
		Sdp:             offer.SDP,
	})

	if err != nil {
		return err
	}

	// wait till peer connection is ready
	if len(offerResult.Sdp) > 0 {
		// Set the remote SessionDescription
		remoteDesc := webrtc.SessionDescription{
			SDP:  offerResult.Sdp,
			Type: webrtc.SDPTypeAnswer,
		}
		if err = peerConnection.SetRemoteDescription(remoteDesc); err != nil {
			log.Printf("SetRemoteDescription error: %v", err)
			return err
		}
		for _, c := range s.pendingCandidates {
			s.sendICECandidate(c)
		}
	} else {
		log.Printf("no more clients")
	}
	done := make(chan bool)
	<-done

	return nil

}

/*
func (s *SliceDownloader) Stop() error {
	if s.peerConnection != nil {
		return s.peerConnection.Close()
	}
	return nil
}
*/

func (s *SliceDownloader) OnICECandidate(candidate *webrtc.ICECandidate) {
	s.candidatesMux.Lock()
	defer s.candidatesMux.Unlock()
	candidateStr := ""
	if candidate != nil {
		candidateStr = candidate.ToJSON().Candidate
	}
	desc := s.peerConnection.RemoteDescription()
	if desc == nil {
		s.pendingCandidates = append(s.pendingCandidates, candidateStr)
	} else {
		s.sendICECandidate(candidateStr)
	}
}

func (s *SliceDownloader) sendICECandidate(candidate string) {
	operationCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if _, err := s.fileClient.SendClientIceCandidate(operationCtx, &pubUserFile.SendIceCandidateRequest{
		ClientIdentity: s.clientID,
		Candidate:      candidate,
	}); err != nil {
		log.Printf("SendIceCandidate error: %v", err)
	}
}

func (s *SliceDownloader) getICECandidate() []string {
	operationCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	result, err := s.fileClient.GetServerIceCandidate(operationCtx, &pubUserFile.GetIceCandidateRequest{
		ClientIdentity: s.clientID,
	})
	if err != nil {
		log.Printf("GetIceCandidate error: %v", err)
		return nil
	}
	return result.Candidate
}
