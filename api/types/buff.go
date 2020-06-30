package types

import "github.com/JoeReid/buffassignment/internal/model"

type Buff struct {
	UUID             string   `json:"buff_id" yaml:"buff_id"`
	VideoStreamUUID  string   `json:"stream_id" yaml:"stream_id"`
	Question         string   `json:"question_text" yaml:"question_text"`
	CorrectAnswer    string   `json:"correct_answer" yaml:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answer" yaml:"incorrect_answer"`
}

func NewBuff(mb model.Buff) Buff {
	b := Buff{
		UUID:            mb.ID.String(),
		VideoStreamUUID: mb.Stream.String(),
		Question:        mb.Question,
	}

	for _, ans := range mb.Answers {
		if ans.Correct {
			b.CorrectAnswer = ans.Text
			continue
		}
		b.IncorrectAnswers = append(b.IncorrectAnswers, ans.Text)
	}
	return b
}

func NewBuffs(mbs []model.Buff) []Buff {
	b := make([]Buff, 0, len(mbs))

	for _, mb := range mbs {
		b = append(b, NewBuff(mb))
	}
	return b
}
