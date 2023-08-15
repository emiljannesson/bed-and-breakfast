package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/emiljannesson/bed-and-breakfast/internal/config"
	"github.com/emiljannesson/bed-and-breakfast/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testAppConfig config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	testAppConfig.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testAppConfig.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testAppConfig.ErrorLog = errorLog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testAppConfig.Session = session

	// the above session setup ensures that our app variable has everything it needs for our tests
	app = &testAppConfig

	os.Exit(m.Run())
}

type testWriter struct{}

func (tw *testWriter) Header() http.Header {
	return http.Header{}
}

func (tw *testWriter) WriteHeader(i int) {
}

func (tw *testWriter) Write(b []byte) (int, error) {
	length := len(b) // we can't return a random int here, so we need to do this
	return length, nil
}
