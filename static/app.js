function ready(fn) {
  if (document.readyState != 'loading'){
    fn();
  } else {
    document.addEventListener('DOMContentLoaded', fn);
  }
}

function onReady() {
  var app = {};

  app.json = function(route, method, body, success, error) {
    success = success || function() {};
    error = error || function() {};

    var request = new XMLHttpRequest();
    request.open(method, route, true);
    request.onload = function() {
        if (request.status >= 200 && request.status < 400) {
          var data = {};
          if (request.responseText &&
              request.responseText.length > 0) {
            var data = JSON.parse(request.responseText);
          }
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
      }

      if (body) {;
        request.send(body);
      } else {
        request.send();
      }
  };

  app.data = function(field) {
    if (!field) {
      return '';
    }

    var el = document.getElementById('app-data');
    console.log(el);
    if (!el) {
      return '';
    }

    var res = el.getAttribute('data-'+field);

    // avoid undefined
    if (!res) {
      return '';
    }

    return res;
  };

  window.app = app;
}

ready(onReady);
