package api

import (
	"gotest.tools/assert"
	"testing"
)

func TestWithCache(t *testing.T) {
	var tests = []struct {
		name string
		option   RequestOption
		expected RequestOptions
	}{
		{
			name: "caching is enabled by default",
			option: nil,
			expected: RequestOptions{
				cache: true,
			},
		},
		{
			name: "enable cache",
			option: WithCache(true),
			expected: RequestOptions{
				cache: true,
			},
		},
		{
			name: "disable cache",
			option: WithCache(false),
			expected: RequestOptions{
				cache: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				opts = make([]RequestOption, 0)
				value *RequestOptions
			)
			if tt.option != nil {
				opts = append(opts, tt.option)
			}
			value = NewRequestOptions(opts)
			assert.DeepEqual(t, value.cache, tt.expected.cache)
		})
	}
}

func TestWithStatus(t *testing.T) {
	var tests = []struct {
		name string
		option   RequestOption
		expected RequestOptions
	}{
		{
			name: "status is empty by default",
			option: nil,
			expected: RequestOptions{
				status: "",
			},
		},
		{
			name: "active members",
			option: WithStatus(memberStatusActive),
			expected: RequestOptions{
				status: memberStatusActive,
			},
		},
		{
			name: "archived members",
			option: WithStatus(memberStatusArchived),
			expected: RequestOptions{
				status: memberStatusArchived,
			},
		},
		{
			name: "all members",
			option: WithStatus(memberStatusAll),
			expected: RequestOptions{
				status: memberStatusAll,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				opts = make([]RequestOption, 0)
				value *RequestOptions
			)
			if tt.option != nil {
				opts = append(opts, tt.option)
			}
			value = NewRequestOptions(opts)
			assert.DeepEqual(t, value.status, tt.expected.status)
		})
	}
}

func TestWithTerm(t *testing.T) {
	var tests = []struct {
		name string
		option   RequestOption
		expected RequestOptions
	}{
		{
			name: "term is empty by default",
			option: nil,
			expected: RequestOptions{
				term: "",
			},
		},
		{
			name: "search by term",
			option: WithTerm("foobar"),
			expected: RequestOptions{
				term: "foobar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				opts = make([]RequestOption, 0)
				value *RequestOptions
			)
			if tt.option != nil {
				opts = append(opts, tt.option)
			}
			value = NewRequestOptions(opts)
			assert.DeepEqual(t, value.term, tt.expected.term)
		})
	}
}
