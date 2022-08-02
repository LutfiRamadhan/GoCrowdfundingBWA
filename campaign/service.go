package campaign

type Service interface {
	FindCampaigns(campaign CampaignInput) ([]Campaign, error)
	FindCampaign(campaign CampaignInput) (Campaign, error)
	CreateCampaign(campaign CampaignInput, userId int) (Campaign, error)
	UpdateCampaign(campaign CampaignInput) (Campaign, error)
	DeleteCampaign(campaign CampaignInput) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) FindCampaigns(campaign CampaignInput) ([]Campaign, error) {
	requestCampaign := Campaign{
		UserId: campaign.UserId,
		Name:   campaign.Name,
	}
	data, err := s.repository.Gets(requestCampaign)
	if err != nil {
		return []Campaign{}, err
	}
	return data, nil
}

func (s *service) FindCampaign(campaign CampaignInput) (Campaign, error) {
	requestCampaign := Campaign{
		ID:   campaign.ID,
		Slug: campaign.Slug,
	}
	data, err := s.repository.Get(requestCampaign)
	if err != nil {
		return Campaign{}, nil
	}
	return data, nil
}

func (s *service) CreateCampaign(campaign CampaignInput, userId int) (Campaign, error) {
	dataCampaign := Campaign{
		Name:             campaign.Name,
		UserId:           userId,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    0,
		Perks:            campaign.Perks,
		BackerCount:      campaign.BackerCount,
		Slug:             campaign.Slug,
	}
	data, err := s.repository.Save(dataCampaign)
	if err != nil {
		return Campaign{}, err
	}
	return data, nil
}

func (s *service) UpdateCampaign(campaign CampaignInput) (Campaign, error) {
	dataCampaign := Campaign{
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Perks:            campaign.Perks,
		BackerCount:      campaign.BackerCount,
		Slug:             campaign.Slug,
	}
	data, err := s.repository.Update(dataCampaign)
	if err != nil {
		return Campaign{}, err
	}
	return data, nil
}

func (s *service) DeleteCampaign(campaign CampaignInput) error {
	dataCampaign := Campaign{
		ID: campaign.ID,
	}
	if err := s.repository.Delete(dataCampaign); err != nil {
		return err
	}
	return nil
}
