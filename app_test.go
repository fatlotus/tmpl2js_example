package exampletmpl2js_test

import (
	"fmt"
	. "github.com/fatlotus/tmpl2js_example"
	"github.com/tebeka/selenium"
	"net"
	"net/http"
	// "os"
	"testing"
)

func ExpectText(t *testing.T, wd selenium.WebDriver, expected string) {
	body, err := wd.FindElement(selenium.ByCSSSelector, "body")
	if err != nil {
		t.Fatal(err)
	}
	actual, err := body.Text()
	if err != nil {
		t.Fatal(err)
	} else if actual != expected {
		t.Fatalf("Expected:\n%s\nActual:\n%s", expected, actual)
	}
}

func Click(t *testing.T, wd selenium.WebDriver, selector string) {
	elem, err := wd.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		t.Fatal(err)
	}
	if err := elem.Click(); err != nil {
		t.Fatal(err)
	}
}

func TestExampleApp(t *testing.T) {
	// Start a local web server
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(listener)
	}
	defer listener.Close()

	// Start running the tests
	server := &http.Server{Handler: &Counter{}}
	go func() {
		server.Serve(listener)
	}()

	// Start a local Selenium Chrome driver instance
	var port = 1337
	service, err := selenium.NewChromeDriverService("chromedriver", port)
	if err != nil {
		t.Fatal(err)
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps,
		fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Quit()

	// Load the basic page
	if err := wd.Get(fmt.Sprintf("http://%s", listener.Addr())); err != nil {
		t.Fatal(err)
	}

	ExpectText(t, wd, "Hello, world!\nThe button has been clicked 0 times.")
	Click(t, wd, "input")
	ExpectText(t, wd, "Hello, world!\nThe button has been clicked 1 time.")
	Click(t, wd, "input")
	ExpectText(t, wd, "Hello, world!\nThe button has been clicked 2 times.")
}
