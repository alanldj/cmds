package main

import (
	"encoding/json"
	"fmt"
	"github.com/cheneylew/gotools/tool"
	"math"
	"sort"
	"strings"
	"time"
)

func main() {
	time.Local = time.FixedZone("PDT", -8*3600)
	//lines, _ := tool.ReadLines("/Users/apple/Downloads/golang_pcp_api_backend_0.log.3")
	lines, _ := tool.ReadLines("/Users/apple/Downloads/golang_pcp_api_backend_0.log.5")
	//lines, _ := tool.ReadLines("/Users/apple/Downloads/golang_pcp_api_backend_0.log.7")
	var logs []Log
	var isDetailLogs []Log
	var isDetailSlowLogs []Log
	var parseSkuLogs []Log
	var parseSkuSlowLogs []Log
	for _, line := range lines {
		var log Log
		err := json.Unmarshal([]byte(line), &log)
		if err != nil || log.Duration == "" {
			//fmt.Println(err)
			continue
		}
		log.Time, _ = time.ParseInLocation("2006-01-02T15:04:05.999-07:00", log.Timestamp, time.UTC)
		log.DurationMs = tool.ToFloat64(strings.TrimRight(log.Duration, "ms"))
		if log.DurationMs == 0 {
			fmt.Println(log)
		}
		if strings.Contains(log.Content, "Detail") {
			isDetailLogs = append(isDetailLogs, log)
			if log.Level == "slow" {
				isDetailSlowLogs = append(isDetailSlowLogs, log)
				continue
			}
		} else {
			parseSkuLogs = append(parseSkuLogs, log)
			if log.Level == "slow" {
				parseSkuSlowLogs = append(parseSkuSlowLogs, log)
				continue
			}
		}
		logs = append(logs, log)
	}
	start := logs[0].Time
	end := logs[len(logs)-1].Time
	seconds := end.Sub(start).Seconds()
	qps := float64(len(logs)) / seconds
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].DurationMs > logs[j].DurationMs
	})
	md := make(map[string]int)
	for _, log := range logs {
		key := fmt.Sprintf("%.0f", math.Ceil(log.DurationMs/100))
		if _, ok := md[key]; !ok {
			md[key] = 0
		}
		md[key]++
	}
	var stats []Stat
	for key, value := range md {
		stats = append(stats, Stat{Ms: tool.ToInt(key), Count: value})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Ms > stats[j].Ms
	})
	fmt.Printf("开始时间：%v\n", start.Add(time.Hour*16).Format("2006-01-02 15:04:05"))
	fmt.Printf("结束时间：%v\n", end.Add(time.Hour*16).Format("2006-01-02 15:04:05"))
	fmt.Printf("QPS:%.2f\n", qps)
	fmt.Printf("总共:%d\n", len(logs))
	fmt.Printf("Detail:%d\n", len(isDetailLogs))
	fmt.Printf("DetailSlow:%d\n", len(isDetailSlowLogs))
	fmt.Printf("DetailSlow占比:%.4f\n", float64(len(isDetailSlowLogs))/float64(len(isDetailLogs)))
	fmt.Printf("ParseSku:%d\n", len(parseSkuLogs))
	fmt.Printf("ParseSkuSlow:%d\n", len(parseSkuSlowLogs))
	fmt.Printf("ParseSkuSlow占比:%.4f\n", float64(len(parseSkuSlowLogs))/float64(len(parseSkuLogs)))
	for _, stat := range stats {
		fmt.Printf("%d00ms次数：%d\n", stat.Ms, stat.Count)
	}
	fmt.Println("结束!")
}

type Log struct {
	Timestamp  string    `json:"@timestamp"`
	Time       time.Time `json:"time"`
	Caller     string    `json:"caller"`
	Content    string    `json:"content"`
	Duration   string    `json:"duration"`
	DurationMs float64   `json:"duration_ms"`
	Level      string    `json:"level"`
	Span       string    `json:"span"`
	Trace      string    `json:"trace"`
}

type Stat struct {
	Ms    int
	Count int
}
