package hub

import (
	"encoding/json"
	"log"
)

func (h *Hub) handleClientMessage(msg rawMessage) {
	log.Println("handler found")
	// msgChannel, ok := msg.getChannelName()

	// if !ok {
	// 	log.Println("Invalid channel name in message received")
	// 	return
	// }
	// log.Println("all channels", h.channels)
	for client, v := range h.clients {
		log.Println("CLIENT:", client, "V", v)
		for _, c := range v.channelsSubscribedTo {
			if c == msg.Payload["channel"].(string) {
				log.Println("found in chan")
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

			} else {
				log.Println("not found in chan")
			}
		}
		// z, err := json.Marshal(msg)
		// if err != nil {
		// 	log.Println(err.Error())
		// }
	}
	// 	select {
	// 	case client.send <- z:
	// 	default:
	// 		close(client.send)
	// 		delete(h.channels[msg.Payload["channel"].(string)], client)
	// 	}
	// }
}

// func (h *Hub) handleJoinNewChannel(msg rawMessage, c *Client) {
// 	fmt.Println("in join new channel")
// 	//h.createChannel(msg.Payload["channel_name"].(string))
// }

// func (h *Hub) handleGetChannels(msg rawMessage) {
// 	h.getChannels()
// }

// func (h *Hub) handleCreateNewChannel(msg rawMessage) {
// 	h.blocker <- true
// 	h.createChannel(msg.Payload["channel_name"].(string))
// }

func (h *Hub) addHandler(msgName string, hander func(rawMessage)) {
	h.msgRouter[msgName] = hander
}

// func (h *Hub) addHandlerClient(msgName string, hander func(rawMessage, *Client)) {
// 	h.msgRouterClient[msgName] = hander
// }

func (h *Hub) findHandler(f string) (func(rawMessage), bool) {
	handler, found := h.msgRouter[f]
	return handler, found
}

// func (h *Hub) findHandlerClients(f string) (func(rawMessage, *Client), bool) {
// 	handler, found := h.msgRouterClient[f]
// 	return handler, found
// }
