package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type scanFile struct {
	Zipfile  string `json:"zipfile"`
	Filename string `json:"filename"`
	Rows     []row  `json:"rows"`
}

type row struct {
	Timestamp int      `json:"timestamp"`
	Location  location `json:"location"`
	Scan      scan     `json:"scan"`
}

type location struct {
	Provider           string  `json:"provider"`
	Timestamp          int     `json:"timestamp"`
	Easting            int     `json:"easting"`
	Northing           int     `json:"northing"`
	Utm_zone           string  `json:"utm_zone"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Accuracy           float64 `json:"accuracy"`
	Speed              float64 `json:"speed"`
	Location_timestamp int     `json:"location_timestamp"`
}

type scan struct {
	Timestamp int        `json:"timestamp"`
	Hotspots  []hotspots `json:"hotspots"`
}

type hotspots struct {
	Id           int    `json:"id"`
	Ssid         string `json:"ssid"`
	Bssid        string `json:"bssid"`
	Cap          string `json:"cap"`
	Signal_level int    `json:"signal_level"`
}

func main() {

	file, err := os.Open("file.24285020.json")
	timeStamps := make(map[string]int)
	times := make(map[string]time.Time)

	var rows []row

	if err != nil {
		fmt.Println(err)
	}

	fileReader := bufio.NewReader(file)
	linesRead := 0
	totalRows := 0
	var buf scanFile

	for {
		line, err := fileReader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println(err)
			break
		}

		err = json.Unmarshal([]byte(line), &buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		for _, row := range buf.Rows {
			totalRows++
			if row.Scan.Timestamp != 0 {
				rows = append(rows, row)
			}
		}
		linesRead++
	}

	for _, row := range rows {

		for _, hs := range row.Scan.Hotspots {

			key := hs.Ssid + hs.Bssid

			if val, ok := timeStamps[key]; !ok {
				timeStamps[key] = row.Scan.Timestamp
			} else {
				if row.Scan.Timestamp > val {

					if key == "HP-Print-d0-LaserJet 1003c77e6553fd0" {
						fmt.Printf("Key: %s. Old: %s. New: %s", key, time.Unix(int64(val/1000), 0), time.Unix(int64(row.Scan.Timestamp/1000), 0))
					}
					timeStamps[key] = row.Scan.Timestamp
				}
			}
		}
	}

	for k, v := range timeStamps {
		times[k] = time.Unix(int64(v/1000), 0)
		if k == "HP-Print-d0-LaserJet 1003c77e6553fd0" {
			fmt.Printf("Key: %s. Time: %s", k, times[k])
		}
	}

	//fmt.Println(timeStamps)
	fmt.Println(linesRead)
	fmt.Println(totalRows)
	fmt.Println(len(rows))
}
