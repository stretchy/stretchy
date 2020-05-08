package configuration

import (
	"fmt"
	"strings"

	"github.com/r3labs/diff"
	"github.com/stretchy/stretchy/pkg/utils"
)

type ChangeType int

const (
	ChangeTypeCreate ChangeType = iota
	ChangeTypeUpdate
	ChangeTypeDelete
)

func (d ChangeType) String() string {
	return [...]string{"CREATE", "UPDATE", "DELETE"}[d]
}

func NewChangeTypeFromDiffType(diffType string) ChangeType {
	switch diffType {
	case diff.CREATE:
		return ChangeTypeCreate
	case diff.UPDATE:
		return ChangeTypeUpdate
	case diff.DELETE:
		return ChangeTypeDelete
	}

	panic("")
}

type Change struct {
	Type ChangeType
	Path []string
	From interface{}
	To   interface{}
}

func (c Change) FullPath() string {
	return strings.Join(c.Path, ".")
}

func (c Change) IsIndexMetadata() (bool, error) {
	return utils.InSlice(
		c.FullPath(),
		[]string{
			"settings.index.creation_date",
			"settings.index.provided_name",
			"settings.index.uuid",
			"settings.index.version",
		},
	)
}

func (c Change) String() string {
	changeFromTo := ""

	fromChange := utils.PrintScalar(c.From)
	toChange := utils.PrintScalar(c.To)

	if fromChange != "" && toChange != "" {
		changeFromTo = fmt.Sprintf(" [From: %s - To: %s]", fromChange, toChange)
	}

	return fmt.Sprintf("%s => %s%s", c.Type.String(), c.FullPath(), changeFromTo)
}

const defaultIgnoreAboveValue float64 = 256

func (c Change) ShouldBeReported() bool {
	pathStart := c.Path[0]
	pathEnd := c.Path[len(c.Path)-1]

	if c.Type == ChangeTypeDelete &&
		pathStart == "mappings" &&
		pathEnd == "ignore_above" &&
		c.From == defaultIgnoreAboveValue &&
		c.To == nil {
		return false
	}

	return true
}

type ChangeCollection []Change
