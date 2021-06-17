package main

const TemplateStr = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link
      href="https://fonts.googleapis.com/css?family=Raleway"
      rel="stylesheet"
      type="text/css"
    />
    <!-- <link rel="stylesheet" href="css/style.css" /> -->
    <title>Spoostys Teilnahmeurkunde</title>
  </head>
  <style>
    body {
      margin: 0;
      padding: 0;
    }

    .container,
    .left-justified-container {
      display: flex;
      height: inherit;
      flex-direction: column;
      justify-content: space-between;
      align-items: center;
    }

    .descriptions {
      text-align: center;
    }

    .border,
    .inner-border,
    .outer-border {
      padding: 20px;
      border: 10px solid #0066ff;
    }

    .outer-border {
      width: 800px;
      height: 1200px;
    }

    .inner-border {
      width: 750px;
      height: 1150px;
      border-width: 5px;
    }

    .default-heading,
    h4,
    h3,
    h2,
    h1 {
      color: #000000;
      font-family: Raleway;
      font-weight: normal;
    }

    h1 {
      font-size: 3em;
      font-weight: bold;
    }

    h2 {
      font-size: 2.5em;
    }

    h3 {
      font-size: 2em;
    }

    h4 {
      font-size: 1.5em;
      font-weight: normal;
    }

    .left-justified-container {
      align-items: flex-start;
      justify-content: flex-start;
      width: 65%;
    }

    img {
      max-width: 200px;
      max-height: 200px;
    }
  </style>
  <body>
    <div class="outer-border">
      <div class="inner-border">
        <div class="container">
          <div class="descriptions">
            <h1>Teilnahmeurkunde</h1>
            <h2>Tri.alone Cup 2021</h2>
            <h3>für {{.Name}}</h3>
            <h3>Kurzdistanz</h3>
            <h3>(10 km run / 40 km bike / 5 km swim)</h3>
          </div>
          <div class="left-justified-container">
            <h4>Run:{{.Time}}</h4>
            <h4>Transition 1:</h4>
            <h4>Bike:</h4>
            <h4>Transition 2:</h4>
            <h4>Swim:</h4>
          </div>
          <!-- <img
            src="https://spoosty.de/wordpress/wp-content/images/spoosty-rgb-200px.png"
            nosend="135"
            alt="spoosty-logo"
          /> -->
        </div>
      </div>
    </div>
  </body>
</html>
`
