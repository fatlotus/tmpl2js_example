package exampletmpl2js

import (
	"html/template"
	"net/http"
	"sync"
	"time"
)

type Counter struct {
	Clicks int

	m sync.Mutex
}

var tmpl = template.Must(template.New("").Parse(`
<h1>Hello, world!</h1>
<p>The button has been clicked {{.Clicks}} time{{if (ne .Clicks 1)}}s{{end}}.</p>
<form method="post"><input type="submit"></form>
`))

func (c *Counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.m.Lock()
	defer c.m.Unlock()

	if r.Method == "POST" {
		c.Clicks += 1
		time.Sleep(1 * time.Second)
	}

	WriteBody(w, r, tmpl, c)
}
