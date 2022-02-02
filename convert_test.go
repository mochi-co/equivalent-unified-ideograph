package eqi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractPairs(t *testing.T) {
	b := strings.NewReader(`
	# invalid line followed by blank line

2E81       ; 5382  #     CJK RADICAL CLIFF
2E82       ; 4E5B  #     CJK RADICAL SECOND ONE
2ECC..2ECE ; 8FB6  # [3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO

2E83       ; 4E5A  #     CJK RADICAL SECOND TWO
2EBE..2EC0 ; 8279  # [3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE
31E1       ; 2010E #     CJK STROKE HZZZG
	`)
	pairs, err := ExtractPairs(b)
	require.NoError(t, err)
	require.Equal(t, []EquivalentPair{
		{Variant: "⺁", Unified: "厂"},
		{Variant: "⺂", Unified: "乛"},
		{Variant: "⻌", Unified: "辶"},
		{Variant: "⻍", Unified: "辶"},
		{Variant: "⻎", Unified: "辶"},
		{Variant: "⺃", Unified: "乚"},
		{Variant: "⺾", Unified: "艹"},
		{Variant: "⺿", Unified: "艹"},
		{Variant: "⻀", Unified: "艹"},
		{Variant: "㇡", Unified: "𠄎"},
	}, pairs)
}

func TestTargetsFromLine(t *testing.T) {
	tt := []struct {
		have string
		t0   string
		t1   string
	}{
		{
			have: `2E81       ; 5382  #     CJK RADICAL CLIFF`, t0: `2E81`, t1: `5382`,
		}, {
			have: `2E87       ; 20628 #     CJK RADICAL TABLE`, t0: `2E87`, t1: `20628`,
		}, {
			have: `31E1       ; 2010E #     CJK STROKE HZZZG`, t0: `31E1`, t1: `2010E`,
		}, {
			have: `31D2..31D3 ; 4E3F  # [2] CJK STROKE P..CJK STROKE SP`, t0: `31D2..31D3`, t1: `4E3F`,
		}, {
			have: `2ECC..2ECE ; 8FB6  # [3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO`, t0: `2ECC..2ECE`, t1: `8FB6`,
		},
	}

	for _, c := range tt {
		t0, t1 := targetsFromLine(c.have)
		require.Equal(t, c.t0, t0)
		require.Equal(t, c.t1, t1)
	}
}

func TestPairsFromTargets(t *testing.T) {
	tt := []struct {
		variant string
		unified string
		pairs   []EquivalentPair
	}{
		{
			variant: `2E81`,
			unified: `5382`,
			pairs:   []EquivalentPair{{"⺁", "厂"}},
		},
		{
			variant: `2FD3`,
			unified: `9F8D`,
			pairs:   []EquivalentPair{{"⿓", "龍"}},
		},
		{
			variant: `2EBE..2EC0`,
			unified: `8279`,
			pairs:   []EquivalentPair{{"⺾", "艹"}, {"⺿", "艹"}, {"⻀", "艹"}},
		},
		{
			variant: `31E1`,
			unified: `2010E`,
			pairs:   []EquivalentPair{{"㇡", "𠄎"}},
		},
	}

	for _, c := range tt {
		pairs, err := pairsFromTargets(c.variant, c.unified)
		require.Equal(t, c.pairs, pairs)
		require.NoError(t, err)
	}
}

func TestIncrementHex(t *testing.T) {
	tt := []struct {
		have string
		want string
	}{
		{have: `2EBE`, want: `2EBF`},
		{have: `2EBF`, want: `2EC0`},
	}

	for _, c := range tt {
		incr, err := incrementHex(c.have)
		require.NoError(t, err)
		require.Equal(t, c.want, incr)
	}
}

func TestIncrementHexError(t *testing.T) {
	_, err := incrementHex("qwerty")
	require.Error(t, err)
}

func TestHexToCode(t *testing.T) {
	tt := []struct {
		have string
		want string
	}{
		{have: `2EBE`, want: `⺾`},
		{have: `2010E`, want: `𠄎`},
	}

	for _, c := range tt {
		h, err := hexToChar(c.have)
		require.NoError(t, err)
		require.Equal(t, c.want, h)
	}
}

func TestHexToCodeError(t *testing.T) {
	_, err := hexToChar("abcdefg")
	require.Error(t, err)
}
