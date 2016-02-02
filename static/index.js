ready(function() {
  var index = {};

  // go is called on click on the main Go button
  // too see the audience of a subreddit.
  index.go = function() {
    var el = document.getElementById('subreddit');

    if (!el) {
      return;
    }

    var entry = el.value;

    if (!entry || entry.length == 0) {
      return;
    }

    document.location.href = '/audiences/' + entry;
  };

  index.onKeyup = function(event) {
    if (event.keyCode == 13) {
      index.go();
    }
  };

  // attach to the app
  app.index = index;
});
