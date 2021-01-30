package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var BattInfoPath string

type Info int

const (
	CapacityInfo Info = iota
	StatusInfo
)

var BattInfo map[Info]string

type Status int

const (
	StatusCharging Status = iota
	StatusDischarging
)

var BattStatus map[Status]string
var DischargeLow, ChargeHigh int

func init() {
	BattInfoPath = "/sys/class/power_supply/BAT0"

	BattInfo = make(map[Info]string)
	BattInfo[CapacityInfo] = "capacity"
	BattInfo[StatusInfo] = "status"

	BattStatus = make(map[Status]string)
	BattStatus[StatusCharging] = "Charging"
	BattStatus[StatusDischarging] = "Discharging"

	DischargeLow = 20
	ChargeHigh = 80
}

func main() {
	checkBatteryDir(BattInfoPath)

	go periodicCheck()

	for {
	}
}

func checkBatteryDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprintf("%s not exist", path))
	}
}

func periodicCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status := read(BattFile(BattInfoPath, BattInfo[StatusInfo]))
			st := strings.Trim(status, "\n")

			capacity := read(BattFile(BattInfoPath, BattInfo[CapacityInfo]))
			caps, err := strconv.Atoi(strings.Trim(capacity, "\n"))
			if err != nil {
				panic(err)
			}

			if st == BattStatus[StatusDischarging] && caps <= DischargeLow {
				msg := fmt.Sprintf("Low battery %d", caps)
				err = beeep.Notify("Battery Info", msg, "")
				if err != nil {
					panic(err)
				}
			}

			if st == BattStatus[StatusCharging] && caps >= ChargeHigh {
				msg := fmt.Sprintf("Batter full %d", caps)
				err = beeep.Notify("Battery Info", msg, "")
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func BattFile(dir, category string) string {
	return fmt.Sprintf("%s/%s", dir, category)
}

func read(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}
