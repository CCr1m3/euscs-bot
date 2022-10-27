package db

import (
	"github.com/haashi/omega-strikers-bot/internal/models"
)

func AddMarkovOccurences(ms []*models.Markov) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	for _, m := range ms {
		res, err := tx.Exec("UPDATE markov SET count=count+1 WHERE word1=? AND word2=? AND word3=?", m.Word1, m.Word2, m.Word3)
		if err != nil {
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			_, err = tx.Exec("INSERT INTO markov (word1,word2,word3,count) VALUES (?,?,?,?)", m.Word1, m.Word2, m.Word3, 1)
			if err != nil {
				return err
			}

		}
	}
	return tx.Commit()
}

func DeleteAllMarkov() error {
	_, err := db.Exec("DELETE from markov")
	return err
}

func GetMarkovOccurencesAndTotal(word1 string, word2 string) ([]*models.Markov, error) {
	markovs := []*models.Markov{}
	err := db.Select(&markovs, "SELECT word1, word2, word3, count, sum(count) OVER (PARTITION BY word1,word2) AS total FROM markov WHERE word1=? AND word2=? ORDER BY total DESC, count DESC", word1, word2)
	if err != nil {
		return nil, err
	}
	return markovs, nil
}

func GetStartingMarkovOccurences() ([]*models.Markov, error) {
	markovs := []*models.Markov{}
	err := db.Select(&markovs, "SELECT word1, word2, sum(count) as count, (SELECT sum(count) FROM markov WHERE word1='__start__') as total FROM markov WHERE word1='__start__' GROUP BY word1,word2 ORDER BY count DESC")
	if err != nil {
		return nil, err
	}
	return markovs, nil
}
