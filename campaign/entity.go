package campaign

import (
	"BWA/user"
	"time"
)

type Campaign struct {
	ID               int              `json:"id"`
	UserId           int              `json:"user_id"`
	Name             string           `json:"name"`
	ShortDescription string           `json:"short_description"`
	Description      string           `json:"description"`
	GoalAmount       int              `json:"goal_amount"`
	CurrentAmount    int              `json:"current_amount"`
	Perks            string           `json:"perks"`
	BackerCount      int              `json:"backer_count"`
	Slug             string           `json:"slug"`
	CampaignImages   []CampaignImages `json:"campaign_images"`
	User             user.User        `json:"user,omitempty"`
	ImagesUrl        string           `json:"images_url" gorm:"-"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type CampaignImages struct {
	ID         int       `json:"id"`
	CampaignID int       `json:"campaign_id"`
	FileName   string    `json:"file_name"`
	IsPrimary  bool      `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
