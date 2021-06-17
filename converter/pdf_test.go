package converter

import "testing"

func TestHTMLToPDF(t *testing.T) {
	tc :=
		`<!DOCTYPE html>
		<html lang="en">
		  <head>
			<meta charset="UTF-8" />
			<meta http-equiv="X-UA-Compatible" content="IE=edge" />
			<meta name="viewport" content="width=device-width, initial-scale=1.0" />
			<title>Document</title>
		  </head>
		  <body>
			<div>Hello World</div>
		  </body>
		</html>
		`

	res, err := HTMLToPDF(tc)
	if err != nil || len(res) == 0 {
		t.Fatal(err)
	}
}
