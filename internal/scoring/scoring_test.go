package scoring

import (
	"testing"
)

func TestRating_String(t *testing.T) {
	tests := []struct {
		rating Rating
		want   string
	}{
		{RatingNone, "None (0)"},
		{RatingPartial, "Partial (1)"},
		{RatingGood, "Good (2)"},
		{RatingExcellent, "Excellent (3)"},
		{Rating(99), "Unknown (99)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.rating.String(); got != tt.want {
				t.Errorf("Rating.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestScore_Total(t *testing.T) {
	tests := []struct {
		name  string
		score Score
		want  int
	}{
		{
			name: "all zeros",
			score: Score{
				Completeness:       Dimension{Rating: RatingNone},
				LintQuality:        Dimension{Rating: RatingNone},
				CodeQuality:        Dimension{Rating: RatingNone},
				OutputValidity:     Dimension{Rating: RatingNone},
				QuestionEfficiency: Dimension{Rating: RatingNone},
			},
			want: 0,
		},
		{
			name: "all excellent",
			score: Score{
				Completeness:       Dimension{Rating: RatingExcellent},
				LintQuality:        Dimension{Rating: RatingExcellent},
				CodeQuality:        Dimension{Rating: RatingExcellent},
				OutputValidity:     Dimension{Rating: RatingExcellent},
				QuestionEfficiency: Dimension{Rating: RatingExcellent},
			},
			want: 15,
		},
		{
			name: "mixed",
			score: Score{
				Completeness:       Dimension{Rating: RatingExcellent},
				LintQuality:        Dimension{Rating: RatingGood},
				CodeQuality:        Dimension{Rating: RatingPartial},
				OutputValidity:     Dimension{Rating: RatingNone},
				QuestionEfficiency: Dimension{Rating: RatingGood},
			},
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.score.Total(); got != tt.want {
				t.Errorf("Score.Total() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestScore_Threshold(t *testing.T) {
	tests := []struct {
		total int
		want  string
	}{
		{0, "Failure"},
		{5, "Failure"},
		{6, "Partial"},
		{9, "Partial"},
		{10, "Success"},
		{12, "Success"},
		{13, "Excellent"},
		{15, "Excellent"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			// Create a score that totals to tt.total
			score := Score{
				Completeness:       Dimension{Rating: Rating(tt.total)},
				LintQuality:        Dimension{Rating: RatingNone},
				CodeQuality:        Dimension{Rating: RatingNone},
				OutputValidity:     Dimension{Rating: RatingNone},
				QuestionEfficiency: Dimension{Rating: RatingNone},
			}
			if got := score.Threshold(); got != tt.want {
				t.Errorf("Score.Threshold() = %q, want %q for total %d", got, tt.want, tt.total)
			}
		})
	}
}

func TestScore_Passed(t *testing.T) {
	tests := []struct {
		total int
		want  bool
	}{
		{0, false},
		{5, false},
		{6, true},
		{10, true},
		{15, true},
	}

	for _, tt := range tests {
		score := Score{
			Completeness: Dimension{Rating: Rating(tt.total)},
		}
		if got := score.Passed(); got != tt.want {
			t.Errorf("Score.Passed() = %v, want %v for total %d", got, tt.want, tt.total)
		}
	}
}

func TestNewScore(t *testing.T) {
	score := NewScore("beginner", "ci-workflow")

	if score.Persona != "beginner" {
		t.Errorf("Persona = %q, want %q", score.Persona, "beginner")
	}
	if score.Scenario != "ci-workflow" {
		t.Errorf("Scenario = %q, want %q", score.Scenario, "ci-workflow")
	}
	if score.Completeness.Name != "Completeness" {
		t.Errorf("Completeness.Name = %q, want %q", score.Completeness.Name, "Completeness")
	}
	if score.LintQuality.Name != "Lint Quality" {
		t.Errorf("LintQuality.Name = %q, want %q", score.LintQuality.Name, "Lint Quality")
	}
}

func TestScoreCompleteness(t *testing.T) {
	tests := []struct {
		expected int
		actual   int
		rating   Rating
	}{
		{0, 0, RatingExcellent},
		{10, 10, RatingExcellent},
		{10, 8, RatingGood},
		{10, 5, RatingPartial},
		{10, 3, RatingNone},
	}

	for _, tt := range tests {
		rating, _ := ScoreCompleteness(tt.expected, tt.actual)
		if rating != tt.rating {
			t.Errorf("ScoreCompleteness(%d, %d) = %v, want %v", tt.expected, tt.actual, rating, tt.rating)
		}
	}
}

func TestScoreLintQuality(t *testing.T) {
	tests := []struct {
		cycles int
		passed bool
		rating Rating
	}{
		{0, true, RatingExcellent},
		{1, true, RatingExcellent},
		{2, true, RatingGood},
		{3, true, RatingPartial},
		{5, true, RatingPartial},
		{0, false, RatingNone},
		{3, false, RatingNone},
	}

	for _, tt := range tests {
		rating, _ := ScoreLintQuality(tt.cycles, tt.passed)
		if rating != tt.rating {
			t.Errorf("ScoreLintQuality(%d, %v) = %v, want %v", tt.cycles, tt.passed, rating, tt.rating)
		}
	}
}

func TestScoreCodeQuality(t *testing.T) {
	tests := []struct {
		issues []string
		rating Rating
	}{
		{nil, RatingExcellent},
		{[]string{}, RatingExcellent},
		{[]string{"issue1"}, RatingGood},
		{[]string{"issue1", "issue2"}, RatingGood},
		{[]string{"issue1", "issue2", "issue3"}, RatingPartial},
		{[]string{"1", "2", "3", "4", "5", "6"}, RatingNone},
	}

	for _, tt := range tests {
		rating, _ := ScoreCodeQuality(tt.issues)
		if rating != tt.rating {
			t.Errorf("ScoreCodeQuality(%v) = %v, want %v", tt.issues, rating, tt.rating)
		}
	}
}

func TestScoreOutputValidity(t *testing.T) {
	tests := []struct {
		errors   int
		warnings int
		rating   Rating
	}{
		{0, 0, RatingExcellent},
		{0, 1, RatingGood},
		{0, 2, RatingGood},
		{0, 3, RatingPartial},
		{1, 0, RatingNone},
		{2, 5, RatingNone},
	}

	for _, tt := range tests {
		rating, _ := ScoreOutputValidity(tt.errors, tt.warnings)
		if rating != tt.rating {
			t.Errorf("ScoreOutputValidity(%d, %d) = %v, want %v", tt.errors, tt.warnings, rating, tt.rating)
		}
	}
}

func TestScoreQuestionEfficiency(t *testing.T) {
	tests := []struct {
		questions int
		rating    Rating
	}{
		{0, RatingExcellent},
		{2, RatingExcellent},
		{3, RatingGood},
		{4, RatingGood},
		{5, RatingPartial},
		{6, RatingPartial},
		{7, RatingNone},
		{10, RatingNone},
	}

	for _, tt := range tests {
		rating, _ := ScoreQuestionEfficiency(tt.questions)
		if rating != tt.rating {
			t.Errorf("ScoreQuestionEfficiency(%d) = %v, want %v", tt.questions, rating, tt.rating)
		}
	}
}

func TestScore_String(t *testing.T) {
	score := NewScore("beginner", "ci-workflow")
	score.Completeness.Rating = RatingExcellent
	score.Completeness.Notes = "All workflows generated"
	score.LintQuality.Rating = RatingGood
	score.CodeQuality.Rating = RatingGood
	score.OutputValidity.Rating = RatingExcellent
	score.QuestionEfficiency.Rating = RatingExcellent

	s := score.String()

	if !contains(s, "Score: 13/15") {
		t.Error("String() should contain total score")
	}
	if !contains(s, "Excellent") {
		t.Error("String() should contain threshold")
	}
	if !contains(s, "beginner") {
		t.Error("String() should contain persona")
	}
	if !contains(s, "ci-workflow") {
		t.Error("String() should contain scenario")
	}
	if !contains(s, "All workflows generated") {
		t.Error("String() should contain notes")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || contains(s[1:], substr)))
}
