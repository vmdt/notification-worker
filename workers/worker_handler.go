package workers

import (
	"encoding/json"

	"github.com/vmdt/notification-worker/server/consumer"
	"github.com/vmdt/notification-worker/shared"
)

func HandleDiscountWorker(msg interface{}, dependencies *shared.DiscountBase) error {
	var discountMsg consumer.DiscountMessage
	if err := json.Unmarshal(msg.([]byte), &discountMsg); err != nil {
		dependencies.Log.Errorf("Error unmarshaling message: %v", err)
		return err
	}

	err := dependencies.Publisher.PublishMessage(discountMsg, "discount", "discount_key")
	if err != nil {
		dependencies.Log.Errorf("Error publishing message: %v", err)
		return err
	}
	return nil
}
