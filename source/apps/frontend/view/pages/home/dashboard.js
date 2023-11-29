//Muze Double Line Chart JavaScript
var options = {
    series: [{
        name: '2021',
        data: [31, 40, 28, 51, 42, 109, 100]
    }, {
        name: '2020',
        data: [11, 32, 80, 45, 75, 80, 41]
    }],
    chart: {
        type: 'line',
        height: 90,
        zoom: {
            enabled: false
        },
        toolbar: {
            show: false,
        }
    },
    dataLabels: {
        enabled: false
    },
    stroke: {
        width: 3,
        colors: ['#008FFB', '#A8CBFE'],
        curve: 'smooth'
    },
    legend: {
        show: false,
    },
    grid: {
        show: false,
        xaxis: {
            lines: {
                show: false
            }
        },
        row: {
            colors: undefined,
            opacity: 0
        },
    },
    tooltip: {
        enabled: true,
        marker: {
            fillColors: ['#008FFB', '#A8CBFE'],
        },
        x: {
            show: false
        },
    },
    markers: {
        colors: ['#008FFB', '#A8CBFE'],
    },
    yaxis: {
        show: false,
    },
    xaxis: {
        labels: {
            show: false,
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false,
        },
        stroke: {
            width: 0,
        },
        tooltip: {
            enabled: false,
        }
    }
};

var chart = new ApexCharts(document.querySelector("#MuzeDoubleLine"), options);
chart.render();

//Muze Single Line Chart JavaScript
var options = {
    series: [{
        name: '2021',
        data: [31, 50, 38, 51, 60, 109, 100]
    }],
    chart: {
        type: 'line',
        height: 90,
        zoom: {
            enabled: false
        },
        toolbar: {
            show: false,
        }
    },
    dataLabels: {
        enabled: false
    },
    stroke: {
        width: 3,
        colors: ['#008FFB'],
        curve: 'straight'
    },
    legend: {
        show: false,
    },
    grid: {
        show: false,
        xaxis: {
            lines: {
                show: false
            }
        },
        row: {
            colors: undefined,
            opacity: 0
        },
    },
    tooltip: {
        enabled: true,
        x: {
            show: false
        },
    },
    yaxis: {
        show: false,
    },
    xaxis: {
        categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep'],
        labels: {
            show: false,
        },
        stroke: {
            width: 0,
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false,
        },
        tooltip: {
            enabled: false,
        }
    }
};

var chart = new ApexCharts(document.querySelector("#MuzeSingleLine"), options);
chart.render();

//Muze Simple Donut Chart JavaScript
var options = {
    series: [50, 50],
    chart: {
        type: 'donut',
        height: 125,
    },
    dataLabels: {
        enabled: false,
    },
    colors: ['#a8cbfe', '#008ffb'],
    stroke: {
        width: 0
    },
    legend: {
        show: false,
    },
    states: {
        hover: {
            filter: {
                type: 'none',
            }
        },
    },
    plotOptions: {
        donut: {
            size: '65%',
            background: 'transparent',
        },
        pie: {
            offsetX: 25
        }
    },
    grid: {
        show: false,
        xaxis: {
            lines: {
                show: false
            }
        },
        row: {
            colors: undefined,
            opacity: 0
        },
    },
    tooltip: {
        enabled: false,
    },
    yaxis: {
        show: false,
    }
};

var chart = new ApexCharts(document.querySelector("#MuzeSimpleDonut"), options);
chart.render();

//Muze Columns Chart JavaScript
var options = {
    series: [{
        name: 'Net Profit',
        data: [40, 40, 40, 40, 40]
    }, {
        name: 'Free Cash Flow',
        data: [70, 70, 70, 70, 70]
    }],
    chart: {
        type: 'bar',
        height: 90,
        zoom: {
            enabled: false
        },
        toolbar: {
            show: false,
        }
    },
    plotOptions: {
        bar: {
            horizontal: false,
            columnWidth: '50%',
            endingShape: 'rounded'
        },
    },
    dataLabels: {
        enabled: false
    },
    legend: {
        show: false,
    },
    grid: {
        show: false,
        xaxis: {
            lines: {
                show: false
            }
        },
    },
    colors: ['#008FFB', '#A8CBFE'],
    stroke: {
        show: false,
    },
    xaxis: {
        labels: {
            show: false,
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false,
        },
        stroke: {
            width: 0,
        },
        tooltip: {
            enabled: false,
        }
    },
    states: {
        hover: {
            filter: {
                type: 'none',
            }
        },
    },
    yaxis: {
        labels: {
            show: false,
        },
    },
    fill: {
        opacity: 1
    },
    tooltip: {
        enabled: false,
    }
};

var chart = new ApexCharts(document.querySelector("#MuzeColumnsChartTwo"), options);
chart.render();

//Muze Single Line Chart JavaScript
var options = {
    series: [{
        name: '2020',
        data: [20, 40, 28, 51, 42, 70, 75, 28, 51, 42, 70, 75]
    }, {
        name: '2021',
        data: [0, 20, 36, 22, 24, 42, 35, 22, 24, 42, 35, 50]
    }],
    chart: {
        type: 'line',
        height: 380,
        zoom: {
            enabled: false
        },
        toolbar: {
            show: false,
        }
    },
    dataLabels: {
        enabled: false,
    },
    stroke: {
        width: 3,
        colors: ['#008FFB', '#A8CBFE'],
        curve: 'smooth'
    },
    legend: {
        show: true,
        position: 'top',
        horizontalAlign: 'left',
        fontSize: '13px',
        fontFamily: 'Open Sans,sans-serif',
        fontWeight: 400,
        labels: {
            colors: '#6C757D',
        },
        markers: {
            width: 12,
            height: 12,
            strokeWidth: 0,
            strokeColor: '#fff',
            fillColors: ['#0D6EFD', '#A8CBFE'],
            radius: 12,
        },
    },
    grid: {
        show: true,
        borderColor: '#E9ECEF',
        xaxis: {
            lines: {
                show: false
            }
        },
        row: {
            colors: undefined,
            opacity: 0
        },
    },
    tooltip: {
        enabled: true,
        marker: {
            fillColors: ['#008FFB', '#A8CBFE'],
        },
        x: {
            show: false
        },
    },
    markers: {
        colors: ['#008FFB', '#A8CBFE'],
    },
    yaxis: {
        show: true,
        labels: {
            style: {
                colors: '#6C757D',
                fontSize: '13px',
                fontFamily: 'Open Sans,sans-serif',
                fontWeight: 400,
            }
        },
    },
    xaxis: {
        categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
        labels: {
            show: true,
            style: {
                colors: '#6C757D',
                fontSize: '13px',
                fontFamily: 'Open Sans,sans-serif',
                fontWeight: 400,
            }
        },
        axisTicks: {
            show: false,
        },
        axisBorder: {
            show: false,
        },
        stroke: {
            width: 0,
        },
        tooltip: {
            enabled: false,
        }
    }
};

var chart = new ApexCharts(document.querySelector("#MuzeDoubleLineTwo"), options);
chart.render();

//Muze Pie Chart JavaScript
Highcharts.chart('MuzePieChartOne', {
    chart: {
        type: 'pie',
        backgroundColor: null,
    },
    title: {
        text: '',
    },
    credits: {
        enabled: false,
    },
    xAxis: {
        lineColor: 'transparent',
        tickLength: 0,
        labels: {
            enabled: false,
        },
    },
    yAxis: {
        gridLineColor: 'transparent',
        title: {
            text: '',
        },
        labels: {
            enabled: false,
        },
    },
    legend: {
        itemStyle: {
            color: '#6C757D',
            fontSize: '12px',
            fontWeight: '500',
            fontFamily: "'Open Sans', sans-serif",
        },
        margin: 30,
        padding: 0,
        symbolWidth: 11,
        symbolHeight: 11,
        itemDistance: 30,
        symbolPadding: 10,
    },
    plotOptions: {
        pie: {
            size: 230,
            borderWidth: 0,
            allowPointSelect: true,
        },
        series: {
            lineWidth: 0,
        },
        column: {
            pointPadding: 0,
            borderWidth: 0,
            pointWidth: 1,
        },
    },
    accessibility: {
        announceNewData: {
            enabled: true,
        },
        point: {
            valueSuffix: '%',
        }
    },
    tooltip: {
        headerFormat: '<span style="font-size:11px">{series.name}</span><br>',
        pointFormat: '<span style="color:{point.color}">{point.name}</span>: <b>{point.y:.2f}%</b> of total<br/>'
    },
    series: [{
        innerSize: '86%',
        dataLabels: [{
            enabled: false,
        }],
        name: 'Browsers',
        showInLegend: true,
        data: [
            {name: 'Email', y: 20, color: '#E6F0FF',},
            {name: 'Refferal', y: 15, color: '#81B4FE',},
            {name: 'Social', y: 36, color: '#3485FD',}],
    }
    ],
});