package utils

import (
	"context"
	"strings"
)

// SplitLyricsToVerses - разбивает текст песни на куплеты
//
// Ориентируется на двойной перенос строки \n\n.
// TODO: учесть максимальную длину символов куплета?
// TODO: тесты
// TODO: бенчмарк
func SplitLyricsToVerses(_ context.Context, lyrics string) (verses []string) {
	for _, verse := range strings.Split(lyrics, "\n\n") {
		if v := strings.TrimSpace(verse); v != "" {
			verses = append(verses, v)
		}
	}
	return verses
}
