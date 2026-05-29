package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateID(t *testing.T) {
	tests := []struct {
		name  string
		uid   string
		toUID string
		want  string
	}{
		{"Success_normal", "u1", "u2", "u1->u2"},
		{"Success_empty_both", "", "", "->"},
		{"Success_empty_uid", "", "u2", "->u2"},
		{"Success_empty_touid", "u1", "", "u1->"},
		{"Success_special_chars", "a->b", "c", "a->b->c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CreateID(tt.uid, tt.toUID))
		})
	}
}

func TestIDGenerate(t *testing.T) {
	id := IDGenerate()
	assert.NotEmpty(t, id)
}

func TestCheckAiKeyWord(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		wantOk       bool
		wantQuestion string
	}{
		{"Success_at_ai_only", "@AI", true, ""},
		{"Success_at_ai_with_question", "@AI what is Go", true, "what is Go"},
		{"Success_at_ai_middle", "hello @AI world", true, "hello  world"},
		{"Success_keyword_111", "111 hello", true, "hello"},
		{"Success_keyword_111_with_context", "prefix 111 suffix", true, "prefix  suffix"},
		{"Fail_no_match", "hello world", false, ""},
		{"Fail_empty_string", "", false, ""},
		{"Fail_case_sensitive", "@ai", false, ""},
		{"Fail_similar_keyword", "@ss", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, question := CheckAiKeyWord(tt.message)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantQuestion, question)
		})
	}
}
