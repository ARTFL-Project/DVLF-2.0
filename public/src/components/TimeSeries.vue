<template>
    <div id="time-series" class="shadow-sm" :class="{'hidden': showTimeSeries == false}">
        <b-card style="margin-top: 15px; min-height: 200px">
            <h4 id="time-series-title">Usage à travers le temps</h4>
            <h6 style="padding: 0 10px">Occurrences du mot
                <b>{{ headword }}</b> pour un million de mots</h6>
            <div>
                <canvas id="line" class="chart" height="200"></canvas>
            </div>
        </b-card>
    </div>
</template>

<script>
import Chart from "chart.js/dist/Chart.js"

export default {
    name: "TimeSeries",
    props: {
        headword: String,
        timeSeries: Array
    },
    data() {
        return {
            showTimeSeries: false,
            chart: null
        }
    },
    created() {
        this.drawChart()
    },
    methods: {
        drawChart() {
            if (this.chart != null) {
                this.chart.destroy()
            }

            let dateList = []
            let counts = []
            this.showTimeSeries = true
            if (this.timeSeries.length > 0) {
                for (var i = 0; i < this.timeSeries.length; i += 1) {
                    dateList.push(this.timeSeries[i][0])
                    counts.push(this.timeSeries[i][1])
                }
                Chart.defaults.global.responsive = true
                Chart.defaults.global.animation.duration = 400
                Chart.defaults.global.tooltipCornerRadius = 0
                Chart.defaults.global.maintainAspectRatio = false
                Chart.defaults.bar.scales.xAxes[0].gridLines.display = false
                this.$nextTick(function() {
                    let chart = document.getElementById("line")
                    let vm = this
                    vm.chart = new Chart(chart, {
                        type: "line",
                        data: {
                            labels: dateList,
                            datasets: [
                                {
                                    borderWidth: 1,
                                    borderColor: "rgb(21, 95, 131)",
                                    backgroundColor: "rgba(21, 95, 131, .3)",
                                    pointBorderWidth: 1,
                                    pointRadius: 2,
                                    pointHoverBorderWidth: 1,
                                    data: counts,
                                    lineTension: 0.2
                                }
                            ]
                        },
                        options: {
                            lineTension: 1,
                            legend: {
                                display: false
                            },
                            scales: {
                                yAxes: [
                                    {
                                        type: "linear",
                                        display: true,
                                        position: "left",
                                        gridLines: {
                                            offsetGridLines: true
                                        },
                                        ticks: {
                                            beginAtZero: true
                                        }
                                    }
                                ]
                            },
                            tooltips: {
                                cornerRadius: 0
                            }
                        }
                    })
                })
                // this.showTimeSeries = true
            } else {
                // this.showTimeSeries = false
            }
        }
    }
}
</script>

<style scoped>
#time-series-title {
    background-color: #f0f0f0;
    border-bottom: 1px solid #eee;
    text-align: center;
    margin-top: 0;
    padding: 7px;
    font-weight: 700;
}
</style>
