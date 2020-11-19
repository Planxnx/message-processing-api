package provider

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

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
	newBytes := make([]byte, 16)
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
	log.Println(providerData)

	return providerData, nil
}

func (pU *ProviderUsercase) CreateNewProvider(ctx context.Context, providerData *ProviderData) (*CreateProviderResult, error) {
	token := pU.getNewToken()
	//TODO: hash token with secret
	providerData.Token = token
	_, err := pU.ProviderCollection.InsertOne(ctx, providerData)
	if err != nil {
		return nil, err
	}
	return &CreateProviderResult{
		ID:    providerData.ID,
		Token: token,
	}, nil
}
