ready(function() {
  var subreddit = app.data('subreddit');

  var audiences = {};

  audiences.draw = function(subreddit) {
    app.json('/api/today/' + subreddit, 'GET', this.on_receive_audience, this.on_error);
  };

  audiences.on_error = function(request) {
    alert('Error while retrieving the data for this subreddit.');
  };

  audiences.on_receive_audience = function(xhr, data) {
    var graph_data = [];

    var values = data.audiences;

    for (var i = 0; i < values.length; i++) {
      graph_data.push({
        x: new Date(values[i].crawl_time).getTime(),
        y: values[i].audience,
      });
    }

    var line_data = [
    {
        area: true,
        values: graph_data,
        key: subreddit,
        color: "#ff7f0e",
        strokeWidth: 2,
    }];

    audiences.draw_graph(line_data);
    audiences.update_labels(data);
  };

  audiences.update_labels = function(data) {
    var average = 'data unavailable';
    if (data) {
      if (data.average) {
        average = data.average;
      }
    }

    document.getElementById('average').innerHTML = average;
  };

  audiences.draw_graph = function(data) {
    nv.addGraph(function() {
      var chart = nv.models.lineChart()
                    .options({
                      useInteractiveGuideline: true,
                      transitionDuration: 350,
                      showLegend: true,
                      showYAxis: true,
                      showXAxis: true,
      });

      chart.xAxis
        .tickPadding(15)
        .tickFormat(function(d) {
          return d3.time.format('%x %H:%M')(new Date(d))
      }).showMaxMin(true);

      chart.yAxis
          .axisLabel(subreddit)
          .tickFormat(d3.format('i'));

      d3.select('#chart').append('svg')
          .datum(data)
          .call(chart);

      //Update the chart when window resizes.
      nv.utils.windowResize(function() { chart.update(); });
      return chart;
    });
  };

  // attach to the app.
  app.audiences = audiences;

  // retrieves and display the data
  audiences.draw(subreddit);
});
