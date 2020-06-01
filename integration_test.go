package integration_test

import (
	"github.com/dennis/hello_go/app"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// TODO: This test is pretty fragile. It assumes that the dummy data exists
// and looks in a very specific way. It should be changed, so that we add
// the data here as a part of setup
func TestMain(t *testing.T) {
	app := app.App{}
	app.Initialize()
	go app.Run()

	// Let's make sure we're protected properly against unauthenticated
	// requests
	t.Run("Unauth: GetMessages", testUnauthenticatedForGetMessages)
	t.Run("Unauth: GetMessage", testUnauthenticatedForGetMessage)
	t.Run("Unauth: UpdateMessage", testUnauthenticatedForUpdateMessage)
	t.Run("Unauth: CreateMessage", testUnauthenticatedForCreateMessage)
	t.Run("Unauth: DeleteMessage", testUnauthenticatedForDeleteMessage)

	t.Run("Auth: GetMessages", testAuthenticatedForGetMessages)
	t.Run("Auth: DeleteMessage", testAuthenticatedForDeleteMessage)
}

func assertUnauthenticated(t *testing.T, r *http.Request, rerr error) {
	if rerr != nil {
		t.Errorf("Error creating request: %v", rerr)
		return
	}

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	} else if res.StatusCode != 401 {
		t.Errorf("Unexpected status code %v, want 401", res.StatusCode)
	}
}

func testUnauthenticatedForGetMessages(t *testing.T) {
	r, rerr := http.NewRequest("GET", "http://localhost:8080/api/messages", nil)

	assertUnauthenticated(t, r, rerr)
}

func testUnauthenticatedForGetMessage(t *testing.T) {
	r, rerr := http.NewRequest("GET", "http://localhost:8080/api/messages/1", nil)

	assertUnauthenticated(t, r, rerr)
}

func testUnauthenticatedForUpdateMessage(t *testing.T) {
	r, rerr := http.NewRequest("PUT", "http://localhost:8080/api/messages/1", nil)

	assertUnauthenticated(t, r, rerr)
}

func testUnauthenticatedForCreateMessage(t *testing.T) {
	r, rerr := http.NewRequest("POST", "http://localhost:8080/api/messages", nil)

	assertUnauthenticated(t, r, rerr)
}

func testUnauthenticatedForDeleteMessage(t *testing.T) {
	r, rerr := http.NewRequest("DELETE", "http://localhost:8080/api/messages/1", nil)

	assertUnauthenticated(t, r, rerr)
}

func testAuthenticatedForGetMessages(t *testing.T) {
	r, rerr := http.NewRequest("GET", "http://localhost:8080/api/messages", nil)
	r.Header.Add("Authorization", "Basic YXV0aHRva2VuZGVubmlzOg==")

	if rerr != nil {
		t.Errorf("Error creating request: %v", rerr)
		return
	}

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	} else if res.StatusCode != 200 {
		t.Errorf("Unexpected status code %v, want 200", res.StatusCode)
	} else {
		body, _ := ioutil.ReadAll(res.Body)

		// Let's look for a known topic
		if strings.Index(string(body), "Lorem lipsum") == -1 {
			t.Errorf("Expected message not present")
		}
	}
}

func testAuthenticatedForDeleteMessage(t *testing.T) {
	r, rerr := http.NewRequest("DELETE", "http://localhost:8080/api/messages/1", nil)
	r.Header.Add("Authorization", "Basic YXV0aHRva2VuZGVubmlzOg==")

	if rerr != nil {
		t.Errorf("Error creating request: %v", rerr)
		return
	}

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	} else if res.StatusCode != 200 {
		t.Errorf("Unexpected status code %v, want 200", res.StatusCode)
	} else {
		r, rerr := http.NewRequest("GET", "http://localhost:8080/api/messages", nil)
		r.Header.Add("Authorization", "Basic YXV0aHRva2VuZGVubmlzOg==")

		if rerr != nil {
			t.Errorf("Error creating request: %v", rerr)
			return
		}

		client := http.Client{}
		res, err := client.Do(r)
		if err != nil {
			t.Errorf("Error sending request: %v", err)
		} else if res.StatusCode != 200 {
			t.Errorf("Unexpected status code %v, want 200", res.StatusCode)
		} else {
			body, _ := ioutil.ReadAll(res.Body)

			// Let's look for a known topic - we have just delete it, so
			// we expect it to be missing
			if strings.Index(string(body), "Lorem lipsum") > -1 {
				t.Errorf("Response contains deleted message")
			}
		}
	}
}
