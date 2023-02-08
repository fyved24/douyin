package utils

import (
	"fmt"
	"time"
)

func NewFileName(userID uint) string {
	now := time.Now()
	return fmt.Sprintf("%d+%s", userID, now.Format("2006-01-02-15h04m05s.999999"))
}
