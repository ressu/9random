package nineRandom

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	fb "github.com/huandu/facebook"
	"github.com/jmcvetta/randutil"

	"appengine"
	"appengine/taskqueue"
	"appengine/urlfetch"
)

type gotGame struct {
	appURL                           string
	siteURL                          string
	enableScrape                     bool
	scrapeURLs                       []string
	facebookAppID, facebookAppSecret string
}

func (g gotGame) randomRedirect(w http.ResponseWriter, r *http.Request) {
	if g.enableScrape && r.URL.Path == "/" && strings.HasPrefix(r.UserAgent(), "facebookexternalhit/") {
		w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Add("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
		http.Redirect(w, r, g.getRandomURL(urls), http.StatusFound)
	} else {
		http.Redirect(w, r, g.siteURL, http.StatusFound)
	}
}

func (g gotGame) getRandomURL(urls weightedList) string {
	item, err := randutil.WeightedChoice(urls)
	if err != nil {
		// In case of error, return self for a clean retry
		return g.appURL
	}
	fmt.Printf("Hit in %d region\n", item.Weight)
	url, err := randutil.ChoiceString(item.Item.([]string))
	if err != nil {
		// In case of error, return self for a clean retry
		return g.appURL
	}
	return url
}

func (g gotGame) refreshFacebook(w http.ResponseWriter, r *http.Request) {
	if !g.enableScrape {
		http.Error(w, "Scrape disabled", http.StatusForbidden)
		return
	}

	var globalApp = fb.New(g.facebookAppID, g.facebookAppSecret)

	session := globalApp.Session(globalApp.AppAccessToken())
	session.HttpClient = urlfetch.Client(appengine.NewContext(r))

	for _, url := range g.scrapeURLs {
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

func (g gotGame) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func (g gotGame) scheduleFacebookRefresh(w http.ResponseWriter, r *http.Request) {
	if !g.enableScrape {
		http.Error(w, "Scrape disabled", http.StatusForbidden)
		return
	}

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
	scrape, err := strconv.ParseBool(os.Getenv("ENABLE_SCRAPE"))
	if err != nil {
		scrape = true
	}

	appID, appSecret := "", ""
	if scrape {
		appID = strings.TrimSpace(os.Getenv("FACEBOOK_APP_ID"))
		appSecret = strings.TrimSpace(os.Getenv("FACEBOOK_APP_SECRET"))

		if len(appID) == 0 || len(appSecret) == 0 {
			fmt.Printf("No valid Facebook credentials")
			os.Exit(1)
		}
	}

	g := gotGame{
		enableScrape:      scrape,
		siteURL:           strings.TrimSpace(os.Getenv("SITE_URL")),
		appURL:            strings.TrimSpace(os.Getenv("APP_URL")),
		scrapeURLs:        strings.Split(strings.TrimSpace(os.Getenv("SCRAPE_URLS")), ","),
		facebookAppID:     appID,
		facebookAppSecret: appSecret,
	}

	http.HandleFunc("/", g.randomRedirect)
	http.HandleFunc("/_internal/scheduleFacebook", g.scheduleFacebookRefresh)
	http.HandleFunc("/_internal/refreshFacebook", g.refreshFacebook)
	http.HandleFunc("/_ah/health", g.healthCheckHandler)
}
