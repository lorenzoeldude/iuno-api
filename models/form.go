package models

type Form struct {

	// =====================================================
	// SURFACE FORM
	// =====================================================

	Form string `json:"form"`

	// =====================================================
	// PART OF SPEECH
	// =====================================================

	Part string `json:"part"`

	// =====================================================
	// NOMINAL FEATURES
	// =====================================================

	Case string `json:"case"`

	Number string `json:"number"`

	Gender string `json:"gender"`

	// =====================================================
	// VERBAL FEATURES
	// =====================================================

	Tense string `json:"tense"`

	Mood string `json:"mood"`

	Voice string `json:"voice"`

	Person int `json:"person"`

	// =====================================================
	// NON-FINITE FEATURES
	// =====================================================

	NonFinite string `json:"non_finite"`
}