package campaign

import "gorm.io/gorm"

type Repository interface {
	Save(campaign Campaign) (Campaign, error)
	Get(campaign Campaign) (Campaign, error)
	Gets(campaign Campaign) ([]Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	Delete(campaign Campaign) error
	CreateImage(image CampaignImages) (CampaignImages, error)
	ResetPrimaryImage(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	if err := r.db.Create(&campaign).Error; err != nil {
		return Campaign{}, err
	}
	return campaign, nil
}

func (r *repository) Get(campaign Campaign) (Campaign, error) {
	var response Campaign
	tx := r.db.Preload("CampaignImages").Preload("User")
	if campaign.ID != 0 {
		tx = tx.Where("id = ?", campaign.ID)
	}
	if campaign.Slug != "" {
		tx = tx.Where("slug = ?", campaign.Slug)
	}
	if campaign.UserId != 0 {
		tx = tx.Where("user_id = ?", campaign.UserId)
	}
	if err := tx.Find(&response).Error; err != nil {
		return Campaign{}, err
	}
	if len(response.CampaignImages) > 1 {
		response.ImagesUrl = response.CampaignImages[0].FileName
	}
	return response, nil
}

func (r *repository) Gets(campaign Campaign) ([]Campaign, error) {
	var response []Campaign
	tx := r.db
	if campaign.Name != "" {
		tx.Where("name LIKE '%?%'", campaign.Name)
	}
	if campaign.UserId != 0 {
		tx.Where("user_id = ?", campaign.UserId)
	}
	if err := tx.Preload("CampaignImages", "is_primary = true").Find(&response).Error; err != nil {
		return []Campaign{}, err
	}
	return response, nil
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	if err := r.db.Save(&campaign).Error; err != nil {
		return Campaign{}, err
	}
	return campaign, nil
}

func (r *repository) Delete(campaign Campaign) error {
	if err := r.db.Delete(&campaign).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateImage(image CampaignImages) (CampaignImages, error) {
	if err := r.db.Create(&image).Error; err != nil {
		return CampaignImages{}, err
	}
	return image, nil
}

func (r *repository) UpdateImage(image CampaignImages) (CampaignImages, error) {
	if err := r.db.Save(&image).Error; err != nil {
		return CampaignImages{}, err
	}
	return image, nil
}

func (r *repository) ResetPrimaryImage(id int) error {
	if err := r.db.Model(&CampaignImages{}).Where("campaign_id = ?", id).Update("is_primary", false).Error; err != nil {
		return err
	}
	return nil
}
