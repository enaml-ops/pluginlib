package pcli

import (
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

type FlagType int

const (
	StringFlag FlagType = iota
	StringSliceFlag
	BoolFlag
	IntFlag
	BoolTFlag
)

type (
	Flag struct {
		Name     string
		Usage    string
		EnvVar   string
		Value    string
		FlagType FlagType
	}
)

func (s Flag) ToCli() interface{} {
	var ret interface{}
	switch s.FlagType {
	case StringFlag:
		ret = StringFlagToCli(s)

	case StringSliceFlag:
		ret = StringSliceFlagToCli(s)

	case BoolFlag:
		ret = BoolFlagToCli(s)

	case BoolTFlag:
		ret = BoolTFlagToCli(s)

	case IntFlag:
		ret = IntFlagToCli(s)

	default:
		lo.G.Error("not sure how to handle this type")
	}
	return ret
}

func StringFlagToCli(s Flag) interface{} {
	return cli.StringFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Value:  s.Value,
		Usage:  s.Usage,
	}
}

func StringSliceFlagToCli(s Flag) interface{} {
	var stringSlice *cli.StringSlice = &cli.StringSlice{}

	if s.Value != "" {
		*stringSlice = strings.Split(s.Value, ",")
	}

	res := cli.StringSliceFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Value:  stringSlice,
		Usage:  s.Usage,
	}
	return res
}

func BoolFlagToCli(s Flag) interface{} {

	return cli.BoolFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Usage:  s.Usage,
	}
}

func IntFlagToCli(s Flag) interface{} {
	intVal, _ := strconv.Atoi(s.Value)
	return cli.IntFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Value:  intVal,
		Usage:  s.Usage,
	}
}

func BoolTFlagToCli(s Flag) interface{} {

	return cli.BoolTFlag{
		Name:   s.Name,
		EnvVar: s.EnvVar,
		Usage:  s.Usage,
	}
}
