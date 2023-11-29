/* global Chart */
window.chartColors = {
    red: 'rgb(255, 99, 132)',
    orange: 'rgb(255, 159, 64)',
    yellow: 'rgb(255, 205, 86)',
    green: 'rgb(75, 192, 192)',
    blue: 'rgb(54, 162, 235)',
    purple: 'rgb(153, 102, 255)',
    grey: 'rgb(201, 203, 207)'
};

/**
 * Chart for report
 * @type type
 */
function chartForReport() {
    var chartLabels = $("#chartData").data("labels");
    var chartImpressions = $("#chartData").data("impressions");
    var chartGrossRevenues = $("#chartData").data("gross_revenues");
    var chartNetRevenues = $("#chartData").data("net_revenues");
    var chartProfits = $("#chartData").data("profits");
    var chartGrossEcpm = $("#chartData").data("gross_ecpms");
    var chartNetEcpm = $("#chartData").data("net_ecpms");
    if (!chartLabels) {
        return;
    }
    var chartData = {
        labels: chartLabels,
        datasets: [
            {
                type: 'line',
                label: 'Impression',
                yAxisID: 'C',
                backgroundColor: window.chartColors.purple,
                borderColor: window.chartColors.purple,
                borderWidth: 2,
                fill: false,
                pointRadius: 0,
                data: chartImpressions
            },
            {
                type: 'line',
                label: 'Gross eCPM',
                yAxisID: 'B',
                backgroundColor: window.chartColors.blue,
                borderColor: window.chartColors.blue,
                borderWidth: 2,
                fill: false,
                pointRadius: 0,
                data: chartGrossEcpm,
                hidden: true
            },
            {
                type: 'line',
                label: 'Net eCPM',
                yAxisID: 'B',
                backgroundColor: window.chartColors.red,
                borderColor: window.chartColors.red,
                borderWidth: 2,
                fill: false,
                pointRadius: 0,
                data: chartNetEcpm,
                hidden: true
            },
            {
                type: 'bar',
                label: 'Revenue',
                yAxisID: 'A',
                backgroundColor: window.chartColors.grey,
                borderColor: window.chartColors.grey,
                borderWidth: 2,
                fill: true,
                data: chartGrossRevenues
            },
            {
                type: 'bar',
                label: 'Cost',
                yAxisID: 'A',
                backgroundColor: window.chartColors.orange,
                borderColor: window.chartColors.orange,
                borderWidth: 2,
                fill: true,
                data: chartNetRevenues
            },
            {
                type: 'bar',
                label: 'Profit',
                yAxisID: 'A',
                backgroundColor: window.chartColors.green,
                borderColor: window.chartColors.green,
                borderWidth: 2,
                fill: true,
                data: chartProfits
            }
        ]

    };
    $('#chartBar').remove(); // this is my <canvas> element
    $('#chartBarWrapper').append('<canvas id="chartBar"><canvas>');
    var reportChartCtx = document.getElementById('chartBar').getContext('2d');
    window.myMixedChart = new Chart(reportChartCtx, {
        type: 'bar',
        data: chartData,
        options: {
            responsive: true,
            maintainAspectRatio: true,
            title: {
                display: true,
                text: 'Revenue Overviews'
            },
            tooltips: {
                mode: 'index',
                intersect: true,
                callbacks: {
                    label: function (tooltipItem, data) {
                        if (tooltipItem.datasetIndex === 0) {
                            return tooltipItem.yLabel.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
                        } else {
                            return "$" + formatMoney(tooltipItem.yLabel, ".", ",");
                        }
                    },
                },
            },
            scales: {
                xAxes: [{
                    // Change here
//                        barPercentage: 1.0,
                    categoryPercentage: 0.6
                }],
                yAxes: [
                    {
                        id: 'A',
                        type: 'linear',
                        position: 'right',
                        ticks: {
                            display: true,
                            beginAtZero: true,
//                            stepSize: 1000,
                            userCallback: function (value, index, values) {
                                return "$" + formatMoney(value, ".", ",");
                            },
                        },
                        gridLines: {
                            display: false
                        }
                    },
                    {
                        id: 'B',
                        type: 'linear',
                        position: 'right',
                        ticks: {
                            display: true,
                            beginAtZero: true,
//                            max: 5,
//                            min: 0,
//                            stepSize: 0.25,
                            userCallback: function (value, index, values) {
                                return "$" + formatMoney(value, ".", ",");
                            }
                        },
                        gridLines: {
                            display: false
                        }
                    },
                    {
                        id: 'C',
                        type: 'linear',
                        position: 'left',
                        ticks: {
                            beginAtZero: true,
                            userCallback: function (value, index, values) {
                                return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
                            },
                        },
                        gridLines: {
                            display: true
                        }
                    }
                ]
            }
        }
    });
}

$(document).ready(function () {
    chartForReport();

    /**
     * Sparkline Revenue chart
     * @type Object
     */
    sparklineChart();
    function sparklineChart() {
        var sparklineRevenue = $("#lineChart").data("revenues");
        $('#lineChart').sparkline(sparklineRevenue, {
            type: 'line',
            height: '70',
            width: '100%',
            lineWidth: '3',
            lineColor: '#ffa534',
            fillColor: 'rgba(255, 165, 52, .14)',
            tooltipFormat: '${{y.2}}'
        });
    }

    /**
     * Impression Chart
     * @type type
     */
    var ctx = document.getElementById('statisticsChart').getContext('2d');
    var impLabels = $("#statisticsChart").data("labels");
    var impRequests = $("#statisticsChart").data("requests");
    var impBids = $("#statisticsChart").data("bids");
    var impImpressions = $("#statisticsChart").data("impressions");
    var statisticsChart = new Chart(ctx, {
        type: 'line',
        axisY: {
            valueFormatString: "#,###" //try properties here
        },
        axisX: {
            valueFormatString: "Sample #"
        },
        data: {
            labels: impLabels,
            datasets: [
                {
                    label: "Impressions",
                    borderColor: '#177dff',
                    pointBackgroundColor: 'rgba(23, 125, 255, 0.6)',
                    pointRadius: 0,
                    backgroundColor: 'rgba(23, 125, 255, 0.4)',
                    legendColor: '#177dff',
                    legendHide: "",
                    fill: true,
                    borderWidth: 2,
                    data: impImpressions
                },
                {
                    label: "Bids",
                    borderColor: '#fdaf4b',
                    pointBackgroundColor: 'rgba(253, 175, 75, 0.6)',
                    pointRadius: 0,
                    backgroundColor: 'rgba(253, 175, 75, 0.4)',
                    legendColor: '#fdaf4b',
                    legendHide: "",
                    fill: true,
                    borderWidth: 2,
                    data: impBids
                },
                {
                    label: "Requests",
                    borderColor: '#f3545d',
                    pointBackgroundColor: 'rgba(243, 84, 93, 0.6)',
                    pointRadius: 0,
                    backgroundColor: 'rgba(243, 84, 93, 0.4)',
                    legendColor: '#f3545d',
                    legendHide: "hidden",
                    fill: true,
                    borderWidth: 2,
                    data: impRequests,
                    hidden: true,
                },
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            legend: {
                display: false
            },
            hover: {
                mode: 'index',
                intersect: true
            },
            tooltips: {
                bodySpacing: 4,
                mode: "index",
//                mode: "nearest",
                intersect: 0,
                position: "nearest",
                xPadding: 10,
                yPadding: 10,
                caretPadding: 10,
                callbacks: {
                    label: function (tooltipItem, data) {
                        return tooltipItem.yLabel.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
                    }
                }

            },
            layout: {
                padding: {left: 5, right: 5, top: 15, bottom: 15}
            },
            scales: {
                yAxes: [{
                    ticks: {
                        callback: function (label, index, labels) {
                            return label / 1000000 + 'M';
                        },
                        fontStyle: "500",
                        beginAtZero: false,
                        maxTicksLimit: 5,
                        padding: 10
                    },
                    gridLines: {
                        drawTicks: false,
                        display: false
                    }
                }],
                xAxes: [{
                    gridLines: {
                        zeroLineColor: "transparent"
                    },
                    ticks: {
                        padding: 10,
                        fontStyle: "500"
                    }
                }]
            },
            legendCallback: function (chart) {
                var text = [];
                text.push('<ul class="' + chart.id + '-legend html-legend">');
                for (var i = 0; i < chart.data.datasets.length; i++) {
                    text.push('<li class="' + chart.data.datasets[i].legendHide + '"><span style="background-color:' + chart.data.datasets[i].legendColor + '"></span>');
                    if (chart.data.datasets[i].label) {
                        text.push(chart.data.datasets[i].label);
                    }
                    text.push('</li>');
                }
                text.push('</ul>');
                return text.join('');
            }
        }
    });
    var myLegendContainer = document.getElementById("myChartLegend");
    // generate HTML legend
    myLegendContainer.innerHTML = statisticsChart.generateLegend();
    var legendItems = myLegendContainer.getElementsByTagName('li');
    for (var i = 0; i < legendItems.length; i += 1) {
        legendItems[i].addEventListener("click", legendClickCallback, false);
    }

    function legendClickCallback(event) {
        event = event || window.event;
        var target = event.target || event.srcElement;
        while (target.nodeName !== 'LI') {
            target = target.parentElement;
        }
        var parent = target.parentElement;
        var chartId = parseInt(parent.classList[0].split("-")[0], 10);
        var chart = Chart.instances[chartId];
        var index = Array.prototype.slice.call(parent.children).indexOf(target);
        chart.legend.options.onClick.call(chart, event, chart.legend.legendItems[index]);
        if (chart.isDatasetVisible(index)) {
            target.classList.remove('hidden');
        } else {
            target.classList.add('hidden');
        }
    }

});

function formatMoney(number, decPlaces, decSep, thouSep) {
    decPlaces = isNaN(decPlaces = Math.abs(decPlaces)) ? 2 : decPlaces,
        decSep = typeof decSep === "undefined" ? "." : decSep;
    thouSep = typeof thouSep === "undefined" ? "," : thouSep;
    var sign = number < 0 ? "-" : "";
    var i = String(parseInt(number = Math.abs(Number(number) || 0).toFixed(decPlaces)));
    var j = (j = i.length) > 3 ? j % 3 : 0;

    return sign +
        (j ? i.substr(0, j) + thouSep : "") +
        i.substr(j).replace(/(\decSep{3})(?=\decSep)/g, "$1" + thouSep) +
        (decPlaces ? decSep + Math.abs(number - i).toFixed(decPlaces).slice(2) : "");
}