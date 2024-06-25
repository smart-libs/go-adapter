package samples

import (
	"context"
	"flag"
	"fmt"
	cliadpt "github.com/smart-libs/go-adapter/cli/lib/pkg"
	"github.com/smart-libs/go-adapter/cli/lib/pkg/condition"
	"github.com/smart-libs/go-adapter/cli/lib/pkg/goflagset"
	"github.com/smart-libs/go-adapter/cli/lib/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Sample_01(t *testing.T) {
	//cliadpt.DebugEnabled = true
	//sdk.DebugEnabled = true
	// Handler Request
	type request struct {
		Indexes            []string   `flag:"i" default:"SELIC" mime-type:"text/csv"`
		StartDate          *time.Time `flag:"start"`
		Filename           string     `flag:"f" mandatory:"true"`
		OutputDir          *string    `flag:"o"`
		ZeroBalanceEnabled *bool      `flag:"z"`
		MerchantID         string     `env:"MERCHANT_ID" mandatory:"true"`
		AssetIDList        []string   `non-flags:"*"`
	}
	// Handler Return
	type response struct {
		FileName string `print:"stdout"`
		Error    error  `error:"panic"`
	}
	// Hook to do assertions
	var (
		requestReceived request
		givenResponse   response
	)
	assertions := func(t *testing.T, req request) response {
		requestReceived = req
		return givenResponse
	}
	// CLI Handler
	sample01Handler := func(ctx context.Context, input request) response {
		return assertions(t, input)
	}
	// GIVEN
	handlerFunc := sample01Handler
	const givenMerchantID = "merch_12345"
	const givenFileName = "file-name"
	const givenOutputDir = "output-dir"
	const givenResponseFileName = "output-file"
	// ENVIRONMENT
	givenEnv := map[string]string{"MERCHANT_ID": givenMerchantID}
	// CLI parameters
	givenFlagSet := flag.NewFlagSet("test", flag.PanicOnError)
	givenFlagSet.String("i", "", "")
	givenFlagSet.String("start", "", "since when")
	givenFlagSet.String("f", "", "file name")
	givenFlagSet.String("o", "", "output dir")
	givenFlagSet.Bool("z", false, "zero balance enabled")
	// Command line arguments
	args := []string{"-start", "20210101", "-f", givenFileName, "-z", "-o", givenOutputDir, "-i", "IPCA,SELIC"}
	givenResponse = response{FileName: givenResponseFileName}

	ctx := context.Background()
	t.Run(fmt.Sprintf("Invoked with args %v", args), func(t *testing.T) {
		// WHEN
		stdout := test.CaptureStdout(func() {
			// CLI adapter config
			givenConfig := cliadpt.Config{
				Bindings: []cliadpt.Binding{
					cliadpt.NewBindingBuilderWithCondition(condition.HasFlag("f")).
						InvokeHandler(handlerFunc).
						Build(),
				},
				FlagSet: goflagset.Wrap(*givenFlagSet),
				EnvGetter: func(s string) (value string, found bool) {
					value, found = givenEnv[s]
					return
				},
			}
			// CLI adapter instance
			givenAdapter, err := cliadpt.NewSingleFlagSetAdapter(givenConfig)
			require.NoError(t, err)
			assert.Equal(t, 0, givenAdapter.Run(ctx, args...))
		})
		// THEN
		if assert.NotNil(t, requestReceived.ZeroBalanceEnabled) {
			assert.Equal(t, true, *requestReceived.ZeroBalanceEnabled)
		}
		assert.Equal(t, givenMerchantID, requestReceived.MerchantID)
		assert.Equal(t, []string{"IPCA", "SELIC"}, requestReceived.Indexes)
		if assert.NotNil(t, requestReceived.StartDate) {
			expectedDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
			assert.Equal(t, expectedDate, *requestReceived.StartDate)
		}
		assert.Equal(t, givenFileName, requestReceived.Filename)
		if assert.NotNil(t, requestReceived.OutputDir) {
			assert.Equal(t, givenOutputDir, *requestReceived.OutputDir)
		}
		assert.Equal(t, givenResponseFileName, stdout)
	})
}
