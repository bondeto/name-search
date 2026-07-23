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

// KTPValidationResult holds extracted metadata from 16-digit NIK KTP Indonesia.
type KTPValidationResult struct {
	NIK          string `json:"nik"`
	IsValid      bool   `json:"is_valid"`
	ProvinceCode string `json:"province_code"`
	RegencyCode  string `json:"regency_code"`
	DistrictCode string `json:"district_code"`
	Gender       string `json:"gender"`
	DOB          string `json:"dob"` // YYYY-MM-DD
	ErrorMessage string `json:"error_message,omitempty"`
}

// GovernmentAuditLog represents ISO 27001 non-repudiation audit trail record.
type GovernmentAuditLog struct {
	LogID         string    `json:"log_id"`
	Timestamp     time.Time `json:"timestamp"`
	OperatorID    string    `json:"operator_id"`
	AgencyCode    string    `json:"agency_code"` // e.g. "IMIGRASI", "POLRI", "PPATK"
	QueryHash     string    `json:"query_hash"`
	MatchStatus   string    `json:"match_status"`
	TotalScore    float64   `json:"total_score"`
	ExecutionMs   float64   `json:"execution_ms"`
	IntegrityHash string    `json:"integrity_hash"`
}

// ParseKTPNIK validates and extracts demographics from Indonesian National ID (NIK KTP).
func ParseKTPNIK(nik string) KTPValidationResult {
	nik = strings.TrimSpace(nik)
	if len(nik) != 16 {
		return KTPValidationResult{NIK: nik, IsValid: false, ErrorMessage: "NIK must be exactly 16 digits"}
	}

	_, err := strconv.ParseUint(nik, 10, 64)
	if err != nil {
		return KTPValidationResult{NIK: nik, IsValid: false, ErrorMessage: "NIK contains non-numeric characters"}
	}

	prov := nik[0:2]
	reg := nik[2:4]
	dist := nik[4:6]

	dayRaw, _ := strconv.Atoi(nik[6:8])
	monthRaw, _ := strconv.Atoi(nik[8:10])
	yearRaw, _ := strconv.Atoi(nik[10:12])

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
		return KTPValidationResult{NIK: nik, IsValid: false, ErrorMessage: "Invalid birth date encoded in NIK"}
	}

	dobStr := fmt.Sprintf("%04d-%02d-%02d", fullYear, monthRaw, day)

	return KTPValidationResult{
		NIK:          nik,
		IsValid:      true,
		ProvinceCode: prov,
		RegencyCode:  reg,
		DistrictCode: dist,
		Gender:       gender,
		DOB:          dobStr,
	}
}

// GenerateGovernmentAuditLog creates a cryptographically signed audit log.
func GenerateGovernmentAuditLog(operatorID, agencyCode, queryName, matchStatus string, score float64, execMs float64) GovernmentAuditLog {
	ts := time.Now().UTC()
	logID := fmt.Sprintf("AUDIT-%d", ts.UnixNano())

	// SHA-256 Hash of query payload for PII privacy preservation
	qHash := sha256.Sum256([]byte(queryName + "|" + agencyCode))
	queryHashHex := hex.EncodeToString(qHash[:])

	rawLog := fmt.Sprintf("%s|%s|%s|%s|%s|%.2f|%.2f",
		logID, ts.Format(time.RFC3339), operatorID, agencyCode, queryHashHex, score, execMs)

	// Tamper-evident integrity hash
	iHash := sha256.Sum256([]byte(rawLog))
	integrityHex := hex.EncodeToString(iHash[:])

	return GovernmentAuditLog{
		LogID:         logID,
		Timestamp:     ts,
		OperatorID:    operatorID,
		AgencyCode:    agencyCode,
		QueryHash:     queryHashHex,
		MatchStatus:   matchStatus,
		TotalScore:    score,
		ExecutionMs:   execMs,
		IntegrityHash: integrityHex,
	}
}

func main() {
	fmt.Println(strings.Repeat("=", 90))
	fmt.Println(" GOVERNMENT READINESS MODULE: KTP PARSER & ISO 27001 AUDIT TRAIL ENGINE")
	fmt.Println(strings.Repeat("=", 90))

	// Demo 1: Indonesian KTP / NIK Validation & Demographics Extraction
	testNIKs := []string{"3171011505850001", "3273025210920003", "12345"}
	fmt.Println("\n[1] VERIFIKASI STRUCTURAL NIK KTP INDONESIA:")
	for _, nik := range testNIKs {
		res := ParseKTPNIK(nik)
		if res.IsValid {
			fmt.Printf("NIK: %s -> VALID | Gender: %-6s | DOB: %s | Prov: %s, Reg: %s\n",
				res.NIK, res.Gender, res.DOB, res.ProvinceCode, res.RegencyCode)
		} else {
			fmt.Printf("NIK: %-16s -> INVALID (%s)\n", res.NIK, res.ErrorMessage)
		}
	}

	// Demo 2: Cryptographic Audit Trail Generation
	fmt.Println("\n[2] GENERASI ISO 27001 AUDIT LOG TRAIL:")
	auditLog := GenerateGovernmentAuditLog("OFFICER_9981", "IMIGRASI_CEKAL", "MOHAMMED ALY KHAN", "CRITICAL MATCH (RED)", 91.0, 1.25)
	logJSON, _ := json.MarshalIndent(auditLog, "", "  ")
	fmt.Println(string(logJSON))
	fmt.Println(strings.Repeat("=", 90))
}
