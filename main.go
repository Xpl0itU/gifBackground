package main

import (
	"image/gif"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/disintegration/imaging"
	"github.com/reujab/wallpaper"
)

func setWallpaper(path string) error {
	return wallpaper.SetFromFile(path)
}

func main() {
	originalWallpaper, err := wallpaper.Get()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP, os.Kill)
	go func() {
		<-c
		err := setWallpaper(originalWallpaper)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	file, err := os.Open("contornos.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	g, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatal(err)
	}

	tmpPath, err := os.MkdirTemp(os.TempDir(), "gifBackground")
	if err != nil {
		log.Fatal(err)
	}
	tmpImgPath := filepath.Join(tmpPath, "temp.jpg")

	for {
		for _, frame := range g.Image {
			img := imaging.Clone(frame)
			err := imaging.Save(img, tmpImgPath)
			if err != nil {
				log.Fatal(err)
			}

			err = setWallpaper(tmpImgPath)
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(250 * time.Millisecond)
		}
	}
}
