package campaign

import "time"

type CampaignFormater struct {
	ID               int       `json:"id"`
	UserId           int       `json:"user_id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	GoalAmount       int       `json:"goal_amount"`
	CurrentAmount    int       `json:"current_amount"`
	Perks            string    `json:"perks"`
	BackerCount      int       `json:"backer_count"`
	Slug             string    `json:"slug"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func FormatCampaign(campaign Campaign) CampaignFormater {
	return CampaignFormater{
		ID:     campaign.ID,
		UserId: campaign.UserId,
	}
}
