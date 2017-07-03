package main

import (
	"net/http"
	"math/rand"
	"strings"
)

func randomRedirect(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.UserAgent(), "facebookexternalhit/") {
		w.Header().Add("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Add("Expires", "Sat, 26 Jul 1997 05:00:00 GMT")
		http.Redirect(w, r, urls[rand.Intn(len(urls))], http.StatusFound)
	} else {
		http.Redirect(w, r, "https://9gag.com/", http.StatusFound)
	}
}

var urls = []string{
	"https://9gag.com/gag/aQe4RzW",
	"https://9gag.com/gag/aMArYjx",
	"https://9gag.com/gag/apQ0Gx5",
	"https://9gag.com/gag/aADreyR",
	"https://9gag.com/gag/aOBb733",
	"https://9gag.com/gag/aq1Nn7P",
	"https://9gag.com/gag/abzyGM8",
	"https://9gag.com/gag/aee0DMm",
	"https://9gag.com/gag/aB83mxP",
	"https://9gag.com/gag/aee0DKb",
	"https://9gag.com/gag/azqE16p",
	"https://9gag.com/gag/aADrgW9",
	"https://9gag.com/gag/aZgx7dV",
	"https://9gag.com/gag/ad9Y5n9",
	"https://9gag.com/gag/ar54G90",
	"https://9gag.com/gag/av7gGzW",
	"https://9gag.com/gag/ax0Z7jD",
	"https://9gag.com/gag/azqE1gN",
	"https://9gag.com/gag/aDz3gRO",
	"https://9gag.com/gag/aGewZv6",
	"https://9gag.com/gag/azqE9ZK",
	"https://9gag.com/gag/a9ALKmm",
	"https://9gag.com/gag/azqE9Mx",
	"https://9gag.com/gag/a6VyOK8",
	"https://9gag.com/gag/aVMWYM8",
	"https://9gag.com/gag/aVMWYDM",
	"https://9gag.com/gag/apQ0mop",
	"https://9gag.com/gag/aDz31wN",
	"https://9gag.com/gag/aOBbrXE",
	"https://9gag.com/gag/a05LK2q",
	"https://9gag.com/gag/a5nPo4y",
	"https://9gag.com/gag/ar54gGK",
	"https://9gag.com/gag/am207Y6",
	"https://9gag.com/gag/aGewe3X",
	"https://9gag.com/gag/aOBbB1N",
	"https://9gag.com/gag/a05L5Ez",
	"https://9gag.com/gag/aRjPjjq",
	"https://9gag.com/gag/aGeweRG",
	"https://9gag.com/gag/aKDPVZO",
	"https://9gag.com/gag/a88YyX6",
	"https://9gag.com/gag/a3M3qN5",
	"https://9gag.com/gag/am20Ym6",
	"https://9gag.com/gag/awQKnNy",
	"https://9gag.com/gag/ayx0LbY",
	"https://9gag.com/gag/aP94MVP",
	"https://9gag.com/gag/a6VyMx2",
	"https://9gag.com/gag/aRjPKL7",
	"https://9gag.com/gag/aOBbd2r",
	"https://9gag.com/gag/aGew0dK",
	"https://9gag.com/gag/a05Ld8L",
	"https://9gag.com/gag/aP94Mwn",
	"https://9gag.com/gag/a05Lp1X",
	"https://9gag.com/gag/av7gnb5",
	"https://9gag.com/gag/aY4WLBv",
	"https://9gag.com/gag/am20zO2",
	"https://9gag.com/gag/a88Y4DQ",
	"https://9gag.com/gag/aKDPqZZ",
	"https://9gag.com/gag/aoO09Rw",
	"https://9gag.com/gag/aMArGE1",
	"https://9gag.com/gag/aoO09We",
}

func main() {
	http.HandleFunc("/", randomRedirect)
	http.ListenAndServe(":8000", nil)
}
