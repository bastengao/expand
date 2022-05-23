package gormadapter

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"
	"gorm.io/gorm"
)

type Scope func(*gorm.DB) *gorm.DB

type Formatter func(expand string) string

type Config struct {
	Formatter Formatter
}

type Option func(config *Config)

var WithCamelCase Option = func(config *Config) {
	config.Formatter = camelCaseExpand
}

func ValidateExpand(expand []string, whitelist any, options ...Option) error {
	var config Config
	for _, opt := range options {
		opt(&config)
	}

	for _, e := range expand {
		if config.Formatter != nil {
			e = config.Formatter(e)
		}

		ok := findExpand(e, whitelist)
		if !ok {
			return fmt.Errorf("invalid expand: %s", e)
		}
	}

	return nil
}

// Inspired by https://stripe.com/docs/expand
func Expand(expand []string, whitelist any, options ...Option) (Scope, error) {
	expand = shrinkExpand(expand)
	err := ValidateExpand(expand, whitelist, options...)
	if err != nil {
		return nil, err
	}

	var config Config
	for _, opt := range options {
		opt(&config)
	}

	return func(db *gorm.DB) *gorm.DB {
		for _, e := range expand {
			if config.Formatter != nil {
				e = config.Formatter(e)
			}

			db = db.Preload(e)
		}
		return db
	}, nil
}

// shrink expand to most deep expand
func shrinkExpand(expand []string) []string {
	var shallowExpand []string

	for i, e := range expand {
		for _, e2 := range expand[i+1:] {
			if strings.HasPrefix(e2, e) {
				shallowExpand = append(shallowExpand, e)
			}
			if strings.HasPrefix(e, e2) {
				shallowExpand = append(shallowExpand, e2)
			}
		}
	}

	var shrunk []string
	for _, e := range expand {
		if !sliceStringContains(shallowExpand, e) {
			shrunk = append(shrunk, e)
		}
	}

	return shrunk
}

func findExpand(expand string, whitelist any) bool {
	if expand == "" {
		return false
	}

	parts := strings.Split(expand, ".")
	return findExpandPart(parts, whitelist)
}

func findExpandPart(parts []string, whitelist any) bool {
	switch v := whitelist.(type) {
	case map[string]any:
		part := parts[0]
		value, ok := v[part]
		if !ok {
			return false
		}
		if len(parts) == 1 { // last one
			return true
		}
		return findExpandPart(popSlice(parts), value)
	case []string:
		part := parts[0]
		ok := sliceStringContains(v, part)
		if !ok {
			return false
		}

		return len(parts) == 1 // last one
	}

	return false
}

func popSlice(s []string) []string {
	if len(s) <= 1 {
		return nil
	}

	return s[1:]
}

func sliceStringContains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}

	return false
}

func camelCaseExpand(expand string) string {
	parts := strings.Split(expand, ".")
	return strings.Join(camelCaseSlice(parts), ".")
}

func camelCaseSlice(s []string) []string {
	n := make([]string, len(s))
	for i := range s {
		n[i] = strcase.UpperCamelCase(s[i])
	}

	return n
}
