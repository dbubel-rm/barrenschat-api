package hub

import (
	"encoding/json"
	"log"
)

func (h *Hub) handleClientMessage(msg rawMessage) {
	for _, v := range h.channelMembers[msg.Payload["channel"].(string)] {

		z, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
		}
		v.send <- z
		// select {
		// case v.send <- z:
		// default:
		// 	close(v.send)
		// 	// delete(h.channels[msg.Payload["channel"].(string)], client)
		// }
	}
}

func (h *Hub) addHandler(msgName string, hander func(rawMessage)) {
	h.msgRouter[msgName] = hander
}

func (h *Hub) findHandler(f string) (func(rawMessage), bool) {
	h.locker <- true
	handler, found := h.msgRouter[f]
	<-h.locker
	return handler, found
}
