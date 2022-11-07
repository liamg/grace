package main

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/aquasecurity/table"
	"github.com/liamg/grace/tracer"
)

func configureSummary(t *tracer.Tracer, w io.Writer, sortKey string) {

	tracker := &tracker{
		counts:    make(map[string]int),
		errors:    make(map[string]int),
		durations: make(map[string]time.Duration),
		starts:    make(map[string]time.Time),
	}

	t.SetSyscallEnterHandler(tracker.recordEnter)
	t.SetSyscallExitHandler(tracker.recordExit)
	t.SetDetachHandler(func(i int) {
		tracker.print(w, sortKey)
	})
}

type tracker struct {
	counts    map[string]int
	errors    map[string]int
	durations map[string]time.Duration
	starts    map[string]time.Time
}

func (t *tracker) recordEnter(s *tracer.Syscall) {
	t.starts[s.Name()] = time.Now()
}

func (t *tracker) recordExit(s *tracer.Syscall) {
	stop := time.Now()
	if start, ok := t.starts[s.Name()]; ok {
		t.durations[s.Name()] += stop.Sub(start)
		delete(t.starts, s.Name())
	}
	if s.Return().Int() < 0 {
		t.errors[s.Name()]++
	}
	t.counts[s.Name()]++
}

func (t *tracker) print(w io.Writer, sortKey string) {

	tab := table.New(w)
	tab.SetRowLines(false)
	tab.AddHeaders("time %", "seconds", "usecs/call", "count", "errors", "syscall")
	tab.SetAlignment(table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignRight, table.AlignLeft)
	tab.SetLineStyle(table.StyleBlue)

	var total time.Duration
	for _, duration := range t.durations {
		total += duration
	}

	type row struct {
		sortKey int
		cols    []string
		name    string
	}

	var rows []row

	for name, count := range t.counts {

		duration := t.durations[name]

		percent := float64(duration) * 100 / float64(total)

		var key int
		switch sortKey {
		case "count":
			key = count
		case "time":
			key = int(percent * 100000)
		case "seconds":
			key = int(duration)
		case "errors":
			key = t.errors[name]

		}

		rows = append(rows, row{
			name:    name,
			sortKey: key,
			cols: []string{
				fmt.Sprintf("%.2f", percent),
				fmt.Sprintf("%.6f", duration.Seconds()),
				fmt.Sprintf("%d", duration.Microseconds()/int64(count)),
				fmt.Sprintf("%d", count),
				fmt.Sprintf("%d", t.errors[name]),
				name,
			},
		})
	}

	if sortKey == "" {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].name < rows[j].name
		})
	} else {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].sortKey > rows[j].sortKey
		})
	}

	for _, row := range rows {
		tab.AddRow(row.cols...)
	}

	tab.Render()
}
