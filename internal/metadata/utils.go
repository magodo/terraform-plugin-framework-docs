package metadata

import (
	"context"
	"strings"
	"unicode"
)

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

func MapOrZero[T any, U any](input T, f func(T) U) U {
	var zero U
	var anyInput any = input
	if anyInput == nil {
		return zero
	}
	output := f(input)
	return output
}

func Sentencefy(s string) string {
	if s == "" {
		return s
	}
	return capitalizeFirstLetter(ensureStringEndsWithDot(s))
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}

	// Convert string to a slice of runes to handle Unicode characters correctly
	runes := []rune(s)

	// Capitalize the first rune
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func ensureStringEndsWithDot(s string) string {
	if s == "" {
		return s
	}
	return strings.TrimRight(s, ".") + "."
}
