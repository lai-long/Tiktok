package main

import (
	react "Tiktok/kitex_gen/react/likeservice"
	"log"
)

func main() {
	svr := react.NewServer(new(LikeServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
