package test

import (
	"context"
	"flag"
	"fmt"
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
	"github.com/smart-libs/go-adapter/cli/lib/pkg/condition"
	"github.com/smart-libs/go-adapter/cli/lib/pkg/goflagset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	TagBasedWithArgsTest[Request any, Response any] struct {
		givenEnv         map[string]string
		givenArgs        []string
		givenFlagSet     func() *flag.FlagSet
		givenAssertions  func(t *testing.T, req Request) Response
		expectedExitCode int
	}
)

func (r TagBasedWithArgsTest[Request, Response]) Run(t *testing.T) {
	ctx := context.Background()
	handler := func(ctx context.Context, input Request) Response {
		return r.givenAssertions(t, input)
	}
	t.Run(fmt.Sprintf("Invoked with args=%v", r.givenArgs), func(t *testing.T) {
		// CLI adapter config
		givenConfig := cliadpt.Config{
			OsArgsUseDisabled: true,
			Bindings: []cliadpt.Binding{
				cliadpt.NewBindingBuilderWithCondition(condition.True()).
					InvokeHandler(handler).
					Build(),
			},
			FlagSet: goflagset.Wrap(*r.givenFlagSet()),
			EnvGetter: func(s string) (value string, found bool) {
				value, found = r.givenEnv[s]
				return
			},
		}
		// CLI adapter instance
		givenAdapter, err := cliadpt.NewSingleFlagSetAdapter(givenConfig)
		require.NoError(t, err)
		// Run CLI adapter
		exitCode := givenAdapter.Run(ctx, r.givenArgs...)
		assert.Equal(t, r.expectedExitCode, exitCode)
	})
}

func Test_InputStringFlag_WithTagAndDefaults(t *testing.T) {
	type request struct {
		StringParam string `flag:"s" default:"SELIC"`
	}

	const sFlagValue = "svalue"
	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{"-s", sFlagValue},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is provided, then handler must receive it", func(t *testing.T) {
				assert.Equal(t, sFlagValue, req.StringParam)
			})
			return nil
		},
	}.Run(t)

	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided and FlagSet has NO default, then handler must receive default value set with the flag tag", func(t *testing.T) {
				assert.Equal(t, "SELIC", req.StringParam)
			})
			return nil
		},
	}.Run(t)

	const sDefaultFlagValue = "default-svalue"
	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", sDefaultFlagValue, "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided, but FlagSet has default, then handler must receive default value", func(t *testing.T) {
				assert.Equal(t, sDefaultFlagValue, req.StringParam)
			})
			return nil
		},
	}.Run(t)

	type MyParam string
	type request2 struct {
		StringParam MyParam `flag:"s"` // the conversion from string to MyParam is performed by go-crosscutting/converter/lib
	}
	const stringRedefinedValue = "XPTO"
	TagBasedWithArgsTest[request2, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{"-s", stringRedefinedValue},
		givenAssertions: func(t *testing.T, req request2) error {
			t.Run("Even if MyParam is a string redefinition, it will be converted by go-crosscutting/converter/lib", func(t *testing.T) {
				assert.Equal(t, MyParam(stringRedefinedValue), req.StringParam)
			})
			return nil
		},
	}.Run(t)

	type Param struct {
		StringParam string `flag:"s"`
	}
	type request3 struct {
		Param // anonymous structure as a field with tagged field inside. This is resolved by go-adapter/cli
	}
	TagBasedWithArgsTest[request3, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{"-s", stringRedefinedValue},
		givenAssertions: func(t *testing.T, req request3) error {
			t.Run("Even if MyParam is a string redefinition, it will be converted by go-crosscutting/converter/lib", func(t *testing.T) {
				assert.Equal(t, stringRedefinedValue, req.StringParam)
			})
			return nil
		},
	}.Run(t)
}

func Test_InputBoolFlag_WithTagAndDefaults(t *testing.T) {
	type request struct {
		BoolParam bool `flag:"s" default:"true"`
	}

	const sFlagValue = true
	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.Bool("s", false, "")
			return result
		},
		givenArgs: []string{"-s", fmt.Sprintf("%v", sFlagValue)},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is provided, then handler must receive it", func(t *testing.T) {
				assert.Equal(t, sFlagValue, req.BoolParam)
			})
			return nil
		},
	}.Run(t)

	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.Bool("s", false, "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided and FlagSet has false as default, then handler must receive default value set by FlagSet", func(t *testing.T) {
				// ************************************
				// Because bool has only 2 values, the false value set as default in the FlagSet object will be
				// returned and the default set in the tag will not. If you want to use the default set by the tag
				// then you must declare it like result.String("s", "", "")
				assert.Equal(t, false, req.BoolParam)
			})
			return nil
		},
	}.Run(t)

	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided and FlagSet is string with empty string as default, then handler must receive default value set by the flag tag", func(t *testing.T) {
				// ************************************
				// Because the empty string is considered NO default, and the SDK converts it to nil, then the
				// flag parameter is nil and so the default value set by the flag tag will be used.
				assert.Equal(t, true, req.BoolParam)
			})
			return nil
		},
	}.Run(t)

	const sDefaultFlagValue = true
	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.Bool("s", sDefaultFlagValue, "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided, but FlagSet has default, then handler must receive default value", func(t *testing.T) {
				assert.Equal(t, sDefaultFlagValue, req.BoolParam)
			})
			return nil
		},
	}.Run(t)
}

func Test_InputStringFlag_WithTagAndNoDefaults(t *testing.T) {
	type request struct {
		OptionalStringParam *string `flag:"s"`
	}

	const sFlagValue = "svalue"
	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{"-s", sFlagValue},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is provided, then handler must receive it", func(t *testing.T) {
				assert.Equal(t, sFlagValue, *req.OptionalStringParam)
			})
			return nil
		},
	}.Run(t)

	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", "", "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided and FlagSet has NO default, then handler must receive nil", func(t *testing.T) {
				assert.Equal(t, (*string)(nil), req.OptionalStringParam)
			})
			return nil
		},
	}.Run(t)

	const sDefaultFlagValue = "default-svalue"
	TagBasedWithArgsTest[request, error]{
		givenFlagSet: func() *flag.FlagSet {
			result := flag.NewFlagSet("test", flag.PanicOnError)
			result.String("s", sDefaultFlagValue, "")
			return result
		},
		givenArgs: []string{},
		givenAssertions: func(t *testing.T, req request) error {
			t.Run("if -s parameter is NOT provided, but FlagSet has default, then handler must receive default value", func(t *testing.T) {
				assert.Equal(t, sDefaultFlagValue, *req.OptionalStringParam)
			})
			return nil
		},
	}.Run(t)
}
