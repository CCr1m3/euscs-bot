package models

type Markov struct {
	Word1 string `db:"word1"`
	Word2 string `db:"word2"`
	Word3 string `db:"word3"`
	Count int    `db:"count"`
	Total int    `db:"total"`
}
