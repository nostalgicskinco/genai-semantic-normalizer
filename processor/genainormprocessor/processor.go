// Copyright 2024 Nostalgic Skin Co.
// SPDX-License-Identifier: Apache-2.0

package genainormprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

type normProc struct {
	logger *zap.Logger
	cfg    *Config
	next   consumer.Traces
	mapper *mapper
}

func newProcessor(
	_ context.Context,
	set processor.Settings,
	cfg *Config,
	next consumer.Traces,
) (*normProc, error) {
	p := &normProc{
		logger: set.Logger,
		cfg:    cfg,
		next:   next,
		mapper: newMapper(cfg),
	}
	set.Logger.Info("genai_semantic_normalizer initialized",
		zap.Int("mapping_count", len(p.mapper.table)),
		zap.Bool("overwrite", cfg.Overwrite),
		zap.Bool("drop_original", cfg.DropOriginal),
	)
	return p, nil
}

func (p *normProc) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (p *normProc) Start(_ context.Context, _ component.Host) error {
	return nil
}

func (p *normProc) Shutdown(_ context.Context) error {
	return nil
}

// ConsumeTraces normalizes attributes on every span and span event.
func (p *normProc) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	rs := td.ResourceSpans()
	for i := 0; i < rs.Len(); i++ {
		scopeSpans := rs.At(i).ScopeSpans()
		for j := 0; j < scopeSpans.Len(); j++ {
			spans := scopeSpans.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				s := spans.At(k)

				// Normalize span attributes
				p.mapper.apply(s.Attributes())

				// Normalize span event attributes (prompt/completion events)
				events := s.Events()
				for e := 0; e < events.Len(); e++ {
					p.mapper.apply(events.At(e).Attributes())
				}
			}
		}
	}
	return p.next.ConsumeTraces(ctx, td)
}
