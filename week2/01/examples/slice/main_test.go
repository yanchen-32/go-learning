package main

import (
	"errors"
	"reflect"
	"testing"
)

func sampleSongs() []Song {
	return []Song{
		{ID: "s-1", Title: "First", Seconds: 100},
		{ID: "s-2", Title: "Second", Seconds: 200},
		{ID: "s-3", Title: "Third", Seconds: 300},
	}
}

func TestRemoveSong(t *testing.T) {
	for _, tt := range []struct {
		name string
		id   string
		want []Song
	}{
		{"first", "s-1", sampleSongs()[1:]},
		{"middle", "s-2", []Song{sampleSongs()[0], sampleSongs()[2]}},
		{"last", "s-3", sampleSongs()[:2]},
		{"missing", "missing", sampleSongs()},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveSong(sampleSongs(), tt.id)
			if tt.name == "missing" {
				if !errors.Is(err, ErrSongNotFound) {
					t.Fatalf("error = %v, want ErrSongNotFound", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("remaining songs = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestAddFindAndDuration(t *testing.T) {
	var songs []Song
	for _, song := range sampleSongs() {
		songs = AddSong(songs, song)
	}
	if !reflect.DeepEqual(songs, sampleSongs()) {
		t.Fatalf("added songs = %+v", songs)
	}
	if song, ok := FindSong(songs, "s-2"); !ok || song != sampleSongs()[1] {
		t.Fatalf("FindSong() = %+v, %v", song, ok)
	}
	if song, ok := FindSong(songs, "missing"); ok || song != (Song{}) {
		t.Fatalf("missing song = %+v, %v", song, ok)
	}
	if got := TotalDuration(songs); got != 600 {
		t.Fatalf("duration = %d, want 600", got)
	}
}

func TestLongerThan(t *testing.T) {
	songs := sampleSongs()
	filtered := LongerThan(songs, 200)
	if !reflect.DeepEqual(filtered, sampleSongs()[2:]) {
		t.Fatalf("LongerThan(200) = %+v", filtered)
	}
	filtered[0].Title = "changed"
	if songs[2].Title != "Third" {
		t.Fatal("filter result shares the original array")
	}
	if got := LongerThan(songs, 300); len(got) != 0 {
		t.Fatalf("LongerThan(300) = %+v, want empty", got)
	}
}

func TestEmptyAndSingleSong(t *testing.T) {
	if TotalDuration(nil) != 0 || len(LongerThan(nil, 0)) != 0 {
		t.Fatal("empty playlist should have zero duration and no matches")
	}
	if got, err := RemoveSong(nil, "missing"); got != nil || !errors.Is(err, ErrSongNotFound) {
		t.Fatalf("remove from empty playlist = %+v, %v", got, err)
	}
	if _, ok := FindSong(nil, "missing"); ok {
		t.Fatal("empty playlist should have no matching song")
	}
	songs := AddSong(nil, Song{ID: "only", Seconds: 1})
	if got, err := RemoveSong(songs, "only"); err != nil || len(got) != 0 {
		t.Fatalf("remove only song = %+v, %v", got, err)
	}
}
