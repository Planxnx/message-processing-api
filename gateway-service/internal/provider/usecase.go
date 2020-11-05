package provider

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ProviderUsercase struct {
	ProviderCollection *qmgo.Collection
}

func New(pC *qmgo.Collection) *ProviderUsercase {
	return &ProviderUsercase{
		ProviderCollection: pC,
	}
}

func (*ProviderUsercase) getNewSecret() string {
	newBytes := make([]byte, 8)
	rand.Read(newBytes)
	return fmt.Sprintf("%x", newBytes)
}

func (pU *ProviderUsercase) GetProviderByID(ctx context.Context, id string) (*ProviderData, error) {
	providerData := &ProviderData{}
	err := pU.ProviderCollection.Find(ctx, bson.M{
		"id": id,
	}).One(providerData)
	if err != nil {
		return nil, err
	}
	return providerData, nil
}

func (pU *ProviderUsercase) CreateNewProvider(ctx context.Context, providerData *ProviderData) (*CreateProviderResult, error) {
	secret := pU.getNewSecret()
	providerData.Secret = secret
	_, err := pU.ProviderCollection.InsertOne(ctx, providerData)
	if err != nil {
		return nil, err
	}
	return &CreateProviderResult{
		ID:     providerData.ID,
		Secret: secret,
	}, nil
}
