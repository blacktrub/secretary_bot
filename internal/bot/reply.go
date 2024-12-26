package bot

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ownerName = "Дмитрий Евгеньевич"
)

var answerOptions = []string{
	"Записала!",
	"На карандаше.",
	"Будет сделано!",
	"Они думали мы не заметим.",
	"Они не знают с кем связались.",
}

type answer struct {
	personal bool
}

func newAnswer(userID int64) answer {
	owner, err := strconv.Atoi(os.Getenv("BOT_OWNER_ID"))

	var personal bool
	if err == nil && userID == int64(owner) {
		personal = true
	}

	return answer{
		personal: personal,
	}
}

func (a answer) Msg() string {
	var b strings.Builder
	if a.personal {
		b.WriteString(fmt.Sprintf("%s, ", ownerName))
	}

	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	m := answerOptions[r.Intn(len(answerOptions)-1)]
	if a.personal {
		m = strings.ToLower(m)
	}

	b.WriteString(m)
	return b.String()
}
