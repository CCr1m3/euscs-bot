package db

import (
	"context"

	"github.com/euscs/euscs-bot/internal/static"
	log "github.com/sirupsen/logrus"
)

type Markov struct {
	Word1 string `db:"word1"`
	Word2 string `db:"word2"`
	Word3 string `db:"word3"`
	Count int    `db:"count"`
	Total int    `db:"total"`
}

func AddMarkovOccurences(ctx context.Context, ms []*Markov) error {
	tx, err := db.Beginx()
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to start transaction")
		return err
	}
	for _, m := range ms {
		res, err := tx.Exec("UPDATE markov SET count=count+1 WHERE word1=? AND word2=? AND word3=?", m.Word1, m.Word2, m.Word3)
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to update markov occurence")
			tx.Rollback()
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.WithFields(log.Fields{
				string(static.UUIDKey):  ctx.Value(static.UUIDKey),
				string(static.ErrorKey): err.Error(),
			}).Error("failed to get affected rows")
			tx.Rollback()
			return err
		}
		if rowsAffected == 0 {
			_, err = tx.Exec("INSERT INTO markov (word1,word2,word3,count) VALUES (?,?,?,?)", m.Word1, m.Word2, m.Word3, 1)
			if err != nil {
				log.WithFields(log.Fields{
					string(static.UUIDKey):  ctx.Value(static.UUIDKey),
					string(static.ErrorKey): err.Error(),
				}).Error("failed to insert markov occurence")
				tx.Rollback()
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

func GetMarkovOccurencesAndTotal(ctx context.Context, word1 string, word2 string) ([]*Markov, error) {
	markovs := []*Markov{}
	err := db.Select(&markovs, "SELECT word1, word2, word3, count, sum(count) OVER (PARTITION BY word1,word2) AS total FROM markov WHERE word1=? AND word2=? ORDER BY total DESC, count DESC", word1, word2)
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get markov occurences")
		return nil, err
	}
	return markovs, nil
}

func GetStartingMarkovOccurences(ctx context.Context) ([]*Markov, error) {
	markovs := []*Markov{}
	err := db.Select(&markovs, "SELECT word1, word2, sum(count) as count, (SELECT sum(count) FROM markov WHERE word1='__start__') as total FROM markov WHERE word1='__start__' GROUP BY word1,word2 ORDER BY count DESC")
	if err != nil {
		log.WithFields(log.Fields{
			string(static.UUIDKey):  ctx.Value(static.UUIDKey),
			string(static.ErrorKey): err.Error(),
		}).Error("failed to get starting markov occurences")
		return nil, err
	}
	return markovs, nil
}
