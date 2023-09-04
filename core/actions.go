package core

import (
	"context"
	"errors"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

type ActionType string

const (
	Click      ActionType = "click"
	FormSubmit ActionType = "formSubmit"
	Scrape	 ActionType = "scrape"
)

type Action struct {
	Type       ActionType
	XPath      string
	Iframe     *string
	SwitchTab  bool
	FormFields map[string]string 
	SubmitBtn  string           
}

func (a *Action) Perform(ctx context.Context, initialTargets []*target.Info) error {
	switch a.Type {
	case Click:
		return performClick(ctx, a.XPath, a.Iframe, a.SwitchTab, initialTargets)
	case FormSubmit:
		return performFormSubmit(ctx, a.XPath, a.FormFields, a.SubmitBtn)
	default:
		return errors.New("unknown action type")
	}
}

func performClick(ctx context.Context, xpath string, iframe *string, switchTab bool, initialTargets []*target.Info) error {
	var err error
	if iframe != nil {
		// Switch to iframe logic here
	}

	// Perform click
	err = chromedp.Run(ctx, chromedp.Click(xpath, chromedp.NodeVisible))
	if err != nil {
		return err
	}

	if switchTab {
		var finalTargets []*target.Info
		err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			finalTargets, err = target.GetTargets().Do(ctx)
			return err
		}))
	
		if err != nil { 
			return err
		}

		var newTab *target.Info
		for _, finalTarget := range finalTargets {
			var isNew = true
			for _, initialTarget := range initialTargets {
				if finalTarget.TargetID == initialTarget.TargetID {
					isNew = false
					break
				}
			}
			if isNew {
				newTab = finalTarget
				break
			}
		}

		if newTab == nil {
			return errors.New("no new tab found")
		}

		newTabCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(newTab.TargetID))
		defer cancel()
		
		// Perform actions in the new tab
		if err := chromedp.Run(newTabCtx,
			// Actions
		); err != nil {
			return err
		}		

		// Perform actions in the new tab 
	}
	return nil
}

func performFormSubmit(ctx context.Context, xpath string, formFields map[string]string, submitBtn string) error {
	tasks := []chromedp.Action{
		chromedp.WaitVisible(xpath),
	}

	for field, value := range formFields {
		tasks = append(tasks, chromedp.SendKeys(field, value))
	}

	tasks = append(tasks, chromedp.Click(submitBtn))

	err := chromedp.Run(ctx, tasks...)
	if err != nil {
		return err
	}

	return nil
}