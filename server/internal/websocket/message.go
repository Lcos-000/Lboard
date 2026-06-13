package websocket

import "encoding/json"

const (
	MessageTypeHello      = "hello"
	MessageTypePing       = "ping"
	MessageTypePong       = "pong"
	MessageTypeJoinRoom   = "join_room"
	MessageTypeJoinAck    = "join_room_ack"
	MessageTypeError      = "error"
	MessageTypeServerInfo = "server_info"
)

// Message 消息
type Message struct {
	Type      string          `json:"type"`
	RequestID string          `json:"requestId,omitempty"`
	RoomID    string          `json:"roomId,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// ErrorPayload 错误负载
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HelloPayload 欢迎负载
type HelloPayload struct {
	UserID string `json:"userId"`
}

// JoinRoomAckPayload 加入房间确认负载
type JoinRoomAckPayload struct {
	RoomID string `json:"roomId"`
	Status string `json:"status"`
}

// MustJSON 快捷方法，将结构体转换为 JSON 字符串，不检查错误，因为marshal几乎
func MustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
