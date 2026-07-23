package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// IdentityRecord represents an entity in the Watchlist/Blacklist.
type IdentityRecord struct {
	ID          string   `json:"id"`
	FullName    string   `json:"full_name"`
	DOB         string   `json:"dob"` // YYYY-MM-DD or YYYY-00-00
	Nationality string   `json:"nationality"`
	PassportNo  string   `json:"passport_no"`
	Aliases     []string `json:"aliases"`
}

// MatchBreakdown contains attribute-level scores.
type MatchBreakdown struct {
	NameScore      float64 `json:"name_score"`
	DOBScore       float64 `json:"dob_score"`
	PassportScore  float64 `json:"passport_score"`
	AliasScore     float64 `json:"alias_score"`
	FellegiLogOdds float64 `json:"fellegi_log_odds"`
	TotalScore     float64 `json:"total_score"`
	Status         string  `json:"status"`
}

// SearchResult represents a candidate match output.
type SearchResult struct {
	Record    IdentityRecord `json:"record"`
	Breakdown MatchBreakdown `json:"breakdown"`
}

// ExpertEngine is an industrial-grade identity resolution and screening system.
type ExpertEngine struct {
	mu          sync.RWMutex
	records     map[string]IdentityRecord
	tokenIDF    map[string]float64
	invertedIdx map[string]set
	aliasGraph  map[string]set // Graph representation of connected alias IDs
	docCount    int
	prefixes    map[string]bool
}

type set map[string]bool

// NewExpertEngine initializes the high-performance screening engine.
func NewExpertEngine() *ExpertEngine {
	engine := &ExpertEngine{
		records:     make(map[string]IdentityRecord),
		tokenIDF:    make(map[string]float64),
		invertedIdx: make(map[string]set),
		aliasGraph:  make(map[string]set),
		prefixes: map[string]bool{
			"AL": true, "EL": true, "ABDUL": true, "BIN": true, "BINTI": true,
			"VON": true, "VAN": true, "DE": true, "DER": true, "SAN": true,
		},
	}
	return engine
}

// StemToken removes common cultural prefixes (Onomastic Normalization).
func (e *ExpertEngine) StemToken(token string) string {
	token = strings.ToUpper(token)
	for prefix := range e.prefixes {
		if strings.HasPrefix(token, prefix+"-") {
			return strings.TrimPrefix(token, prefix+"-")
		}
	}
	return token
}

// IndexRecord adds an entity to the Inverted Index, IDF Dictionary, and Alias Graph.
func (e *ExpertEngine) IndexRecord(rec IdentityRecord) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.records[rec.ID] = rec
	e.docCount++

	allNames := append([]string{rec.FullName}, rec.Aliases...)
	tokensInDoc := make(map[string]bool)

	for _, name := range allNames {
		clean := Normalize(name)
		for _, token := range strings.Fields(clean) {
			stemmed := e.StemToken(token)
			tokensInDoc[stemmed] = true

			soundexKey := Soundex(stemmed)
			if e.invertedIdx[soundexKey] == nil {
				e.invertedIdx[soundexKey] = make(set)
			}
			e.invertedIdx[soundexKey][rec.ID] = true
		}
	}

	for token := range tokensInDoc {
		e.tokenIDF[token] += 1.0
	}

	if e.aliasGraph[rec.ID] == nil {
		e.aliasGraph[rec.ID] = make(set)
	}
	for _, alias := range rec.Aliases {
		aliasID := "ALIAS_" + Normalize(alias)
		e.aliasGraph[rec.ID][aliasID] = true
		if e.aliasGraph[aliasID] == nil {
			e.aliasGraph[aliasID] = make(set)
		}
		e.aliasGraph[aliasID][rec.ID] = true
	}
}

// GetIDF returns the Inverse Document Frequency weight for a token.
func (e *ExpertEngine) GetIDF(token string) float64 {
	e.mu.RLock()
	defer e.mu.RUnlock()

	count, exists := e.tokenIDF[token]
	if !exists || count == 0 {
		return math.Log(float64(e.docCount + 1))
	}
	return math.Log(float64(e.docCount+1) / count)
}

// CalculateDOBScore computes fuzzy date of birth distance with Gaussian decay & transposition handling.
func CalculateDOBScore(dob1, dob2 string) float64 {
	if dob1 == "" || dob2 == "" {
		return 0.5
	}
	if dob1 == dob2 {
		return 1.0
	}

	t1, err1 := time.Parse("2006-01-02", dob1)
	t2, err2 := time.Parse("2006-01-02", dob2)

	if err1 == nil && err2 == nil {
		daysDelta := math.Abs(t1.Sub(t2).Hours() / 24.0)

		if t1.Year() == t2.Year() && t1.Month() == time.Month(t2.Day()) && t1.Day() == int(t2.Month()) {
			return 0.95
		}

		sigma := 30.0
		return math.Exp(-0.5 * math.Pow(daysDelta/sigma, 2))
	}

	parts1 := strings.Split(dob1, "-")
	parts2 := strings.Split(dob2, "-")

	if len(parts1) == 3 && len(parts2) == 3 {
		if parts1[0] == parts2[0] {
			if parts1[1] == "00" || parts2[1] == "00" || parts1[1] == parts2[1] {
				return 0.85
			}
			return 0.70
		}
	}

	return 0.0
}

// CalculateEntropyWeightedNameScore scores two names using token-level IDF entropy weighting.
func (e *ExpertEngine) CalculateEntropyWeightedNameScore(name1, name2 string) float64 {
	tokens1 := strings.Fields(Normalize(name1))
	tokens2 := strings.Fields(Normalize(name2))

	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	matcher := NewNameMatcher(0.3, 0.4, 0.3)

	var totalWeight float64
	var weightedScore float64

	for _, t1 := range tokens1 {
		stem1 := e.StemToken(t1)
		idf1 := e.GetIDF(stem1)

		maxTokenSimilarity := 0.0
		for _, t2 := range tokens2 {
			stem2 := e.StemToken(t2)
			sim := matcher.CalculateJaroWinkler(stem1, stem2)
			if Soundex(stem1) == Soundex(stem2) {
				sim = math.Max(sim, 0.9)
			}
			if sim > maxTokenSimilarity {
				maxTokenSimilarity = sim
			}
		}

		weightedScore += maxTokenSimilarity * idf1
		totalWeight += idf1
	}

	if totalWeight == 0 {
		return 0.0
	}

	baseScore := weightedScore / totalWeight
	tokenSortScore := matcher.CalculateTokenSortScore(name1, name2)

	return (baseScore * 0.7) + (tokenSortScore * 0.3)
}

// Search performs candidate retrieval via Inverted LSH Index and multi-attribute probabilistic scoring.
func (e *ExpertEngine) Search(queryIdentity IdentityRecord, limit int) []SearchResult {
	e.mu.RLock()
	defer e.mu.RUnlock()

	queryTokens := strings.Fields(Normalize(queryIdentity.FullName))
	candidateIDs := make(set)

	for _, token := range queryTokens {
		stemmed := e.StemToken(token)
		soundexKey := Soundex(stemmed)
		if ids, found := e.invertedIdx[soundexKey]; found {
			for id := range ids {
				candidateIDs[id] = true
			}
		}
	}

	if len(candidateIDs) == 0 {
		for id := range e.records {
			candidateIDs[id] = true
		}
	}

	var results []SearchResult

	for id := range candidateIDs {
		rec := e.records[id]

		nameScore := e.CalculateEntropyWeightedNameScore(queryIdentity.FullName, rec.FullName)

		aliasScore := 0.0
		for _, alias := range rec.Aliases {
			aScore := e.CalculateEntropyWeightedNameScore(queryIdentity.FullName, alias)
			if aScore > aliasScore {
				aliasScore = aScore
			}
		}
		effectiveNameScore := math.Max(nameScore, aliasScore)

		dobScore := CalculateDOBScore(queryIdentity.DOB, rec.DOB)

		passportScore := 0.0
		if queryIdentity.PassportNo != "" && rec.PassportNo != "" {
			if queryIdentity.PassportNo == rec.PassportNo {
				passportScore = 1.0
			}
		}

		logOdds := 0.0
		if effectiveNameScore > 0.85 {
			logOdds += 4.5
		} else if effectiveNameScore < 0.5 {
			logOdds -= 3.0
		}

		if dobScore > 0.90 {
			logOdds += 3.5
		} else if dobScore < 0.3 {
			logOdds -= 2.0
		}

		if passportScore == 1.0 {
			logOdds += 8.0
		}

		composite := (effectiveNameScore * 0.50) + (dobScore * 0.30) + (passportScore * 0.20)
		compositePercent := math.Round(composite*10000) / 100.0

		var status string
		if compositePercent >= 88.0 || passportScore == 1.0 {
			status = "CRITICAL MATCH (RED)"
		} else if compositePercent >= 70.0 {
			status = "POTENTIAL MATCH (YELLOW)"
		} else {
			status = "NO MATCH (CLEAR)"
		}

		results = append(results, SearchResult{
			Record: rec,
			Breakdown: MatchBreakdown{
				NameScore:      math.Round(nameScore * 10000) / 100.0,
				DOBScore:       math.Round(dobScore * 10000) / 100.0,
				PassportScore:  math.Round(passportScore * 10000) / 100.0,
				AliasScore:     math.Round(aliasScore * 10000) / 100.0,
				FellegiLogOdds: math.Round(logOdds*100) / 100.0,
				TotalScore:     compositePercent,
				Status:         status,
			},
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Breakdown.TotalScore > results[j].Breakdown.TotalScore
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func main() {
	engine := NewExpertEngine()

	watchlist := []IdentityRecord{
		{
			ID:          "WL-001",
			FullName:    "Mohamad Ali Al-Khan",
			DOB:         "1982-05-14",
			Nationality: "ID",
			PassportNo:  "A12345678",
			Aliases:     []string{"Abu Ali", "KHAN, MOHAMMED ALY"},
		},
		{
			ID:          "WL-002",
			FullName:    "Usama Bin Laden",
			DOB:         "1957-03-10",
			Nationality: "SA",
			PassportNo:  "P98765432",
			Aliases:     []string{"Osama Bin Laden", "Abu Abdallah"},
		},
		{
			ID:          "WL-003",
			FullName:    "Vladimir Vladimirovich Putin",
			DOB:         "1952-10-07",
			Nationality: "RU",
			PassportNo:  "R55443322",
			Aliases:     []string{"VLADYMYR PUTYN"},
		},
	}

	for _, rec := range watchlist {
		engine.IndexRecord(rec)
	}

	searchQueries := []IdentityRecord{
		{
			FullName:   "MOHAMMED ALY KHAN",
			DOB:        "1982-14-05",
			PassportNo: "A12345678",
		},
		{
			FullName: "OSAMA BIN LADIN",
			DOB:      "1957-03-10",
		},
		{
			FullName: "Budi Santoso",
			DOB:      "1990-01-01",
		},
	}

	fmt.Println(strings.Repeat("=", 105))
	fmt.Printf("%-24s | %-22s | %-8s | %-8s | %-7s | %s\n",
		"QUERY NAME", "MATCHED RECORD", "NAME %", "DOB %", "TOTAL %", "STATUS")
	fmt.Println(strings.Repeat("=", 105))

	for _, query := range searchQueries {
		results := engine.Search(query, 1)
		if len(results) > 0 {
			res := results[0]
			b := res.Breakdown
			fmt.Printf("%-24s | %-22s | %6.1f%%  | %6.1f%%  | %5.1f%%  | %s\n",
				query.FullName, res.Record.FullName, b.NameScore, b.DOBScore, b.TotalScore, b.Status)
		}
	}
	fmt.Println(strings.Repeat("=", 105))
}
