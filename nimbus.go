// Package nimbus provides utilities for quickly creating [Pulumi] stacks.
//
// [Pulumi]: https://www.pulumi.com/
package nimbus

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Option[I pulumi.Input] func(*Params[I])

type Params[I pulumi.Input] struct {
	Name    string
	Input   I
	Options []pulumi.ResourceOption
}

type Builder[I pulumi.Input, R pulumi.Resource] struct {
	new    func(*pulumi.Context, string, *I, ...pulumi.ResourceOption) (R, error)
	params Params[I]
}

func (b *Builder[I, R]) Configure(opts ...Option[I]) *Builder[I, R] {
	for _, apply := range opts {
		apply(&b.params)
	}

	return b
}

func (b Builder[I, R]) Register(ctx *pulumi.Context) (R, error) {
	return b.new(ctx, b.params.Name, &b.params.Input, b.params.Options...)
}

// Build initializes a new builder with the provided resource factory function.
func Build[I pulumi.Input, R pulumi.Resource](new func(*pulumi.Context, string, *I, ...pulumi.ResourceOption) (R, error)) *Builder[I, R] {
	return &Builder[I, R]{new: new}
}

func WithName[I pulumi.Input](name string) Option[I] {
	return func(p *Params[I]) {
		p.Name = name
	}
}

type Config[I pulumi.Input] interface {
	Configure(input *I)
}

func WithConfig[I pulumi.Input](config Config[I]) Option[I] {
	return func(p *Params[I]) {
		config.Configure(&p.Input)
	}
}

func WithResourceOptions[I pulumi.Input](opts ...pulumi.ResourceOption) Option[I] {
	return func(p *Params[I]) {
		p.Options = append(p.Options, opts...)
	}
}
