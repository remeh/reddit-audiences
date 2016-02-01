function ready(fn) {
  if (document.readyState != 'loading'){
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

function onReady() {
  var app = {};

  app.json = function(route, method, success, error) {
    var request = new XMLHttpRequest();
    request.open(method, route, true);
    request.onload = function() {
        if (request.status >= 200 && request.status < 400) {
          var data = JSON.parse(request.responseText);
          success(request, data);
        } else {
          if (error) {
            error(request, request.status);
          }
        }
      };
      request.onerror = function() {
        if (error) {
          error(request, request.status);
        }
      };
      request.send();
  }

  window.app = app;
}

ready(onReady);
