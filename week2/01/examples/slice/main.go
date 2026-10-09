package main

import (
	"errors"
	"fmt"
)

var ErrSongNotFound = errors.New("song not found")

type Song struct {
	ID      string
	Title   string
	Seconds int
}

// AddSong 返回追加后的 slice，调用方必须接收返回值。
func AddSong(songs []Song, song Song) []Song {
	return append(songs, song)
}

// FindSong 返回歌曲副本；bool 区分找到与未找到。
func FindSong(songs []Song, id string) (Song, bool) {
	for _, song := range songs {
		if song.ID == id {
			return song, true
		}
	}
	return Song{}, false
}

func TotalDuration(songs []Song) int {
	total := 0
	for _, song := range songs {
		total += song.Seconds
	}
	return total
}

// RemoveSong 根据 ID 删除歌曲，并返回更新后的 slice。
// 删除会前移底层数组中的元素，其他共享该数组的 slice 也会看到变化。
func RemoveSong(songs []Song, id string) ([]Song, error) {
	for index, song := range songs {
		if song.ID == id {
			copy(songs[index:], songs[index+1:])
			// 清空不再使用的尾元素，释放其持有的字符串引用。
			songs[len(songs)-1] = Song{}
			return songs[:len(songs)-1], nil
		}
	}

	return songs, ErrSongNotFound
}

// LongerThan 返回一个新的 slice，避免调用方修改结果时影响原 slice。
func LongerThan(songs []Song, seconds int) []Song {
	result := make([]Song, 0)
	for _, song := range songs {
		if song.Seconds > seconds {
			result = append(result, song)
		}
	}
	return result
}

func main() {
	songs := []Song{
		{ID: "s-1", Title: "Morning", Seconds: 180},
		{ID: "s-2", Title: "Journey", Seconds: 260},
	}

	// append 可能更换底层数组，所以必须接收返回的新 slice。
	songs = AddSong(songs, Song{ID: "s-3", Title: "Night", Seconds: 320})

	fmt.Printf("歌曲数=%d，容量=%d\n", len(songs), cap(songs))
	for index, song := range songs {
		fmt.Printf("%d: %s (%d 秒)\n", index, song.Title, song.Seconds)
	}

	longSongs := LongerThan(songs, 200)
	fmt.Println("超过 200 秒的歌曲:", longSongs)
	fmt.Println("总时长（秒）:", TotalDuration(songs))
	if song, exists := FindSong(songs, "s-2"); exists {
		fmt.Println("找到歌曲:", song)
	}

	var err error
	songs, err = RemoveSong(songs, "s-2")
	if err != nil {
		fmt.Println("删除失败:", err)
		return
	}
	fmt.Println("删除后的播放列表:", songs)
	if _, err := RemoveSong(songs, "missing"); err != nil {
		fmt.Println("删除不存在的歌曲:", err)
	}
}
