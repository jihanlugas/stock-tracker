package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Generate PDF menggunakan chromedp dan cdproto/page
func GeneratePDFWithChromedp(inputHTML string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	htmlFile := fmt.Sprintf("file://%s/%s", fullPath, inputHTML)
	var pdfBuffer []byte

	var contentHeight float64

	// Menggunakan cdproto/page untuk generate PDF
	err = chromedp.Run(ctx,
		chromedp.Navigate(htmlFile),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(`document.body.scrollHeight`, &contentHeight),
		chromedp.ActionFunc(func(ctx context.Context) error {
			const pxPerInch = 96.0
			heightInInch := contentHeight / pxPerInch
			buf, _, err := page.PrintToPDF().
				WithPaperWidth(5.9).
				WithPaperHeight(heightInInch).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPreferCSSPageSize(false).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuffer = buf
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
