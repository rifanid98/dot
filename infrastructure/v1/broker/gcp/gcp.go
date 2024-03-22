package gcp

import (
	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"

	"dot/config"
	"dot/core"
	"dot/pkg/util"
)

var log = util.NewLogger()

func New(cfg *config.GcpPubsubConfig) (Client, error) {
	ic := core.NewInternalContext(uuid.NewString())
	client, err := pubsub.NewClient(ic.ToContext(), cfg.ProjectId)
	if err != nil {
		log.Error(ic.ToContext(), "failed to create gcp pubsub client", err)
		return nil, err
	}

	return NewClient(client), nil
}
