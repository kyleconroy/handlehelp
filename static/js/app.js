var search = document.getElementById("search");
var box = document.getElementById("query");
var results = document.getElementById("results");
var source = null;

search.addEventListener("submit", function(e) {
  e.preventDefault();

  results.innerHTML = '';

  if (source !== null) {
    source.close();
  }

  source = new EventSource('/search?handle=' + encodeURI(query.value));

  source.onerror = function(event) {
    source.close()
  };

  source.onmessage = function(event) {
    var result = JSON.parse(event.data);

    var item = document.createElement('li');

    if (!result.Available) {
      var link = document.createElement('a');
      link.href = result.Profile;
      link.appendChild(document.createTextNode(result.Site.Name));
      item.appendChild(link);
    } else {
      item.appendChild(document.createTextNode(result.Site.Name));
    }

    if (result.Available) {
      item.className = "available";

      if (result.Site.RegisterURL) {
        var signup = document.createElement('a');
        signup.href = result.Site.RegisterURL;
        signup.className = 'signup'
        signup.appendChild(document.createTextNode('Sign Up'));
        item.appendChild(signup);
      }
    }

    results.appendChild(item);
  };

  return false;
});
