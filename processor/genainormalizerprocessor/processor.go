package genainormalizerprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

type tracesProcessor struct {
	logger *zap.Logger
	next   consumer.Traces

	overwrite    bool
	dropOriginal bool
	mappings     map[string]string
}

func newTracesProcessor(_ context.Context, settings processor.CreateSettings, cfg *Config, next consumer.Traces) (*tracesProcessor, error) {
	if next == nil {
		return nil, fmt.Errorf("next consumer is nil")
	}

	m := map[string]string{}
	if cfg.EnableDefaults {
		for k, v := range defaultMappings {
			m[k] = v
		}
	}
	for k, v := range cfg.Mappings {
		m[k] = v
	}

	return &tracesProcessor{
		logger:       settings.Logger,
		next:         next,
		overwrite:    cfg.Overwrite,
		dropOriginal: cfg.DropOriginal,
		mappings:     m,
	}, nil
}

func (p *tracesProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (p *tracesProcessor) Start(context.Context, component.Host) error { return nil }
func (p *tracesProcessor) Shutdown(context.Context) error              { return nil }

func (p *tracesProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	rs := td.ResourceSpans()
	for i := 0; i < rs.Len(); i++ {
		ss := rs.At(i).ScopeSpans()
		for j := 0; j < ss.Len(); j++ {
			spans := ss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				p.applyMappings(span.Attributes())

				events := span.Events()
				for e := 0; e < events.Len(); e++ {
					p.applyMappings(events.At(e).Attributes())
				}
			}
		}
	}

	return p.next.ConsumeTraces(ctx, td)
}

func (p *tracesProcessor) applyMappings(attrs pcommon.Map) {
	// Iterate over mappings instead of attributes to avoid iterator invalidation
	// when deleting keys.
	for src, dst := range p.mappings {
		val, ok := attrs.Get(src)
		if !ok {
			continue
		}

		if src == dst {
			continue
		}

		if _, exists := attrs.Get(dst); exists && !p.overwrite {
			// Destination already present and overwrite disabled.
			continue
		}

		// Copy value.
		val.CopyTo(attrs.PutEmpty(dst))

		// Optionally drop original.
		if p.dropOriginal {
			attrs.Remove(src)
		}
	}
}
