package morphology

// import "iuno-api/models"

type Word struct {
    Lemma string
    Type  string
    Gender string
    Stem string
}

type Form struct {
    Form    string `json:"form"`

    Part    string `json:"part"`

    Case    string `json:"case,omitempty"`
    Number  string `json:"number,omitempty"`

    Gender  string `json:"gender,omitempty"`

    Tense   string `json:"tense,omitempty"`
    Mood    string `json:"mood,omitempty"`
    Voice   string `json:"voice,omitempty"`
    Person  int    `json:"person,omitempty"`
}