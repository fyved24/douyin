package utils

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
	"time"
)

func NewFileName(userID uint) string {
	now := time.Now()
	return fmt.Sprintf("%d+%s", userID, now.Format("2006-01-02-15h04m05s.999999"))
}

func CutFirstFrameOfVideo(videoURL string) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoURL).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 5)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()

	// 结果显示
	if err != nil {
		log.Fatalln("截取图片失败", err)
		return nil
	}
	log.Printf("截取图片成功")
	return buf
}
