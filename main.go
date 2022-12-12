package main

import (
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	limit := "9.6TiB"
	if len(args) == 2 {
		limit = args[1]
	}
	limitBytes := HumanReadableToInt(limit)
	log.Println("limit traffic:", limit, limitBytes)

	output, err := exec.Command("/bin/bash", "-c", "vnstat --oneline").Output()
	if err != nil {
		log.Println(err)
		return
	}

	lines := strings.Split(string(output), ";")
	tx := lines[9]
	usedBytes := HumanReadableToInt(tx)
	log.Println("上传:", tx, "byte:", usedBytes)
	if usedBytes >= limitBytes {
		log.Printf("大于>%s,停止服务: %s\n", limit, tx)
		output, err = exec.Command("/root/stop.sh").Output()
		if err != nil {
			log.Println(">停止服务失败:", err)
			return
		}
		log.Println(">停止服务成功:", string(output))
	}
}

var HumanizeSuffixes = [6]string{"KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

var SizeSuffix = map[string]int{
	"kb": 1024,
	"mb": 1 << (10 * 2),
	"gb": 1 << (10 * 3),
	"tb": 1 << (10 * 4),
}

// HumanReadableToInt Converts a human-readable size to int, param value: A string such as "10MB"
func HumanReadableToInt(value string) int {
	value = strings.ToLower(value)
	value = strings.Replace(value, "i", "", 1)
	l := len(value)
	suffix := value[l-2 : l]
	if multiplier, ok := SizeSuffix[suffix]; ok && l > 2 {
		value = strings.TrimSpace(value[0 : l-2])
		if size, err := strconv.ParseFloat(value, 32); err == nil {
			n := int(math.Floor(float64(multiplier) * size))
			return n
		}
	}
	log.Println("Invalid size value:", value)
	return 0
}
