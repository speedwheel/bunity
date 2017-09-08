/*
* Name:        Oscar
* Written by: 	Unifato - (http://unifato.com)
* Version:     1.0.0
*/

(function($) {
  'use strict';

  var Custom = {
    init: function() {
      this.warningModal();
      this.cancelModal();
      this.warningModal();
      this.infoModal();
      this.imageModal();

      //home
      this.dialKnob();
      this.extraAreaMorrisChart();
      this.productLineHomeMorrisChart();
      this.siteVisitMorrisChart();
      this.barMorrisDashboardChart();
      this.homeSparkline1();
      this.homeSparkline2();

      //flot
      this.lineFlotChart();
      this.movingLineFlotChart();
      this.barFlotChart();
      this.salesBarFlotChart();

      //morris
      this.lineMorrisChart();
      this.donuteMorrisChart();
      this.productLineMorrisChart();
      this.barMorrisChart();

      //knob
      this.knobChart();

      //Sparkline
      this.sparklineChart();

      //Chart Js
      this.lineChartJs();
      this.barChartJs();
      this.pieChartJs();
      this.doughnutChartJs();
      this.polarAreaChartJs();
      this.radarChartJs();

      //Widgets
      this.doughnutChartWidget();
      this.barMorrisChartWidget();
      this.lineChartWidget();

      //Demos
      this.barMorrisChartDemo();
      this.barMorrisChartUniversityDemo();
      this.donuteMorrisChartDemo();
      this.doughnutChartJsDemo();

      //Toastr
      this.enableToastr();

      //Animations
      this.testAnimation();
    },


    warningModal: function() {
      var el = $('.modal-alert-warning');
      if ( !el.length ) return;
      el.on('click', function () {
        swal({
          title: 'Are you sure?',
          text: "You won't be able to revert this!",
          type: 'warning',
          showCancelButton: true,
          confirmButtonClass: 'btn btn-warning',
          confirmButtonText: 'Yes, delete it!'
        }).then(function () {
          swal({
            title: 'Deleted!',
            text: 'Your file has been deleted.',
            type: 'success',
            confirmButtonClass: 'btn btn-success'
          });
        });
      });
    },

    cancelModal: function() {
      var el = $('.modal-alert-cancel');
      if ( !el.length ) return;
      el.on('click', function () {
        swal({
          title: 'Are you sure?',
          text: "You won't be able to recover this imaginary file!",
          type: 'error',
          showCancelButton: true,
          confirmButtonClass: 'btn btn-danger',
          confirmButtonText: 'Yes, delete it!'
        }).then(function () {
          swal({
            title: 'Deleted!',
            text: 'Your imaginary file has been deleted.',
            type: 'success',
            confirmButtonClass: 'btn btn-success'
          });
        });
      });
    },

    infoModal: function() {
      var el = $('.modal-alert-info');
      if ( !el.length ) return;
      el.on('click', function () {
        swal({
          title: 'Are you sure?',
          text: "You won't be able to recover this imaginary file!",
          type: 'info',
          showCancelButton: true,
          confirmButtonText: 'Info'
        });
      });
    },

    imageModal: function () {
      var el = $('.modal-alert-image');
      if ( !el.length ) return;
      el.on('click', function () {
        swal({
          title: 'Sweet!',
          text: "Here's a custom image.",
          confirmButtonClass: 'btn btn-success',
          imageUrl: 'https://unsplash.it/400/200',
          imageWidth: 400,
          imageHeight: 200
        });
      });
    },


    /*************** Flot Charts ****************/
    lineFlotChart: function(){
      var ctx = document.getElementById("lineFlot");
      if ( ctx === null ) return;

      var offset = 0;
      plot();

      function plot() {
        var sin = [],
        cos = [];
        for (var i = 0; i < 7; i += 0.2) {
          sin.push([i, Math.sin(i + offset)]);
          cos.push([i, Math.cos(i + offset)]);
        }

        var options = {
          series: {
            lines: {
              show: true
            },
            points: {
              show: true
            }
          },
          grid: {
            hoverable: true, //IMPORTANT! this is needed for tooltip to work
            color: "#AFAFAF",
            borderWidth: 0,
            backgroundColor: '#FFF'
          },
          yaxis: {
            min: -1.2,
            max: 1.2
          },
          colors: ["#007bb6", "#f0ad4e"],
          tooltip: true,
          tooltipOpts: {
            content: "'%s' of %x.1 is %y.4",
            shifts: {
              x: -60,
              y: 25
            }
          }
        };

        var plotObj = function() {
          $.plot($("#lineFlot"), [{
            data: sin,
            label: "&nbsp; sin(x)",
          }, {
            data: cos,
            label: "&nbsp; cos(x)"
          }],
          options);
        }

        plotObj();
        $(document).on('SIDEBAR_CHANGED_WIDTH', plotObj);
        $(window).on('resize', plotObj);
      }
    },

    movingLineFlotChart: function() {
      var ctx = document.getElementById("movingLineFlot");
      if ( ctx === null ) return;

      var container = $("#movingLineFlot");
      // Determine how many data points to keep based on the placeholder's initial size;
      // this gives us a nice high-res plot while avoiding more than one point per pixel.
      var maximum = container.outerWidth() / 2 || 300;
      var data = [];
      function getRandomData() {
        if (data.length) {
          data = data.slice(1);
        }
        while (data.length < maximum) {
          var previous = data.length ? data[data.length - 1] : 50;
          var y = previous + Math.random() * 10 - 5;
          data.push(y < 0 ? 0 : y > 100 ? 100 : y);
        }
        // zip the generated y values with the x values
        var res = [];
        for (var i = 0; i < data.length; ++i) {
          res.push([i, data[i]]);
        }
        return res;
      }
      var series = [{
        data: getRandomData(),
        lines: {
          fill: true
        }
      }];

      var plot = null;

      var FlotFunc = function() {
        plot = $.plot(container, series, {
          colors: ["#e6614f"],
          grid: {
            borderWidth: 0,
            minBorderMargin: 20,
            labelMargin: 10,
            backgroundColor: {
              colors: ["#fff", "#fff"]
            },
            margin: {
              top: 8,
              bottom: 20,
              left: 0
            },

            markings: function(axes) {
              var markings = [];
              var xaxis = axes.xaxis;
              for (var x = Math.floor(xaxis.min); x < xaxis.max; x += xaxis.tickSize * 1) {
                markings.push({
                  xaxis: {
                    from: x,
                    to: x + xaxis.tickSize
                  },
                  color: "#fff"
                });
              }
              return markings;
            }
          },
          xaxis: {
            show: false,
            tickFormatter: function() {
              return "";
            }
          },
          yaxis: {
            min: 0,
            max: 110
          },
          legend: {
            show: true
          }
        });

      }

      FlotFunc();

      // Update the random dataset at 25FPS for a smoothly-animating chart
      setInterval(function updateRandom() {
        series[0].data = getRandomData();
        plot.setData(series);
        plot.draw();
      }, 40);

      $(document).on('SIDEBAR_CHANGED_WIDTH', FlotFunc);
      $(window).on('resize', FlotFunc);

    },

    barFlotChart: function() {
      var ctx = document.getElementById("barFlot");
      if ( ctx === null ) return;

      var barOptions = {
        series: {
          bars: {
            show: true,
            barWidth: 0.5
          }
        },
        xaxis: {
          mode: "time",
          timeformat: "%m/%d",
          minTickSize: [2, "day"]
        },
        grid: {
          hoverable: true,
          color: "#AFAFAF",
          borderWidth: 0,
          backgroundColor: '#FFF'
        },
        legend: {
          show: false
        },
        tooltip: true,
        tooltipOpts: {
          content: "x: %x, y: %y"
        }
      };
      var barData = {
        label: "bar",
        color: "#51d2b7",
        data: [
          [1, 5000],
          [2, 3000],
          [3, 6000],
          [4, 2000],
          [5, 4000],
          [6, 7000]
        ]
      };
      var FlotFunc = function() {
        $.plot($("#barFlot"), [barData], barOptions);
      }
      FlotFunc();
      $(document).on('SIDEBAR_CHANGED_WIDTH', FlotFunc);
      $(window).on('resize', FlotFunc);
    },

    salesBarFlotChart: function() {
      var ctx = document.getElementById("salesBarFlot");
      if ( ctx === null ) return;

      var d1 = [];
      var i = 0;
      // for (i = 0; i <= 5; i += 1)
      // d1.push([i, parseInt(Math.random() * 60)]);

      var d2 = [];
      for (i = 0; i <= 5; i += 1)
      d2.push([i, parseInt(Math.random() * 60)]);

      var d3 = [];
      for (i = 0; i <= 5; i += 1)
      d3.push([i, parseInt(Math.random() * 40)]);

      var ds = [];

      ds.push({
        label : "Data One",
        data : d2,
        bars : {
          order : 2
        }
      });
      ds.push({
        label : "Data Two",
        data : d3,
        bars : {
          order : 3
        }
      });

      var stack = 0,
      bars = true,
      lines = true,
      steps = true;

      var options = {
        bars : {
          show : true,
          barWidth : 0.4,
          fill : 1
        },
        grid : {
          show : true,
          aboveData : false,
          labelMargin : 5,
          axisMargin : 0,
          borderWidth : 1,
          minBorderMargin : 5,
          clickable : true,
          hoverable : true,
          autoHighlight : false,
          mouseActiveRadius : 20,
          borderColor : '#f5f5f5'
        },
        series : {
          stack : stack
        },
        legend : {
          position : "ne",
          margin : [0, 0],
          noColumns : 0,
          labelBoxBorderColor : null,
          labelFormatter : function(label, series) {
            // just add some space to labes
            return '&nbsp;' + label + '&nbsp;&nbsp;&nbsp;';
          },
          width : 30,
          height : 5
        },
        yaxis : {
          tickColor : '#f5f5f5',
          font : {
            color : '#bdbdbd'
          }
        },
        xaxis : {
          tickColor : '#f5f5f5',
          font : {
            color : '#bdbdbd'
          }
        },
        colors : ["#5867c3", "#00c5dc"],
        tooltip : true, //activate tooltip
        tooltipOpts : {
          content : "%s : %y.0",
          shifts : {
            x : -30,
            y : -50
          },
        }
      };

      var FlotFunc = function() {
        $.plot($("#salesBarFlot"), ds, options);
      }

      FlotFunc();
      $(document).on('SIDEBAR_CHANGED_WIDTH', FlotFunc);
      $(window).on('resize', FlotFunc);
    },


    /*************** Morris charts ****************/
    productLineMorrisChart: function() {
      var ctx = document.getElementById("productLineMorris");
      if ( ctx === null ) return;

      var chart = Morris.Area({
        element: 'productLineMorris',
        data: [{
          period: '2010',
          dell: 70,
          lenovo: 110,
          sony: 80
        }, {
          period: '2011',
          dell: 75,
          lenovo: 95,
          sony: 90
        }, {
          period: '2012',
          dell: 85,
          lenovo: 102,
          sony: 75
        }, {
          period: '2013',
          dell: 75,
          lenovo: 90,
          sony: 104
        }, {
          period: '2014',
          dell: 95,
          lenovo: 105,
          sony: 75
        }, {
          period: '2015',
          dell: 87,
          lenovo: 95,
          sony: 80
        },
        {
          period: '2016',
          dell: 92,
          lenovo: 108,
          sony: 70
        }],
        xkey: 'period',
        ykeys: ['dell', 'lenovo', 'sony'],
        labels: ['Dell', 'Lenovo', 'Sony'],
        ymax: 'auto[110]',
        ymin: 'auto[70]',
        pointSize: 3,
        fillOpacity: 0,
        pointStrokeColors:['#00bfc7', '#fb9678', '#9675ce'],
        behaveLikeLine: true,
        gridLineColor: '#e0e0e0',
        lineWidth: 1,
        hideHover: 'auto',
        lineColors: ['#00bfc7', '#fb9678', '#9675ce'],
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },


    /*************** Site Visit Morris Chart ****************/
    siteVisitMorrisChart: function() {
      var ctx = document.getElementById("siteVisitMorris");
      if ( ctx === null ) return;

      var chart = Morris.Area({
        element: 'siteVisitMorris',
        data: [{
          period: '2010',
          SiteA: 0,
          SiteB: 0,

        }, {
          period: '2011',
          SiteA: 130,
          SiteB: 100,

        }, {
          period: '2012',
          SiteA: 60,
          SiteB: 80,

        }, {
          period: '2013',
          SiteA: 180,
          SiteB: 200,

        }, {
          period: '2014',
          SiteA: 280,
          SiteB: 100,

        }, {
          period: '2015',
          SiteA: 170,
          SiteB: 150,
        },
        {
          period: '2016',
          SiteA: 200,
          SiteB: 80,

        },{
          period: '2017',
          SiteA: 0,
          SiteB: 0,

        }],
        xkey: 'period',
        ykeys: ['SiteA', 'SiteB'],
        labels: ['Site A', 'Site B'],
        pointSize: 0,
        fillOpacity: 1,
        pointStrokeColors:['#5867c3', '#00c5dc'],
        behaveLikeLine: true,
        gridLineColor: '#e0e0e0',
        lineWidth: 0,
        smooth: false,
        hideHover: 'auto',
        lineColors: ['#5867c3', '#00c5dc'],
        resize: true
      });
      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    lineMorrisChart: function() {
      var ctx = document.getElementById("lineMorris");
      if ( ctx === null ) return;

      var line = new Morris.Line({
        element: 'lineMorris',
        resize: true,
        data: [
          {y: '2011 Q1', item1: 5666},
          {y: '2011 Q2', item1: 3778},
          {y: '2011 Q3', item1: 4912},
          {y: '2011 Q4', item1: 1767},
          {y: '2012 Q1', item1: 9810},
          {y: '2012 Q2', item1: 6670},
          {y: '2012 Q3', item1: 9820},
          {y: '2012 Q4', item1: 8073},
          {y: '2013 Q1', item1: 4687},
          {y: '2013 Q2', item1: 9432}
        ],
        xkey: 'y',
        ykeys: ['item1'],
        labels: ['Item 1'],
        gridLineColor: '#eef0f2',
        lineColors: ['#1d7c2a'],
        lineWidth: 1,
        hideHover: 'auto'
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', line.resizeHandler);
    },

    donuteMorrisChart: function() {
      var ctx = document.getElementById("donuteMorris");
      if ( ctx === null ) return;

      var chart = Morris.Donut({
        element: 'donuteMorris',
        data: [{
          label: "Dell",
          value: 12,

        }, {
          label: "Samsung",
          value: 30
        }, {
          label: "Apple",
          value: 20
        },{
          label: "Lenovo",
          value: 24,

        },{
          label: "Sony",
          value: 37,
        }],
        resize: true,
        colors:['#99d683', '#13dafe', '#6164c1', '#ffb600', '#ff6261']
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    barMorrisChart: function() {
      var ctx = document.getElementById("barMorris");
      if ( ctx === null ) return;

      var chart = Morris.Bar({
        element: 'barMorris',
        data: [{
          y: '2006',
          a: 100,
          b: 90
        }, {
          y: '2007',
          a: 75,
          b: 65
        }, {
          y: '2008',
          a: 50,
          b: 40
        }, {
          y: '2009',
          a: 75,
          b: 65
        }, {
          y: '2010',
          a: 50,
          b: 40
        }, {
          y: '2011',
          a: 75,
          b: 65
        }, {
          y: '2012',
          a: 100,
          b: 90
        }],
        xkey: 'y',
        ykeys: ['a', 'b'],
        labels: ['A', 'B'],
        barColors:['#00a9ff', '#5867c3'],
        barOpacity: 0.5,
        barSizeRatio: 0.5,
        hideHover: 'auto',
        gridLineColor: '#eef0f2',
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },


    /*************** Homepage ****************/
    productLineHomeMorrisChart: function() {
      var ctx = document.getElementById("productLineHomeMorris");
      if ( ctx === null ) return;

      var chart = Morris.Area({
        element: 'productLineHomeMorris',
        data: [{
          period: '2010',
          iphone: 70,
          ipad: 110,
          itouch: 80
        }, {
          period: '2011',
          iphone: 75,
          ipad: 95,
          itouch: 90
        }, {
          period: '2012',
          iphone: 85,
          ipad: 102,
          itouch: 75
        }, {
          period: '2013',
          iphone: 75,
          ipad: 90,
          itouch: 104
        }, {
          period: '2014',
          iphone: 95,
          ipad: 105,
          itouch: 75
        }, {
          period: '2015',
          iphone: 87,
          ipad: 95,
          itouch: 80
        },
        {
          period: '2016',
          iphone: 92,
          ipad: 108,
          itouch: 70
        }],
        xkey: 'period',
        ykeys: ['iphone', 'ipad', 'itouch'],
        labels: ['iPhone', 'iPad', 'iPod Touch'],
        ymax: 'auto[110]',
        ymin: 'auto[70]',
        pointSize: 3,
        fillOpacity: 0,
        pointStrokeColors:['#00bfc7', '#fb9678', '#9675ce'],
        behaveLikeLine: true,
        gridLineColor: '#e0e0e0',
        lineWidth: 1,
        hideHover: 'auto',
        lineColors: ['#00bfc7', '#fb9678', '#9675ce'],
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    barMorrisDashboardChart: function() {
      var ctx = document.getElementById("barMorrisDashboard");
      if ( ctx === null ) return;

      var chart = Morris.Bar({
        element: 'barMorrisDashboard',
        data: [{
          y: '2006',
          a: 100,
          b: 90,
          c: 60
        }, {
          y: '2007',
          a: 75,
          b: 65,
          c: 40
        }, {
          y: '2008',
          a: 50,
          b: 40,
          c: 30
        }, {
          y: '2009',
          a: 75,
          b: 65,
          c: 40
        },{
          y: '2010',
          a: 45,
          b: 95,
          c: 30
        }],
        xkey: 'y',
        ykeys: ['a', 'b', 'c'],
        labels: ['Twitter', 'Facebook', 'Linkedin'],
        barColors:['#b8edf0', '#b4c1d7', '#fcc9ba'],
        barSizeRatio:0.55,
        hideHover: 'auto',
        gridLineColor: '#eef0f2',
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    extraAreaMorrisChart: function() {
      var ctx = document.getElementById("extraAreaMorris");
      if ( ctx === null ) return;

      function generateData() {
        var ele = null;
        var data = [];
        for( var i=0; i<10; i++) {
          ele = {
            period: (2010 + i) + '',
            lenovo: Math.ceil( Math.random() * 100 ),
            dell: Math.ceil( Math.random() * 100 ),
            hp: Math.ceil( Math.random() * 100 ),
          }
          data.push( ele );
        }
        return data;
      }

      var chart = Morris.Area({
        element: 'extraAreaMorris',
        data: generateData(),
        lineColors: ['#6DC9DE','#49C1D0', '#9BD992'],
        xkey: 'period',
        ykeys: ['lenovo','dell','hp'],
        labels: ['Lenovo','Dell','HP'],
        lineWidth: 0,
        fillOpacity: 1.0,
        pointSize: 0,
        pointFillColors: ['#fff','#fff','#fff'],
        pointStrokeColors: ["#cfedcb","#86d6e0","#abe0ec"],
        behaveLikeLine: false,
        smooth: true,
        gridLineColor: '#e0e0e0',
        resize: true,
        hideHover: 'auto'
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    dialKnob: function() {
      var el = $(".knobHomePage");
      if ( !el.length ) return;
      el.knob({
        'format' : function (value) {
          return value + '%';
        }
      });
    },

    homeSparkline1: function() {
      var el = $("#homeSparkline1");
      if ( !el.length ) return;

      var sparklineFunc = function() {
        el.sparkline([5,6,2,8,9,4,7,10,11,12,10,4,7,10,6,4,6,9,2], {
          type: 'bar',
          height: '45',
          barWidth: 4,
          barSpacing: 2,
          barColor: '#fb9678'
        });
      }

      sparklineFunc();
      $(window).on('resize', sparklineFunc);
      $(document).on('SIDEBAR_CHANGED_WIDTH', sparklineFunc);
    },

    homeSparkline2: function() {
      var el = $("#homeSparkline2");
      if ( !el.length ) return;

      var sparklineFunc = function() {
        el.sparkline([10,7,4,10,12,11,10,7,4,9,8,2,6,5,9,4,7,10], {
          type: 'bar',
          height: '45',
          barWidth: 4,
          barSpacing: 2,
          barColor: '#03a9f3'
        });
      };

      sparklineFunc();
      $(window).on('resize', sparklineFunc);
      $(document).on('SIDEBAR_CHANGED_WIDTH', sparklineFunc);
    },


    /************** Knob Charts ****************/
    knobChart: function() {
      var el = $('[data-plugin="knob"]');
      if ( !el.length ) return;
      el.knob();
    },

    /************** Sparkline Charts ****************/
    sparklineChart: function() {
      var el = $(".sparklineChart");
      if ( !el.length ) return;

      $(document).ready(function(){
        var sparklineLogin = function() {

          $("#sparkline1").sparkline([5,6,2,8,9,4,7,10,11,12,10], {
            type: 'bar',
            height: '45',
            barWidth: 7,
            barSpacing: 4,
            barColor: '#99d683'
          });

          $('#sparkline2').sparkline([20, 40, 30], {
            type: 'pie',
            width: '50',
            height: '45',
            resize: true,
            sliceColors: ['#13dafe', '#6164c1', '#f1f2f7']
          });

          $('#sparkline3').sparkline([5, 6, 2, 9, 4, 7, 10, 12], {
            type: 'bar',
            height: '164',
            barWidth: '7',
            resize: true,
            barSpacing: '5',
            barColor: '#f96262'
          });

          $("#sparkline4").sparkline([0, 23, 43, 35, 44, 45, 56, 37, 40, 45, 56, 7, 10], {
            type: 'line',
            width: '120',
            height: '45',
            lineColor: '#fb6d9d',
            fillColor: 'transparent',
            spotColor: '#fb6d9d',
            minSpotColor: undefined,
            maxSpotColor: undefined,
            highlightSpotColor: undefined,
            highlightLineColor: undefined
          });

          $('#sparkline5').sparkline([15, 23, 55, 35, 54, 45, 66, 47, 30], {
            type: 'line',
            width: '100%',
            height: '160',
            chartRangeMax: 50,
            resize: true,
            lineColor: '#13dafe',
            fillColor: 'rgba(19, 218, 254, 0.3)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
          });

          $('#sparkline5').sparkline([0, 13, 10, 14, 15, 10, 18, 20, 0], {
            type: 'line',
            width: '100%',
            height: '160',
            chartRangeMax: 40,
            lineColor: '#6164c1',
            fillColor: 'rgba(97, 100, 193, 0.3)',
            composite: true,
            resize: true,
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
          });

          $('#sparkline6').sparkline([5, 6, 2, 8, 9, 4, 7, 10, 11, 12, 10], {
            type: 'bar',
            height: '45',
            barWidth: '7',
            barSpacing: '4',
            barColor: '#13dafe'
          });

          $("#sparkline7").sparkline([0,2,8,6,8,5,6,4,8,6,4,2 ], {
            type: 'line',
            width: '100%',
            height: '50',
            lineColor: '#ffca4a',
            fillColor: '#ffca4a',
            highlightLineColor: 'rgba(0, 0, 0, 0.2)',
            highlightSpotColor: '#4f4f4f'
          });

          $("#sparkline8").sparkline([0,6,8,10,7,4,6,8,4,3,1,0 ], {
            type: 'line',
            width: '100%',
            height: '50',
            lineColor: '#ff7676',
            fillColor: '#ff7676',
            maxSpotColor: '#ff7676',
            highlightLineColor: 'rgba(0, 0, 0, 0.2)',
            highlightSpotColor: '#ff7676'
          });

          $("#sparkline9").sparkline([0,6,8,10,7,4,6,8,4,3,1,0 ], {
            type: 'line',
            width: '100%',
            height: '50',
            lineColor: '#2cabe3',
            fillColor: '#2cabe3',
            minSpotColor:'#2cabe3',
            maxSpotColor: '#2cabe3',
            highlightLineColor: 'rgba(0, 0, 0, 0.2)',
            highlightSpotColor: '#2cabe3'
          });

          $("#sparkline10").sparkline([0,6,8,10,7,4,6,8,4,3,1,0 ], {
            type: 'line',
            width: '100%',
            height: '50',
            lineColor: '#6164c1',
            fillColor: '#6164c1',
            maxSpotColor: '#6164c1',
            highlightLineColor: 'rgba(0, 0, 0, 0.2)',
            highlightSpotColor: '#6164c1'
          });

          $('#sparkline11').sparkline([20, 40, 30, 50, 70], {
            type: 'pie',
            height: '200',
            resize: true,
            sliceColors: ['#38d57a','#9c27b0', '#e6614f', '#03a9f3', '#007bb6' ]
          });

          $("#sparkline12").sparkline([10,4,6,12,5,9,3,17,20,10,15,9,4,11], {
            type: 'bar',
            height: '200',
            barWidth: 10,
            barSpacing: 7,
            barColor: '#007bb6'
          });

          $('#sparkline13').sparkline([10,4,6,12,5,9,3,17,20,10,15,9,4,11], {
            type: 'bar',
            height: '200',
            barWidth: '10',
            resize: true,
            barSpacing: '7',
            barColor: '#51d2b7'
          });

          $('#sparkline13').sparkline([10,4,6,12,5,9,3,17,20,10,15,9,4,11], {
            type: 'line',
            height: '200',
            lineColor: '#51d2b7',
            spotRadius: 0,
            fillColor: 'transparent',
            composite: true,
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)'
          });

          $("#sparkline14").sparkline([0, 23, 43, 35, 44, 45, 56, 37, 40, 45, 56, 7, 10], {
            type: 'line',
            width: '100%',
            height: '200',
            lineColor: '#fff',
            fillColor: 'transparent',
            spotColor: '#fff',
            minSpotColor: undefined,
            maxSpotColor: undefined,
            highlightSpotColor: undefined,
            highlightLineColor: undefined
          });

          $('#sparkline15').sparkline([5, 6, 2, 8, 9, 4, 7, 10, 11, 12, 10, 9, 4, 7], {
            type: 'bar',
            height: '200',
            barWidth: '10',
            barSpacing: '10',
            barColor: '#13dafe'
          });

          $('#sparkline16').sparkline([25, 55, 35, 75, 64, 55, 86, 47, 50], {
            type: 'line',
            width: '100%',
            height: '200',
            chartRangeMax: 50,
            resize: true,
            lineColor: '#e6614f',
            fillColor: 'rgba(255,170,170,0.3)',
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
          });

          $('#sparkline16').sparkline([0, 23, 30, 55, 43, 80, 38, 40, 0], {
            type: 'line',
            width: '100%',
            height: '200',
            chartRangeMax: 40,
            lineColor: '#007f3f',
            fillColor: 'rgba(86, 255, 170, 0.3)',
            composite: true,
            resize: true,
            highlightLineColor: 'rgba(0,0,0,.1)',
            highlightSpotColor: 'rgba(0,0,0,.2)',
          });

          $('#sparklinedash').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12], {
            type: 'line',
            width: '100%',
            height: '50',
            lineWidth: 2,
            spotRadius: 0,
            lineColor: '#03a9f3',
            fillColor: '#03a9f3',
            highlightLineColor: 0
          });

          $('#sparklinedash2').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12], {
            type: 'line',
            width: '100%',
            height: '50',
            lineWidth: 2,
            spotRadius: 0,
            lineColor: '#e6614f',
            fillColor: '#e6614f',
            highlightLineColor: 0
          });

          $('#sparklinedash3').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12], {
            type: 'line',
            width: '100%',
            height: '50',
            lineWidth: 2,
            spotRadius: 0,
            lineColor: '#51d2b7',
            fillColor: '#51d2b7',
            highlightLineColor: 0
          });

          $('#sparklinedash4').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12], {
            type: 'line',
            width: '100%',
            height: '50',
            lineWidth: 2,
            spotRadius: 0,
            lineColor: '#fb9678',
            fillColor: '#fb9678',
            highlightLineColor: 0
          });

          $('#sparklinedashbar').sparkline([ 12, 11,10, 9, 8, 7, 6, 5, 6, 7, 8, 9, 10,11, 12,11, 10, 9, 8, 7, 6, 5, 6, 7, 8,9, 10, 12 ], {
            type: 'line',
            height: '50',
            width: '100%',
            chartRangeMin: 3,
            lineColor: "#03a9f3",
            fillColor: "#03a9f3"
          });

          $('#sparklinedashbar2').sparkline([ 12, 11,10, 9, 8, 7, 6, 5, 6, 7, 8, 9, 10,11, 12,11, 10, 9, 8, 7, 6, 5, 6, 7, 8,9, 10, 12 ], {
            type: 'line',
            height: '50',
            width: '100%',
            chartRangeMin: 3,
            lineColor: "#e6614f",
            fillColor: "#e6614f"
          });

          $('#sparklinedashbar3').sparkline([ 12, 11,10, 9, 8, 7, 6, 5, 6, 7, 8, 9, 10,11, 12,11, 10, 9, 8, 7, 6, 5, 6, 7, 8,9, 10, 12 ], {
            type: 'line',
            height: '50',
            width: '100%',
            chartRangeMin: 3,
            fillColor: "#51d2b7",
            lineColor: "#51d2b7"
          });

          $('#sparklinedashbar4').sparkline([ 12, 11,10, 9, 8, 7, 6, 5, 6, 7, 8, 9, 10,11, 12,11, 10, 9, 8, 7, 6, 5, 6, 7, 8,9, 10, 12 ], {
            type: 'line',
            height: '50',
            width: '100%',
            chartRangeMin: 3,
            lineColor: "#fb9678",
            fillColor: "#fb9678"
          });

          $('#sparklinedashbar5').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12,6, 5, 6, 8, 10, 9 ], {
            type: 'bar',
            height: '50',
            barSpacing: 5,
            chartRangeMin: 1,
            barColor: "#03a9f3"
          });

          $('#sparklinedashbar6').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12,6, 5, 6, 8, 10, 9 ], {
            type: 'bar',
            height: '50',
            barSpacing: 5,
            chartRangeMin: 1,
            barColor: "#e6614f"
          });

          $('#sparklinedashbar7').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12,6, 5, 6, 8, 10, 9], {
            type: 'bar',
            height: '50',
            barSpacing: 5,
            chartRangeMin: 1,
            barColor: "#51d2b7"
          });

          $('#sparklinedashbar8').sparkline([ 8, 6, 4, 4, 3, 5, 6, 8, 10, 9, 10, 8, 9, 8, 10, 9, 12,6, 5, 6, 8, 10, 9], {
            type: 'bar',
            height: '50',
            barSpacing: 5,
            chartRangeMin: 1,
            barColor: "#fb9678"
          });

          $('#sparklineWidget1').sparkline('html', {
            type: 'bar',
            height: '30',
            barWidth: '4',
            resize: true,
            barSpacing: '1',
            barColor: '#ccc',
            zeroAxis: true,
          });

        };
        sparklineLogin();
        $(window).on('resize',sparklineLogin);
        $(document).on('SIDEBAR_CHANGED_WIDTH',sparklineLogin);
      });

    },


    /*************** Chart Js ******************/
    lineChartJs: function() {
      var ctx = document.getElementById("chartJsLine");
      if ( ctx === null ) return;

      var ctx1 = document.getElementById("chartJsLine").getContext("2d");
      var data1 = {
        labels: ["January", "February", "March", "April", "May", "June", "July"],
        datasets: [
          {
            label: "First",
            backgroundColor: "rgba(88, 103, 195, 0.4)",
            borderColor: "rgba(88, 103, 195, 0.4)",
            pointStrokeColor: "#fff",
            pointBorderColor: "#fff",
            pointHoverBackgroundColor: "#fff",
            pointHoverBorderColor: "rgba(88, 103, 195, 0.4)",
            data: [0, 80, 40, 58, 20, 55, 10]
          },
          {
            label: "Second",
            backgroundColor: "rgba(255, 128, 111, 0.4)",
            borderColor: "rgba(255, 128, 111, 0.4)",
            pointStrokeColor: "#fff",
            pointBorderColor: "#fff",
            pointHoverBackgroundColor: "#fff",
            pointHoverBorderColor: "rgba(255, 128, 111, 0.4)",
            data: [28, 48, 60, 19, 86, 27, 90]
          }
        ]
      };

      var chartJsLine = new Chart(ctx1, {
        type: "line",
        data: data1,
        responsive: true,
        options: {
          legend: {
            display: false
          },
          tooltips: {
            mode: 'index',
            intersect: false,
            titleFontColor: "#000",
            titleMarginBottom: 10,
            backgroundColor: "rgba(255,255,255,.9)",
            bodyFontColor: "#000",
            borderColor: "#e9e9e9",
            bodySpacing: 10,
            borderWidth: 3,
            xPadding: 10,
            yPadding: 10,
          },
          scales: {
            xAxes: [{
              gridLines: {
                display:false
              }
            }],
            yAxes: [{
              gridLines: {
                display:false
              }
            }]
          }
        }
      });
    },

    barChartJs: function() {
      var ctx = document.getElementById("chartJsBar");
      if ( ctx === null ) return;

      var ctx2 = document.getElementById("chartJsBar").getContext("2d");
      var data2 = {
        labels: ["January", "February", "March", "April", "May", "June", "July"],
        datasets: [
          {
            label: "My First dataset",
            backgroundColor: "#5867c3",
            strokeColor: "#5867c3",
            data: [16, 60, 30, 50, 26, 70, 40]
          },
          {
            label: "My Second dataset",
            backgroundColor: "#00cedc",
            strokeColor: "#00cedc",
            data: [28, 48, 40, 19, 80, 27, 50]
          }
        ]
      };

      var chartJsBar = new Chart(ctx2, {
        type: "bar",
        data: data2,
        options: {
          legend: {
            display: false
          },
          tooltips: {
            mode: 'index',
            intersect: false,
            titleFontColor: "#000",
            titleMarginBottom: 10,
            backgroundColor: "rgba(255,255,255,.9)",
            bodyFontColor: "#000",
            borderColor: "#e9e9e9",
            bodySpacing: 10,
            borderWidth: 3,
            xPadding: 10,
            yPadding: 10,
          },
          scales: {
            xAxes: [{
              gridLines: {
                display:false
              }
            }],
            yAxes: [{
              gridLines: {
                display:false
              }
            }]
          }
        },
        responsive: true
      });
    },

    pieChartJs: function() {
      var ctx = document.getElementById("chartJsPie");
      if ( ctx === null ) return;

      var ctx3 = document.getElementById("chartJsPie").getContext("2d");
      var data3 = {
        labels: [
          "Site A",
          "Site B",
          "Site C",
          "Site D",
          "Site E",
          "Site F",
          "Site G",
          "Site H"
        ],
        datasets: [
          {
            data: [150,30,30,30,30,30,30, 30],
            backgroundColor: [
              "#5793f3",
              "#dd4d79",
              "#bd3b47",
              "#dd4444",
              "#fd9c35",
              "#fec42c",
              "#d4df5a",
              "#5578c2"
            ],
            hoverBackgroundColor: [
              "#5793f3",
              "#dd4d79",
              "#bd3b47",
              "#dd4444",
              "#fd9c35",
              "#fec42c",
              "#d4df5a",
              "#5578c2"
            ]
          }]
        };
        var chartJsPie = new Chart(ctx3,{
          type: "pie",
          data: data3,
          responsive: true,
          options: {
            legend: {
              display: false
            },
            tooltips: {
              mode: 'index',
              intersect: false,
              titleFontColor: "#000",
              titleMarginBottom: 0,
              backgroundColor: "rgba(255,255,255,.9)",
              bodyFontColor: "#000",
              borderColor: "#e9e9e9",
              bodySpacing: 0,
              borderWidth: 0,
              xPadding: 10,
              yPadding: 10,
            },
          }
        });
      },

    doughnutChartJs: function() {
      var ctx = document.getElementById("chartJsDoughnut");
      if ( ctx === null ) return;

      var ctx4 = document.getElementById("chartJsDoughnut").getContext("2d");
      var data4 = {
        labels: [
          "Green",
          "Blue",
          "Orange",
          "Red",

        ],
        datasets: [
          {
            data: [300,90,100,170],
            backgroundColor: [
              "#05d05b",
              "#00c1dd",
              "#ffb701",
              "#ff6161",

            ],
            hoverBackgroundColor: [
              "#05d05b",
              "#00c1dd",
              "#ffb701",
              "#ff6161",

            ]
          }]
        };
        var chartJsDoughnut = new Chart(ctx4, {
          type: "doughnut",
          data: data4,
          responsive: true,
          options: {
            legend: {
              display: false
            },
            tooltips: {
              mode: 'index',
              intersect: false,
              titleFontColor: "#000",
              titleMarginBottom: 10,
              backgroundColor: "rgba(255,255,255,.9)",
              bodyFontColor: "#000",
              borderColor: "#e9e9e9",
              bodySpacing: 10,
              borderWidth: 3,
              xPadding: 10,
              yPadding: 10,
            },
          }
        });
      },

    polarAreaChartJs: function() {
      var ctx = document.getElementById("chartJsPolarArea");
      if ( ctx === null ) return;

      var ctx5 = document.getElementById("chartJsPolarArea").getContext("2d");
      var data5 =  {
        labels: [
          "Blue",
          "Megna",
          "Orange",
          "Grey"
        ],
        datasets: [
          {
            data: [300,250,200,100],
            backgroundColor: [
              "rgba(138,43,226,.6)",
              "rgba(0,250,154,.6)",
              "rgba(255,165,0,.6)",
              "rgba(169,169,169,.6)"
            ],
            hoverBackgroundColor: [
              "rgba(138,43,226,.6)",
              "rgba(0,250,154,.6)",
              "rgba(255,165,0,.6)",
              "rgba(169,169,169,.6)"
            ]
          }]
        };
        var chartJsPolarArea = new Chart(ctx5, {
          type: "polarArea",
          data: data5,
          options: {
            legend: {
              display: false
            },
            scale: {
              display:true,
            },
            tooltips: {
              mode: 'index',
              intersect: false,
              titleFontColor: "#000",
              titleMarginBottom: 10,
              backgroundColor: "rgba(255,255,255,.9)",
              bodyFontColor: "#000",
              borderColor: "#e9e9e9",
              bodySpacing: 10,
              borderWidth: 3,
              xPadding: 10,
              yPadding: 10,
            },
          },
        });
      },

    radarChartJs: function() {
      var ctx = document.getElementById("chartJsRadar");
      if ( ctx === null ) return;

      var ctx6 = document.getElementById("chartJsRadar").getContext("2d");
      var data6 = {
        labels: ["Eating", "Drinking", "Sleeping", "Designing", "Coding", "Cycling", "Running"],
        datasets: [
          {

            label: "My Second dataset",
            backgroundColor: "rgba(255,118,118,0.8)",
            strokeColor: "rgba(255,118,118,1)",
            pointColor: "rgba(255,118,118,1)",
            pointStrokeColor: "#fff",
            pointHighlightFill: "#fff",
            pointHighlightStroke: "rgba(255,118,118,1)",
            data: [30, 38, 50, 29, 86, 37, 110]

          },
          {
            label: "My First dataset",
            backgroundColor: "rgba(44,171,227,0.8)",
            strokeColor: "rgba(44,171,227,1)",
            pointColor: "rgba(44,171,227,1)",
            pointStrokeColor: "#fff",
            pointHighlightFill: "#fff",
            pointHighlightStroke: "rgba(19,218,254,1)",
            data: [55, 99, 70, 91, 66, 35, 40]
          }
        ]
      };

      var chartJsRadar = new Chart(ctx6, {
        type: "radar",
        data: data6,
        options: {
          legend: {
            display: false
          },
          tooltips: {
            mode: 'index',
            intersect: false,
            titleFontColor: "#000",
            titleMarginBottom: 10,
            backgroundColor: "rgba(255,255,255,.9)",
            bodyFontColor: "#000",
            borderColor: "#e9e9e9",
            bodySpacing: 10,
            borderWidth: 3,
            xPadding: 10,
            yPadding: 10,
            displayColors: true
          },
        },
        responsive: true
      });
    },

    /******************** Widgets *********************/
    doughnutChartWidget: function() {
      var ctx = document.getElementById("doughnutWidget");
      if ( ctx === null ) return;
      var data4 = {
        labels: [
          "Blue",
          "Green",
          "Yellow",
          "Red",

        ],
        datasets: [
          {
            data: [300,50,120,100],
            backgroundColor: [
              "#03a9f3",
              "#38d57a",
              "#ffcc02",
              "#e6614f"
            ],
            hoverBackgroundColor: [
              "#03a9f3",
              "#38d57a",
              "#ffcc02",
              "#e6614f"
            ]
          }]
        };
        var chartJsDoughnut = new Chart(ctx, {
          type: "doughnut",
          data: data4,
          responsive: true,
          options: {
            maintainAspectRatio: false,
            legend: {
              display: false
            },
            cutoutPercentage: 80,
            tooltips: {
              mode: 'index',
              intersect: false,
              titleFontColor: "#000",
              titleMarginBottom: 10,
              backgroundColor: "rgba(255,255,255,.9)",
              bodyFontColor: "#000",
              borderColor: "#e9e9e9",
              bodySpacing: 10,
              borderWidth: 3,
              xPadding: 10,
              yPadding: 10,
            },
          }
        });
      },

    barMorrisChartWidget: function() {
      var ctx = document.getElementById("barMorrisWidget");
      if ( ctx === null ) return;

      var chart = Morris.Bar({
        element: 'barMorrisWidget',
        data: [{
          y: '2006',
          a: 100,
          b: 90,
          c: 60
        }, {
          y: '2007',
          a: 75,
          b: 65,
          c: 40
        }, {
          y: '2008',
          a: 50,
          b: 40,
          c: 30
        }, {
          y: '2009',
          a: 75,
          b: 65,
          c: 40
        },{
          y: '2010',
          a: 45,
          b: 95,
          c: 30
        }],
        xkey: 'y',
        ykeys: ['a', 'b', 'c'],
        labels: ['Twitter', 'Facebook', 'Linkedin'],
        barColors:['#b8edf0', '#b4c1d7', '#fcc9ba'],
        barSizeRatio:0.55,
        hideHover: 'auto',
        grid: false,
        axes: false,
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    lineChartWidget: function() {
        var ctx = document.getElementById("lineWidget");
        if ( ctx === null ) return;
        var data1 = {
          labels: ["January", "February", "March", "April", "May", "June"],
          datasets: [
            {
              label: "First",
              backgroundColor: "#8ee2d0",
              borderColor: "#31bfa1",
              pointStrokeColor: "#000",
              pointRadius: 4,
              pointBorderWidth: 2,
              pointBorderColor: "#fff",
              pointBackgroundColor: "#2caa90",
              pointHoverBackgroundColor: "#2caa90",
              pointHoverBorderColor: "#fff",
              pointHoverRadius: 6,
              pointHoverBorderWidth: 2,
              data: [13, 28, 19, 24, 43, 49]
            },

          ]
        };

        var lineWidget = new Chart(ctx, {
          type: "line",
          data: data1,
          responsive: true,
          options: {
            maintainAspectRatio: false,
            legend: {
              display: false
            },
            tooltips: {
              mode: 'index',
              intersect: false,
              titleFontColor: "#000",
              titleMarginBottom: 10,
              backgroundColor: "rgba(255,255,255,.9)",
              bodyFontColor: "#000",
              borderColor: "#e9e9e9",
              bodySpacing: 10,
              borderWidth: 3,
              xPadding: 10,
              yPadding: 10,
            },
            scales: {
              xAxes: [{
                display: false,
                gridLines: {
                  display:false
                }
              }],
              yAxes: [{
                display: false,
                gridLines: {
                  display:false
                }
              }]
            }
          }
        });
      },


    /********************* Demos ***********************/
    barMorrisChartDemo: function() {
      var ctx = document.getElementById("barMorrisDemo");
      if ( ctx === null ) return;

      var chart = Morris.Bar({
        element: 'barMorrisDemo',
        data: [{
          y: '2006',
          a: 100,
          b: 90
        }, {
          y: '2007',
          a: 75,
          b: 65
        }, {
          y: '2008',
          a: 50,
          b: 40
        }, {
          y: '2009',
          a: 75,
          b: 65
        }, {
          y: '2010',
          a: 50,
          b: 40
        }, {
          y: '2011',
          a: 75,
          b: 65
        }, {
          y: '2012',
          a: 100,
          b: 90
        }],
        xkey: 'y',
        ykeys: ['a', 'b',],
        labels: ['A', 'B'],
        barColors:['#8ae1e9', '#1e83c3'],
        hideHover: 'auto',
        barSizeRatio: .4,
        gridLineColor: '#eef0f2',
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    barMorrisChartUniversityDemo: function() {
      var ctx = document.getElementById("barMorrisUniversityDemo");
      if ( ctx === null ) return;

      var chart = Morris.Line({
        element: 'barMorrisUniversityDemo',
        data: [{
          y: '2006',
          a: 100,
          b: 90,
          c: 60
        }, {
          y: '2007',
          a: 75,
          b: 65,
          c: 40
        }, {
          y: '2008',
          a: 50,
          b: 40,
          c: 30
        }, {
          y: '2009',
          a: 75,
          b: 65,
          c: 40
        }, {
          y: '2010',
          a: 50,
          b: 40,
          c: 30
        }, {
          y: '2011',
          a: 75,
          b: 65,
          c: 40
        }, {
          y: '2012',
          a: 100,
          b: 90,
          c: 40
        }, {
          y: '2013',
          a: 75,
          b: 55,
          c: 40
        }, {
          y: '2014',
          a: 75,
          b: 45,
          c: 30
        }, {
          y: '2015',
          a: 75,
          b: 65,
          c: 40
        }, {
          y: '2016',
          a: 90,
          b: 65,
          c: 40
        }],
        xkey: 'y',
        ykeys: ['a', 'b', 'c'],
        labels: ['A', 'B', 'C'],
        lineColors:['#03A9F3', '#E6614F', '#51D2B7'],
        hideHover: 'auto',
        lineWidth: 2,
        pointSize: 5,
        barSizeRatio: .5,
        gridLineColor: '#eef0f2',
        resize: true
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    donuteMorrisChartDemo: function() {
      var ctx = document.getElementById("donuteMorrisDemo");
      if ( ctx === null ) return;

      var chart = Morris.Donut({
        element: 'donuteMorrisDemo',
        data: [{
          label: "Property For Sale",
          value: 12,

        }, {
          label: "Property For Rent",
          value: 30
        }, {
          label: "Property Sold",
          value: 20
        }],
        resize: true,
        colors:['#8ae1e9', '#1e83c3', '#38d57a']
      });

      $(document).on('SIDEBAR_CHANGED_WIDTH', chart.resizeHandler);
    },

    doughnutChartJsDemo: function() {
      var ctx = document.getElementById("chartJsDoughnutDemo");
      if ( ctx === null ) return;

      var ctx4 = document.getElementById("chartJsDoughnutDemo").getContext("2d");
      var data4 = {
        labels: [
          "Delivered",
          "Order",
          "Pending"
        ],
        datasets: [
          {
            data: [300,120,100],
            backgroundColor: [
              "#5793f3",
              "#38d57a",
              "#fb9678"
            ],
            hoverBackgroundColor: [
              "#5793f3",
              "#38d57a",
              "#fb9678"
            ]
          }]
        };
        var chartJsDoughnut = new Chart(ctx4, {
          type: "doughnut",
          data: data4,
          responsive: false,
          options: {
            legend: {
              display: false
            },
            maintainAspectRatio: false,
            tooltips: {
              mode: 'index',
              intersect: false,
              titleFontColor: "#000",
              titleMarginBottom: 10,
              backgroundColor: "rgba(255,255,255,.9)",
              bodyFontColor: "#000",
              borderColor: "#e9e9e9",
              bodySpacing: 10,
              borderWidth: 3,
              xPadding: 10,
              yPadding: 10,
            },
          }
        });
      },

    /*********************** Animations **************************/
    testAnimation: function() {
      var ctx = document.getElementById("animationSandbox");
      if ( ctx === null ) return;

      function testAnim(x) {
        $('#animationSandbox').removeClass().addClass(x + ' animated').one('webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend', function () {
          $(this).removeClass();
        });
      }

      $(document).ready(function () {
        $('.js--triggerAnimation').click(function (e) {
          e.preventDefault();
          var anim = $('.js--animations').val();
          testAnim(anim);
        });
        $('.js--animations').change(function () {
          var anim = $(this).val();
          testAnim(anim);
        });
      });
    },

    /*********************** Toastr **************************/
    enableToastr: function() {
      var toastInfo = $('.toast-info');
      if( toastInfo.length ) {
        toastInfo.click(function(){
          $.toast({
            heading: 'Welcome to Oscar Admin',
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-right',
            icon: 'info',
            stack: false
          });
        });
      }

      var toastSuccess = $('.toast-success');
      if( toastSuccess.length ) {
        toastSuccess.click(function(){
          $.toast({
            heading: 'Welcome to Oscar Admin',
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-right',
            icon: 'success',
            stack: false
          });
        });
      }


      var toastDanger = $('.toast-danger');
      if( toastDanger.length ) {
        toastDanger.click(function(){
          $.toast({
            heading: 'Welcome to Oscar Admin',
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-right',
            icon: 'error',
            stack: false
          });
        });
      }


      var toastWarning = $('.toast-warning');
      if( toastWarning.length ) {
        toastWarning.click(function(){
          $.toast({
            heading: 'Welcome to Oscar Admin',
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-right',
            icon: 'warning',
            stack: false
          });
        });
      }

      var toastBasic = $('.toast-basic');
      if( toastBasic.length ) {
        toastBasic.click(function(){
          $.toast({
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-right',
            stack: false
          });
        });
      }

      var toastArray = $('.toast-array');
      if( toastArray.length ) {
        toastArray.click(function(){
          $.toast({
            heading: 'Read documentation and check:',
            text:["Connections","Cables & accessories","Display & Touch"],
            position: 'top-right',
            stack: false
          });
        });
      }


      var toastHtml = $('.toast-html');
      if( toastHtml.length ) {
        toastHtml.click(function(){
          $.toast({
            heading: 'Welcome to <i>Oscar Admin</i>',
            text: 'This alert needs your attention. <a href="#">Check this out</a>.',
            position: 'top-right',
            stack: false
          });
        });
      }


      var toastTransitions = $('.toast-transitions');
      if( toastTransitions.length ) {
        toastTransitions.click(function(){
          $.toast({
            heading: 'Slide Transitions',
            text: 'Set the `showHideTransition` property to fade|plain|slide to achieve different transitions',
            position: 'top-right',
            showHideTransition:"slide",
            stack: false
          });
        });
      }

      var toastTopCenter = $('.toast-top-center');
      if( toastTopCenter.length ) {
        toastTopCenter.click(function(){
          $.toast({
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-center',
            stack: false
          });
        });
      }

      var toastBottomRight = $('.toat-bottom-right');
      if( toastBottomRight.length ) {
        toastBottomRight.click(function(){
          $.toast({
            text: 'This alert needs your attention, but it is not super important.',
            position: 'bottom-right',
            stack: false
          });
        });
      }

      var toastBottomCenter = $('.toast-bottom-center');
      if( toastBottomCenter.length ) {
        toastBottomCenter.click(function(){
          $.toast({
            text: 'This alert needs your attention, but it is not super important.',
            position: 'bottom-center',
            stack: false
          });
        });
      }

      var toastSticky = $('.toast-sticky');
      if( toastSticky.length ) {
        toastSticky.click(function(){
          $.toast({
            text: 'This alert needs your attention, but it is not super important.',
            position: 'top-right',
            hideAfter: false,
            stack: false
          });
        });
      }

    },
  };
  $(document).ready( function() {
    Custom.init();
	UserListTable();
	//BusinessesListTable();
	
	
  });
})(jQuery);


function UserListTable() {
	var showKey = "showBusinesses";
	var detailRows = [];
	var $userListTable = $('#userListTable');
	
	//init table
	var dt = $userListTable.DataTable( {
        "processing": true,
		"serverSide": true,
		"responsive": true,
		"autoWidth": false,
		"stateSave": false,
		"aaSorting": [[1, "asc"]],
        "ajax": "/userlist",
		"iDisplayLength": 5,
		"lengthMenu": [ [5, 10, 25, 50, 100, -1], [5, 10, 25, 50, 100, "All"] ],
		"rowId": function(b) {
			  return b.id;
		},
		/*'createdRow': function( row, data, dataIndex ) {
				$(row).attr('data-id', data.id);
		},*/
		/*"dom": "<'row'<'col-sm-3'l><'col-sm-6 text-left'B><'col-sm-3'f>>" +
    "<'row'<'col-sm-12'tr>>" +
    "<'row'<'col-sm-5'i><'col-sm-7'p>>",
		"buttons": [
            {
                text: ' Show Businesses',
				className: "hideShowBtn",
                action: function ( e, dt, node, config ) {
					var showBusinesses;
					
					if (typeof(Storage) !== "undefined") {
						if (sessionStorage.getItem(showKey) == "true") {
							sessionStorage.setItem(showKey, false);
						} else {
							sessionStorage.setItem(showKey, true);
						}
						showBusinesses = sessionStorage.getItem(showKey);
						var hideShow = $(".hideShowBtn");
						dt.rows().every( function () {
							var row = this;
							var rowData = this.data();
							if(rowData.Business.length > 0) {
								if(showBusinesses == "true") {
									row.child.show();
									hideShow.addClass("showBusiness");
									hideShow.find("span").text(" Hide Businesses");
									$("#"+rowData.id).addClass( 'details' );
								} else {
									row.child.hide();
									hideShow.removeClass("showBusiness");
									hideShow.find("span").text(" Show Businesses");
									$("#"+rowData.id).removeClass( 'details' );
								}
							}
						});
					}
                }
            }
        ],*/
		"columns": [
			{
                "class":          "details-control",
                "orderable":      false,
                "data":           null,
                "defaultContent": "",
				"width": "20"
            },
			//{ "name": "id", "data": "id" },
			{ "name": "firstname", "data": "firstname" },
			{ "name": "lastname", "data": "lastname" },
			{ "name": "email", "data": "email" },
			{
                "class":          "actionBtns",
                "orderable":      false,
                "data":           null,
                "defaultContent": '<div class="btn-group btn-group-sm actionBtns">'+
									'<button type="button" class="userEditBtn btn btn-sm btn-default"><i class="fa fa-pencil" aria-hidden="true"></i></button>' +
									'<button type="button" class="userDeleteBtn btn btn-sm btn-default"><i class="fa fa-trash" aria-hidden="true"></i></button>' +
									
									'</div>',
				"width": "80"
            },
		],
		"drawCallback": function () {
			dt.rows().every( function () {
				var rowS = this;
				var rowData = this.data();
				var hideShow = $(".hideShowBtn");
				if(rowData.Business.length < 1) {
					
					$userListTable.find("#"+rowData.id+" .details-control").removeClass("details-control");
				} else {
					//show by default all the businesses expanded
					rowS.child( format( rowData ) );
					if (sessionStorage.getItem(showKey) == "true") {
						$("#"+rowData.id).addClass( 'details' );
						hideShow.addClass("showBusiness");
						rowS.child.show();
					} else {
						rowS.child.hide();
					}
					
					
				}
			});
		}		
        //"deferLoading": 57
    });
	

	
	//expand business button
	$userListTable.on( 'click', 'tr td.details-control', function () {
        var tr = $(this).closest('tr');
        var row = dt.row( tr );
		
        var idx = $.inArray( tr.attr('id'), detailRows );
 
        if ( row.child.isShown() ) {
            tr.removeClass( 'details' );
            row.child.hide();
 
            // Remove from the 'open' array
            detailRows.splice( idx, 1 );
        }
        else {
            tr.addClass( 'details' );
            row.child.show();
 
            // Add to the 'open' array
            if ( idx === -1 ) {
                detailRows.push( tr.attr('id') );
            }
        }
    });
	
	//trigger click when chainging view
	dt.on( 'draw', function () {
        $.each( detailRows, function ( i, id ) {
            $('#'+id+' td.details-control').trigger( 'click' );
        } );
    } );
	
	
	var bizID;
	var userID;
	
	//deleteBusiness
	$("#userListTable").on("click", ".businessDeleteBtn", function() {
		userID = $(this).parents('tr').last().prev().attr("id");
		bizID = $(this).closest('tr').data("id");
		$('.popupContainer').html(confirmDeleteModal);
		var bName = $(this).closest('tr').find(".tdBusinessName").text();
		$(".deleteBusinessName strong").text(bName);
		
		$('.confirmDeleteModal').modal();
	});
	
	$(".popupContainer").on("click", ".confirmBDelete", function() {
		$.ajax({
			type:"DELETE",
			url:"/business/"+bizID, //+"/"+bizID,
			success: function(data) {
				if (data.success === true) {
					dt.ajax.reload();
					$('.confirmDeleteModal').modal("hide");
				}
			},
			dataType: 'json',
		  });
	});
	
	
	//hide modal event
	$(".popupContainer").on('hidden.bs.modal', '.modal', function (e) {
		$(this).remove();
	})
	
	//business edit action
	
	$("#userListTable").on("click", ".businessEditBtn", function() {
		//$(".businessEditModal").remove();
		$('.popupContainer').html(businessEditModal);
		
		bizID = $(this).closest('tr').data("id");
		userID = $(this).parents('tr').last().prev().attr("id");
		$.ajax({
			type:"GET",
			url:"/business/"+bizID,
			success: function(data) {
				var b = data.business;
				$("#businessEditModalLabel span").text(b.name);
				$(".businessName").val(b.name);
				$(".businessPhone").val(b.phone);
				$(".businessCountry").val(b.country);
				$(".businessState").val(b.state);
				$(".businessArea").val(b.area);
				$(".businessCity").val(b.city);
				$(".businessAddress").val(b.address);
				$('.businessEditModal').modal();
			},
			dataType: 'json',
		  });
	})
	
	//business update event
	$(".popupContainer").on("click", ".businessUpdate", function() {
		var formSelector = $(".editBusinessForm");
		var b = {};
		$.each($(formSelector).serializeArray(), function(i, field) {
			b[field.name] = field.value;
		});
		b.userid = userID;
		$.ajax({
		type:"POST",
		url:"/business/"+bizID,
			data: {b},
			success: function(data) {
				dt.ajax.reload();
			},
			dataType: 'json',
		  });
	});
	
	$(".hideShowBtn").on("click", function() {
	var showBusinesses;
		if (typeof(Storage) !== "undefined") {
			if (sessionStorage.getItem(showKey) == "true") {
				sessionStorage.setItem(showKey, false);
			} else {
				sessionStorage.setItem(showKey, true);
			}
			showBusinesses = sessionStorage.getItem(showKey);
			var hideShow = $(".hideShowBtn");
			dt.rows().every( function () {
				var row = this;
				var rowData = this.data();
				if(rowData.Business.length > 0) {
					if(showBusinesses == "true") {
						row.child.show();
						hideShow.addClass("showBusiness");
						$("#"+rowData.id).addClass( 'details' );
					} else {
						row.child.hide();
						hideShow.removeClass("showBusiness");

						$("#"+rowData.id).removeClass( 'details' );
					}
				}
			});
		}
	});
}


//business 
var businessEditModal = '<div class="modal fade businessEditModal" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">'+
	 '<div class="modal-dialog modal-lg">'+
		'<div class="modal-content">'+
			'<div class="modal-header">'+
				'<h5 class="modal-title" id="businessEditModalLabel">Edit Business (<span></span>)</h5>'+
			'</div>'+
			'<div class="modal-body">'+
				'<form class="editBusinessForm">' +
					'<div class="tabs">'+
						'<ul class="nav nav-tabs">'+
							'<li class="nav-item" aria-expanded="false"><a class="nav-link active" href="#home-tab" data-toggle="tab" aria-expanded="false">Name & Address</a></li>'+
							'<li class="nav-item"><a class="nav-link" href="#profile-tab" data-toggle="tab" aria-expanded="true">Details</a></li>'+
							'<li class="nav-item"><a class="nav-link" href="#messages-tab" data-toggle="tab" aria-expanded="false">Description</a></li>'+
							'<li class="nav-item"><a class="nav-link" href="#settings-tab" data-toggle="tab" aria-expanded="false">Social Accounts</a>'+
							'<li class="nav-item"><a class="nav-link" href="#images-tab" data-toggle="tab" aria-expanded="false">Images</a></li>'+
						'</ul>'+
						'<div class="tab-content">'+
							'<div class="tab-pane active" id="home-tab" aria-expanded="false">'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">Business Name</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessName" type="text" value="" name="name">'+
									'</div>'+
								'</div>'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">Business Phone</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessPhone" type="text" value="" name="phone">'+
									'</div>'+
								'</div>'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">Country</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessCountry" type="text" value="" name="country">'+
									'</div>'+
								'</div>'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">Area</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessArea" type="text" value="" name="area">'+
									'</div>'+
								'</div>'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">State</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessState" type="text" value="" name="state">'+
									'</div>'+
								'</div>'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">City</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessCity" type="text" value="" name="city">'+
									'</div>'+
								'</div>'+
								'<div class="form-group row">'+
									'<label for="example-text-input" class="col-2 col-form-label">Address</label>'+
										'<div class="col-10">'+
											'<input class="form-control businessAddress" type="text" value="" name="address">'+
									'</div>'+
								'</div>'+
							'</div>'+
							'<div class="tab-pane" id="profile-tab" aria-expanded="true">'+
								
							'</div>'+
							'<div class="tab-pane" id="messages-tab" aria-expanded="false">'+
								
							'</div>'+
							'<div class="tab-pane" id="settings-tab" aria-expanded="false">'+
								
							'</div>'+
							'<div class="tab-pane" id="images-tab" aria-expanded="false">'+
								
							'</div>'+
						'</div>'+
					'</div>'+
					
				'</form>' +
			'</div>'+
			'<div class="modal-footer">'+
				'<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>'+
				'<button type="button" class="btn btn-success businessUpdate">Update</button>'+
			'</div>'+
		'</div>'+
	  '</div>'+
	'</div>';
	
var confirmDeleteModal = '<div class="modal fade confirmDeleteModal" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">'+
	 '<div class="modal-dialog modal-md">'+
		'<div class="modal-content">'+
			'<div class="modal-header">'+
				'<h5 class="modal-title" id="confirmDeleteModalLabel">Confirm Delete <span></span></h5>'+
			'</div>'+
			'<div class="modal-body">'+
				'<p class="deleteBusinessName">You are about to delete a business: <strong></strong>, this procedure is irreversible.</p>'+
				'<p>Do you want to proceed?</p>'+
			'</div>'+
			'<div class="modal-footer">'+
				'<button type="button" class="btn btn-secondary" data-dismiss="modal">Cancel</button>'+
				'<button type="button" class="btn btn-danger confirmBDelete">Delete</button>'+
			'</div>'+
		'</div>'+
	  '</div>'+
	'</div>';
	
function format ( d ) {
	var b = d.Business;
	var html = '<table><tbody>';
	for (var i=0; i< b.length; i++) {
		html += '<tr data-id="'+b[i].id+'"><td class="tdBusinessName">'+b[i].name+'</td><td><div class="btn-group btn-group-sm actionBtns">'+
									'<button type="button" class="businessEditBtn btn btn-sm btn-default"><i class="fa fa-pencil" aria-hidden="true"></i></button>' +
									'<button type="button" class="businessDeleteBtn btn btn-sm btn-default"><i class="fa fa-trash" aria-hidden="true"></i></button>' +		
									'</div></td></tr>';
	}
	html += '</tbody></table>';
    return html;
}

/*function BusinessesListTable() {
	var $businessesListTable = $('#businessesListTable');
	var dt = $businessesListTable.DataTable( {
        "processing": true,
		"serverSide": true,
		"responsive": true,
		//"pagingType": "input",
		"aaSorting": [[0, "asc"]],
        "ajax": "/businesslist",
		"iDisplayLength": 5,
		"lengthMenu": [ [5, 10, 25, 50, 100, -1], [5, 10, 25, 50, 100, "All"] ],
		'createdRow': function( row, data, dataIndex ) {
				$(row).attr('data-id', data.id);
		},
		"columns": [
			{ "name": "name", "data": "name" },
			{ "name": "phone", "data": "phone" },
			{ "name": "country", "data": "country" },
			{ "name": "website", "data": "website" },
			{
                "class":          "actionBtns",
                "orderable":      false,
                "data":           null,
                "defaultContent": '<div class="btn-group btn-group-sm actionBtns">'+
									'<button type="button" class="businessEditBtn btn btn-sm btn-default"><i class="fa fa-pencil" aria-hidden="true"></i></button>' +
									'<button type="button" class="businessDeleteBtn btn btn-sm btn-default"><i class="fa fa-trash" aria-hidden="true"></i></button>' +
									'</div>',
				"width": "80"
            },
		]
    });
	
	var bizID;
	$("#userListTable").on("click", ".businessEditBtn", function() {
		//$(".businessEditModal").remove();
		$('.popupContainer').html(businessEditModal);
		
		bizID = $(this).closest('tr').data("id");
		$.ajax({
			type:"GET",
		url:"/business/"+bizID,
			data: {"bizID":bizID},
			success: function(data) {
				var b = data.business;
				$("#businessEditModalLabel span").text(b.name);
				$(".businessName").val(b.name);
				$(".businessPhone").val(b.phone);
				$(".businessCountry").val(b.country);
				$(".businessState").val(b.state);
				$(".businessArea").val(b.area);
				$(".businessCity").val(b.city);
				$(".businessAddress").val(b.address);
				$('.businessEditModal').modal();
			},
			dataType: 'json',
		  });
	})
	
	$("body").on('hidden.bs.modal', '.businessEditModal', function (e) {
		$(this).remove();
	})
	
	$("body").on("click", ".businessUpdate", function() {
		var formSelector = $(".editBusinessForm");
		var b = {};
		$.each($(formSelector).serializeArray(), function(i, field) {
			b[field.name] = field.value;
		});
		$.ajax({
		type:"POST",
		url:"/business/"+bizID,
			data: {b},
			success: function(data) {
				
			},
			dataType: 'json',
		  });
	});
}
*/