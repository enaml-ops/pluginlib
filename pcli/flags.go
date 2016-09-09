package pcli

import (
	"strconv"
	"strings"

	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
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

	if s.EnvVar == "" {
		s.EnvVar = createEnvVar(s.Name)
	}

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

func createEnvVar(flagname string) string {
	return "OMG_" + strings.Replace(strings.ToUpper(flagname), "-", "_", -1)
}

func StringFlagToCli(s Flag) interface{} {
	return &cli.StringFlag{
		Name:    s.Name,
		EnvVars: []string{s.EnvVar},
		Value:   s.Value,
		Usage:   s.Usage,
	}
}

func StringSliceFlagToCli(s Flag) interface{} {
	var stringSlice *cli.StringSlice

	if s.Value != "" {
		stringSlice = cli.NewStringSlice(strings.Split(s.Value, ",")...)
	} else {
		stringSlice = cli.NewStringSlice()
	}

	res := &cli.StringSliceFlag{
		Name:    s.Name,
		EnvVars: []string{s.EnvVar},
		Value:   stringSlice,
		Usage:   s.Usage,
	}
	return res
}

func BoolFlagToCli(s Flag) interface{} {
	return &cli.BoolFlag{
		Name:    s.Name,
		EnvVars: []string{s.EnvVar},
		Usage:   s.Usage,
	}
}

func IntFlagToCli(s Flag) interface{} {
	intVal, _ := strconv.Atoi(s.Value)
	return &cli.IntFlag{
		Name:    s.Name,
		EnvVars: []string{s.EnvVar},
		Value:   intVal,
		Usage:   s.Usage,
	}
}

func BoolTFlagToCli(s Flag) interface{} {
	return &cli.BoolFlag{
		Value:   true,
		Name:    s.Name,
		EnvVars: []string{s.EnvVar},
		Usage:   s.Usage,
	}
}
