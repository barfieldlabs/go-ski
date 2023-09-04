package examples

import (
	"context"

	"github.com/barfieldlabs/go-ski/core"
	"github.com/chromedp/cdproto/target"
)

func CanvasExample() {
	proc := core.NewProcedures()

	proc.Actions = []core.Action{
		{
			Type:  core.Click,
			XPath: "some_xpath",
		},
		{
			Type:       core.FormSubmit,
			XPath:      "form_xpath",
			FormFields: map[string]string{"username": "user", "password": "pass"},
			SubmitBtn:  "submit_btn_xpath",
		},
		{
			Type:  core.Scrape,
			XPath: "scrape_xpath",
		},
	}

	ctx := context.Background()

	var targetInfo []*target.Info

	err := proc.Execute(ctx, targetInfo)
	if err != nil {
		panic(err)
	}
}
