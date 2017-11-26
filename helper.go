package exampletmpl2js

import (
	"bytes"
	"encoding/json"
	"github.com/fatlotus/tmpl2js"
	"html/template"
	"net/http"
)

type contentForShim struct {
	Template    template.JS
	Prerendered template.HTML
}

var shim = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<title>Monster clicker</title>
</head>
<body>
{{.Prerendered}}
<script>
	var _template = {{.Template}};
	var rebind = function() {
		document.querySelector("input").onclick = function() {
			this.disabled = true;
			var xhr = new XMLHttpRequest;
			xhr.open("POST", "/", true);
			xhr.setRequestHeader("x-tmpl2js", "true");
			xhr.onreadystatechange = function() {
				if (xhr.readyState == 4) {
					document.body.innerHTML = _template(
						JSON.parse(xhr.responseText));
					rebind();
				}
			}
			xhr.send(null);
		};
	}
	rebind();
</script>
</body>
</html>
`))

func WriteBody(w http.ResponseWriter, r *http.Request, tmpl *template.Template, context interface{}) {
	if r.Header.Get("x-tmpl2js") == "true" {
		err := json.NewEncoder(w).Encode(context)
		if err != nil {
			panic(err)
		}
	} else {
		js, err := tmpl2js.ConvertHTML(tmpl, context, nil)
		if err != nil {
			panic(err)
		}

		buf := bytes.Buffer{}
		if err := tmpl.Execute(&buf, context); err != nil {
			panic(err)
		}
		cnt := &contentForShim{
			Prerendered: template.HTML(buf.Bytes()),
			Template:    template.JS(js),
		}
		if err := shim.Execute(w, cnt); err != nil {
			panic(err)
		}
	}
}
