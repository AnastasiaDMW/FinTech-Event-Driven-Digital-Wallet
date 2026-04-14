package kafka

const UserChangedTopic = "user.changed"

type UserChangedEvent struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
