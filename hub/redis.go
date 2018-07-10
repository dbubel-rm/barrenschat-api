package hub

import "log"

func (h *hub) listenForNewMessages() {
	go func() {
		defer h.RedisClient.Close()
		pSub := h.RedisClient.Subscribe("datapipe")
		for {
			msg, err := pSub.ReceiveMessage()
			if err != nil {
				log.Println(err.Error())
			}

			for _, j := range h.RoomList {
				for i := 0; i < len(j); i++ {
					j[i].newMsgChan <- msg.Payload
				}
			}
		}
	}()
}

func (h *hub) getChannels() {
	channelList := h.RedisClient.PubSubChannels("*")
	for _, j := range channelList.Val() {
		log.Println(j)
	}
}
