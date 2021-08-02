package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"learn_gqlgen/graph/generated"
	"learn_gqlgen/graph/model"
	"log"

	"github.com/thanhpk/randstr"
)

func (r *mutationResolver) SetMessage(ctx context.Context, input model.NewMessage) (*model.Message, error) {
	m := &model.Message{
		Message: input.Message,
	}
	r.message = m
	log.Println("message channel length: ", len(r.messageChannels))
	for key, ch := range r.messageChannels {
		log.Println("pushed message to channel", key)
		ch <- m
	}
	return m, nil
}

func (r *queryResolver) Message(ctx context.Context) (*model.Message, error) {
	return r.message, nil
}

func (r *subscriptionResolver) MessagesSubscription(ctx context.Context) (<-chan *model.Message, error) {
	log.Println("reached here")
	token := randstr.Hex(16)
	mc := make(chan *model.Message, 1)
	r.messageChannels[token] = mc

	go func() {
		<-ctx.Done()
		delete(r.messageChannels, token)
		log.Println("Subscription deleted: ", token)
	}()

	log.Println("Subscription created: ", token)

	return mc, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
