package hub

import (
	"encoding/json"
	"log"
)

func (h *Hub) handleClientMessage(msg rawMessage) {
	for _, v := range h.channelMembers {

		z, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
		}

		select {
		case v.send <- z:
		default:
			close(v.send)
			// delete(h.channels[msg.Payload["channel"].(string)], client)
		}
	}
}

func (h *Hub) addHandler(msgName string, hander func(rawMessage)) {
	h.msgRouter[msgName] = hander
}

func (h *Hub) findHandler(f string) (func(rawMessage), bool) {
	handler, found := h.msgRouter[f]
	return handler, found
}
