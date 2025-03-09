package topic

import (
	"github.com/bitcoin-sv/go-sdk/overlay"
)

type TopicManager interface {
	IdentifyAdmissableOutputs(ctx overlay.SubmitContext) (*overlay.Admittance, error)
	IdentifyNeededInputs(ctx overlay.SubmitContext) ([]*overlay.Outpoint, error)
	GetDependencies() []string
	GetDocumentation() string
	GetMetaData() overlay.MetaData
}

type BaseTopicManager struct{}

func (b *BaseTopicManager) IdentifyAdmissableOutputs(ctx overlay.SubmitContext) (*overlay.Admittance, error) {
	return nil, nil
}
func (b *BaseTopicManager) IdentifyNeededInputs(ctx overlay.SubmitContext) (*overlay.Admittance, error) {
	return nil, nil
}
func (b *BaseTopicManager) GetDependencies() []string {
	return []string{}
}
