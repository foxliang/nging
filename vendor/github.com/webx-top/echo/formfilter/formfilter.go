package formfilter

import (
	"strings"

	"github.com/webx-top/echo"
)

type (
	Filter  func(*Data)
	Filters map[string][]Filter
	Options func() (key string, filter Filter)
	Data    struct {
		Key           string
		Value         []string
		normalizedKey string
	}
)

const All = `*`

func (d *Data) NormalizedKey() string {
	return d.normalizedKey
}

func Build(options ...Options) echo.FormDataFilter {
	filterMap := Filters{}
	for _, opt := range options {
		key, filter := opt()
		key = strings.Title(key)
		if _, ok := filterMap[key]; !ok {
			filterMap[key] = []Filter{}
		}
		filterMap[key] = append(filterMap[key], filter)
	}
	return echo.FormDataFilter(func(k string, v []string) (string, []string) {
		if k == All {
			return ``, nil
		}
		key := strings.Title(k)
		filters, ok := filterMap[key]
		filtersAll, okAll := filterMap[All]
		if !ok && !okAll {
			return k, v
		}
		data := &Data{Key: k, Value: v, normalizedKey: key}
		if ok {
			for _, filter := range filters {
				filter(data)
				if len(data.Key) == 0 {
					return data.Key, data.Value
				}
			}
		}
		if okAll {
			for _, filter := range filtersAll {
				filter(data)
				if len(data.Key) == 0 {
					return data.Key, data.Value
				}
			}
		}
		return data.Key, data.Value
	})
}