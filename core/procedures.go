package core

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

type Procedures struct {
	Actions []Action
	UseGUI  bool // New field to control GUI
}

func NewProcedures(useGUI bool) *Procedures {
	return &Procedures{UseGUI: useGUI}
}

func (p *Procedures) Execute(ctx context.Context, initialTargets []*target.Info) error {
	// Create a new context for chromedp
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)

	// If UseGUI is false, add the Headless option
	if !p.UseGUI {
		opts = append(opts, chromedp.Headless)
	} else {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	var contextStack []context.Context // Initialize a stack for context management

	currentCtx := ctx // Initialize with the chromedp context

	for _, action := range p.Actions {
		newCtx, err := action.Perform(currentCtx, initialTargets, &contextStack)
		if err != nil {
			log.Fatalf("Failed to perform action: %v", err)
			return err
		}

		// Update the current context for the next action
		currentCtx = newCtx
	}

	return nil
}
