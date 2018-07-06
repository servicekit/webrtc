package mediasoup

type RtpCapabilities struct {
	Codecs []struct {
		Kind       string `json:"kind"`
		Name       string `json:"name"`
		MimeType   string `json:"mimeType"`
		ClockRate  int    `json:"clockRate"`
		Channels   int    `json:"channels,omitempty"`
		Parameters struct {
			Useinbandfec int `json:"useinbandfec"`
			Minptime     int `json:"minptime"`
		} `json:"parameters"`
		RtcpFeedback         []interface{} `json:"rtcpFeedback,omitempty"`
		PreferredPayloadType int           `json:"preferredPayloadType"`
	} `json:"codecs"`
	HeaderExtensions []struct {
		Kind             string `json:"kind"`
		URI              string `json:"uri"`
		PreferredID      int    `json:"preferredId"`
		PreferredEncrypt bool   `json:"preferredEncrypt"`
	} `json:"headerExtensions"`
	FecMechanisms []interface{} `json:"fecMechanisms"`
}

type QueryRoomData struct {
	RtpCapabilities            RtpCapabilities `json:"rtpCapabilities"`
	MandatoryCodecPayloadTypes []interface{}   `json:"mandatoryCodecPayloadTypes"`
}

type JoinRoomData struct {
	Method          string          `json:"method"`
	Target          string          `json:"target"`
	PeerName        string          `json:"peerName"`
	RtpCapabilities RtpCapabilities `json:"rtpCapabilities"`
	AppData         struct {
		DisplayName string `json:"displayName"`
		Device      struct {
			Flag    string `json:"flag"`
			Name    string `json:"name"`
			Version string `json:"version"`
			Bowser  struct {
				Name      string `json:"name"`
				Chrome    bool   `json:"chrome"`
				Version   string `json:"version"`
				Blink     bool   `json:"blink"`
				Windows   bool   `json:"windows"`
				Osname    string `json:"osname"`
				Osversion string `json:"osversion"`
				A         bool   `json:"a"`
			} `json:"bowser"`
		} `json:"device"`
	} `json:"appData"`
}
