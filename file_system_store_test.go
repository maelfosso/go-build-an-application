package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

func TestFileSystemStore(t *testing.T) {
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})

	t.Run("league from reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, _ := NewFileSystemPlayerStore(database)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, _ := NewFileSystemPlayerStore(database)

		got := store.GetPlayerScore("Chris")

		want := 33

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, _ := NewFileSystemPlayerStore(database)
		store.RecordWin("Chris")

		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, _ := NewFileSystemPlayerStore(database)
		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		want := 1
		assertScoreEquals(t, got, want)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("score is not equal, got %d want %d", got, want)
	}
}
