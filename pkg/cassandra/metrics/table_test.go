package metrics

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger/pkg/testutils"
)

func TestTableEmit(t *testing.T) {
	testCases := []struct {
		err    error
		counts map[string]int64
		gauges map[string]int64
	}{
		{
			err: nil,
			counts: map[string]int64{
				"a_table.attempts": 1,
				"a_table.inserts":  1,
			},
			gauges: map[string]int64{
				"a_table.latency-ok.P999": 51,
				"a_table.latency-ok.P50":  51,
				"a_table.latency-ok.P75":  51,
				"a_table.latency-ok.P90":  51,
				"a_table.latency-ok.P95":  51,
				"a_table.latency-ok.P99":  51,
			},
		},
		{
			err: errors.New("some error"),
			counts: map[string]int64{
				"a_table.attempts": 1,
				"a_table.errors":   1,
			},
			gauges: map[string]int64{
				"a_table.latency-err.P999": 51,
				"a_table.latency-err.P50":  51,
				"a_table.latency-err.P75":  51,
				"a_table.latency-err.P90":  51,
				"a_table.latency-err.P95":  51,
				"a_table.latency-err.P99":  51,
			},
		},
	}
	for _, tc := range testCases {
		mf := metrics.NewLocalFactory(time.Second)
		tm := NewTable(mf, "a_table")
		tm.Emit(tc.err, 50*time.Millisecond)
		counts, gauges := mf.Snapshot()
		assert.Equal(t, tc.counts, counts)
		assert.Equal(t, tc.gauges, gauges)
		mf.Stop()
	}
}

func TestTableExec(t *testing.T) {
	testCases := []struct {
		q      insertQuery
		log    bool
		counts map[string]int64
	}{
		{
			q: insertQuery{},
			counts: map[string]int64{
				"a_table.attempts": 1,
				"a_table.inserts":  1,
			},
		},
		{
			q: insertQuery{
				str: "SELECT * FROM something",
				err: errors.New("failed"),
			},
			counts: map[string]int64{
				"a_table.attempts": 1,
				"a_table.errors":   1,
			},
		},
		{
			q: insertQuery{
				str: "SELECT * FROM something",
				err: errors.New("failed"),
			},
			log: true,
			counts: map[string]int64{
				"a_table.attempts": 1,
				"a_table.errors":   1,
			},
		},
	}

	for _, tc := range testCases {
		mf := metrics.NewLocalFactory(0)
		tm := NewTable(mf, "a_table")
		logger, logBuf := testutils.NewLogger(false)

		useLogger := logger
		if !tc.log {
			useLogger = nil
		}
		err := tm.Exec(tc.q, useLogger)
		if tc.q.err == nil {
			assert.NoError(t, err)
			assert.Len(t, logBuf.Bytes(), 0)
		} else {
			assert.Error(t, err, tc.q.err.Error())
			if tc.log {
				assert.Equal(t,
					"[E] Failed to exec query query=SELECT * FROM something error=failed\n",
					logBuf.String())
			} else {
				assert.Len(t, logBuf.Bytes(), 0)
			}
		}
		counts, _ := mf.Snapshot()
		assert.Equal(t, tc.counts, counts)
	}
}

type insertQuery struct {
	err error
	str string
}

func (q insertQuery) Exec() error {
	return q.err
}

func (q insertQuery) String() string {
	return q.str
}
