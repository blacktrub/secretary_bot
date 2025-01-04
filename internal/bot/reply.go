package bot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var answerOptions = []string{
	"Записала!",
	"На карандаше.",
	"Будет сделано!",
	"Они думали мы не заметим.",
	"Они не знают с кем связались.",
}

type answerBuilder struct {
	aliaser aliaser
}

func newAnswer(aliaser aliaser) answerBuilder {
	return answerBuilder{
		aliaser: aliaser,
	}
}

func (a answerBuilder) Msg(userID int64) string {
	var b strings.Builder
	var personal bool
	if alias := a.aliaser.Alias(userID); alias != "" {
		b.WriteString(fmt.Sprintf("%s, ", alias))
		personal = true
	}

	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	m := answerOptions[r.Intn(len(answerOptions)-1)]
	if personal {
		m = strings.ToLower(m)
	}

	b.WriteString(m)
	return b.String()
}
