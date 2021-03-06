package http

import (
	"testing"

	"github.com/zoncoen/scenarigo/context"
	"github.com/zoncoen/yaml"
)

func TestExpect_Build(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		tests := map[string]struct {
			vars   interface{}
			expect *Expect
			result *result
		}{
			"default": {
				expect: &Expect{},
				result: &result{
					status: "200 OK",
				},
			},
			"status code": {
				expect: &Expect{
					Code: "404",
				},
				result: &result{
					status: "404 Not Found",
				},
			},
			"status code string": {
				expect: &Expect{
					Code: "Not Found",
				},
				result: &result{
					status: "404 Not Found",
				},
			},
			"assert body": {
				expect: &Expect{
					Body: yaml.MapSlice{
						yaml.MapItem{
							Key:   "foo",
							Value: "bar",
						},
					},
				},
				result: &result{
					status: "200 OK",
					body:   map[string]string{"foo": "bar"},
				},
			},
			"with vars": {
				vars: map[string]string{"foo": "bar"},
				expect: &Expect{
					Body: yaml.MapSlice{
						yaml.MapItem{
							Key:   "foo",
							Value: "{{vars.foo}}",
						},
					},
				},
				result: &result{
					status: "200 OK",
					body:   map[string]string{"foo": "bar"},
				},
			},
		}
		for name, test := range tests {
			test := test
			t.Run(name, func(t *testing.T) {
				ctx := context.FromT(t)
				if test.vars != nil {
					ctx = ctx.WithVars(test.vars)
				}
				assertion, err := test.expect.Build(ctx)
				if err != nil {
					t.Fatalf("failed to build assertion: %s", err)
				}
				if err := assertion.Assert(test.result); err != nil {
					t.Errorf("got assertion error: %s", err)
				}
			})
		}
	})
	t.Run("ng", func(t *testing.T) {
		tests := map[string]struct {
			expect            *Expect
			result            *result
			expectBuildError  bool
			expectAssertError bool
		}{
			"wrong status code": {
				expect: &Expect{},
				result: &result{
					status: "404 Not Found",
				},
				expectAssertError: true,
			},
			"wrong body": {
				expect: &Expect{
					Body: yaml.MapSlice{
						yaml.MapItem{
							Key:   "foo",
							Value: "bar",
						},
					},
				},
				result: &result{
					status: "200 OK",
				},
				expectAssertError: true,
			},
			"failed to execute template": {
				expect: &Expect{
					Body: yaml.MapSlice{
						yaml.MapItem{
							Key:   "foo",
							Value: "{{vars.foo}}",
						},
					},
				},
				expectBuildError: true,
			},
		}
		for name, test := range tests {
			test := test
			t.Run(name, func(t *testing.T) {
				ctx := context.FromT(t)
				assertion, err := test.expect.Build(ctx)
				if test.expectBuildError && err == nil {
					t.Fatal("succeeded building assertion")
				}
				if !test.expectBuildError && err != nil {
					t.Fatalf("failed to build assertion: %s", err)
				}
				if err != nil {
					return
				}

				err = assertion.Assert(test.result)
				if test.expectAssertError && err == nil {
					t.Errorf("no assertion error")
				}
				if !test.expectAssertError && err != nil {
					t.Errorf("got assertion error: %s", err)
				}
			})
		}
	})
}
