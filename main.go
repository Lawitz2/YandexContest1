package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
программа по расшифровке логов
*/

type flight struct {
	m      int
	status string
}

type flights []flight

type rockets struct {
	f   []flight
	dur int
}

type list struct {
	id  int
	dur int
}

type lists []list

// реализация сортировки
func (a flights) Len() int           { return len(a) }
func (a flights) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a flights) Less(i, j int) bool { return a[i].m < a[j].m }

func (a lists) Len() int           { return len(a) }
func (a lists) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a lists) Less(i, j int) bool { return a[i].id < a[j].id }

func main() {
	logs := make(map[int]rockets)
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("error opening the file: %v", err)
		return
	}

	defer file.Close()

	fo, err := os.Create("output.txt")
	if err != nil {
		fmt.Printf("error creating the file: %v", err)
		return
	}

	writer := bufio.NewWriter(fo)
	scanner := bufio.NewScanner(file)
	scanner.Scan()

	for scanner.Scan() {
		dataslice := strings.Split(scanner.Text(), " ")
		mins, _ := strconv.Atoi(dataslice[0])
		mins2, _ := strconv.Atoi(dataslice[1])

		// переходим от дни/часы/минуты к минутам
		mins = mins*60*24 + mins2*60
		mins2, _ = strconv.Atoi(dataslice[2])
		mins += mins2

		idx, _ := strconv.Atoi(dataslice[3])
		r := logs[idx]

		if r.f == nil { // если это первая запись для ракеты idx - создаем слайс
			r.f = make([]flight, 0, 128) // в слайсе хранятся данные о полётах конкретной ракеты
		}

		var fl flight
		fl.status = dataslice[4]
		fl.m = mins
		r.f = append(r.f, fl)
		logs[idx] = r
	}

	res := make(lists, 0, len(logs))

	// подсчёт общего налёта каждый ракеты
	for i, r := range logs {
		sort.Sort(flights(r.f))
		delta := 0
		for _, d := range r.f {
			if d.status == "B" {
				continue
			}
			if d.status == "A" {
				delta = d.m
			} else {
				r.dur += d.m - delta
			}
		}
		res = append(res, list{id: i, dur: r.dur})
	}

	sort.Sort(res)

	for _, i := range res {
		fmt.Fprintf(writer, "%d ", i.dur)
	}

	writer.Flush()
}
