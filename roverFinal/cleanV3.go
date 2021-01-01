package main

import (
        "log"
        "net/http"
        "strings"
        "time"

        "github.com/gorilla/websocket"

        "github.com/stianeikeland/go-rpio/v4"

        "github.com/pion/mediadevices"
        "github.com/pion/mediadevices/examples/internal/signal"
        "github.com/pion/mediadevices/pkg/frame"
        "github.com/pion/mediadevices/pkg/prop"
        "github.com/pion/webrtc/v3"

        // If you don't like x264, you can also use vpx by importing as below
        // "github.com/pion/mediadevices/pkg/codec/vpx" // This is required to use VP8/VP9 video encoder
        // or you can also use openh264 for alternative h264 implementation
        // "github.com/pion/mediadevices/pkg/codec/openh264"
        // or if you use a raspberry pi like, you can use mmal for using its hardware encoder
        "github.com/pion/mediadevices/pkg/codec/mmal"
        //"github.com/pion/mediadevices/pkg/codec/opus" // This is required to use opus audio encoder
        //"github.com/pion/mediadevices/pkg/codec/x264" // This is required to use h264 video encoder

        // Note: If you don't have a camera or microphone or your adapters are not supported,
        //       you can always swap your adapters with our dummy adapters below.
        // _ "github.com/pion/mediadevices/pkg/driver/videotest"
        // _ "github.com/pion/mediadevices/pkg/driver/audiotest"
        _ "github.com/pion/mediadevices/pkg/driver/camera" // This is required to register camera adapter
        //_ "github.com/pion/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter
)

var upgrader = websocket.Upgrader{} // use default options for upgrader


func echo(w http.ResponseWriter, r *http.Request) {
        c, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Print("upgrade:", err)
                return
        }
        defer c.Close()

        //webrtc stuffffffffff

        config := webrtc.Configuration{
                ICEServers: []webrtc.ICEServer{
                        {
                                URLs: []string{"stun:stun.l.google.com:19302"},
                        },
                },
        }

        // Create a new RTCPeerConnection
        mmalParams, err := mmal.NewParams()
        if err != nil {
                panic(err)
        }
        mmalParams.BitRate = 500_000 // 1000kbps
        //mmalParams.BitRate = 0

        // opusParams, err := opus.NewParams()
        //  if err != nil {
        //          panic(err)
        //  }

        codecSelector := mediadevices.NewCodecSelector(
                mediadevices.WithVideoEncoders(&mmalParams),
                //mediadevices.WithAudioEncoders(&opusParams),
        )

        mediaEngine := webrtc.MediaEngine{}
        codecSelector.Populate(&mediaEngine)
        //if err := mediaEngine.PopulateFromSDP(offer); err != nil {
        //        panic(err)
        //}
        api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
        peerConnection, err := api.NewPeerConnection(config)
        if err != nil {
                panic(err)
        }

        // Set the handler for ICE connection state
        // This will notify you when the peer has connected/disconnected
        peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
                log.Printf("Connection State has changed %s \n", connectionState.String())
        })

        s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
                Video: func(c *mediadevices.MediaTrackConstraints) {
                        c.FrameFormat = prop.FrameFormat(frame.FormatYUY2)
                        //c.FrameFormat = prop.FrameFormatExact(frame.FormatI420)
                        c.Width = prop.Int(352)
                        c.Height = prop.Int(288)
                },
                //Audio: func(c *mediadevices.MediaTrackConstraints) {
                //},
                Codec: codecSelector,
        })
        if err != nil {
                panic(err)
        }

        for _, track := range s.GetTracks() {
                track.OnEnded(func(err error) {
                        log.Printf("Track (ID: %s) ended with error: %v\n",
                                track.ID(), err)
                })

                // In Pion/webrtc v3, bind will be called automatically after SDP negotiation
                //webrtcTrack, err := track.Bind(peerConnection)
                //if err != nil {
                //        panic(err)
                //}

                _, err = peerConnection.AddTransceiverFromTrack(track,
                        webrtc.RtpTransceiverInit{
                                Direction: webrtc.RTPTransceiverDirectionSendonly,
                        },
                )
                if err != nil {
                        panic(err)
                }
        }

        // Create an offer
        offer, err := peerConnection.CreateOffer(nil)
        if err != nil {
                panic(err)
        }

        //Create a channel that is blocked until ICE Gathering is complete
        gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

        // Sets the LocalDescription, looks for ICE candidates
        err = peerConnection.SetLocalDescription(offer)
        if err != nil {
                panic(err)
        }

        //Wait for ICE gathering to complete (non-trickle ICE)
        <-gatherComplete

        //dt = time.Now()
        //log.Print(dt.String())

        //Output the SDP with the final ICE candidate
        log.Println(signal.Encode(*peerConnection.LocalDescription()))

        //Send the SDP with the final ICE candidate to the browser as our offer
        err = c.WriteMessage(1, []byte(signal.Encode(*peerConnection.LocalDescription()))) //write message back to browser
        if err != nil {
                log.Println("write:", err)
        }

        //Wait for the browser to return an answer (its SDP)
        msgType, message, err2 := c.ReadMessage() //ReadMessage blocks until message received
        if err2 != nil {
                log.Println("read:", err)
        }

        answer := webrtc.SessionDescription{}

        signal.Decode(string(message), &answer) //set offer to the decoded SDP
        log.Print(answer)

        // Set the remote SessionDescription
        err = peerConnection.SetRemoteDescription(answer)
        if err != nil {
                panic(err)
        }

        //log.Printf("Type: %s", msgType)
        //log.Printf("Message: %s", message)
        log.Printf("%s sent: %s\n", c.RemoteAddr(), string(message), msgType)

        //====================Now Listen For User Rover Controls========================
        var control string
        var lastTurn string //used to rember what the last direction of turning was

        for {
                msgType, message, err2 = c.ReadMessage() //ReadMessage blocks until message received
                if err2 != nil {
                        log.Println("read:", err)
                        break
                }

                control = string(message)

                if strings.HasPrefix(control, "!") {
                        //Turning
                        if control == "!l" {
                                //go left if right pin is up
                                pinR.Low()
                                pinL.High()
                                lastTurn = "l"
                        } else if control == "!r" {
                                //go right if left pin is up
                                pinL.Low()
                                pinR.High()
                                lastTurn = "r"
                        }

                        //Forwards or Backwards
                        if control == "!f" {
                                //go forwards if back pin is up
                                pinB.Low()
                                pinF.High()
                        } else if control == "!b" {
                                //go backwards if forwards pin is up
                                pinF.Low()
                                pinB.High()
                        }

                        //Stop Movement || Stop Turning
                        if control == "!keyUpT" {
                                //pull turn pins down
                                pinR.Low()
                                pinL.Low()

                                //turn back to mid
                                if lastTurn == "r" {
                                        pinL.High()
                                        time.Sleep(time.Second / 15)
                                        pinL.Low()
                                } else {
                                        pinR.High()
                                        time.Sleep(time.Second / 15)
                                        pinR.Low()
                                }
                        }

                        if control == "!keyUpM" {
                                //pull move pins down
                                pinF.Low()
                                pinB.Low()
                        }

                }

        }

}

//==================Pin Setup===================================================
//Initialize pins, Pin refers to the bcm2835 pin, not the physical pin on the raspberry pi header
var (
        pinB  = rpio.Pin(13) //Backwards
        pinF  = rpio.Pin(12) //Forwards
        pinR  = rpio.Pin(5)  //Steer Right
        pinL  = rpio.Pin(6)  //Steer Left
        pinV  = rpio.Pin(0)  //Vcc pin to tell controller what the logic voltage is
        pinPB = rpio.Pin(15) //PWMB set always on
        pinPA = rpio.Pin(18) //PWMA set always on
)

//==============================================================================

func main() {
        //Open memory range for GPIO access in /dev/mem
        err := rpio.Open()
        if err != nil {
                panic(err) //panic logs error and exits program
        }
        // Unmap gpio memory when main() exits
        defer rpio.Close()

        //Set Pins to mode Output
        pinB.Output()
        pinF.Output()
        pinR.Output()
        pinL.Output()

        pinV.Output()
        pinV.High()
        pinPB.Output()
        pinPB.High()
        pinPA.Output()
        pinPA.High()

        //flag.Parse()
        //log.SetFlags(0)

        fileServer := http.FileServer(http.Dir("./public"))
        http.HandleFunc("/echo", echo) //this request comes from webrtc.html
        http.Handle("/", fileServer)

        err = http.ListenAndServe(":80", nil) //Http server blocks
        if err != nil {
                log.Fatal(err)
        }

}
