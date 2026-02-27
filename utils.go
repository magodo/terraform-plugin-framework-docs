package tfproviderdocs

import "context"

type DescriptionProvider interface {
	GetDescription() string
	GetMarkdownDescription() string
}

func DescriptionOf(d DescriptionProvider) string {
	if v := d.GetMarkdownDescription(); v != "" {
		return v
	}
	return d.GetDescription()
}

type DescriptionCtxProvider interface {
	Description(context.Context) string
	MarkdownDescription(context.Context) string
}

func DescriptionCtxOf(ctx context.Context, d DescriptionCtxProvider) string {
	if v := d.MarkdownDescription(ctx); v != "" {
		return v
	}
	return d.Description(ctx)
}

func MapSlice[T any, U any](input []T, f func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = f(v)
	}
	return result
}

func MapOrNil[T any, U any](input T, f func(T) U) *U {
	var anyInput any = input
	if anyInput == nil {
		return nil
	}
	output := f(input)
	return &output
}
