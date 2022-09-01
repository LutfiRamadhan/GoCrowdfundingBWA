package transaction

import "time"

type Transaction struct {
	ID         int
	CampaignId int
	UserId     int
	Amount     int
	status     string
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
