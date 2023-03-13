package db

import (
	"context"
	"testing"
)

func Test_db_Markov(t *testing.T) {
	Clear()
	Init()
	ctx := context.TODO()
	m1 := &Markov{Word1: "__start__", Word2: "i", Word3: "love"}
	m2 := &Markov{Word1: "__start__", Word2: "i", Word3: "must"}
	err := AddMarkovOccurences(ctx, []*Markov{m1, m1, m1, m2})
	if err != nil {
		t.Errorf("failed to add markov occurences: " + err.Error())
	}
	ms, err := GetMarkovOccurencesAndTotal(ctx, "", "")
	if err != nil {
		t.Errorf("failed to get markov occurences: " + err.Error())
	}
	if len(ms) != 0 {
		t.Errorf("found unexisting occurences")
	}

	ms, err = GetMarkovOccurencesAndTotal(ctx, "__start__", "i")
	if err != nil {
		t.Errorf("failed to get markov occurences: " + err.Error())
	}
	if len(ms) > 2 {
		t.Errorf("found unexisting occurences")
	}
	if len(ms) < 2 {
		t.Errorf("missing occurences")
	}
	if ms[0].Count != 3 || ms[0].Total != 4 {
		t.Errorf("sums are wrong")
	}

	ms, err = GetStartingMarkovOccurences(ctx)
	if err != nil {
		t.Errorf("failed to get markov occurences: " + err.Error())
	}
	if len(ms) > 1 {
		t.Errorf("found unexisting occurences")
	}
	if len(ms) < 1 {
		t.Errorf("missing occurences")
	}
	if ms[0].Count != 4 || ms[0].Total != 4 {
		t.Errorf("sums are wrong")
	}

	err = DeleteAllMarkov()
	if err != nil {
		t.Errorf("failed to delete markov occurences: " + err.Error())
	}
}
