import style from './style'

export default `
<!doctype html>
<html ⚡ class="amp-border-box">
<head>
  <meta charset="utf-8">
  <link rel="canonical" href="#">
  <title>Post</title>
  <meta name="description" content="Post entry">
  <meta name="viewport" content="width=device-width,minimum-scale=1">
  ${style}
  <script async src="https://cdn.ampproject.org/v0.js"></script>
  <script async custom-element="amp-form" src="https://cdn.ampproject.org/v0/amp-form-0.1.js"></script>
  <script async custom-element="amp-bind" src="https://cdn.ampproject.org/v0/amp-bind-0.1.js"></script>
</head>
<body>
<div class="container">
  <h1 id="title"></h1>
  <p id="author"></p>
  <p id="date"></p>
  <p id="views"></p>
  <p id="content"></p>
  <a id="mainsite" href="#">view without amp</a>
</div>
</body>
</html>
`
