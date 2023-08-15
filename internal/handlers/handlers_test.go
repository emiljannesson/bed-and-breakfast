package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{
		name:               "home",
		url:                "/",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "about",
		url:                "/about",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "contact",
		url:                "/contact",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "generals-quarters",
		url:                "/generals-quarters",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "majors-suite",
		url:                "/majors-suite",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "make-reservation",
		url:                "/make-reservation",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "search-availability",
		url:                "/search-availability",
		method:             "GET",
		params:             []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "post-search-availability",
		url:    "/search-availability",
		method: "POST",
		params: []postData{
			{
				key:   "start",
				value: "2023-08-15",
			},
			{
				key:   "end",
				value: "2023-08-16",
			},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "post-search-availability-json",
		url:    "/search-availability-json",
		method: "POST",
		params: []postData{
			{
				key:   "start",
				value: "2023-08-15",
			},
			{
				key:   "end",
				value: "2023-08-16",
			},
		},
		expectedStatusCode: http.StatusOK,
	},
	{
		name:   "post-make-reservation",
		url:    "/make-reservation",
		method: "POST",
		params: []postData{
			{
				key:   "first_name",
				value: "John",
			},
			{
				key:   "last_name",
				value: "Smith",
			},
			{
				key:   "email",
				value: "test@example.com",
			},
			{
				key:   "phone",
				value: "888-555-5555",
			},
		},
		expectedStatusCode: http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("test failed: %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range test.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("test failed: %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
