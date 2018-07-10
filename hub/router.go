package hub

import (
	"log"
)

func (h *hub) handleClientMessage(msg map[string]interface{}) {
	if redisClient.Publish("datapipe", msg[MESSAGE_TEXT]).Err() != nil {
		log.Println("Error publishing message to datapipe")
	}
}

func (h *hub) handleUpdateClientInfo(msg map[string]interface{}) {
	log.Println("Revs client info")
}

func (h *hub) handleWhoCommand(msg map[string]interface{}) {
	log.Println("who command recv")
}

func (h *hub) handleMsg(msgName string, hander func(map[string]interface{})) {
	h.Router[msgName] = hander
}

func (h *hub) findHandler(f string) (func(map[string]interface{}), bool) {
	handler, found := h.Router[f]
	return handler, found
}
