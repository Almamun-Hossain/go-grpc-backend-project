package service

import (
	"context"

	pb "website-service/api/website/v1"
	"website-service/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type WebsiteService struct {
	pb.UnimplementedWebsiteServer

	uc *biz.WebsiteUsecase
}

func NewWebsiteService(uc *biz.WebsiteUsecase) *WebsiteService {
	return &WebsiteService{uc: uc}
}

func (s *WebsiteService) GetSiteSettings(ctx context.Context, req *pb.GetSiteSettingsRequest) (*pb.GetSiteSettingsReply, error) {
	settings, err := s.uc.GetSiteSettings(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetSiteSettingsReply{
		Settings: &pb.SiteSettings{
			Id:              settings.ID,
			Title:           settings.Title,
			Description:     settings.Description,
			Keywords:        settings.Keywords,
			LogoUrl:         settings.LogoURL,
			ContactEmail:    settings.ContactEmail,
			ContactPhone:    settings.ContactPhone,
			MaintenanceMode: settings.MaintenanceMode,
			CreatedAt:       timestamppb.New(settings.CreatedAt),
			UpdatedAt:       timestamppb.New(settings.UpdatedAt),
		},
	}, nil
}

func (s *WebsiteService) UpdateSiteSettings(ctx context.Context, req *pb.UpdateSiteSettingsRequest) (*pb.UpdateSiteSettingsReply, error) {
	settings, err := s.uc.UpdateSiteSettings(
		ctx,
		req.Title,
		req.Description,
		req.Keywords,
		req.ContactEmail,
		req.ContactPhone,
		req.MaintenanceMode,
	)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSiteSettingsReply{
		Settings: &pb.SiteSettings{
			Id:              settings.ID,
			Title:           settings.Title,
			Description:     settings.Description,
			Keywords:        settings.Keywords,
			LogoUrl:         settings.LogoURL,
			ContactEmail:    settings.ContactEmail,
			ContactPhone:    settings.ContactPhone,
			MaintenanceMode: settings.MaintenanceMode,
			CreatedAt:       timestamppb.New(settings.CreatedAt),
			UpdatedAt:       timestamppb.New(settings.UpdatedAt),
		},
	}, nil
}

func (s *WebsiteService) UpdateSiteTitle(ctx context.Context, req *pb.UpdateSiteTitleRequest) (*pb.UpdateSiteTitleReply, error) {
	title, updatedAt, err := s.uc.UpdateSiteTitle(ctx, req.Title)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSiteTitleReply{
		Title:     title,
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *WebsiteService) UploadLogo(ctx context.Context, req *pb.UploadLogoRequest) (*pb.UploadLogoReply, error) {
	logoURL, uploadedAt, err := s.uc.UploadLogo(ctx, req.FileData, req.FileName, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &pb.UploadLogoReply{
		LogoUrl:    logoURL,
		UploadedAt: timestamppb.New(uploadedAt),
	}, nil
}

func (s *WebsiteService) GetLogo(ctx context.Context, req *pb.GetLogoRequest) (*pb.GetLogoReply, error) {
	logoURL, uploadedAt, err := s.uc.GetLogo(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetLogoReply{
		LogoUrl:    logoURL,
		UploadedAt: timestamppb.New(uploadedAt),
	}, nil
}

func (s *WebsiteService) DeleteLogo(ctx context.Context, req *pb.DeleteLogoRequest) (*pb.DeleteLogoReply, error) {
	err := s.uc.DeleteLogo(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteLogoReply{
		Success: true,
	}, nil
}

func (s *WebsiteService) GetDynamicSettings(ctx context.Context, req *pb.GetDynamicSettingsRequest) (*pb.GetDynamicSettingsReply, error) {
	settings, err := s.uc.GetDynamicSettings(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetDynamicSettingsReply{
		Settings: settings,
	}, nil
}

func (s *WebsiteService) UpdateDynamicSetting(ctx context.Context, req *pb.UpdateDynamicSettingRequest) (*pb.UpdateDynamicSettingReply, error) {
	setting, err := s.uc.UpdateDynamicSetting(ctx, req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateDynamicSettingReply{
		Setting: &pb.DynamicSetting{
			Key:       setting.Key,
			Value:     setting.Value,
			CreatedAt: timestamppb.New(setting.CreatedAt),
			UpdatedAt: timestamppb.New(setting.UpdatedAt),
		},
	}, nil
}

func (s *WebsiteService) DeleteDynamicSetting(ctx context.Context, req *pb.DeleteDynamicSettingRequest) (*pb.DeleteDynamicSettingReply, error) {
	err := s.uc.DeleteDynamicSetting(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteDynamicSettingReply{
		Success: true,
	}, nil
}

