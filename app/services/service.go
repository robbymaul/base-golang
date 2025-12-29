package services

import (
	"context"
	"fmt"
	"path"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
)

type Service struct {
	ctx        context.Context
	config     *config.Config
	repository *repositories.RepositoryContext
}

func NewService(ctx context.Context, r *repositories.RepositoryContext, cfg *config.Config) *Service {

	return &Service{ctx: ctx, config: cfg, repository: r}
}

func (s *Service) mapToConfigurationResponse(conf *models.Configuration) *web.ResponseConfiguration {
	var aggregator *web.AggregatorResponse

	if conf.Aggregator != nil {
		aggregator = s.mapToAggregatorResponse(conf.Aggregator)
	}

	return &web.ResponseConfiguration{
		Id:           conf.Id,
		AggregatorId: conf.AggregatorId,
		Aggregator:   aggregator,
		ConfigName:   conf.ConfigName,
		ConfigValue:  conf.ConfigValue,
		IsActive:     conf.IsActive,
		ConfigJson:   s.mapToConfigJsonResponse(conf.ConfigJson),
		CreatedAt:    conf.CreatedAt,
		UpdatedAt:    conf.UpdatedAt,
	}
}

func (s *Service) mapToAggregatorResponse(aggregator *models.Aggregator) *web.AggregatorResponse {
	return &web.AggregatorResponse{
		Id:          aggregator.Id,
		Name:        aggregator.Name,
		Slug:        aggregator.Slug,
		Description: aggregator.Description,
		IsActive:    aggregator.IsActive,
		Currency:    aggregator.Currency,
		CreatedAt:   aggregator.CreatedAt,
		UpdatedAt:   aggregator.UpdatedAt,
	}
}

func (s *Service) mapToConfigJsonResponse(json models.ConfigJson) *web.ConfigJson {
	return &web.ConfigJson{
		SandboxBaseUrl:               json.SandboxBaseUrl,
		ProductionBaseUrl:            json.ProductionBaseUrl,
		SandboxMerchantId:            helpers.DecryptAES(json.SandboxMerchantId),
		ProductionMerchantId:         helpers.DecryptAES(json.ProductionMerchantId), // json.ProductionMerchantId,
		SandboxMerchantCode:          helpers.DecryptAES(json.SandboxMerchantCode),
		ProductionMerchantCode:       helpers.DecryptAES(json.ProductionMerchantCode),
		SandboxMerchantName:          helpers.DecryptAES(json.SandboxMerchantName),
		ProductionMerchantName:       helpers.DecryptAES(json.ProductionMerchantName),
		SandboxApiKey:                helpers.DecryptAES(json.SandboxApiKey),
		ProductionApiKey:             helpers.DecryptAES(json.ProductionApiKey),
		SandboxServerKey:             helpers.DecryptAES(json.SandboxServerKey),       // json.SandboxServerKey,
		ProductionServerKey:          helpers.DecryptAES(json.ProductionServerKey),    // json.ProductionServerKey,
		SandboxSecretKey:             helpers.DecryptAES(json.SandboxSecretKey),       // json.SandboxSecretKey,
		ProductionSecretKey:          helpers.DecryptAES(json.ProductionSecretKey),    // json.ProductionSecretKey,
		SandboxClientKey:             helpers.DecryptAES(json.SandboxClientKey),       // json.SandboxClientKey,
		ProductionClientKey:          helpers.DecryptAES(json.ProductionClientKey),    // json.ProductionClientKey,
		SandboxSignatureKey:          helpers.DecryptAES(json.SandboxSignatureKey),    // json.SandboxSignatureKey,
		ProductionSignatureKey:       helpers.DecryptAES(json.ProductionSignatureKey), // json.ProductionSignatureKey,
		SandboxCredentialPassword:    helpers.DecryptAES(json.SandboxCredentialPassword),
		ProductionCredentialPassword: helpers.DecryptAES(json.ProductionCredentialPassword),
		ReturnUrl:                    json.ReturnUrl,
	}
}

func (s *Service) mapToDetailPlatformResponse(platform *models.Platforms) *web.DetailPlatformResponse {
	return &web.DetailPlatformResponse{
		Id:              platform.Id,
		Code:            platform.Code,
		Name:            platform.Name,
		Description:     platform.Description,
		ApiKey:          platform.ApiKey,
		SecretKey:       platform.SecretKey,
		IsActive:        platform.IsActive,
		NotificationUrl: platform.NotificationURL,
		CreatedAt:       platform.CreatedAt,
		UpdatedAt:       platform.UpdatedAt,
	}
}

func (s *Service) mapToDetailChannelResponse(channel *models.Channel) *web.DetailChannelResponse {
	channelImages := make([]*web.ImageResponse, 0, len(channel.ChannelImage))

	for _, channelImage := range channel.ChannelImage {
		imageResponse := s.mapToChannelImageResponse(channelImage)
		if imageResponse != nil {
			channelImages = append(channelImages, imageResponse)
		}
	}

	return &web.DetailChannelResponse{
		Id:            channel.Id,
		Code:          channel.Code,
		Name:          channel.Name,
		PaymentMethod: channel.PaymentMethod,
		//TransactionType: channel.TransactionType,
		//Provider:        channel.Provider,
		Currency:      channel.Currency,
		FeeType:       channel.FeeType,
		FeeFixed:      channel.FeeAmount,
		FeePercentage: channel.FeePercentage,
		IsActive:      channel.IsActive,
		ProductName:   channel.ProductName,
		ProductCode:   channel.ProductCode,
		BankName:      channel.BankName,
		BankCode:      channel.BankCode,
		Images:        channelImages,
		Instruction:   channel.Instruction,
	}
}

func (s *Service) mapToChannelImageResponse(channelImage *models.ChannelImage) *web.ImageResponse {
	assetPath, err := s.GetAssetPath(enums.ASSET_LOGO, channelImage.FileName)
	if err != nil {
		return nil
	}

	return &web.ImageResponse{
		FileUrl:   s.urlFile(assetPath),
		SizeType:  channelImage.SizeType,
		Geometric: channelImage.GeometricType,
	}
}

func (s *Service) GetAssetPath(assetType enums.AssetType, name string) (string, error) {
	assetDir := ""

	switch assetType {
	case enums.ASSET_LOGO:
		assetDir = s.config.S3AssetLogo
	default:
		return "", fmt.Errorf("unknown asset type and asset path")
	}

	return path.Join(assetDir, name), nil
}

func (s *Service) urlFile(assetPath string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.S3BucketName, s.config.S3Region, assetPath)
}

func (s *Service) mapToChannelResponseWithCurrency(channels []*models.Channel) map[enums.Currency]map[enums.PaymentMethod][]*web.DetailChannelResponse {
	responseData2 := map[enums.Currency]map[enums.PaymentMethod][]*web.DetailChannelResponse{
		enums.CURRENCY_MYR: map[enums.PaymentMethod][]*web.DetailChannelResponse{
			enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_QRIS:            []*web.DetailChannelResponse{},
			enums.PAYMENT_METEHOD_K_WALLET:       []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_CREDIT_CARD:     []*web.DetailChannelResponse{},
		},
		enums.CURRENCY_IDR: map[enums.PaymentMethod][]*web.DetailChannelResponse{
			enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_QRIS:            []*web.DetailChannelResponse{},
			enums.PAYMENT_METEHOD_K_WALLET:       []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_CREDIT_CARD:     []*web.DetailChannelResponse{},
		},
	}

	for _, ch := range channels {
		channel := s.mapToDetailChannelResponse(ch)

		//if _, ok := responseData[ch.PaymentMethod]; ok {
		//	responseData[ch.PaymentMethod] = append(responseData[ch.PaymentMethod], s.mapToDetailChannelResponse(ch))
		//} else {
		//	responseData[enums.PAYMENT_METHOD_MULTI_PAYMENT] = append(responseData[enums.PAYMENT_METHOD_MULTI_PAYMENT], s.mapToDetailChannelResponse(ch))
		//}

		//response = append(response, s.mapToDetailChannelResponse(ch))

		if _, ok := responseData2[ch.Currency]; ok {
			//responseData2[ch.Currency] = map[enums.PaymentMethod][]*web.DetailChannelResponse{
			//	enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
			//	enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
			//}

			switch ch.PaymentMethod {
			case enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT:
				responseData2[ch.Currency][enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT] = append(responseData2[ch.Currency][enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT], channel)
			case enums.PAYMENT_METHOD_QRIS:
				responseData2[ch.Currency][enums.PAYMENT_METHOD_QRIS] = append(responseData2[ch.Currency][enums.PAYMENT_METHOD_QRIS], channel)
			case enums.PAYMENT_METEHOD_K_WALLET:
				responseData2[ch.Currency][enums.PAYMENT_METEHOD_K_WALLET] = append(responseData2[ch.Currency][enums.PAYMENT_METEHOD_K_WALLET], channel)
			case enums.PAYMENT_METHOD_CREDIT_CARD:
				responseData2[ch.Currency][enums.PAYMENT_METHOD_CREDIT_CARD] = append(responseData2[ch.Currency][enums.PAYMENT_METHOD_CREDIT_CARD], channel)

			default:
				responseData2[ch.Currency][enums.PAYMENT_METHOD_MULTI_PAYMENT] = append(responseData2[ch.Currency][enums.PAYMENT_METHOD_MULTI_PAYMENT], channel)
			}

		}
	}

	return responseData2
}

func (s *Service) mapToChannelResponseWithoutCurrency(channels []*models.Channel, kWallets []*models.KWallet) map[enums.PaymentMethod][]*web.DetailChannelResponse {
	responseData := map[enums.PaymentMethod][]*web.DetailChannelResponse{
		enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
		enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
		enums.PAYMENT_METEHOD_K_WALLET:       []*web.DetailChannelResponse{},
		enums.PAYMENT_METHOD_CREDIT_CARD:     []*web.DetailChannelResponse{},
		enums.PAYMENT_METHOD_QRIS:            []*web.DetailChannelResponse{},
	}

	for _, ch := range channels {
		//response = append(response, s.mapToDetailChannelResponse(ch))

		channel := s.mapToDetailChannelResponse(ch)

		if _, ok := responseData[ch.PaymentMethod]; ok {
			if ch.PaymentMethod == enums.PAYMENT_METEHOD_K_WALLET {
				for _, kWallet := range kWallets {
					responseData[ch.PaymentMethod] = append(responseData[ch.PaymentMethod], s.mapToDetailChannelKWalletResponse(ch, kWallet))
				}
			} else {
				responseData[ch.PaymentMethod] = append(responseData[ch.PaymentMethod], channel)
			}
		} else {
			responseData[enums.PAYMENT_METHOD_MULTI_PAYMENT] = append(responseData[enums.PAYMENT_METHOD_MULTI_PAYMENT], channel)
		}
	}

	return responseData
}

func (s *Service) mapToDetailChannelKWalletResponse(ch *models.Channel, kWallet *models.KWallet) *web.DetailChannelResponse {
	channelImageResponse := make([]*web.ImageResponse, 0, len(ch.ChannelImage))

	for _, channelImage := range ch.ChannelImage {
		channelImageResponse = append(channelImageResponse, s.mapToChannelImageResponse(channelImage))
	}

	return &web.DetailChannelResponse{
		Id:            ch.Id,
		Code:          ch.Code,
		Name:          ch.Name,
		PaymentMethod: ch.PaymentMethod,
		//TransactionType: ch.TransactionType,
		//Provider:        ch.Provider,
		Currency:      ch.Currency,
		FeeType:       ch.FeeType,
		FeeFixed:      ch.FeeAmount,
		FeePercentage: ch.FeePercentage,
		IsActive:      ch.IsActive,
		MemberID:      kWallet.MemberID,
		FullName:      kWallet.FullName,
		NoRekening:    kWallet.NoRekening,
		GenVa:         kWallet.GenVA,
		Balance:       kWallet.Balance.IntPart(),
		Symbol:        kWallet.Symbol,
		Status:        kWallet.Status,
		ProductName:   ch.ProductName,
		ProductCode:   ch.ProductCode,
		BankName:      ch.BankName,
		BankCode:      ch.BankCode,
		Images:        channelImageResponse,
	}
}

func (s *Service) mapToDetailPaymentResponse(payments *models.Payments) *web.DetailPaymentResponse {
	var platform *web.DetailPlatformResponse
	var channel *web.DetailChannelResponse
	var aggregator *web.AggregatorResponse

	if payments.Channel != nil {
		channel = s.mapToDetailChannelResponse(payments.Channel)
	}

	if payments.Platform != nil {
		platform = s.mapToDetailPlatformResponse(payments.Platform)
	}

	if payments.Aggregator != nil {
		aggregator = s.mapToAggregatorResponse(payments.Aggregator)
	}

	return &web.DetailPaymentResponse{
		Id:                   payments.Id,
		TransactionId:        payments.TransactionId,
		OrderId:              payments.OrderId,
		PlatformId:           payments.PlatformId,
		PaymentMethodId:      payments.PaymentMethodId,
		Amount:               payments.Amount.IntPart(),
		FeeAmount:            payments.FeeAmount.IntPart(),
		TotalAmount:          payments.TotalAmount.IntPart(),
		Currency:             payments.Currency,
		Status:               payments.Status,
		CustomerId:           payments.CustomerId,
		CustomerName:         payments.CustomerName,
		CustomerEmail:        payments.CustomerEmail,
		CustomerPhone:        payments.CustomerPhone,
		ReferenceId:          payments.ReferenceId,
		ReferenceType:        payments.ReferenceType,
		GatewayTransactionId: payments.GatewayTransactionId,
		GatewayReference:     payments.GatewayReference,
		GatewayResponse:      payments.GatewayResponse,
		CallbackUrl:          payments.CallbackUrl,
		ReturnUrl:            payments.ReturnUrl,
		ExpiredAt:            payments.ExpiredAt,
		PaidAt:               payments.PaidAt,
		CreatedAt:            payments.CreatedAt,
		UpdatedAt:            payments.UpdatedAt,
		Platform:             platform,
		Channel:              channel,
		Aggregator:           aggregator,
	}
}

func (s *Service) mapToCheckStatusPaymentResponse(payment *models.Payments) *web.CheckStatusPaymentResponse {
	return &web.CheckStatusPaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
	}
}

func (s *Service) mapToChannelResponse(channel *models.Channel) *web.ChannelResponse {
	var methodGroup enums.PaymentMethod

	switch channel.PaymentMethod {
	case enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT:
		methodGroup = enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT
	case enums.PAYMENT_METHOD_CREDIT_CARD:
		methodGroup = enums.PAYMENT_METHOD_CREDIT_CARD
	case enums.PAYMENT_METHOD_QRIS:
		methodGroup = enums.PAYMENT_METHOD_QRIS
	case enums.PAYMENT_METHOD_GOPAY:
		methodGroup = enums.PAYMENT_METHOD_GOPAY
	default:
		methodGroup = enums.PAYMENT_METHOD_MULTI_PAYMENT
	}

	return &web.ChannelResponse{
		Id:            channel.Id,
		Name:          channel.Name,
		MethodGroup:   methodGroup,
		PaymentMethod: channel.PaymentMethod,
		Currency:      channel.Currency,
		FeeType:       channel.FeeType,
		FeeFixed:      channel.FeeAmount,
		FeePercentage: channel.FeePercentage,
		Bank:          channel.BankName,
		Instruction:   channel.Instruction,
		Images:        nil,
	}
}
