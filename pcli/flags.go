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

type Flag struct {
	Name     string
	Usage    string
	EnvVar   string
	Value    string
	FlagType FlagType
}

// CreateStringFlag creates a string flag with an optional default value.
func CreateStringFlag(name, usage string, value ...string) Flag {
	res := Flag{FlagType: StringFlag, Name: name, Usage: usage, EnvVar: makeEnvVarName(name)}
	if len(value) > 0 {
		res.Value = value[0]
	}
	return res
}

// CreateBoolFlag creates a bool flag that is false by default.
func CreateBoolFlag(name, usage string) Flag {
	return Flag{FlagType: BoolFlag, Name: name, Usage: usage, EnvVar: makeEnvVarName(name)}
}

// CreateBoolTFlag creates a bool flag that is true by default.
func CreateBoolTFlag(name, usage string) Flag {
	return Flag{FlagType: BoolTFlag, Name: name, Usage: usage, EnvVar: makeEnvVarName(name)}
}

// CreateIntFlag creates an int flag with an optional default value.
func CreateIntFlag(name, usage string, value ...string) Flag {
	res := Flag{FlagType: IntFlag, Name: name, Usage: usage, EnvVar: makeEnvVarName(name)}
	if len(value) > 0 {
		res.Value = value[0]
	}
	return res
}

// CreateStringSliceFlag creates a string slice flag with an optional default value.
func CreateStringSliceFlag(name, usage string, value ...string) Flag {
	res := Flag{FlagType: StringSliceFlag, Name: name, Usage: usage, EnvVar: makeEnvVarName(name), Value: strings.Join(value, ",")}
	return res
}

func makeEnvVarName(flagName string) string {
	return strings.Replace(strings.ToUpper(flagName), "-", "_", -1)
}

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
		Value:   s.Value == "true",
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
		Value:   s.Value != "false",
		Name:    s.Name,
		EnvVars: []string{s.EnvVar},
		Usage:   s.Usage,
	}
}
