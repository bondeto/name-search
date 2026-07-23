import unittest
from matcher import NameMatcher

class TestNameMatcher(unittest.TestCase):
    def setUp(self):
        self.matcher = NameMatcher(w_phonetic=0.3, w_distance=0.4, w_token=0.3)

    def test_exact_match(self):
        result = self.matcher.match("John Doe", "John Doe")
        self.assertEqual(result["composite_score"], 100.0)
        self.assertEqual(result["status"], "CRITICAL MATCH")

    def test_case_and_punctuation_insensitivity(self):
        result = self.matcher.match("john-doe", "JOHN DOE!")
        self.assertEqual(result["composite_score"], 100.0)

    def test_word_order_swap(self):
        result = self.matcher.match("Osama Bin Laden", "Laden Bin Osama")
        self.assertGreaterEqual(result["composite_score"], 80.0)
        self.assertIn(result["status"], ["POTENTIAL MATCH", "CRITICAL MATCH"])

    def test_transliteration_and_typo(self):
        result = self.matcher.match("Muhamad Ali", "Mohammed Ally")
        self.assertGreaterEqual(result["composite_score"], 75.0)

    def test_disparate_names(self):
        result = self.matcher.match("Alexander", "Zulqarnain")
        self.assertLess(result["composite_score"], 50.0)
        self.assertEqual(result["status"], "NO MATCH")

if __name__ == "__main__":
    unittest.main()
