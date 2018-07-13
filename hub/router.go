package hub

import (
	"log"
)

func (h *Hub) handleClientMessage(msg map[string]interface{}) {
	// if redisClient.Publish("datapipe", msg[MESSAGE_TEXT]).Err() != nil {
	// 	log.Println("Error publishing message to datapipe")
	// }
}

func (h *Hub) handleUpdateClientInfo(msg map[string]interface{}) {
	log.Println("Revs client info")
}

func (h *Hub) handleWhoCommand(msg map[string]interface{}) {
	log.Println("who command recv")
}

func (h *Hub) handleMsg(msgName string, hander func(map[string]interface{})) {
	// h.Router[msgName] = hander
}

// func (h *Hub) findHandler(f string) (func(map[string]interface{}), bool) {
// 	// handler, found := h.Router[f]
// 	// return handler, found
// }
