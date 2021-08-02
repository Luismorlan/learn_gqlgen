package graph

import "learn_gqlgen/graph/model"

//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	message         *model.Message
	messageChannels map[string]chan *model.Message
}

func NewResolver() *Resolver {
	return &Resolver{
		message:         nil,
		messageChannels: map[string]chan *model.Message{},
	}
}
