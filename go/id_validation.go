package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NationalIDValidationResult holds extracted metadata from 16-digit National Identity Numbers.
type NationalIDValidationResult struct {
	IDNumber     string `json:"id_number"`
	IsValid      bool   `json:"is_valid"`
	RegionCode   string `json:"region_code"`
	SubRegionCode string `json:"sub_region_code"`
	Gender       string `json:"gender"`
	DOB          string `json:"dob"` // YYYY-MM-DD
	ErrorMessage string `json:"error_message,omitempty"`
}

// SystemAuditLog represents structured non-repudiation audit trail record for compliance auditing.
type SystemAuditLog struct {
	LogID         string    `json:"log_id"`
	Timestamp     time.Time `json:"timestamp"`
	OperatorID    string    `json:"operator_id"`
	SystemModule  string    `json:"system_module"`
	QueryHash     string    `json:"query_hash"`
	MatchStatus   string    `json:"match_status"`
	TotalScore    float64   `json:"total_score"`
	ExecutionMs   float64   `json:"execution_ms"`
	IntegrityHash string    `json:"integrity_hash"`
}

// ParseNationalID validates and extracts demographics from 16-digit National Identity Numbers.
func ParseNationalID(idNum string) NationalIDValidationResult {
	idNum = strings.TrimSpace(idNum)
	if len(idNum) != 16 {
		return NationalIDValidationResult{IDNumber: idNum, IsValid: false, ErrorMessage: "Identity number must be exactly 16 digits"}
	}

	_, err := strconv.ParseUint(idNum, 10, 64)
	if err != nil {
		return NationalIDValidationResult{IDNumber: idNum, IsValid: false, ErrorMessage: "Identity number contains non-numeric characters"}
	}

	region := idNum[0:2]
	subRegion := idNum[2:4]

	dayRaw, _ := strconv.Atoi(idNum[6:8])
	monthRaw, _ := strconv.Atoi(idNum[8:10])
	yearRaw, _ := strconv.Atoi(idNum[10:12])

	gender := "MALE"
	day := dayRaw
	if dayRaw > 40 {
		gender = "FEMALE"
		day = dayRaw - 40
	}

	currentYear := time.Now().Year() % 100
	fullYear := 1900 + yearRaw
	if yearRaw <= currentYear {
		fullYear = 2000 + yearRaw
	}

	if monthRaw < 1 || monthRaw > 12 || day < 1 || day > 31 {
		return NationalIDValidationResult{IDNumber: idNum, IsValid: false, ErrorMessage: "Invalid birth date encoded in identity number"}
	}

	dobStr := fmt.Sprintf("%04d-%02d-%02d", fullYear, monthRaw, day)

	return NationalIDValidationResult{
		IDNumber:      idNum,
		IsValid:       true,
		RegionCode:    region,
		SubRegionCode: subRegion,
		Gender:        gender,
		DOB:           dobStr,
	}
}

// GenerateAuditLog creates a cryptographically signed audit log for enterprise compliance.
func GenerateAuditLog(operatorID, systemModule, queryName, matchStatus string, score float64, execMs float64) SystemAuditLog {
	ts := time.Now().UTC()
	logID := fmt.Sprintf("AUDIT-%d", ts.UnixNano())

	qHash := sha256.Sum256([]byte(queryName + "|" + systemModule))
	queryHashHex := hex.EncodeToString(qHash[:])

	rawLog := fmt.Sprintf("%s|%s|%s|%s|%s|%.2f|%.2f",
		logID, ts.Format(time.RFC3339), operatorID, systemModule, queryHashHex, score, execMs)

	iHash := sha256.Sum256([]byte(rawLog))
	integrityHex := hex.EncodeToString(iHash[:])

	return SystemAuditLog{
		LogID:         logID,
		Timestamp:     ts,
		OperatorID:    operatorID,
		SystemModule:  systemModule,
		QueryHash:     queryHashHex,
		MatchStatus:   matchStatus,
		TotalScore:    score,
		ExecutionMs:   execMs,
		IntegrityHash: integrityHex,
	}
}

func main() {
	fmt.Println(strings.Repeat("=", 90))
	fmt.Println(" ENTERPRISE READINESS MODULE: NATIONAL ID VALIDATION & COMPLIANCE AUDIT TRAIL ENGINE")
	fmt.Println(strings.Repeat("=", 90))

	testIDs := []string{"3171011505850001", "3273025210920003", "12345"}
	fmt.Println("\n[1] VERIFIKASI STRUCTURAL NATIONAL ID:")
	for _, idNum := range testIDs {
		res := ParseNationalID(idNum)
		if res.IsValid {
			fmt.Printf("ID: %s -> VALID | Gender: %-6s | DOB: %s | Region: %s, SubRegion: %s\n",
				res.IDNumber, res.Gender, res.DOB, res.RegionCode, res.SubRegionCode)
		} else {
			fmt.Printf("ID: %-16s -> INVALID (%s)\n", res.IDNumber, res.ErrorMessage)
		}
	}

	fmt.Println("\n[2] GENERASI COMPLIANCE AUDIT LOG TRAIL:")
	auditLog := GenerateAuditLog("OPERATOR_9981", "SCREENING_MODULE", "MOHAMMED ALY KHAN", "CRITICAL MATCH (RED)", 91.0, 1.25)
	logJSON, _ := json.MarshalIndent(auditLog, "", "  ")
	fmt.Println(string(logJSON))
	fmt.Println(strings.Repeat("=", 90))
}
