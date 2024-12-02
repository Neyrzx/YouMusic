package utils_test

import (
	"context"
	"testing"

	"github.com/neyrzx/youmusic/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSplitLyricsToVerses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		lyrics         string
		expectedVerses []string
	}{
		{
			"case: #1",
			"Verse1\n\nVerse2\n\n",
			[]string{"Verse1", "Verse2"},
		},
		{
			"case: #3",
			"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight  ",
			[]string{
				"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?",
				"Ooh\nYou set my soul alight\nOoh\nYou set my soul alight",
			},
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actualVerses := utils.SplitLyricsToVerses(ctx, test.lyrics)

			assert.Equal(t, test.expectedVerses, actualVerses)
		})
	}
}
