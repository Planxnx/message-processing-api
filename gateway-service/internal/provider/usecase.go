package provider

import (
	"context"
	"crypto/rand"
	"log"
	"time"

	"encoding/base64"

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

func (*ProviderUsercase) getNewToken() string {
	newBytes := make([]byte, 64)
	rand.Read(newBytes)
	return base64.StdEncoding.EncodeToString(newBytes)
}

func (pU *ProviderUsercase) GetProviderByID(ctx context.Context, id string) (*ProviderData, error) {
	providerData := &ProviderData{}
	err := pU.ProviderCollection.Find(ctx, bson.M{
		"id": id,
	}).One(providerData)
	if err != nil {
		return nil, err
	}
	log.Println(providerData)

	return providerData, nil
}

func (pU *ProviderUsercase) CreateNewProvider(ctx context.Context, providerData *ProviderData) (*CreateProviderResult, error) {
	token := pU.getNewToken()
	//TODO: hash token with secret
	providerData.Token = token

	providerData.CreatedAt = time.Now()
	providerData.UpdatedAt = providerData.CreatedAt

	_, err := pU.ProviderCollection.InsertOne(ctx, providerData)
	if err != nil {
		return nil, err
	}
	return &CreateProviderResult{
		ID:    providerData.ID,
		Token: token,
	}, nil
}
