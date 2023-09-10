package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	net_url "net/url"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

type ActionType string

const (
	Click      			ActionType = "click"
	FormSubmit 			ActionType = "formSubmit"
	Scrape     			ActionType = "scrape"
	Navigate   			ActionType = "navigate"
	Sleep      			ActionType = "sleep"
	SwitchToIframe 		ActionType = "switchToIframe"
	ReturnToRoot 		ActionType = "returnToRoot"
)

type FormField struct {
	XPath string
	Value string
}

type FormDetails struct {
	Fields  []FormField  
	Submit  string         
	Delay   time.Duration  
}

type Action struct {
	Type        ActionType
	Delay       time.Duration
	XPath       string
	IframeXPath string
	SearchIframe   bool
	URL		 string
	SwitchTab   bool
	FormDetails *FormDetails 
}

func (a *Action) Perform(ctx context.Context, initialTargets []*target.Info, contextStack *[]context.Context) (context.Context, error) {
	var err error
	var newCtx context.Context

	switch a.Type {
	case Click:
		newCtx, err = performClick(ctx, a.XPath, a.SwitchTab, initialTargets)
		if err != nil {
			return nil, fmt.Errorf("Failed to perform click: %v", err)
		}
		return newCtx, nil
	case FormSubmit:
		if a.FormDetails == nil {
			return nil, errors.New("FormDetails cannot be nil for FormSubmit action")
		}
		err = performFormSubmit(ctx, a.FormDetails)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	case Navigate:
		err = performNavigate(ctx, a.URL)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	case SwitchToIframe:
		newCtx, err := performSwitchToIframe(ctx, a.XPath, contextStack)
		if err != nil {
			return nil, err
		}
		return newCtx, nil
	case ReturnToRoot:
		newCtx, err := performReturnToRoot(contextStack)
		if err != nil {
			return nil, err
		}
		return newCtx, nil
	case Sleep:
		err = chromedp.Run(ctx, chromedp.Sleep(a.Delay))
		if err != nil {
			return nil, err
		}
		return ctx, nil
	case Scrape:
		scrapedData, err := performScrape(ctx, a.XPath)
		if err != nil {
			return nil, err
		}
		fmt.Println("Scraped Data:", scrapedData)
		return ctx, nil
	default:
		return nil, errors.New("unknown action type")
	}
}

func performClick(ctx context.Context, xpath string, switchTab bool, initialTargets []*target.Info) (context.Context, error) {
	var err error

	// Perform the click
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(xpath),
		chromedp.Click(xpath, chromedp.NodeVisible),
	)
	if err != nil {
		return nil, err
	}

	if switchTab {
		// Get the list of current targets (tabs)
		var currentTargets []*target.Info
		err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			currentTargets, err = target.GetTargets().Do(ctx)
			return err
		}))
		if err != nil {
			return nil, fmt.Errorf("Failed to get current targets: %v", err)
		}

		// Find the new tab by comparing with initial targets
		var newTarget *target.Info
		for _, curr := range currentTargets {
			isNew := true
			for _, init := range initialTargets {
				if curr.TargetID == init.TargetID {
					isNew = false
					break
				}
			}
			if isNew {
				newTarget = curr
				break
			}
		}

		if newTarget == nil {
			return nil, errors.New("No new tab found")
		}

		// Switch to the new tab
		newCtx, cancel := context.WithCancel(ctx)
		defer cancel()  

		err = chromedp.Run(newCtx, chromedp.ActionFunc(func(ctx context.Context) error {
			return target.ActivateTarget(newTarget.TargetID).Do(ctx)
		}))
		if err != nil {
			return nil, fmt.Errorf("Failed to switch to new tab: %v", err)
		}

		return newCtx, nil
	}

	return ctx, nil
}

func performFormSubmit(ctx context.Context, details *FormDetails) error {
	if details == nil {
		log.Println("FormDetails is nil")
		return errors.New("FormDetails cannot be nil")
	}

	tasks := []chromedp.Action{}

	// Wait for the form to be visible (you can choose any field or the submit button)
	if details.Submit != "" {
		log.Println("Waiting for submit button to be visible")
		tasks = append(tasks, chromedp.WaitVisible(details.Submit))
	}

	// Fill in the fields
	for _, field := range details.Fields {
		log.Printf("Waiting for field with XPath %s to be visible", field.XPath)
		tasks = append(tasks, chromedp.WaitVisible(field.XPath))
		log.Printf("Waiting for field with XPath %s to be enabled", field.XPath)
		tasks = append(tasks, chromedp.WaitEnabled(field.XPath))
		log.Printf("Sending keys to field with XPath %s", field.XPath)
		tasks = append(tasks, chromedp.SendKeys(field.XPath, field.Value))
	}

	// Add optional delay
	if details.Delay > 0 {
		log.Printf("Sleeping for %v", details.Delay)
		tasks = append(tasks, chromedp.Sleep(details.Delay))
	}

	// Click the submit button or hit Enter
	if details.Submit != "" {
		log.Println("Clicking the submit button")
		tasks = append(tasks, chromedp.Click(details.Submit))
	} else {
		log.Println("Sending Enter key event")
		tasks = append(tasks, chromedp.KeyEvent("\r"))
	}

	// Run the tasks
	log.Println("Running tasks")
	err := chromedp.Run(ctx, tasks...)
	if err != nil {
		log.Printf("Error running tasks: %v", err)
		return err
	}

	log.Println("Tasks completed successfully")
	return nil
}

func performNavigate(ctx context.Context, url string) error {
	log.Println("Starting to navigate to URL:", url)
    // Validate URL
    _, err := net_url.Parse(url)
    if err != nil {
        return fmt.Errorf("Invalid URL '%s': %v", url, err)
    }

    err = chromedp.Run(ctx, chromedp.Navigate(url))
    if err != nil {
        return fmt.Errorf("Failed to navigate to URL '%s': %v. Is the URL correct?", url, err)
    }

	if err != nil {
		log.Println("Failed to navigate:", err)
		return err
	}
	log.Println("Successfully navigated to URL:", url)


    return nil
}

func performSwitchToIframe(ctx context.Context, xpath string, contextStack *[]context.Context) (context.Context, error) {
	log.Println("Starting to switch to iframe with XPath:", xpath)

	var iframes []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Nodes(xpath, &iframes, chromedp.ByQueryAll))
	if err != nil {
		return nil, err
	}

	if len(iframes) == 0 {
		return nil, errors.New("no iframes found")
	}

	// Assume the iframe we need is the last one in the list
	targetIframe := iframes[len(iframes)-1]

	// Create a new context for the iframe
	iframeCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	// Switch to the iframe context
	err = chromedp.Run(iframeCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		// Wait for an element inside the iframe to become visible
		err := chromedp.WaitVisible(xpath, chromedp.FromNode(targetIframe)).Do(ctx)
		if err != nil {
			return fmt.Errorf("Failed to wait for element inside iframe: %v", err)
		}
		log.Println("Switched to iframe")
		return nil
	}))

	if err != nil {
		log.Println("Failed to switch to iframe:", err)
		return nil, err
	}

	// Add the iframe context to the context stack
	*contextStack = append(*contextStack, iframeCtx)

	log.Println("Successfully switched to iframe with XPath:", xpath)
	return iframeCtx, nil
}

func performReturnToRoot(contextStack *[]context.Context) (context.Context, error) {
	if len(*contextStack) == 0 {
		return nil, errors.New("no context to return to")
	}

	lastIndex := len(*contextStack) - 1
	rootCtx := (*contextStack)[lastIndex]

	*contextStack = (*contextStack)[:lastIndex]

	return rootCtx, nil
}

func performScrape(ctx context.Context, xpath string) (string, error) {
	var scrapedText string

	// Define the tasks to perform
	tasks := []chromedp.Action{
		chromedp.WaitVisible(xpath), // Wait until the element is visible
		chromedp.Text(xpath, &scrapedText, chromedp.NodeVisible), // Scrape the text content
	}

	fmt.Println("Starting to scrape element with XPath:", xpath) // Log the start

	// Run the tasks
	if err := chromedp.Run(ctx, tasks...); err != nil {
		fmt.Println("Failed to perform scrape:", err) // Log the failure
		return "", fmt.Errorf("Failed to perform scrape: %v", err)
	}

	fmt.Println("Successfully scraped text:", scrapedText) // Log the success

	return scrapedText, nil
}