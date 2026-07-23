package main

import (
	"testing"
)

func TestExactMatch(t *testing.T) {
	matcher := NewNameMatcher(0.3, 0.4, 0.3)
	res := matcher.Match("John Doe", "John Doe")
	if res.CompositeScore != 100.0 {
		t.Errorf("Expected 100.0, got %f", res.CompositeScore)
	}
	if res.Status != "CRITICAL MATCH" {
		t.Errorf("Expected CRITICAL MATCH, got %s", res.Status)
	}
}

func TestWordOrderSwap(t *testing.T) {
	matcher := NewNameMatcher(0.3, 0.4, 0.3)
	res := matcher.Match("Osama Bin Laden", "Laden Bin Osama")
	if res.CompositeScore < 75.0 {
		t.Errorf("Expected >= 75.0 for swapped names, got %f", res.CompositeScore)
	}
}

func TestDisparateNames(t *testing.T) {
	matcher := NewNameMatcher(0.3, 0.4, 0.3)
	res := matcher.Match("Alexander", "Zulqarnain")
	if res.CompositeScore >= 50.0 {
		t.Errorf("Expected < 50.0 for completely different names, got %f", res.CompositeScore)
	}
	if res.Status != "NO MATCH" {
		t.Errorf("Expected NO MATCH, got %s", res.Status)
	}
}

func BenchmarkNameMatching(b *testing.B) {
	matcher := NewNameMatcher(0.3, 0.4, 0.3)
	name1 := "Mohamad Ali Khan"
	name2 := "KHAN, MOHAMMED ALY"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = matcher.Match(name1, name2)
	}
}
