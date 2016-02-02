ready(function() {
  var subreddit = app.data('subreddit');

  var audiences = {};

  audiences.draw = function(subreddit) {
    app.json('/api/today/' + subreddit, 'GET', this.on_receive_audience, this.on_error);
  }

  audiences.on_error = function(request) {
    alert('Error while retrieving the data for this subreddit.');
  }

  audiences.on_receive_audience = function(xhr, data) {
    var graph_data = [];

    for (var i = 0; i < data.length; i++) {
      graph_data.push({
        x: new Date(data[i].CrawlTime).getTime(),
        y: data[i].Audience,
      });
    }

    var line_data = [
    {
        area: true,
        values: graph_data,
        key: "Audience",
        color: "#ff7f0e",
        strokeWidth: 4,
    }];

    audiences.draw_graph(line_data);
  }

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
          console.log(d);
          return d3.time.format('%x %H:%M')(new Date(d))
      }).showMaxMin(true);

      chart.yAxis
          .axisLabel(subreddit)
          .tickFormat(d3.format('.02f'));

      d3.select('#chart').append('svg')
          .datum(data)
          .call(chart);

      //Update the chart when window resizes.
      nv.utils.windowResize(function() { chart.update(); });
      return chart;
    });
  }

  // retrieves and display the data
  audiences.draw(subreddit);
});
