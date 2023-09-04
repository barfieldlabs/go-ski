package core

import (
	"context"

	"github.com/chromedp/cdproto/target"
)

type Procedures struct {
	Actions []Action
}

func NewProcedures() *Procedures {
	return &Procedures{}
}

func (p *Procedures) Execute(ctx context.Context, initialTargets []*target.Info) error {
	for _, action := range p.Actions {
		err := action.Perform(ctx, initialTargets)
		if err != nil {
			return err
		}
	}
	return nil
}
