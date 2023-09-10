package examples

import (
	"context"
	"time"

	"github.com/barfieldlabs/go-ski/core"
	"github.com/chromedp/cdproto/target"
)

func LoginExample() {
	proc := core.NewProcedures(true)

	// Define actions
	proc.Actions = []core.Action{
		{
			Type:  core.Navigate,
			URL: "https://example.com/login",
		},
		{
			Type: core.Sleep,
			Delay: 2 * time.Second,
		},
		{
			Type: core.FormSubmit,
			FormDetails: &core.FormDetails{
				Fields: []core.FormField{
					{XPath: "/html/form/p[1]/input", Value: "example_user"},
					{XPath: "/html/form/p[2]/input", Value: "example_password"},
				},
				Submit: "/html/body/main/div[2]/fieldset/form/input",
				Delay:  1 * time.Second, 
			},
		},
		{
			Type: core.Click,
			XPath: "/html/button",
		},
		{
			Type: core.Sleep,
			Delay: 5 * time.Second,
		},
		{
			Type: core.SwitchToIframe,
			IframeXPath: "/html/body/iframe",
		},
		{
			Type: core.Sleep,
			Delay: 5 * time.Second,
		},
		{
			Type: core.Scrape,
			XPath: "/html/body/h1",
		},
	}

	// Create a context
	ctx := context.Background()

	var targetInfo []*target.Info

	// Execute actions
	err := proc.Execute(ctx, targetInfo)
	if err != nil {
		panic(err)
	}
}
