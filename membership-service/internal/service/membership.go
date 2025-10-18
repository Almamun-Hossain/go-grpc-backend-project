package service

import (
	"context"

	pb "membership-service/api/membership/v1"
	"membership-service/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type MembershipService struct {
	pb.UnimplementedMembershipServer

	uc *biz.MembershipUsecase
}

func NewMembershipService(uc *biz.MembershipUsecase) *MembershipService {
	return &MembershipService{uc: uc}
}

func (s *MembershipService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	user, accessToken, refreshToken, err := s.uc.Register(ctx, req.Email, req.Password, req.Username, req.ReferralCode)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterReply{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(user.CreatedAt.Add(24 * 3600 * 1000000000)), // 24 hours
	}, nil
}

func (s *MembershipService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, accessToken, refreshToken, err := s.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginReply{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(user.UpdatedAt.Add(24 * 3600 * 1000000000)),
	}, nil
}

func (s *MembershipService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenReply, error) {
	accessToken, refreshToken, err := s.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokenReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.Now(),
	}, nil
}

func (s *MembershipService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileReply, error) {
	// TODO: Extract user ID from JWT token in context
	userID := "user_123"
	
	user, err := s.uc.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileReply{
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			Username:     user.Username,
			Phone:        user.Phone,
			AvatarUrl:    user.AvatarURL,
			ReferralCode: user.ReferralCode,
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (s *MembershipService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	// TODO: Extract user ID from JWT token in context
	userID := "user_123"
	
	user, err := s.uc.UpdateProfile(ctx, userID, req.Username, req.Phone, req.AvatarUrl)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProfileReply{
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			Username:     user.Username,
			Phone:        user.Phone,
			AvatarUrl:    user.AvatarURL,
			ReferralCode: user.ReferralCode,
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (s *MembershipService) GetAssets(ctx context.Context, req *pb.GetAssetsRequest) (*pb.GetAssetsReply, error) {
	// TODO: Extract user ID from JWT token in context
	userID := "user_123"
	
	assets, totalBalance, err := s.uc.GetAssets(ctx, userID)
	if err != nil {
		return nil, err
	}

	pbAssets := make([]*pb.Asset, 0, len(assets))
	for _, asset := range assets {
		pbAssets = append(pbAssets, &pb.Asset{
			Id:        asset.ID,
			UserId:    asset.UserID,
			AssetType: asset.AssetType,
			Balance:   asset.Balance,
			CreatedAt: timestamppb.New(asset.CreatedAt),
			UpdatedAt: timestamppb.New(asset.UpdatedAt),
		})
	}

	return &pb.GetAssetsReply{
		Assets:       pbAssets,
		TotalBalance: totalBalance,
	}, nil
}

func (s *MembershipService) UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.UpdateAssetReply, error) {
	asset, err := s.uc.UpdateAsset(ctx, req.AssetId, req.Amount, req.Operation)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateAssetReply{
		Asset: &pb.Asset{
			Id:        asset.ID,
			UserId:    asset.UserID,
			AssetType: asset.AssetType,
			Balance:   asset.Balance,
			CreatedAt: timestamppb.New(asset.CreatedAt),
			UpdatedAt: timestamppb.New(asset.UpdatedAt),
		},
	}, nil
}

func (s *MembershipService) GetReferrals(ctx context.Context, req *pb.GetReferralsRequest) (*pb.GetReferralsReply, error) {
	// TODO: Extract user ID from JWT token in context
	userID := "user_123"
	
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	referrals, total, totalEarnings, err := s.uc.GetReferrals(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	pbReferrals := make([]*pb.Referral, 0, len(referrals))
	for _, ref := range referrals {
		pbReferrals = append(pbReferrals, &pb.Referral{
			Id:               ref.ID,
			ReferredUserId:   ref.ReferredUserID,
			ReferredUsername: ref.ReferredUsername,
			Earnings:         ref.Earnings,
			Status:           ref.Status,
			CreatedAt:        timestamppb.New(ref.CreatedAt),
		})
	}

	return &pb.GetReferralsReply{
		Referrals:     pbReferrals,
		Total:         int32(total),
		TotalEarnings: totalEarnings,
	}, nil
}

func (s *MembershipService) RequestWithdrawal(ctx context.Context, req *pb.WithdrawalRequest) (*pb.WithdrawalReply, error) {
	// TODO: Extract user ID from JWT token in context
	userID := "user_123"
	
	withdrawal, err := s.uc.RequestWithdrawal(ctx, userID, req.Amount, req.PaymentMethod, req.PaymentDetails)
	if err != nil {
		return nil, err
	}

	return &pb.WithdrawalReply{
		Withdrawal: &pb.Withdrawal{
			Id:            withdrawal.ID,
			UserId:        withdrawal.UserID,
			Amount:        withdrawal.Amount,
			PaymentMethod: withdrawal.PaymentMethod,
			Status:        withdrawal.Status,
			CreatedAt:     timestamppb.New(withdrawal.CreatedAt),
			UpdatedAt:     timestamppb.New(withdrawal.UpdatedAt),
		},
	}, nil
}

func (s *MembershipService) GetWithdrawals(ctx context.Context, req *pb.GetWithdrawalsRequest) (*pb.GetWithdrawalsReply, error) {
	// TODO: Extract user ID from JWT token in context
	userID := "user_123"
	
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	withdrawals, total, err := s.uc.GetWithdrawals(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	pbWithdrawals := make([]*pb.Withdrawal, 0, len(withdrawals))
	for _, w := range withdrawals {
		pbWithdrawals = append(pbWithdrawals, &pb.Withdrawal{
			Id:            w.ID,
			UserId:        w.UserID,
			Amount:        w.Amount,
			PaymentMethod: w.PaymentMethod,
			Status:        w.Status,
			CreatedAt:     timestamppb.New(w.CreatedAt),
			UpdatedAt:     timestamppb.New(w.UpdatedAt),
		})
	}

	return &pb.GetWithdrawalsReply{
		Withdrawals: pbWithdrawals,
		Total:       int32(total),
	}, nil
}

