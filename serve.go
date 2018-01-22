package nineRandom

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	fb "github.com/huandu/facebook"
	"github.com/jmcvetta/randutil"

	"appengine"
	"appengine/taskqueue"
	"appengine/urlfetch"
)

func randomRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && strings.HasPrefix(r.UserAgent(), "facebookexternalhit/") {
		w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Add("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
		http.Redirect(w, r, getRandomURL(urls), http.StatusFound)
	} else {
		http.Redirect(w, r, strings.TrimSpace(os.Getenv("SITE_URL")), http.StatusFound)
	}
}

func appURL() string {
	return strings.TrimSpace(os.Getenv("APP_URL"))
}

func scrapeURLs() []string {
	return strings.Split(strings.TrimSpace(os.Getenv("SCRAPE_URLS")), ",")
}

func getRandomURL(urls weightedList) string {
	item, err := randutil.WeightedChoice(urls)
	if err != nil {
		// In case of error, return self for a clean retry
		return appURL()
	}
	fmt.Printf("Hit in %d region\n", item.Weight)
	url, err := randutil.ChoiceString(item.Item.([]string))
	if err != nil {
		// In case of error, return self for a clean retry
		return appURL()
	}
	return url
}

func refreshFacebook(w http.ResponseWriter, r *http.Request) {
	appID := strings.TrimSpace(os.Getenv("FACEBOOK_APP_ID"))
	appSecret := strings.TrimSpace(os.Getenv("FACEBOOK_APP_SECRET"))

	if len(appID) == 0 || len(appSecret) == 0 {
		fmt.Printf("No valid Facebook credentials")
		return
	}

	var globalApp = fb.New(appID, appSecret)

	session := globalApp.Session(globalApp.AppAccessToken())
	session.HttpClient = urlfetch.Client(appengine.NewContext(r))

	for _, url := range scrapeURLs() {
		res, _ := session.Post("/", fb.Params{
			"id":     url,
			"scrape": "true",
		})

		if res.Err() != nil {
			fmt.Printf("Error occurred with scrape: %s", res.Err())
			return
		}
	}
}

type weightedList []randutil.Choice

var urls = weightedList{
	randutil.Choice{
		Weight: 70,
		Item: []string{
			"https://9gag.com/gag/aKDqdBZ",
			"https://9gag.com/gag/aRjA0ZM",
			"https://9gag.com/gag/aMAGY4X",
			"https://9gag.com/gag/ad9j5Od",
			"https://9gag.com/gag/aVMP0E8",
		},
	},
	randutil.Choice{
		Weight: 30,
		Item: []string{
			"https://9gag.com/gag/a9A7M1Z",
			"https://9gag.com/gag/a4GY9x6",
			"https://9gag.com/gag/a88MbQe",
			"https://9gag.com/gag/aNznoEb",
			"https://9gag.com/gag/aAD1gVZ",
			"https://9gag.com/gag/aVMP0g8",
			"https://9gag.com/gag/a4GYR6p",
			"https://9gag.com/gag/aWq6E7q",
			"https://9gag.com/gag/ax0j73b",
		},
	},
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func scheduleFacebookRefresh(w http.ResponseWriter, r *http.Request) {
	t := taskqueue.NewPOSTTask("/_internal/refreshFacebook", url.Values{"refresh": {"true"}})
	batches := 3
	count := 50
	tasks := make([]*taskqueue.Task, count)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = t
	}
	c := appengine.NewContext(r)

	for x := 0; x < batches; x++ {
		if _, err := taskqueue.AddMulti(c, tasks, "facebook-refresh"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func init() {
	// go refreshFacebook()
	http.HandleFunc("/", randomRedirect)
	http.HandleFunc("/_internal/scheduleFacebook", scheduleFacebookRefresh)
	http.HandleFunc("/_internal/refreshFacebook", refreshFacebook)
	http.HandleFunc("/_ah/health", healthCheckHandler)
}
