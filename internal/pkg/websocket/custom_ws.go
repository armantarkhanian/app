package websocket

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/centrifugal/centrifuge"
	"github.com/centrifugal/protocol"
	"nhooyr.io/websocket"
)

type customWebsocketHandler struct {
	node *centrifuge.Node
}

func newWebsocketHandler(node *centrifuge.Node) http.Handler {
	return &customWebsocketHandler{node}
}

const websocketTransportName = "websocket"

type customWebsocketTransport struct {
	mu      sync.RWMutex
	closed  bool
	closeCh chan struct{}

	conn      *websocket.Conn
	protoType centrifuge.ProtocolType
	request   *http.Request
}

func newWebsocketTransport(conn *websocket.Conn, protoType centrifuge.ProtocolType) *customWebsocketTransport {
	return &customWebsocketTransport{
		conn:      conn,
		protoType: protoType,
		closeCh:   make(chan struct{}),
	}
}

// Name implementation.
func (t *customWebsocketTransport) Name() string {
	return websocketTransportName
}

// Protocol implementation.
func (t *customWebsocketTransport) Protocol() centrifuge.ProtocolType {
	return t.protoType
}

// Encoding implementation.
func (t *customWebsocketTransport) Encoding() centrifuge.EncodingType {
	return centrifuge.EncodingTypeJSON
}

// Unidirectional implementation.
func (t *customWebsocketTransport) Unidirectional() bool {
	return false
}

// DisabledPushFlags ...
func (t *customWebsocketTransport) DisabledPushFlags() uint64 {
	return centrifuge.PushFlagDisconnect
}

// Write ...
func (t *customWebsocketTransport) Write(messages ...[]byte) error {
	select {
	case <-t.closeCh:
		return nil
	default:
		var messageType = websocket.MessageText
		protoType := protocol.TypeJSON

		if t.Protocol() == centrifuge.ProtocolTypeProtobuf {
			messageType = websocket.MessageBinary
			protoType = protocol.TypeProtobuf
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		encoder := protocol.GetDataEncoder(protoType)
		defer protocol.PutDataEncoder(protoType, encoder)
		for i := range messages {
			_ = encoder.Encode(messages[i])
		}
		return t.conn.Write(ctx, messageType, encoder.Finish())
	}
}

// Close ...
func (t *customWebsocketTransport) Close(disconnect *centrifuge.Disconnect) error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	close(t.closeCh)
	t.mu.Unlock()

	if disconnect != nil {
		return t.conn.Close(websocket.StatusCode(disconnect.Code), disconnect.CloseText())
	}
	return t.conn.Close(websocket.StatusNormalClosure, "")
}

func (s *customWebsocketHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Accept(rw, r, &websocket.AcceptOptions{})
	if err != nil {
		s.node.Log(centrifuge.NewLogEntry(centrifuge.LogLevelDebug, "websocket upgrade error", map[string]interface{}{"error": err.Error()}))
		return
	}

	var protoType = centrifuge.ProtocolTypeJSON
	if r.URL.Query().Get("format") == "protobuf" {
		protoType = centrifuge.ProtocolTypeProtobuf
	}

	transport := newWebsocketTransport(conn, protoType)

	select {
	case <-s.node.NotifyShutdown():
		_ = transport.Close(centrifuge.DisconnectShutdown)
		return
	default:
	}

	c, closeFn, err := centrifuge.NewClient(r.Context(), s.node, transport)
	if err != nil {
		s.node.Log(centrifuge.NewLogEntry(centrifuge.LogLevelError, "error creating client", map[string]interface{}{"transport": websocketTransportName}))
		return
	}
	defer func() { _ = closeFn() }()
	s.node.Log(centrifuge.NewLogEntry(centrifuge.LogLevelDebug, "client connection established", map[string]interface{}{"client": c.ID(), "transport": websocketTransportName}))
	defer func(started time.Time) {
		s.node.Log(centrifuge.NewLogEntry(centrifuge.LogLevelDebug, "client connection completed", map[string]interface{}{"client": c.ID(), "transport": websocketTransportName, "duration": time.Since(started)}))
	}(time.Now())

	for {
		_, data, err := conn.Read(context.Background())
		if err != nil {
			return
		}
		ok := c.Handle(data)
		if !ok {
			return
		}
	}
}
