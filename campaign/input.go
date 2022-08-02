package campaign

type CampaignInput struct {
	ID               int
	Name             string `json:"name" form:"name" binding:"required"`
	UserId           int    `form:"user_id"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int
	Perks            string `json:"perks"`
	BackerCount      int    `json:"backer_count"`
	Slug             string `json:"slug" binding:"required"`
}

type CampaignImageInput struct {
	CampaignId int    `json:"campaign_id"`
	FileName   string `json:"file_name"`
	IsPrimary  bool   `json:"is_primary"`
}
