package hub

import (
	"log"
)

func (h *Hub) HandleClientMessage(msg map[string]interface{}) {
	redisClient.Publish("datapipe", msg["msgText"]).Err()
}
func (h *Hub) HandleUpdateClientInfo(msg map[string]interface{}) {
	log.Println("Revs client info")
}
func (h *Hub) HandleWhoCommand(msg map[string]interface{}) {
	log.Println("who command recv")
}
func (h *Hub) HandleMsg(msgName string, hander func(map[string]interface{})) {
	h.Router[msgName] = hander
}
func (h *Hub) findHandler(f string) (func(map[string]interface{}), bool) {
	handler, found := h.Router[f]
	return handler, found
}
