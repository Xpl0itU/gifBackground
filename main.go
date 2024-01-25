package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/reujab/wallpaper"
)

func setWallpaper(path string) error {
	return wallpaper.SetFromFile(path)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func main() {
	originalWallpaper, err := wallpaper.Get()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP, os.Kill)

	frameFolder := "frames"
	imageFiles, err := filepath.Glob(filepath.Join(frameFolder, "*.png"))
	if err != nil {
		log.Fatal(err)
	}
	tempFolder, err := os.MkdirTemp(os.TempDir(), "gifWallpaper")
	if err != nil {
		log.Fatal(err)
	}

	for _, imagePath := range imageFiles {
		if err := copyFile(imagePath, filepath.Join(tempFolder, filepath.Base(imagePath))); err != nil {
			log.Fatal(err)
		}
	}

	for {
		for _, imagePath := range imageFiles {
			select {
			case <-c:
				err := setWallpaper(originalWallpaper)
				if err != nil {
					log.Fatal(err)
				}
				os.RemoveAll(tempFolder)
				os.Exit(0)
			default:
				if err := setWallpaper(filepath.Join(tempFolder, filepath.Base(imagePath))); err != nil {
					log.Fatal(err)
				}

				time.Sleep(250 * time.Millisecond)
			}
		}
	}
}
