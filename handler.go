package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/jalasoft/vodni-svet-kolin-metrics-exporter/vodnisvetkolin"
)

const pattern = `# HELP vodni_svet_kolin_occupancy Actual number of people swimming in vodni svet kolin
# TYPE vodni_svet_kolin_occupancy gauge
vodni_svet_kolin_occupancy {{.Occupancy}}
# HELP vodni_svet_kolin_visited Number of people who visited vodni svet kolin in given period
# TYPE vodni_svet_kolin_visited gauge
vodni_svet_kolin_visited{when="today"} {{.VisitedToday}}
vodni_svet_kolin_visited{when="week"} {{.VisitedThisWeek}}
vodni_svet_kolin_visited{when="month"} {{.VisitedThisMonth}}
vodni_svet_kolin_visited{when="year"} {{.VisitedThisYear}}
`

var (
	metricsTemplate              *template.Template
	VodniSvetKolinMetricsHandler http.Handler
)

func init() {
	metricsTemplate = template.New("metrics-template")
	metricsTemplate.Parse(pattern)

	VodniSvetKolinMetricsHandler = http.HandlerFunc(handler)
}

func handler(writer http.ResponseWriter, request *http.Request) {

	info, err := vodnisvetkolin.ReadStatistics()

	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	if err := metricsTemplate.Execute(writer, info); err != nil {
		fmt.Fprintf(writer, "%s\n", err.Error())
	}
}
