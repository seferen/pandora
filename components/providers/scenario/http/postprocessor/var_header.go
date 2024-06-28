package postprocessor

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/yandex/pandora/lib/str"
)

type VarHeaderPostprocessor struct {
	Mapping map[string]string
}

func (p *VarHeaderPostprocessor) ReturnedParams() []string {
	result := make([]string, len(p.Mapping))
	for k := range p.Mapping {
		result = append(result, k)
	}
	return result
}

func (p *VarHeaderPostprocessor) Process(resp *http.Response, _ io.Reader) (map[string]any, error) {
	if len(p.Mapping) == 0 {
		return nil, nil
	}
	result := make(map[string]any, len(p.Mapping))
	for k, v := range p.Mapping {
		headerVal, modifier, err := p.parseValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse value %s: %w", v, err)
		}
		val := resp.Header.Get(headerVal)
		if val == "" {
			continue
		}
		result[k] = modifier(val)
	}
	return result, nil
}

func (p *VarHeaderPostprocessor) parseValue(v string) (value string, modifier func(in string) string, err error) {
	vals := strings.Split(v, "|")
	if len(vals) == 1 {
		return vals[0], func(in string) string { return in }, nil
	}

	value = vals[0]
	modifier = func(in string) string { return in }

	for _, modStr := range vals[1:] {
		mod, err := p.parseModifier(modStr)
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse modifier %s: %w", modStr, err)
		}
		previousModifier := modifier
		modifier = func(in string) string {
			return mod(previousModifier(in))
		}
	}

	return value, modifier, nil
}

func (p *VarHeaderPostprocessor) parseModifier(s string) (func(in string) string, error) {
	name, args, err := str.ParseStringFunc(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function name %s: %w", s, err)
	}
	switch name {
	case "lower":
		return func(in string) string { return strings.ToLower(in) }, nil
	case "upper":
		return func(in string) string { return strings.ToUpper(in) }, nil
	case "substr":
		return p.substr(args)
	case "replace":
		if len(args) != 2 {
			return nil, fmt.Errorf("replace modifier requires 2 arguments")
		}
		return func(in string) string {
			return strings.ReplaceAll(in, args[0], args[1])
		}, nil
	}
	return nil, fmt.Errorf("unknown modifier %s", name)
}

func (p *VarHeaderPostprocessor) substr(args []string) (func(in string) string, error) {
	if len(args) == 0 || len(args) > 2 {
		return nil, fmt.Errorf("substr modifier requires one or two arguments")
	}
	start, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("substr modifier requires integer as first argument, got %s", args[0])
	}
	end := 0
	if len(args) == 2 {
		end, err = strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("substr modifier requires integer as second argument, got %s", args[1])
		}
	}
	return func(in string) string {
		l := len(in)
		if start < 0 {
			start = l + start
		}
		if end <= 0 {
			end = l + end
		}
		if end > l {
			end = l
		}
		if start > end {
			start, end = end, start
		}
		return in[start:end]
	}, nil
}
