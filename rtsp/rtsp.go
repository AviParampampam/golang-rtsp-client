package rtsp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// RTSP Methods
const (
	OPTIONS      = "OPTIONS"
	DESCRIBE     = "DESCRIBE"
	ANNOUNCE     = "ANNOUNCE"
	SETUP        = "SETUP"
	PLAY         = "PLAY"
	PAUSE        = "PAUSE"
	TEARDOWN     = "TEARDOWN"
	GETPARAMETER = "GET_PARAMETER"
	SETPARAMETER = "SET_PARAMETER"
	REDIRECT     = "REDIRECT"
	RECORD       = "RECORD"
)

// RTSP Status codes
const (
	Continue                      = 100
	OK                            = 200
	Created                       = 201
	LowOnStorageSpace             = 250
	MultipleChoices               = 300
	MovedPermanently              = 301
	MovedTemporarily              = 302
	SeeOther                      = 303
	UseProxy                      = 305
	BadRequest                    = 400
	Unauthorized                  = 401
	PaymentRequired               = 402
	Forbidden                     = 403
	NotFound                      = 404
	MethodNotAllowed              = 405
	NotAcceptable                 = 406
	ProxyAuthenticationRequired   = 407
	RequestTimeout                = 408
	Gone                          = 410
	LengthRequired                = 411
	PreconditionFailed            = 412
	RequestEntityTooLarge         = 413
	RequestURITooLong             = 414
	UnsupportedMediaType          = 415
	Invalidparameter              = 451
	IllegalConferenceIdentifier   = 452
	NotEnoughBandwidth            = 453
	SessionNotFound               = 454
	MethodNotValidInThisState     = 455
	HeaderFieldNotValid           = 456
	InvalidRange                  = 457
	ParameterIsReadOnly           = 458
	AggregateOperationNotAllowed  = 459
	OnlyAggregateOperationAllowed = 460
	UnsupportedTransport          = 461
	DestinationUnreachable        = 462
	InternalServerError           = 500
	NotImplemented                = 501
	BadGateway                    = 502
	ServiceUnavailable            = 503
	GatewayTimeout                = 504
	RTSPVersionNotSupported       = 505
	OptionNotsupport              = 551
)

// Client - rtsp client
type Client struct {
	UserAgent string
	Sessions  []Session
}

// Session - ..
type Session struct {
	Host    string
	CSeq    int
	conn    net.Conn
	Session string
}

type request struct {
	method  string
	url     string
	version string
	header  map[string][]string
	body    []byte
}

// Response - ...
type Response struct {
	Version    string
	StatusCode int
	Status     string
	Header     map[string][]string
	Body       []byte
	String     func() string
}

// NewClient - creating new rtsp client
func NewClient() Client {
	return Client{
		"SV .1 Local Client",
		make([]Session, 4),
	}
}

// NewSession - creating new sessions
func (client *Client) NewSession(host string, session string) (Session, error) {
	s := Session{
		host,
		0,
		nil,
		session,
	}
	client.Sessions = append(client.Sessions, s)

	err := s.Connect()
	return s, err
}

func (req request) string() string {
	s := fmt.Sprintf("%s %s %s\r\n", req.method, req.url, req.version)

	for k, v := range req.header {
		for _, v := range v {
			s += fmt.Sprintf("%s: %s\r\n", k, v)
		}
	}
	s += "\r\n" + string(req.body)

	return s
}

// Connect - connecting to server
func (session *Session) Connect() error {
	conn, err := net.Dial("tcp", session.Host)
	if err != nil {
		return err
	}
	session.conn = conn
	return nil
}

// Disconnect - disconnect from server
func (session *Session) Disconnect() error {
	return session.conn.Close()
}

func (session *Session) sendRequest(req request) error {
	_, err := io.WriteString(session.conn, req.string())
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) getResponse() (Response, error) {
	res := Response{}

	// Receiving response
	reader := bufio.NewReader(session.conn)
	buf := make([]byte, 4096)
	_, err := reader.Read([]byte(buf))
	if err != nil {
		return res, err
	}
	s := string(bytes.Trim(buf, "\x00"))
	lines := strings.Split(s, "\n")

	// Making response
	res.Header = make(map[string][]string)

	firstLineWords := strings.Split(lines[0], " ")
	res.Version = firstLineWords[0]
	res.StatusCode, err = strconv.Atoi(firstLineWords[1])
	if err != nil {
		return res, err
	}
	res.Status = firstLineWords[2]

	var indexSplitLine int
	for i, l := range lines {
		if len(l) == 0 {
			indexSplitLine = i
		}
	}
	headLines := lines[1 : indexSplitLine-1]
	for _, line := range headLines {
		l := strings.Split(line, ":")
		if len(l) > 1 {
			res.Header[l[0]] = strings.Split(l[1], ",")
		} else {
			res.Header[l[0]] = []string{}
		}
	}

	bodyLines := lines[indexSplitLine:]
	res.Body = []byte(strings.Join(bodyLines, "\n"))
	res.String = func() string { return s }

	return res, nil
}

func (session *Session) nextCSeq() string {
	session.CSeq++
	return strconv.Itoa(session.CSeq)
}

// Describe - send DESCRIBE method
func (session *Session) Describe(url string) (Response, error) {
	hs := map[string][]string{
		"CSeq":    {session.nextCSeq()},
		"Session": {session.Session},
	}
	req := request{DESCRIBE, url, "RTSP/1.0", hs, nil}

	err := session.sendRequest(req)
	if err != nil {
		return Response{}, err
	}
	return session.getResponse()
}

// Options - send OPTIONS method
func (session *Session) Options(url string) (Response, error) {
	hs := map[string][]string{
		"CSeq":    {session.nextCSeq()},
		"Session": {session.Session},
	}
	req := request{OPTIONS, url, "RTSP/1.0", hs, nil}

	err := session.sendRequest(req)
	if err != nil {
		return Response{}, err
	}
	return session.getResponse()
}

// Setup - send SETUP method
func (session *Session) Setup(url string) (Response, error) {
	hs := map[string][]string{
		"CSeq":      {session.nextCSeq()},
		"Session":   {session.Session},
		"Transport": {"RTP/AVP;unicast;client_port=8010-8011"},
	}
	req := request{SETUP, url, "RTSP/1.0", hs, nil}

	err := session.sendRequest(req)
	if err != nil {
		return Response{}, err
	}
	return session.getResponse()
}
