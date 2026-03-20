package metadata

import (
	"fmt"
	"slices"
	"strings"
)

type FunctionFields []FunctionField

type FunctionObjects map[string]FunctionObject

type FunctionObject struct {
	functionKey           string
	fields                map[string]FunctionField
	customTypeDescription string
}

func (r FunctionObject) CustomTypeDescription() string {
	return Sentencefy(r.customTypeDescription)
}

type FunctionField struct {
	parents  []string
	name     string
	dataType DataType

	allowNull    bool
	allowUnknown bool
	isVariadic   bool

	description           string
	customTypeDescription string

	validators []string
	isObject   bool
}

func (r FunctionField) Parents() []string {
	return r.parents
}

func (r FunctionField) Name() string {
	return r.name
}

func (r FunctionField) DataType() DataType {
	return r.dataType
}

func (r FunctionField) AllowNull() bool {
	return r.allowNull
}

func (r FunctionField) AllowUnknown() bool {
	return r.allowUnknown
}

func (r FunctionField) Description() string {
	return r.description
}

func (r FunctionField) CustomTypeDescription() string {
	return Sentencefy(r.customTypeDescription)
}

func (r FunctionField) Validators() []string {
	return MapSlice(r.validators, Sentencefy)
}

func (field FunctionField) Traits() string {
	var traits []string
	traits = append(traits, field.DataType().String())
	if field.AllowNull() {
		traits = append(traits, "Nullable")
	}
	if field.AllowUnknown() {
		traits = append(traits, "Unknownable")
	}
	return strings.Join(traits, ", ")
}

func (field FunctionField) nestedKey() string {
	return strings.Join(slices.Concat(field.Parents(), []string{field.Name()}), ".")
}

func (field FunctionField) NestedLink() string {
	if field.isObject {
		return fmt.Sprintf("See the nested schema [here](#nested--%s).", field.nestedKey())
	}
	return ""
}
