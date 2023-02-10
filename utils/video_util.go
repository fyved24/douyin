package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func NewFileName(userID uint) string {
	now := time.Now()
	return fmt.Sprintf("%d+%s", userID, now.Format("2006-01-02-15h04m05s.999999"))
}

func CutFirstFrameOfVideo(coverPath, videoPath string) error {
	cuDir, _ := os.Getwd()
	log.Printf("cuDir %v", cuDir)
	cmdArguments := []string{"-i", videoPath, "-vf", "select=eq(n\\,0)",
		"-vframes", "1", coverPath}
	cmd := exec.Command("./utils/ffmpeg.exe", cmdArguments...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		log.Println(errOut.String())
	}
	return err
}
