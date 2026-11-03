package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/puppe1990/matheuspuppe-blog/internal/ogcard"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	out := flag.String("out", "web/static/og.png", "PNG path")
	host := flag.String("host", "sintaxe.apps.gestaobem.com", "hostname painted on the card")
	title := flag.String("title", "", "if set, paint an article card with this headline")
	kicker := flag.String("kicker", "", "article kicker (caderno)")
	byline := flag.String("byline", "", "article byline")
	flag.Parse()

	card := ogcard.SiteCard(*host)
	if *title != "" {
		card = ogcard.Card{Title: *title, Kicker: *kicker, Byline: *byline, Host: *host}
	}

	f, err := os.Create(*out)
	if err != nil {
		return fmt.Errorf("create %s: %w", *out, err)
	}
	defer func() { _ = f.Close() }()
	if err := ogcard.WritePNG(f, card); err != nil {
		return fmt.Errorf("write png: %w", err)
	}
	return nil
}
