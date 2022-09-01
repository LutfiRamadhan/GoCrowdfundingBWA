package campaign

import "errors"

type Service interface {
	FindCampaigns(campaign CampaignInput) ([]Campaign, error)
	FindCampaign(campaign CampaignInput) (Campaign, error)
	CreateCampaign(campaign CampaignInput, userId int) (Campaign, error)
	UpdateCampaign(campaign CampaignInput) (Campaign, error)
	DeleteCampaign(campaign CampaignInput) error
	ExportCampaign(campaign CampaignExportFilter) ([]Campaign, error)
	CreateImage(campaignImage CampaignImageInput) (CampaignImages, error)
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
		return Campaign{}, err
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
	requestCampaign := Campaign{
		ID:     campaign.ID,
		UserId: campaign.UserId,
	}
	check, err := s.repository.Get(requestCampaign)
	if err != nil {
		return Campaign{}, err
	}
	if check.ID == 0 {
		return Campaign{}, errors.New("Not an owner of the campaign")
	}

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

func (s *service) ExportCampaign(campaign CampaignExportFilter) ([]Campaign, error) {
	requestCampaign := Campaign{
		UserId: campaign.UserId,
		Perks:  campaign.Perks,
	}
	data, err := s.repository.Gets(requestCampaign)
	if err != nil {
		return []Campaign{}, err
	}
	return data, nil
}

func (s *service) CreateImage(campaignImage CampaignImageInput) (CampaignImages, error) {
	requestCampaign := Campaign{
		ID: campaignImage.CampaignId,
	}
	data, err := s.repository.Get(requestCampaign)
	if err != nil {
		return CampaignImages{}, err
	}

	if data.UserId != campaignImage.UserId {
		return CampaignImages{}, errors.New("Not an owner of the campaign")
	}

	dataImage := CampaignImages{
		CampaignID: campaignImage.CampaignId,
		FileName:   campaignImage.FileName,
		IsPrimary:  campaignImage.IsPrimary,
	}
	if dataImage.IsPrimary {
		if err := s.repository.ResetPrimaryImage(dataImage.CampaignID); err != nil {
			return CampaignImages{}, err
		}
	}
	result, err := s.repository.CreateImage(dataImage)
	if err != nil {
		return CampaignImages{}, err
	}
	return result, nil
}
