package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

// MatchResult holds the scoring output for a name comparison.
type MatchResult struct {
	InputName      string  `json:"input_name"`
	CandidateName  string  `json:"candidate_name"`
	PhoneticScore  float64 `json:"phonetic_score"`
	DistanceScore  float64 `json:"distance_score"`
	TokenScore     float64 `json:"token_score"`
	CompositeScore float64 `json:"composite_score"`
	Status         string  `json:"status"`
}

// NameMatcher handles hybrid fuzzy name matching.
type NameMatcher struct {
	WPhonetic float64
	WDistance float64
	WToken    float64
}

// NewNameMatcher initializes a matcher with default or custom weights.
func NewNameMatcher(wPhonetic, wDistance, wToken float64) *NameMatcher {
	total := wPhonetic + wDistance + wToken
	return &NameMatcher{
		WPhonetic: wPhonetic / total,
		WDistance: wDistance / total,
		WToken:    wToken / total,
	}
}

// Normalize sanitizes string into uppercase alphanumeric tokens.
func Normalize(s string) string {
	s = strings.ToUpper(s)
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}
	fields := strings.Fields(builder.String())
	return strings.Join(fields, " ")
}

// Soundex implementation for zero-dependency Go phonetic matching.
func Soundex(s string) string {
	s = strings.ToUpper(s)
	var clean strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			clean.WriteRune(r)
		}
	}
	str := clean.String()
	if len(str) == 0 {
		return "Z000"
	}

	getDigit := func(r byte) byte {
		switch r {
		case 'B', 'F', 'P', 'V':
			return '1'
		case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
			return '2'
		case 'D', 'T':
			return '3'
		case 'L':
			return '4'
		case 'M', 'N':
			return '5'
		case 'R':
			return '6'
		default:
			return '0'
		}
	}

	res := []byte{str[0]}
	lastDigit := getDigit(str[0])

	for i := 1; i < len(str); i++ {
		digit := getDigit(str[i])
		if digit != '0' && digit != lastDigit {
			res = append(res, digit)
			lastDigit = digit
		}
		if len(res) == 4 {
			break
		}
	}

	for len(res) < 4 {
		res = append(res, '0')
	}

	return string(res)
}

// CalculatePhoneticScore calculates Jaccard index over token Soundex codes.
func (m *NameMatcher) CalculatePhoneticScore(s1, s2 string) float64 {
	tokens1 := strings.Fields(s1)
	tokens2 := strings.Fields(s2)

	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, t := range tokens1 {
		set1[Soundex(t)] = true
	}
	for _, t := range tokens2 {
		set2[Soundex(t)] = true
	}

	if len(set1) == 0 || len(set2) == 0 {
		return 0.0
	}

	intersection := 0
	unionMap := make(map[string]bool)

	for k := range set1 {
		unionMap[k] = true
		if set2[k] {
			intersection++
		}
	}
	for k := range set2 {
		unionMap[k] = true
	}

	return float64(intersection) / float64(len(unionMap))
}

// CalculateJaroWinkler calculates the Jaro-Winkler distance metric.
func (m *NameMatcher) CalculateJaroWinkler(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	len1 := len(s1)
	len2 := len(s2)
	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	matchDistance := int(math.Max(float64(len1), float64(len2))/2) - 1
	if matchDistance < 0 {
		matchDistance = 0
	}

	s1Matches := make([]bool, len1)
	s2Matches := make([]bool, len2)

	matches := 0
	transpositions := 0

	for i := 0; i < len1; i++ {
		start := int(math.Max(0, float64(i-matchDistance)))
		end := int(math.Min(float64(i+matchDistance+1), float64(len2)))

		for j := start; j < end; j++ {
			if s2Matches[j] || s1[i] != s2[j] {
				continue
			}
			s1Matches[i] = true
			s2Matches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	k := 0
	for i := 0; i < len1; i++ {
		if !s1Matches[i] {
			continue
		}
		for !s2Matches[k] {
			k++
		}
		if s1[i] != s2[k] {
			transpositions++
		}
		k++
	}

	jaro := (float64(matches)/float64(len1) +
		float64(matches)/float64(len2) +
		(float64(matches)-float64(transpositions)/2.0)/float64(matches)) / 3.0

	prefix := 0
	maxPrefix := int(math.Min(4, math.Min(float64(len1), float64(len2))))
	for i := 0; i < maxPrefix; i++ {
		if s1[i] == s2[i] {
			prefix++
		} else {
			break
		}
	}

	return jaro + float64(prefix)*0.1*(1.0-jaro)
}

// CalculateTokenSortScore computes score after sorting word tokens alphabetically.
func (m *NameMatcher) CalculateTokenSortScore(s1, s2 string) float64 {
	t1 := strings.Fields(s1)
	t2 := strings.Fields(s2)

	sort.Strings(t1)
	sort.Strings(t2)

	sorted1 := strings.Join(t1, " ")
	sorted2 := strings.Join(t2, " ")

	return m.CalculateJaroWinkler(sorted1, sorted2)
}

// Match executes the scoring logic.
func (m *NameMatcher) Match(inputName, candidateName string) MatchResult {
	norm1 := Normalize(inputName)
	norm2 := Normalize(candidateName)

	sPhonetic := m.CalculatePhoneticScore(norm1, norm2)
	sDistance := m.CalculateJaroWinkler(norm1, norm2)
	sToken := m.CalculateTokenSortScore(norm1, norm2)

	composite := (m.WPhonetic * sPhonetic) + (m.WDistance * sDistance) + (m.WToken * sToken)
	compositeScore := math.Round(composite*10000) / 100.0

	var status string
	if compositeScore >= 90.0 {
		status = "CRITICAL MATCH"
	} else if compositeScore >= 75.0 {
		status = "POTENTIAL MATCH"
	} else {
		status = "NO MATCH"
	}

	return MatchResult{
		InputName:      inputName,
		CandidateName:  candidateName,
		PhoneticScore:  math.Round(sPhonetic*10000) / 100.0,
		DistanceScore:  math.Round(sDistance*10000) / 100.0,
		TokenScore:     math.Round(sToken*10000) / 100.0,
		CompositeScore: compositeScore,
		Status:         status,
	}
}

func main() {
	matcher := NewNameMatcher(0.3, 0.4, 0.3)

	testCases := []struct {
		input     string
		candidate string
	}{
		{"Mohamad Ali Khan", "KHAN, MOHAMMED ALY"},
		{"Osama Bin Laden", "USAMA BIN LADIN"},
		{"Vladimir Putin", "VLADYMYR PUTYN"},
		{"John Smith", "JONATHAN SMITH"},
		{"Budi Santoso", "SANTOSO, BUDI"},
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-22s | %-22s | %-7s | %s\n", "INPUT NAME", "CANDIDATE NAME", "SCORE", "STATUS")
	fmt.Println(strings.Repeat("=", 80))

	for _, tc := range testCases {
		res := matcher.Match(tc.input, tc.candidate)
		fmt.Printf("%-22s | %-22s | %5.1f%% | %s\n", res.InputName, res.CandidateName, res.CompositeScore, res.Status)
	}

	fmt.Println(strings.Repeat("=", 80))
}
