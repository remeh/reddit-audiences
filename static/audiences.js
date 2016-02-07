ready(function() {
  var subreddit = app.data('subreddit');

  var audiences = {};

  audiences.templates = {};

  audiences.draw = function(subreddit) {
    app.json('/api/today/' + subreddit, 'GET', this.on_receive_audience, this.on_error);
  };

  audiences.on_error = function(request) {
    alert('Error while retrieving the data for this subreddit.');
  };

  audiences.compile_templates = function() {
    audiences.templates['article'] = _.template(document.getElementById('template_article').innerHTML);
  };

  audiences.on_receive_audience = function(xhr, data) {
    var graph_data = [];
    var ranking_data = [];

    // audiences
    // ----------------------

    var values = data.audiences;
    for (var i = 0; i < values.length; i++) {
      graph_data.push({
        x: new Date(values[i].crawl_time).getTime(),
        y: values[i].audience,
      });
    }
    var lines_data = [
    {
        area: true,
        values: graph_data,
        key: subreddit,
        color: '#ff7f0e',
        strokeWidth: 2,
    }];

    audiences.draw_graph(lines_data, '#chart');
    audiences.update_labels(data);

    // articles
    // ---------------------- 

    var articles_container = document.getElementById('articles_container');

    for (var i = 0; i < data.articles.length; i++) {
      var article = data.articles[i];
      var rendered_article = audiences.templates['article'](article);
      articles_container.insertAdjacentHTML('beforeend', rendered_article);
    }
  };

  audiences.update_labels = function(data) {
    var average = 'data unavailable.';
    var highest = 'data unavailable.';
    var lowest = 'data unavailable.';

    if (data) {
      if (data.average) {
        average = data.average;
      }

      if (data.highest_audience) {
        highest = '' + data.highest_audience.audience + ' at ' + moment(new Date(data.highest_audience.crawl_time)).format('HH:mm a');
      }

      if (data.lowest_audience) {
        lowest = '' + data.lowest_audience.audience + ' at ' + moment(new Date(data.lowest_audience.crawl_time)).format('HH:mm a');
      }
    }

    document.getElementById('average').innerHTML = average;
    document.getElementById('lowest').innerHTML = lowest;
    document.getElementById('highest').innerHTML = highest;
  };

  audiences.draw_graph = function(data, domNodeSelector) {
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

      d3.select(domNodeSelector).append('svg')
          .datum(data)
          .call(chart);

      //Update the chart when window resizes.
      nv.utils.windowResize(function() { chart.update(); });
      return chart;
    });
  };

  // attach to the app.
  app.audiences = audiences;

  // compile the templates from 
  audiences.compile_templates();

  // retrieves and display the data
  audiences.draw(subreddit);
});
