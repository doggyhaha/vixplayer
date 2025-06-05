package main

import (
	"fmt"

	"github.com/doggyhaha/vixplayer"
)

func main() {
	player := vixplayer.NewVixPlayer("https://vixsrc.to", nil)
	playerData, err := player.GetMovieHLS("786892", "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Player Data for movie:", playerData)
	// Process playerData as needed

	playerData, err = player.GetShowHLS("1396", 1, 1, "")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Player Data for show:", playerData)

}
