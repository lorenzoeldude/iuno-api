package models

type Form struct {
    Form     string `json:"form"`
    Part     string `json:"part"`

    Case     string `json:"case"`
    Number   string `json:"number"`

    Gender   string `json:"gender"`

    Tense    string `json:"tense"`
    Mood     string `json:"mood"`
    Voice    string `json:"voice"`

    Person   int    `json:"person"`
}