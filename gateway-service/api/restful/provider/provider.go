package provider

import (
	providerusecase "github.com/Planxnx/message-processing-api/gateway-service/internal/provider"
)

type ProviderHandler struct {
	ProviderUsecase *providerusecase.ProviderUsercase
}

func New(pU *providerusecase.ProviderUsercase) *ProviderHandler {
	return &ProviderHandler{
		ProviderUsecase: pU,
	}
}
