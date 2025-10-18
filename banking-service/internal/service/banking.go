package service

import (
	"context"

	pb "banking-service/api/banking/v1"
	"banking-service/internal/biz"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type BankingService struct {
	pb.UnimplementedBankingServer

	uc *biz.BankingUsecase
}

func NewBankingService(uc *biz.BankingUsecase) *BankingService {
	return &BankingService{uc: uc}
}

func (s *BankingService) GetCardInfo(ctx context.Context, req *pb.GetCardInfoRequest) (*pb.GetCardInfoReply, error) {
	card, err := s.uc.GetCardInfo(ctx, req.CardId)
	if err != nil {
		return nil, err
	}

	return &pb.GetCardInfoReply{
		Card: &pb.Card{
			Id:               card.ID,
			UserId:           card.UserID,
			CardNumberMasked: card.CardNumberMasked,
			CardType:         card.CardType,
			CardBrand:        card.CardBrand,
			Status:           card.Status,
			ExpiryDate:       timestamppb.New(card.ExpiryDate),
			CreatedAt:        timestamppb.New(card.CreatedAt),
		},
	}, nil
}

func (s *BankingService) ListCards(ctx context.Context, req *pb.ListCardsRequest) (*pb.ListCardsReply, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	cards, total, err := s.uc.ListCards(ctx, req.UserId, page, pageSize)
	if err != nil {
		return nil, err
	}

	pbCards := make([]*pb.Card, 0, len(cards))
	for _, card := range cards {
		pbCards = append(pbCards, &pb.Card{
			Id:               card.ID,
			UserId:           card.UserID,
			CardNumberMasked: card.CardNumberMasked,
			CardType:         card.CardType,
			CardBrand:        card.CardBrand,
			Status:           card.Status,
			ExpiryDate:       timestamppb.New(card.ExpiryDate),
			CreatedAt:        timestamppb.New(card.CreatedAt),
		})
	}

	return &pb.ListCardsReply{
		Cards: pbCards,
		Total: int32(total),
	}, nil
}

func (s *BankingService) GetCreditCardBalance(ctx context.Context, req *pb.GetCreditCardBalanceRequest) (*pb.GetCreditCardBalanceReply, error) {
	balance, err := s.uc.GetCardBalance(ctx, req.CardId)
	if err != nil {
		return nil, err
	}

	return &pb.GetCreditCardBalanceReply{
		CardId:             balance.CardID,
		AvailableBalance:   balance.AvailableBalance,
		CreditLimit:        balance.CreditLimit,
		OutstandingBalance: balance.OutstandingBalance,
		Currency:           balance.Currency,
		LastUpdated:        timestamppb.New(balance.LastUpdated),
	}, nil
}

func (s *BankingService) GetPhysicalCardDetails(ctx context.Context, req *pb.GetPhysicalCardDetailsRequest) (*pb.GetPhysicalCardDetailsReply, error) {
	details, err := s.uc.GetPhysicalCardDetails(ctx, req.CardId)
	if err != nil {
		return nil, err
	}

	return &pb.GetPhysicalCardDetailsReply{
		CardId:           details.CardID,
		CardNumberMasked: details.CardNumberMasked,
		CardholderName:   details.CardholderName,
		ExpiryDate:       details.ExpiryDate,
		CvvMasked:        details.CVVMasked,
		CardType:         details.CardType,
		CardBrand:        details.CardBrand,
	}, nil
}

func (s *BankingService) InitiateRemittance(ctx context.Context, req *pb.InitiateRemittanceRequest) (*pb.InitiateRemittanceReply, error) {
	remittance, err := s.uc.InitiateRemittance(
		ctx,
		req.FromAccount,
		req.ToAccount,
		req.Amount,
		req.FromCurrency,
		req.ToCurrency,
		req.BeneficiaryName,
		req.BeneficiaryBank,
		req.Purpose,
	)
	if err != nil {
		return nil, err
	}

	totalAmount := remittance.Amount + remittance.Fee

	return &pb.InitiateRemittanceReply{
		TransactionId: remittance.TransactionID,
		Status:        remittance.Status,
		Amount:        remittance.Amount,
		ExchangeRate:  remittance.ExchangeRate,
		Fee:           remittance.Fee,
		TotalAmount:   totalAmount,
		CreatedAt:     timestamppb.New(remittance.CreatedAt),
	}, nil
}

func (s *BankingService) GetRemittanceStatus(ctx context.Context, req *pb.GetRemittanceStatusRequest) (*pb.GetRemittanceStatusReply, error) {
	remittance, err := s.uc.GetRemittanceStatus(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	return &pb.GetRemittanceStatusReply{
		Remittance: &pb.Remittance{
			TransactionId:   remittance.TransactionID,
			FromAccount:     remittance.FromAccount,
			ToAccount:       remittance.ToAccount,
			Amount:          remittance.Amount,
			FromCurrency:    remittance.FromCurrency,
			ToCurrency:      remittance.ToCurrency,
			ExchangeRate:    remittance.ExchangeRate,
			Fee:             remittance.Fee,
			Status:          remittance.Status,
			BeneficiaryName: remittance.BeneficiaryName,
			BeneficiaryBank: remittance.BeneficiaryBank,
			CreatedAt:       timestamppb.New(remittance.CreatedAt),
			UpdatedAt:       timestamppb.New(remittance.UpdatedAt),
		},
	}, nil
}

func (s *BankingService) ListRemittances(ctx context.Context, req *pb.ListRemittancesRequest) (*pb.ListRemittancesReply, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	remittances, total, err := s.uc.ListRemittances(ctx, req.UserId, page, pageSize, req.Status)
	if err != nil {
		return nil, err
	}

	pbRemittances := make([]*pb.Remittance, 0, len(remittances))
	for _, r := range remittances {
		pbRemittances = append(pbRemittances, &pb.Remittance{
			TransactionId:   r.TransactionID,
			FromAccount:     r.FromAccount,
			ToAccount:       r.ToAccount,
			Amount:          r.Amount,
			FromCurrency:    r.FromCurrency,
			ToCurrency:      r.ToCurrency,
			ExchangeRate:    r.ExchangeRate,
			Fee:             r.Fee,
			Status:          r.Status,
			BeneficiaryName: r.BeneficiaryName,
			BeneficiaryBank: r.BeneficiaryBank,
			CreatedAt:       timestamppb.New(r.CreatedAt),
			UpdatedAt:       timestamppb.New(r.UpdatedAt),
		})
	}

	return &pb.ListRemittancesReply{
		Remittances: pbRemittances,
		Total:       int32(total),
	}, nil
}

func (s *BankingService) GetExchangeRate(ctx context.Context, req *pb.GetExchangeRateRequest) (*pb.GetExchangeRateReply, error) {
	rate, err := s.uc.GetExchangeRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, err
	}

	return &pb.GetExchangeRateReply{
		FromCurrency: rate.FromCurrency,
		ToCurrency:   rate.ToCurrency,
		Rate:         rate.Rate,
		Timestamp:    timestamppb.New(rate.Timestamp),
	}, nil
}

