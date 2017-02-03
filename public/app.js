window.addEventListener('DOMContentLoaded', function() {
  var ul = document.querySelector('#tweets');
  function update() {
    ul.innerHTML = '';
    window.fetch('/search', {
		method: 'POST',
		body: new FormData(document.querySelector('#form'))
	}).then(function(response) {
      return response.json();
    }).then(function(tweets) {
      for (var tweet of tweets) {
        var a = document.createElement('a')
        a.textContent = tweet.text;
        a.href = 'https://twitter.com/statuses/' + tweet.tweet_id;
        a.target = '_blank';
        var li = document.createElement('li')
        li.appendChild(a);
        ul.appendChild(li);
      }
    }); 
  }

  document.querySelector('#q').addEventListener('keydown', function(e) {
    if (e.which == 13) {
      e.preventDefault();
      update();
    }
  });
  document.querySelector('#b').addEventListener('click', update, false);
}, false);
