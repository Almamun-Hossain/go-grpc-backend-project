package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// SiteSettings represents the main site settings
type SiteSettings struct {
	ID              string
	Title           string
	Description     string
	Keywords        string
	LogoURL         string
	ContactEmail    string
	ContactPhone    string
	MaintenanceMode bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// DynamicSetting represents a key-value setting
type DynamicSetting struct {
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WebsiteRepo defines the interface for website data operations
type WebsiteRepo interface {
	// Site Settings
	GetSiteSettings(ctx context.Context) (*SiteSettings, error)
	UpdateSiteSettings(ctx context.Context, settings *SiteSettings) (*SiteSettings, error)
	
	// Logo Management
	SaveLogoURL(ctx context.Context, url string) error
	GetLogoURL(ctx context.Context) (string, time.Time, error)
	DeleteLogo(ctx context.Context) error
	
	// Dynamic Settings
	GetDynamicSettings(ctx context.Context) (map[string]string, error)
	GetDynamicSetting(ctx context.Context, key string) (*DynamicSetting, error)
	UpdateDynamicSetting(ctx context.Context, setting *DynamicSetting) (*DynamicSetting, error)
	DeleteDynamicSetting(ctx context.Context, key string) error
}

// WebsiteUsecase handles website business logic
type WebsiteUsecase struct {
	repo WebsiteRepo
	log  *log.Helper
}

// NewWebsiteUsecase creates a new WebsiteUsecase
func NewWebsiteUsecase(repo WebsiteRepo, logger log.Logger) *WebsiteUsecase {
	return &WebsiteUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// GetSiteSettings retrieves the site settings
func (uc *WebsiteUsecase) GetSiteSettings(ctx context.Context) (*SiteSettings, error) {
	return uc.repo.GetSiteSettings(ctx)
}

// UpdateSiteSettings updates the site settings
func (uc *WebsiteUsecase) UpdateSiteSettings(ctx context.Context, title, description, keywords, contactEmail, contactPhone string, maintenanceMode bool) (*SiteSettings, error) {
	settings := &SiteSettings{
		Title:           title,
		Description:     description,
		Keywords:        keywords,
		ContactEmail:    contactEmail,
		ContactPhone:    contactPhone,
		MaintenanceMode: maintenanceMode,
		UpdatedAt:       time.Now(),
	}
	
	return uc.repo.UpdateSiteSettings(ctx, settings)
}

// UpdateSiteTitle updates only the site title
func (uc *WebsiteUsecase) UpdateSiteTitle(ctx context.Context, title string) (string, time.Time, error) {
	settings, err := uc.repo.GetSiteSettings(ctx)
	if err != nil {
		return "", time.Time{}, err
	}
	
	settings.Title = title
	settings.UpdatedAt = time.Now()
	
	updated, err := uc.repo.UpdateSiteSettings(ctx, settings)
	if err != nil {
		return "", time.Time{}, err
	}
	
	return updated.Title, updated.UpdatedAt, nil
}

// UploadLogo handles logo upload
func (uc *WebsiteUsecase) UploadLogo(ctx context.Context, fileData []byte, fileName, contentType string) (string, time.Time, error) {
	// TODO: Implement actual file upload to storage (S3, local, etc.)
	// For now, return a placeholder URL
	logoURL := "https://example.com/logos/" + fileName
	
	err := uc.repo.SaveLogoURL(ctx, logoURL)
	if err != nil {
		return "", time.Time{}, err
	}
	
	return logoURL, time.Now(), nil
}

// GetLogo retrieves the logo URL
func (uc *WebsiteUsecase) GetLogo(ctx context.Context) (string, time.Time, error) {
	return uc.repo.GetLogoURL(ctx)
}

// DeleteLogo removes the logo
func (uc *WebsiteUsecase) DeleteLogo(ctx context.Context) error {
	// TODO: Delete the actual file from storage
	return uc.repo.DeleteLogo(ctx)
}

// GetDynamicSettings retrieves all dynamic settings
func (uc *WebsiteUsecase) GetDynamicSettings(ctx context.Context) (map[string]string, error) {
	return uc.repo.GetDynamicSettings(ctx)
}

// UpdateDynamicSetting updates a dynamic setting
func (uc *WebsiteUsecase) UpdateDynamicSetting(ctx context.Context, key, value string) (*DynamicSetting, error) {
	setting := &DynamicSetting{
		Key:       key,
		Value:     value,
		UpdatedAt: time.Now(),
	}
	
	// Check if setting exists
	existing, err := uc.repo.GetDynamicSetting(ctx, key)
	if err == nil && existing != nil {
		setting.CreatedAt = existing.CreatedAt
	} else {
		setting.CreatedAt = time.Now()
	}
	
	return uc.repo.UpdateDynamicSetting(ctx, setting)
}

// DeleteDynamicSetting deletes a dynamic setting
func (uc *WebsiteUsecase) DeleteDynamicSetting(ctx context.Context, key string) error {
	return uc.repo.DeleteDynamicSetting(ctx, key)
}

