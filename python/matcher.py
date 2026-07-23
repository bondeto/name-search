"""
Hybrid Fuzzy Name Matcher Engine (Python Implementation)
Provides multi-layered scoring: Phonetic Matching, Jaro-Winkler Distance, and Token Sorting.
"""

from typing import Dict, Tuple, Set
import re

try:
    from rapidfuzz import distance, fuzz
    HAS_RAPIDFUZZ = True
except ImportError:
    HAS_RAPIDFUZZ = False

try:
    import doublemetaphone
    HAS_METAPHONE = True
except ImportError:
    HAS_METAPHONE = False


class NameMatcher:
    """
    Intelligence-grade Hybrid Name Matching Scorer.
    Combines:
    1. Phonetic matching (Double Metaphone)
    2. Jaro-Winkler Edit Distance
    3. Token Permutation (handles word order swaps)
    """

    def __init__(self, w_phonetic: float = 0.3, w_distance: float = 0.4, w_token: float = 0.3):
        total_weight = w_phonetic + w_distance + w_token
        self.w_phonetic = w_phonetic / total_weight
        self.w_distance = w_distance / total_weight
        self.w_token = w_token / total_weight

    def _normalize(self, name: str) -> str:
        """Clean and normalize name string."""
        name = name.upper()
        name = re.sub(r'[^A-Z0-9\s]', ' ', name)
        return re.sub(r'\s+', ' ', name).strip()

    def _simple_soundex(self, name: str) -> str:
        """Fallback Soundex encoder if Double Metaphone library is not installed."""
        name = re.sub(r'[^A-Z]', '', name.upper())
        if not name:
            return "Z000"
        
        mapping = {
            'B': '1', 'F': '1', 'P': '1', 'V': '1',
            'C': '2', 'G': '2', 'J': '2', 'K': '2', 'Q': '2', 'S': '2', 'X': '2', 'Z': '2',
            'D': '3', 'T': '3',
            'L': '4',
            'M': '5', 'N': '5',
            'R': '6'
        }
        
        first_letter = name[0]
        codes = [first_letter]
        
        for char in name[1:]:
            code = mapping.get(char, '0')
            if code != '0' and code != codes[-1]:
                codes.append(code)
                
        soundex = "".join(codes).replace('0', '')
        return (soundex + "0000")[:4]

    def _get_phonetic_keys(self, token: str) -> Set[str]:
        """Extract phonetic representation for a single token."""
        if HAS_METAPHONE:
            primary, secondary = doublemetaphone.doublemetaphone(token)
            return {k for k in [primary, secondary] if k}
        else:
            return {self._simple_soundex(token)}

    def calculate_phonetic_score(self, name1: str, name2: str) -> float:
        """Calculate phonetic token overlap score."""
        tokens1 = name1.split()
        tokens2 = name2.split()

        keys1 = {k for token in tokens1 for k in self._get_phonetic_keys(token)}
        keys2 = {k for token in tokens2 for k in self._get_phonetic_keys(token)}

        if not keys1 or not keys2:
            return 0.0

        intersection = keys1.intersection(keys2)
        union = keys1.union(keys2)
        return len(intersection) / len(union)

    def calculate_jaro_winkler(self, s1: str, s2: str) -> float:
        """Calculate Jaro-Winkler similarity."""
        if HAS_RAPIDFUZZ:
            return distance.JaroWinkler.similarity(s1, s2)
        
        # Pure Python fallback Jaro-Winkler calculation
        if s1 == s2:
            return 1.0
        
        len1, len2 = len(s1), len(s2)
        if len1 == 0 or len2 == 0:
            return 0.0

        match_distance = max(len1, len2) // 2 - 1
        if match_distance < 0:
            match_distance = 0

        s1_matches = [False] * len1
        s2_matches = [False] * len2

        matches = 0
        transpositions = 0

        for i in range(len1):
            start = max(0, i - match_distance)
            end = min(i + match_distance + 1, len2)

            for j in range(start, end):
                if s2_matches[j]:
                    continue
                if s1[i] != s2[j]:
                    continue
                s1_matches[i] = True
                s2_matches[j] = True
                matches += 1
                break

        if matches == 0:
            return 0.0

        k = 0
        for i in range(len1):
            if not s1_matches[i]:
                continue
            while not s2_matches[k]:
                k += 1
            if s1[i] != s2[k]:
                transpositions += 1
            k += 1

        jaro = (matches / len1 + matches / len2 + (matches - transpositions / 2) / matches) / 3.0

        # Winkler prefix scaling
        prefix = 0
        for i in range(min(4, min(len1, len2))):
            if s1[i] == s2[i]:
                prefix += 1
            else:
                break

        return jaro + prefix * 0.1 * (1 - jaro)

    def calculate_token_sort_score(self, s1: str, s2: str) -> float:
        """Calculate score independent of word order."""
        if HAS_RAPIDFUZZ:
            return fuzz.token_sort_ratio(s1, s2) / 100.0
        
        sorted1 = " ".join(sorted(s1.split()))
        sorted2 = " ".join(sorted(s2.split()))
        return self.calculate_jaro_winkler(sorted1, sorted2)

    def match(self, input_name: str, candidate_name: str) -> Dict[str, any]:
        """Perform hybrid name matching and return detailed scoring breakdown."""
        norm1 = self._normalize(input_name)
        norm2 = self._normalize(candidate_name)

        s_phonetic = self.calculate_phonetic_score(norm1, norm2)
        s_distance = self.calculate_jaro_winkler(norm1, norm2)
        s_token = self.calculate_token_sort_score(norm1, norm2)

        composite = (
            (self.w_phonetic * s_phonetic) +
            (self.w_distance * s_distance) +
            (self.w_token * s_token)
        )
        composite_percent = round(composite * 100, 2)

        if composite_percent >= 90.0:
            status = "CRITICAL MATCH"
        elif composite_percent >= 75.0:
            status = "POTENTIAL MATCH"
        else:
            status = "NO MATCH"

        return {
            "input_name": input_name,
            "candidate_name": candidate_name,
            "phonetic_score": round(s_phonetic * 100, 2),
            "distance_score": round(s_distance * 100, 2),
            "token_score": round(s_token * 100, 2),
            "composite_score": composite_percent,
            "status": status
        }


if __name__ == "__main__":
    matcher = NameMatcher()
    
    test_cases = [
        ("Mohamad Ali Khan", "KHAN, MOHAMMED ALY"),
        ("Osama Bin Laden", "USAMA BIN LADIN"),
        ("Vladimir Putin", "VLADYMYR PUTYN"),
        ("John Smith", "JONATHAN SMITH"),
        ("Budi Santoso", "SANTOSO, BUDI"),
    ]

    print("=" * 80)
    print(f"{'INPUT NAME':<22} | {'CANDIDATE NAME':<22} | {'SCORE':<7} | STATUS")
    print("=" * 80)

    for query, target in test_cases:
        res = matcher.match(query, target)
        print(f"{res['input_name']:<22} | {res['candidate_name']:<22} | {res['composite_score']:>5.1f}% | {res['status']}")
    
    print("=" * 80)
