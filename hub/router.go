package hub

import (
	"encoding/json"
	"log"
)

func (h *Hub) handleClientMessage(msg rawMessage) {

	msgChannel, ok := msg.getChannelName()

	if !ok {
		log.Println("Invalid channel name in message received")
		return
	}

	for client := range h.channels[msgChannel] {
		z, err := json.Marshal(msg)
		if err != nil {
			log.Println(err.Error())
		}
		select {
		case client.send <- z:
		default:
			close(client.send)
			delete(h.channels[msg.Payload["channel"].(string)], client)
		}
	}
}
func (h *Hub) handleCreateNewChannel(msg rawMessage) {
	h.createChannel(msg.Payload["channel_name"].(string))
}

// func (h *Hub) handleUpdateClientInfo(msg map[string]interface{}) {
// 	log.Println("Revs client info")
// }

// func (h *Hub) handleWhoCommand(msg map[string]interface{}) {
// 	log.Println("who command recv")
// }

func (h *Hub) addHandler(msgName string, hander func(rawMessage)) {
	h.msgRouter[msgName] = hander
}

func (h *Hub) findHandler(f string) (func(rawMessage), bool) {
	handler, found := h.msgRouter[f]
	return handler, found
}
