package tool

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"
)

type ChanData struct {
	Success bool
	Message string
	Items []interface{}
}

func Conc(jobs []interface{}, process func(a interface{}))  {
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(job interface{}) {
			process(job)
			wg.Done()
		}(job)
	}
	wg.Wait()
}

//按指定并发执行，如果process遇到错误，则停止后续协程，终止
func ConcQueueV2(maxConc uint8, items interface{}, sliceSize int, process func(items[]interface{}) ([]interface{}, error)) ([]interface{}, error) {
	if reflect.TypeOf(items).Kind() != reflect.Slice {
		return nil, errors.New("items请传入一个slice类型")
	}
	groupChan := make(chan []interface{})
	resultChan := make(chan []interface{})
	var exitError error
	go func(items interface{}) {
		defer close(groupChan) //全部传完后关掉
		//分组传入通道缓冲区，异步不阻塞
		var sliceItems []interface{}
		s := reflect.ValueOf(items)
		for i:=0; i<s.Len(); i++ {
			sliceItems = append(sliceItems, s.Index(i).Interface())
			if len(sliceItems) == sliceSize {
				groupChan<-sliceItems
				sliceItems = []interface{}{}
			}
		}
		if len(sliceItems) > 0 {
			groupChan <- sliceItems
		}
	}(items)
	go func() {
		var wg sync.WaitGroup
		maxConcChan := make(chan struct{}, maxConc) //控制最大并发数量
		defer close(resultChan)
		defer close(maxConcChan)
		var lock sync.Mutex
		for sliceItems := range groupChan {
			maxConcChan<- struct{}{}
			if exitError != nil {
				<-maxConcChan
				continue //让分组数组走完，能够正常关闭通道
			}
			//每组开启一个协程处理任务
			wg.Add(1)
			go func(sliceItems []interface{}) {
				defer wg.Done()
				result, err := process(sliceItems)
				if err != nil {
					lock.Lock()
					exitError = errors.New(fmt.Sprintf("error:%v array:%v", err, sliceItems))
					lock.Unlock()
				} else {
					resultChan <- result
				}
				<-maxConcChan
			}(sliceItems)
		}
		//等到全部完成
		wg.Wait() //全部执行完了
	}()
	var results []interface{}
	for result := range resultChan {
		results = append(results, result...)
	}
	return results, exitError
}

func ConcQueue(jobs []interface{}, process func(jobsSlice []interface{}) []interface{}, maxConc uint8, sliceSize int) []interface{} {
	perPageNum := sliceSize
	pageTotalCount := int(math.Ceil(float64(len(jobs)) / float64(perPageNum)))
	cd := make(chan ChanData)
	maxConcurrencyNum := maxConc
	var curConNum uint8
	curConNum = 0
	curIndex := 0
	finished := 0
	addedCount := 0
	var totalResults []interface{}
	for {
		if curIndex < pageTotalCount && curConNum < maxConcurrencyNum {
			var endIndex int
			if curIndex*perPageNum+perPageNum >= len(jobs) {
				endIndex = len(jobs)
			} else {
				endIndex = curIndex*perPageNum + perPageNum
			}
			targetJobItems := jobs[curIndex*perPageNum : endIndex]

			go func(cd chan ChanData, pageId int, targetJobItems []interface{}) {
				results := process(targetJobItems)
				cd <- ChanData{
					Success: true,
					Message: fmt.Sprintf("%d已完成", pageId),
					Items:   results,
				}
			}(cd, curIndex, targetJobItems)
			curConNum++
			curIndex++
			addedCount++
		}
		//并发满了，等待完成后进入下一个;或者并发未满，但并发池子里还有未完成的任务，等待所有完成才退出。
		if curConNum == maxConcurrencyNum || (addedCount == pageTotalCount && finished != pageTotalCount) {
			next := false
			exit := false
			select {
			case data := <-cd:
				if data.Success && data.Items != nil {
					totalResults = append(totalResults, data.Items...)
				}
				finished++
				curConNum--
				next = true
				if finished == pageTotalCount {
					exit = true
					break
				}
			}
			if exit {
				break
			}
			if next {
				continue
			}
		}
	}

	close(cd)

	return totalResults
}
