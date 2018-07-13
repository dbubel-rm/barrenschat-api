package hub

import "log"

func (h *Hub) pubSubListen(pubSub string) {
	// defer h.RedisClient.Close()
	pSub := redisClient.Subscribe(pubSub)
	for {
		msg, err := pSub.ReceiveMessage()
		if err != nil {
			log.Println(err.Error())
		}
		h.pubSubRecv <- []byte(msg.Payload)
	}
}

func (h *Hub) getChannels() {
	channelList := redisClient.PubSubChannels("*")
	for _, j := range channelList.Val() {
		log.Println(j)
	}
}
