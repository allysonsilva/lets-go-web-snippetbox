package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	t.Parallel()

	// // Create a new instance of our application struct. For now, this just
	// // contains a couple of mock loggers (which discard anything written to
	// // them).
	// app := &application{
	// 	errorLog: log.New(ioutil.Discard, "", 0),
	// 	infoLog:  log.New(ioutil.Discard, "", 0),
	// }

	// // We then use the httptest.NewTLSServer() function to create a new test
	// // server, passing in the value returned by our app.routes() method as the
	// // handler for the server. This starts up a HTTPS server which listens on a
	// // randomly-chosen port of your local machine for the duration of the test.
	// // Notice that we defer a call to ts.Close() to shutdown the server when
	// // the test finishes.
	// ts := httptest.NewTLSServer(app.routes())
	// defer ts.Close()

	// // The network address that the test server is listening on is contained
	// // in the ts.URL field. We can use this along with the ts.Client().Get()
	// // method to make a GET /ping request against the test server. This
	// // returns a http.Response struct containing the response.
	// rs, err := ts.Client().Get(ts.URL + "/ping")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // We can then check the value of the response status code and body using
	// // the same code as before.
	// if rs.StatusCode != http.StatusOK {
	// 	t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	// }

	// defer rs.Body.Close()

	// body, err := ioutil.ReadAll(rs.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if string(body) != "OK" {
	// 	t.Errorf("want body to equal %q", "OK")
	// }

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	t.Parallel()

	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up some table-driven tests to check the responses sent by our
	// application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

//  go test -v -run="TestSignupUser" -count 1 ./cmd/web/
func TestSignupUser(t *testing.T) {
	t.Parallel()

	// Create the application struct containing our mocked dependencies and set
	// up the test server for running and end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	defer ts.Close()

	// Make a GET /user/signup request and then extract the CSRF token from the response body.
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestCreateSnippetForm(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/create")
		if code != http.StatusSeeOther {
			t.Errorf("want %d; got %d", http.StatusSeeOther, code)
		}

		if headers.Get("Location") != "/user/login" {
			t.Errorf("want %s; got %s", "/user/login", headers.Get("Location"))
		}
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Authenticate the user...
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)

		// Then check that the authenticated user is shown the create snippet form.
		code, _, body := ts.get(t, "/snippet/create")
		if code != 200 {
			t.Errorf("want %d; got %d", 200, code)
		}

		formTag := "<form action='/snippet/create' method='POST'>"
		if !bytes.Contains(body, []byte(formTag)) {
			t.Errorf("want body %s to contain %q", body, formTag)
		}
	})
}
