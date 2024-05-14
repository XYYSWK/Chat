package retry

import (
	"fmt"
	"github.com/XYYSWK/Rutils/pkg/utils"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	t.Parallel()
	end := int(utils.RandomInt(1, 5))
	start := 0
	testFunc := func() error {
		start++
		if start == end {
			return nil
		}
		return fmt.Errorf("%d", start)
	}
	log.Println("start")
	report := <-NewTry("test", testFunc, 100*time.Millisecond, end).Run()
	require.True(t, report.Result)
	require.Equal(t, report.Times, end)
	require.Len(t, report.Errs, end-1)
	for i, err := range report.Errs {
		require.Equal(t, fmt.Sprintf("%d", i+1), err.Error())
	}
}
