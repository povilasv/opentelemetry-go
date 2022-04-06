package views

import (
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

type (
	View struct {
		cfg Config
	}

	Config struct {
		// Matchers for the instrument
		instrumentName       string
		instrumentNameRegexp *regexp.Regexp
		instrumentKind       sdkinstrument.Kind
		numberKind           number.Kind
		library              instrumentation.Library

		// Properties of the view
		keys        []attribute.Key // nil implies all keys, []attribute.Key{} implies none
		name        string
		description string
		aggregation aggregation.Kind
		temporality aggregation.Temporality
		acfg        aggregator.Config
	}

	Option func(cfg *Config)
)

const (
	unsetKind       = sdkinstrument.Kind(-1)
	unsetNumberKind = number.Kind(-1)
)

// Matchers

func MatchInstrumentName(name string) Option {
	return func(cfg *Config) {
		cfg.instrumentName = name
	}
}

func MatchInstrumentNameRegexp(re *regexp.Regexp) Option {
	return func(cfg *Config) {
		cfg.instrumentNameRegexp = re
	}
}

func MatchKind(k sdkinstrument.Kind) Option {
	return func(cfg *Config) {
		cfg.instrumentKind = k
	}
}

func MatchNumberKind(k number.Kind) Option {
	return func(cfg *Config) {
		cfg.numberKind = k
	}
}

func MatchInstrumentationLibrary(lib instrumentation.Library) Option {
	return func(cfg *Config) {
		cfg.library = lib
	}
}

// Properties

func WithKeys(keys []attribute.Key) Option {
	return func(cfg *Config) {
		if keys == nil {
			cfg.keys = nil
		} else {
			cfg.keys = append(cfg.keys, keys...)
		}
	}
}

func WithName(name string) Option {
	return func(cfg *Config) {
		cfg.name = name
	}
}

func WithDescription(desc string) Option {
	return func(cfg *Config) {
		cfg.description = desc
	}
}

func WithAggregation(kind aggregation.Kind) Option {
	return func(cfg *Config) {
		cfg.aggregation = kind
	}
}

func WithTemporality(tempo aggregation.Temporality) Option {
	return func(cfg *Config) {
		cfg.temporality = tempo
	}
}

func WithAggregatorConfig(acfg aggregator.Config) Option {
	return func(cfg *Config) {
		cfg.acfg = acfg
	}
}

func New(opts ...Option) View {
	cfg := Config{
		instrumentKind: unsetKind,
		numberKind:     unsetNumberKind,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return View{
		cfg: cfg,
	}
}

// IsSingleInstrument is a requirement when HasName().
func (v View) IsSingleInstrument() bool {
	return v.cfg.instrumentName != ""
}

// HasName implies IsSingleInstrument SHOULD be required.
func (v View) HasName() bool {
	return v.cfg.name != ""
}

func (v View) Name() string {
	return v.cfg.name
}

func (v View) Keys() []attribute.Key {
	return v.cfg.keys
}

func (v View) Description() string {
	return v.cfg.description
}

func (v View) Aggregation() aggregation.Kind {
	return v.cfg.aggregation
}

func (v View) Temporality() aggregation.Temporality {
	return v.cfg.temporality
}

func (v View) AggregatorConfig() aggregator.Config {
	return v.cfg.acfg
}

func stringMismatch(test, value string) bool {
	return test != "" && test != value
}

func ikindMismatch(test, value sdkinstrument.Kind) bool {
	return test != unsetKind && test != value
}

func nkindMismatch(test, value number.Kind) bool {
	return test != unsetNumberKind && test != value
}

func regexpMismatch(test *regexp.Regexp, value string) bool {
	return test != nil && test.MatchString(value)
}

func (v View) Matches(lib instrumentation.Library, desc sdkinstrument.Descriptor) bool {
	return !stringMismatch(v.cfg.library.Name, lib.Name) &&
		!stringMismatch(v.cfg.library.Version, lib.Version) &&
		!stringMismatch(v.cfg.library.SchemaURL, lib.SchemaURL) &&
		!stringMismatch(v.cfg.instrumentName, desc.Name) &&
		!ikindMismatch(v.cfg.instrumentKind, desc.Kind) &&
		!nkindMismatch(v.cfg.numberKind, desc.NumberKind) &&
		!regexpMismatch(v.cfg.instrumentNameRegexp, desc.Name)
}
