package eqi

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// EquivalentPair is a pair of equiavalent unicode characters, specifying
// a variant (eg. CHK Radical, CJK Stroke, or Kangxi Radical) and its
// Unified CJK equivalent.
type EquivalentPair struct {
	VariantName string
	Variant     string
	Unified     string
}

// Replace iterates through the runes ofa string and replaces them with
// their CJK Unified Equivalents if applicable.
func Replace(s string) []byte {
	n := []byte{}
	for _, r := range s {
		if v, ok := MappedPairs[string(r)]; ok {
			n = append(n, []byte(v)...)
			continue
		}
		n = append(n, []byte(string(r))...)
	}

	return n
}

// BufferedReplace iterates through the runes of an io.Reader and replaces
// them with their CJK Unified Equivalents if applicable.
func BufferedReplace(r io.Reader) (out *bytes.Buffer, err error) {
	out = bytes.NewBuffer([]byte{})

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanRunes)
	for s.Scan() {
		rn := s.Text()
		if v, ok := MappedPairs[string(rn)]; ok {
			out.WriteString(v)
			continue
		}
		out.WriteString(rn)
	}

	return out, nil
}

// MapPairs returns the pairs in a map, keyed on variant with the
// unified character as the value.
func MapPairs(pairs []EquivalentPair) map[string]string {
	m := map[string]string{}
	for _, v := range pairs {
		m[v.Variant] = v.Unified
	}

	return m
}

// ExtractPairs extracts the unicode pairs from an io.Reader to a properly
// formaetted EquivalentUnifiedIdeograph.txt file (or equivalent).
func ExtractPairs(r io.Reader) (pairs []EquivalentPair, err error) {
	pairs = []EquivalentPair{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if line == "" { // Skip blank lines.
			continue
		}

		if strings.HasPrefix(line, "#") { // Skip Commented lines.
			continue
		}

		variant, unified, variantName := targetsFromLine(line)
		extracted, err := pairsFromTargets(variant, unified, variantName)
		if err != nil {
			return pairs, err
		}

		pairs = append(pairs, extracted...)
	}

	return
}

// targetsFromLine extracts a pair of unicode targets from an index line.
func targetsFromLine(line string) (variant, unified, variantName string) {
	codePoints := strings.Split(line, "#")
	variantName = strings.TrimSpace(codePoints[1])
	targets := strings.Split(codePoints[0], ";")
	variant = strings.TrimSpace(targets[0])
	unified = strings.TrimSpace(targets[1])

	return
}

// hexToChar converts a hex value to a unicode characer. If the
// hex value is more than 4 characters, the value will be extended to the
// correct representation for decoding.
func hexToChar(hex string) (char string, err error) {
	k := fmt.Sprintf("\\u%s", hex)
	if len(hex) > 4 {
		k = fmt.Sprintf("\\U%s%s", strings.Repeat("0", 8-len(hex)), hex)
	}

	char, err = strconv.Unquote(`"` + k + `"`)
	if err != nil {
		return
	}

	return
}

// pairsFromTargets returns a slice of equivalent pairs from a target string.
// If the target specifies a range (eg. xxxx..xxxx), multiple pairs covering
// this range will be returned.
func pairsFromTargets(variant, unified, variantName string) (pairs []EquivalentPair, err error) {
	pairs = []EquivalentPair{}
	variants := []string{}
	if strings.Contains(variant, "..") {
		vr := strings.Split(variant, "..")
		v := vr[0]
		var i = 0
		for {
			variants = append(variants, v)

			if i == 10 || v == vr[1] {
				break
			}

			i++
			v, err = incrementHex(v)
			if err != nil {
				return pairs, err
			}
		}
	} else {
		variants = append(variants, variant)
	}

	unifiedChar, err := hexToChar(unified)
	if err != nil {
		return
	}

	var variantChar string
	for _, variant := range variants {
		variantChar, err = hexToChar(variant)
		if err != nil {
			return pairs, err
		}

		pairs = append(pairs, EquivalentPair{
			VariantName: variantName,
			Variant:     variantChar,
			Unified:     unifiedChar,
		})
	}

	return pairs, nil
}

// incrementHex increments a hex value.
func incrementHex(hex string) (incr string, err error) {
	dec, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return
	}

	return strings.ToUpper(strconv.FormatInt(dec+1, 16)), nil
}
