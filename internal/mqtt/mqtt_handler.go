package mqtt

import (
	"encoding/json"
	"log"

	"gps-backend/internal/model"

	"github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	Client          mqtt.Client
	OnTrackReceived func(track model.Track) // callback при получении трека
}

func NewHandler(broker string) *MQTTHandler {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	client := mqtt.NewClient(opts)
	return &MQTTHandler{
		Client: client,
	}
}

func (h *MQTTHandler) Connect() error {
	if token := h.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Println("MQTT connected")
	return nil
}

func (h *MQTTHandler) Subscribe(topic string) {
	h.Client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var track model.Track
		if err := json.Unmarshal(msg.Payload(), &track); err != nil {
			log.Println("MQTT unmarshal error:", err)
			return
		}
		log.Printf("Received track: %+v\n", track)

		if h.OnTrackReceived != nil {
			h.OnTrackReceived(track) // вызываем callback
		}
	})
}
