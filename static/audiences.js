ready(function() {
  var subreddit = app.data('subreddit');

  var audiences = {};

  audiences.article_info = {};
  audiences.templates = {};

  audiences.draw = function(subreddit, duration) {
    var t = '';
    if (duration && duration.length > 0) {
      t = '?t=36h';
      if (duration === '7d') {
        t = '?t=7d';
      }
    }
    app.json('/api/today/' + subreddit + t, 'GET', undefined, this.on_receive_audience, this.on_error);
  };

  audiences.on_error = function(request) {
    console.error('Error while retrieving the data for this subreddit.');
  };

  // compile_templates compiles the template
  // used in the audiences view.
  audiences.compile_templates = function() {
    audiences.templates['article'] = _.template(document.getElementById('template_article').innerHTML);
    audiences.templates['article_info'] = _.template(document.getElementById('template_article_info').innerHTML);
  };

  // article_details is called when someone click on
  // an article to display details.
  // It calls the backend to retrieve articles details
  // and then draw the new information if the user
  // has the correct rights.
  audiences.article_details = function(event, div) {
    var id = div.getAttribute("data-id");
    if (!id) {
      return;
    }

    var container = document.getElementById('article-info-' + id);
    var select = document.getElementById('article_select_' + id);

    if (container.innerHTML.trim() !== '') {
      return;
    }

    app.json('/api/article/' + subreddit + '/' + id, 'GET', undefined, this.on_receive_article, this.on_error);
  };

  audiences.close_all_details = function() {
    var nodes = document.querySelectorAll('.article-info');
    for (var i = 0; i < nodes.length; i++) {
      var node = nodes[i];
      node.innerHTML = '';
    }
  };

  // on_receive_article is called when the backend has
  // returned all the article information.
  // If data.demo_mode is set true, no data have been
  // returned because the user isn't logged in.
  audiences.on_receive_article = function(xhr, data) {
    if (!data.id) {
      return;
    }

    // demo mode, display some appealing message
    if (data.demo_mode_message) {
      // close all others
      audiences.close_all_details();
    }

    var rendered = audiences.templates['article_info']({
      demo_mode_message: data.demo_mode_message,
      id: data.id,
    });

    var container = document.getElementById('article-info-' + data.id);
    container.insertAdjacentHTML('beforeend', rendered);
    container.style.display = '';

    // if not in demo mode, renders the graph.
    if (!data.demo_mode_message) {
      // store the data for when we'll change the select
      audiences.article_info = data;
      audiences.render_article_details(data.id, data.ranks, 'Rank');
    }
  };

  audiences.select_article_info_graph = function(select, event) {
    var id = select.getAttribute('data-id');

    if (!id) {
      return;
    }

    switch (+select.selectedIndex) {
      case 0:
        audiences.render_article_details(id, audiences.article_info.ranks, 'Rank');
        break;
      case 1:
        audiences.render_article_details(id, audiences.article_info.scores, 'Score');
        break;
      case 2:
        audiences.render_article_details(id, audiences.article_info.comments, 'Comments');
        break;
    }
  };

  audiences.render_article_details = function(id, data, label) {
    if (!data) {
      return;
    }

    var lines_data = [];

    // article details data
    // ----------------------

    var graph_data = [];

    for (var i = 0; i < data.length; i++) {
      var v = data[i];
      graph_data.push({
        x: new Date(v.time).getTime(),
        y: v.value,
      });
    }

    lines_data.push({
      area: false,
      values: graph_data,
      key: label,
      color: '#ff7f0e',
      strokeWidth: 2,
    });

    // graph rendering.
    // ----------------------

    audiences.draw_graph(lines_data, '#article_chart_' + id, true);
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

    audiences.draw_graph(lines_data, '#chart', false);
    audiences.update_labels(data);

    // articles
    // NOTE(remy): should use d3 to append data
    // ---------------------- 

    d3.select('#articles_container').html(''); // remove previous articles
    var articles_container = document.getElementById('articles_container');

    for (var i = 0; i < data.articles.length; i++) {
      var article = data.articles[i];
      var rendered_article = audiences.templates['article'](article);
      articles_container.insertAdjacentHTML('beforeend', rendered_article);
    }

    // hide removed articles
    if (document.getElementById('hide-removed').checked) {
      audiences.toggle_visibility(false);
    }

    // show demo mode message ?
    // ----------------------

    demo_mode_msg = document.getElementById("demomessage");
    if (demo_mode_msg) {
      if (data.demo_mode_message) {
        // show
        demo_mode_msg.style.display = '';
        // force selection of 36h
        var time = document.getElementById('time');
        if (time) {
          time.selectedIndex = 0;
        }
      } else {
        demo_mode_msg.style.display = 'none';
      }
    }

  };

  audiences.toggle_removed = function(checkbox) {
    if (checkbox.checked) {
      audiences.toggle_visibility(false);
    } else {
      audiences.toggle_visibility(true);
    }
  };

  audiences.on_time_selection_change = function(event) {
    if (!event) {
      return;
    }

    t = undefined;
    if (event.target.selectedIndex == 1) {
      t = '7d';
    }
    // retrieves and display the data
    audiences.draw(subreddit, t);
  };

  audiences.toggle_visibility = function(show) {
    var nodes = document.querySelectorAll('.article-removed');
    for (var i = 0; i < nodes.length; i++) {
      var node = nodes[i];
      if (show) {
        node.style.display = '';
      } else {
        node.style.display = 'none';
      }
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
        highest = '' + data.highest_audience.audience + ' on ' + moment(new Date(data.highest_audience.crawl_time)).format('Y-MM-DD') + ' at ' + moment(new Date(data.highest_audience.crawl_time)).format('HH:mm a');
      }

      if (data.lowest_audience) {
        lowest = '' + data.lowest_audience.audience + ' on ' + moment(new Date(data.lowest_audience.crawl_time)).format('Y-MM-DD') + ' at ' + moment(new Date(data.lowest_audience.crawl_time)).format('HH:mm a');
      }
    }

    document.getElementById('average').innerHTML = average;
    document.getElementById('lowest').innerHTML = lowest;
    document.getElementById('highest').innerHTML = highest;
  };

  audiences.draw_graph = function(data, domNodeSelector, margin) {
    // clear previous content
    d3.select(domNodeSelector).html('');

    var m = {};
    if (margin) {
      m.right = 100;
    }

    // create the graph
    nv.addGraph(function() {
      var chart = nv.models.lineChart()
        .options({
          useInteractiveGuideline: true,
          margin: m,
          transitionDuration: 350,
          showLegend: true,
          showYAxis: true,
          showXAxis: true,
      });

      chart.xAxis
        .tickPadding(3)
        .tickFormat(function(d) {
          return d3.time.format('%x %H:%M')(new Date(d))
      }).showMaxMin(true);

      chart.yAxis
          .axisLabel(subreddit)
          .tickFormat(d3.format('i'));

      d3.select(domNodeSelector).append('svg')
          .datum(data)
          .call(chart);

      // attach the click event for anotation
      // ---------------------- 

      chart.lines.dispatch.on("elementClick", function(e) {
        // NOTE(remy): cut off atm
        return;

        if (!e || e.length == 0) {
          return;
        }

        var point = e[0].point;

        var message = window.prompt('Annotation ?');
        if (message === null || message == '') {
          return;
        }

        var body = {
          t: new Date(point.x),
          m: message
        };

        var d = JSON.stringify(body);
        app.json('/api/annotate/' + subreddit, 'POST', d, function() { app.audiences.draw(subreddit); }, undefined);
      });

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
