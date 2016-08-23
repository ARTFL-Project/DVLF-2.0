(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('timeSeries', timeSeries);

    function timeSeries($rootScope, $timeout, $location) {
        var buildTimeSeries = function(scope) {
            var dateList = [];
            var counts = [];

            scope.showTimeSeries = true;
            if (typeof(scope.Results.results.timeSeries) !== 'undefined') {
                for (var i=0; i  < scope.Results.results.timeSeries.length; i +=1) {
                    dateList.push(scope.Results.results.timeSeries[i][0]);
                    counts.push(scope.Results.results.timeSeries[i][1]);
                }
                Chart.defaults.global.responsive = true;
                Chart.defaults.global.animation.duration = 400;
                Chart.defaults.global.tooltipCornerRadius = 0;
                Chart.defaults.global.maintainAspectRatio = false;
                Chart.defaults.bar.scales.xAxes[0].gridLines.display = false;
                var chart = angular.element("#line");
                scope.myBarChart = new Chart(chart, {
                    type: 'line',
                    data: {
                        labels: dateList,
                        datasets: [{
                            borderWidth: 1,
                            borderColor: "rgb(21, 95, 131)",
                            backgroundColor: "rgba(21, 95, 131, .3)",
                            pointBorderWidth: 1,
                            pointRadius: 2,
                            pointHoverBorderWidth: 1,
                            data: counts,
                            lineTension: 0.2
                        }],
                    },
                    options: {
                        lineTension: 1,
                        legend: {
                            display: false,
                        },
                        scales: {
                            yAxes: [{
                                type: "linear",
                                display: true,
                                position: "left",
                                gridLines: {
                                    offsetGridLines: true
                                },
                                ticks: {
                                    beginAtZero: true
                                }
                            }]
                        },
                        tooltips: {
                            cornerRadius: 0,
                        }
                    }
                });
            } else {
                scope.showTimeSeries = false;
            }
        }
        return {
            templateUrl: "app/components/timeSeries/timeSeries.html",
            link: function(scope) {
                scope.$on('resultsUpdate', function () {
                    $timeout(function() {
                        buildTimeSeries(scope);
                    });
                });
            }
        }
    }
})();
