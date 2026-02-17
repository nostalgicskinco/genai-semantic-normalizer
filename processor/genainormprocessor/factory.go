package genainormprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

const (
	typeStr   = "genai_semantic_normalizer"
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for the genai semantic normalizer processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType(typeStr),
		func() component.Config { return createDefaultConfig() },
		processor.WithTraces(createTracesProcessor, stability),
	)
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	pCfg := cfg.(*Config)
	return newNormalizerProcessor(set.Logger, pCfg, nextConsumer), nil
}