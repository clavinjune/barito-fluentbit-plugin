package barito

import (
	"structs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tt := []struct {
		_           structs.HostLayout
		Name        string
		Input       string
		Expected    string
		ExpectedErr error
	}{
		{
			Name:        "empty string",
			Input:       "",
			Expected:    "",
			ExpectedErr: errInvalidURL,
		},
		{
			Name:        "no scheme",
			Input:       "github.com/baritolog/",
			Expected:    "",
			ExpectedErr: errInvalidURL,
		},
		{
			Name:        "with port",
			Input:       "https://github.com:443/baritolog/",
			Expected:    "https://github.com:443/baritolog/",
			ExpectedErr: nil,
		},
		{
			Name:        "IP and Port",
			Input:       "https://12.34.56.78:443/baritolog/",
			Expected:    "https://12.34.56.78:443/baritolog/",
			ExpectedErr: nil,
		},
		{
			Name:        "https",
			Input:       "https://github.com",
			Expected:    "https://github.com",
			ExpectedErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			actual, err := parseURL(tc.Input)
			if tc.ExpectedErr == nil {
				r.NoError(err)
			} else {
				r.ErrorIs(err, tc.ExpectedErr)
			}

			r.Equal(tc.Expected, actual)
		})
	}
}

func TestParseExtraLabels(t *testing.T) {
	tt := []struct {
		_        structs.HostLayout
		Name     string
		Input    string
		Expected map[string]string
	}{
		{
			Name:     "emptystring",
			Input:    "",
			Expected: map[string]string{},
		},
		{
			Name:  "whitespace",
			Input: "  a=2 ",
			Expected: map[string]string{
				"a": "2",
			},
		},
		{
			Name:  "trailing comma",
			Input: "  a=2 ,",
			Expected: map[string]string{
				"a": "2",
			},
		},
		{
			Name:  "empty value",
			Input: "  a=2 ,b=",
			Expected: map[string]string{
				"a": "2",
				"b": "",
			},
		},
		{
			Name:  "empty value and trailing comma",
			Input: "  a=2 ,b=,",
			Expected: map[string]string{
				"a": "2",
				"b": "",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			r.Equal(tc.Expected, parseExtraLabels(tc.Input))
		})
	}
}
