package nats

import (
	"github.com/goodluck0107/gcore/gencoding/json"
	"github.com/goodluck0107/gcore/geventbus"
	"github.com/goodluck0107/gcore/gutils/gconv"
	"github.com/goodluck0107/gcore/gutils/gtime"
	"github.com/goodluck0107/gcore/gutils/guuid"
	"github.com/goodluck0107/gcore/gwrap/value"
)

type data struct {
	ID        string `json:"id"`        // 事件ID
	Topic     string `json:"topic"`     // 事件主题
	Payload   string `json:"payload"`   // 事件载荷
	Timestamp int64  `json:"timestamp"` // 事件时间
}

// 序列化
func serialize(topic string, payload interface{}) ([]byte, error) {
	return json.Marshal(&data{
		ID:        guuid.UUID(),
		Topic:     topic,
		Payload:   gconv.String(payload),
		Timestamp: gtime.Now().UnixNano(),
	})
}

// 反序列化
func deserialize(v []byte) (*geventbus.Event, error) {
	d := &data{}

	err := json.Unmarshal(v, d)
	if err != nil {
		return nil, err
	}

	return &geventbus.Event{
		ID:        d.ID,
		Topic:     d.Topic,
		Payload:   value.NewValue(d.Payload),
		Timestamp: gtime.UnixNano(d.Timestamp),
	}, nil
}
