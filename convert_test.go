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
		{Variant: "⺁", Unified: "厂", VariantName: "CJK RADICAL CLIFF"},
		{Variant: "⺂", Unified: "乛", VariantName: "CJK RADICAL SECOND ONE"},
		{Variant: "⻌", Unified: "辶", VariantName: "[3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO"},
		{Variant: "⻍", Unified: "辶", VariantName: "[3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO"},
		{Variant: "⻎", Unified: "辶", VariantName: "[3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO"},
		{Variant: "⺃", Unified: "乚", VariantName: "CJK RADICAL SECOND TWO"},
		{Variant: "⺾", Unified: "艹", VariantName: "[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE"},
		{Variant: "⺿", Unified: "艹", VariantName: "[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE"},
		{Variant: "⻀", Unified: "艹", VariantName: "[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE"},
		{Variant: "㇡", Unified: "𠄎", VariantName: "CJK STROKE HZZZG"},
	}, pairs)
}

func TestTargetsFromLine(t *testing.T) {
	tt := []struct {
		have        string
		variantName string
		variant     string
		unified     string
	}{
		{
			have:        `2E81       ; 5382  #     CJK RADICAL CLIFF`,
			variantName: "CJK RADICAL CLIFF",
			variant:     `2E81`,
			unified:     `5382`,
		}, {
			have:        `2E87       ; 20628 #     CJK RADICAL TABLE`,
			variantName: "CJK RADICAL TABLE",
			variant:     `2E87`,
			unified:     `20628`,
		}, {
			have:        `31E1       ; 2010E #     CJK STROKE HZZZG`,
			variantName: "CJK STROKE HZZZG",
			variant:     `31E1`,
			unified:     `2010E`,
		}, {
			have:        `31D2..31D3 ; 4E3F  # [2] CJK STROKE P..CJK STROKE SP`,
			variantName: "[2] CJK STROKE P..CJK STROKE SP",
			variant:     `31D2..31D3`,
			unified:     `4E3F`,
		}, {
			have:        `2ECC..2ECE ; 8FB6  # [3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO`,
			variantName: "[3] CJK RADICAL SIMPLIFIED WALK..CJK RADICAL WALK TWO",
			variant:     `2ECC..2ECE`,
			unified:     `8FB6`,
		},
	}

	for _, c := range tt {
		variant, unified, name := targetsFromLine(c.have)
		require.Equal(t, c.variant, variant)
		require.Equal(t, c.unified, unified)
		require.Equal(t, c.variantName, name)
	}
}

func TestPairsFromTargets(t *testing.T) {
	tt := []struct {
		variantName string
		variant     string
		unified     string
		pairs       []EquivalentPair
	}{
		{
			variantName: "CJK RADICAL CLIFF",
			variant:     `2E81`,
			unified:     `5382`,
			pairs:       []EquivalentPair{{"CJK RADICAL CLIFF", "⺁", "厂"}},
		},
		{
			variantName: "KANGXI RADICAL DRAGON",
			variant:     `2FD3`,
			unified:     `9F8D`,
			pairs:       []EquivalentPair{{"KANGXI RADICAL DRAGON", "⿓", "龍"}},
		},
		{
			variantName: "[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE",
			variant:     `2EBE..2EC0`,
			unified:     `8279`,
			pairs: []EquivalentPair{
				{"[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE", "⺾", "艹"},
				{"[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE", "⺿", "艹"},
				{"[3] CJK RADICAL GRASS ONE..CJK RADICAL GRASS THREE", "⻀", "艹"},
			},
		},
		{
			variantName: "CJK STROKE HZZZG",
			variant:     `31E1`,
			unified:     `2010E`,
			pairs:       []EquivalentPair{{"CJK STROKE HZZZG", "㇡", "𠄎"}},
		},
	}

	for _, c := range tt {
		pairs, err := pairsFromTargets(c.variant, c.unified, c.variantName)
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
