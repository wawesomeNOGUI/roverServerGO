// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	//"html/template"
	"log"
	"net/http"
	//"fmt"
	//"sync"

	"github.com/gorilla/websocket"


/*
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
	//_ "github.com/pion/mediadevices/pkg/driver/camera"     // This is required to register camera adapter
	//_ "github.com/pion/mediadevices/pkg/driver/microphone" // This is required to register microphone adapter

*/
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options for upgrader

var SDP string   //SDP set in echo function


func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	//for {   //infinite loop to wait for browser's SDP

		 msgType, message, err2 := c.ReadMessage()   //ReadMessage blocks until message received
			if err2 != nil {
				log.Println("read:", err)
				//break
			}

			SDP = string(message)

		//log.Printf("Type: %s", msgType)
		//log.Printf("Message: %s", message)
		log.Printf("%s sent: %s\n", c.RemoteAddr(), string(message), msgType)

	//}


/*
		//webrtc stuffffffffff

		        config := webrtc.Configuration{
		                ICEServers: []webrtc.ICEServer{
		                        {
		                                URLs: []string{"stun:stun.l.google.com:19302"},
		                        },
		                },
		        }

		        // Wait for the offer to be pasted
		        offer := webrtc.SessionDescription{}
		        //signal.Decode(signal.MustReadStdin(), &offer)
						signal.Decode(SDP, &offer)   //set offer to the decoded SDP

		        // Create a new RTCPeerConnection
		        mmalParams, err := mmal.NewParams()
		        if err != nil {
		                panic(err)
		        }
		        mmalParams.BitRate = 500_000 // 500kbps


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
		        if err := mediaEngine.PopulateFromSDP(offer); err != nil {
		                panic(err)
		        }
		        api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
		        peerConnection, err := api.NewPeerConnection(config)
		        if err != nil {
		                panic(err)
		        }

		        // Set the handler for ICE connection state
		        // This will notify you when the peer has connected/disconnected
		        peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		                fmt.Printf("Connection State has changed %s \n", connectionState.String())
		        })

		        s, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		                Video: func(c *mediadevices.MediaTrackConstraints) {
		                        c.FrameFormat = prop.FrameFormat(frame.FormatYUY2)
		                        c.Width = prop.Int(640)
		                        c.Height = prop.Int(480)
		                },
		                Audio: func(c *mediadevices.MediaTrackConstraints) {
		                },
		                Codec: codecSelector,
		        })
		        if err != nil {
		                panic(err)
		        }

		        for _, tracker := range s.GetTracks() {
		                tracker.OnEnded(func(err error) {
		                        fmt.Printf("Track (ID: %s) ended with error: %v\n",
		                                tracker.ID(), err)
		                })

		                // In Pion/webrtc v3, bind will be called automatically after SDP negotiation
		                webrtcTrack, err := tracker.Bind(peerConnection)
		                if err != nil {
		                        panic(err)
		                }

		                _, err = peerConnection.AddTransceiverFromTrack(webrtcTrack,
		                        webrtc.RtpTransceiverInit{
		                                Direction: webrtc.RTPTransceiverDirectionSendonly,
		                        },
		                )
		                if err != nil {
		                        panic(err)
		                }
		        }

		        // Set the remote SessionDescription
		        err = peerConnection.SetRemoteDescription(offer)
		        if err != nil {
		                panic(err)
		        }

		        // Create an answer
		        answer, err := peerConnection.CreateAnswer(nil)
		        if err != nil {
		                panic(err)
		        }

		        // Sets the LocalDescription, and starts our UDP listeners
		        err = peerConnection.SetLocalDescription(answer)
		        if err != nil {
		                panic(err)
		        }

		        // Output the answer in base64 so we can paste it in browser
		        log.Println(signal.Encode(answer))

						err = c.WriteMessage(msgType, signal.Encode(answer))  //write message back to browser
							if err != nil {
								log.Println("write:", err)
								//break
							}

*/


/*
		err = c.WriteMessage(msgType, message)  //write message back to browser
			if err != nil {
				log.Println("write:", err)
				break
			}
*/
	//}
}


func home(w http.ResponseWriter, r *http.Request) {
	//homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
	http.ServeFile(w, r, "./public/webrtc.html")
}

func httpServerAndWebsockets() {

	log.Fatal(http.ListenAndServe(*addr, nil))

	//runtime.Goexit()  //exits a go routine
}



func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo) //this request comes from webrtc.html
	http.HandleFunc("/", home)

	//wg := new(sync.WaitGroup)
	//wg.Add(1)  //one wait group that will keep the program running until the go routine
						//httpServerAndWebsockets exits,  I'll put the blocking part at the end of func main()
	go httpServerAndWebsockets()




					select {}  //wait until go routines done??


	//wg.Wait()  //the main function won't exit until httpServerAndWebsockets is exited
}


/*
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
*/
