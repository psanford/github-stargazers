package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

var authToken = flag.String("auth_token", "", "Auth token for better rate limits")

func main() {
	flag.Parse()
	ctx := context.Background()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <user/repo>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	ownerRepo := strings.Split(args[0], "/")
	if len(ownerRepo) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <user/repo>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var hc *http.Client
	if *authToken != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *authToken},
		)
		hc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(hc)

	var (
		owner = ownerRepo[0]
		repo  = ownerRepo[1]
		w     = csv.NewWriter(os.Stdout)

		opts github.ListOptions
	)

	w.Write([]string{"login", "id", "starred_at"})

	for {
		gazers, resp, err := client.Activity.ListStargazers(ctx, owner, repo, &opts)
		if err != nil {
			panic(err)
		}

		for _, g := range gazers {
			w.Write([]string{*g.User.Login, strconv.Itoa(int(*g.User.ID)), g.StarredAt.String()})
			w.Flush()
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
}
