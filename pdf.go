package main

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"time"
)

func (service *Service) PdfHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	err := printPDF(r.Context(), &buf)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func printPDF(ctx context.Context, buf *[]byte) error {
	url := "https://en.wikipedia.org/wiki/Special:Random"
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()
	var title string
	err := chromedp.Run(ctx,
		chromedp.Tasks{
			enableLifeCycleEvents(),
			chromedp.Navigate(url),
			waitFor("networkIdle"),
			chromedp.Title(&title),
			printToPdf(buf, &title),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func printToPdf(res *[]byte, title *string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var err error
		log.Printf("Printing page \"%s\"", *title)
		*res, _, err = page.PrintToPDF().Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func enableLifeCycleEvents() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		err := page.Enable().Do(ctx)
		if err != nil {
			return err
		}
		err = page.SetLifecycleEventsEnabled(true).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

// waitFor blocks until eventName is received.
// Events you can wait for:
// init, DOMContentLoaded, firstPaint,
// firstContentfulPaint, firstImagePaint,
// firstMeaningfulPaintCandidate,
// load, networkAlmostIdle, firstMeaningfulPaint, networkIdle
func waitFor(eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		ch := make(chan struct{})
		cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		chromedp.ListenTarget(cctx, func(ev interface{}) {
			switch e := ev.(type) {
			case *page.EventLifecycleEvent:
				if e.Name == eventName {
					cancel()
					close(ch)
				}
			}
		})
		select {
		case <-ch:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-cctx.Done():
			return nil
		}
	}
}
