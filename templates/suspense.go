package templates

import (
	"context"
	"io"
	"sync"

	"github.com/a-h/templ"
)

type SuspenseCtxKey int

const (
	SuspenseKey SuspenseCtxKey = iota
)

type SuspenseCtx struct {
	out chan templ.Component
	wg  *sync.WaitGroup
}

func (s *SuspenseCtx) Stream(ctx context.Context, w io.Writer) {
	go func() {
		s.wg.Wait()
		close(s.out)
	}()

	for c := range s.out {
		c.Render(ctx, w)
	}
}

func NewSuspenseCtx() *SuspenseCtx {
	return &SuspenseCtx{
		out: make(chan templ.Component),
		wg:  &sync.WaitGroup{},
	}
}

func WithSuspenseCtx(ctx context.Context, susCtx *SuspenseCtx) context.Context {
	return context.WithValue(ctx, SuspenseKey, susCtx)
}

func GetSuspenseCtx(ctx context.Context) *SuspenseCtx {
	return ctx.Value(SuspenseKey).(*SuspenseCtx)
}

func Suspense[T any](
	placeholder templ.Component,
	data func() (T, error),
	component func(T) templ.Component,
	errorComponent func(error) templ.Component,
) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		susCtx := GetSuspenseCtx(ctx)

		susCtx.wg.Add(1)
		go func() {
			defer susCtx.wg.Done()
			d, err := data()
			if err != nil {
				susCtx.out <- errorComponent(err)
				return
			}
			susCtx.out <- component(d)
		}()

		return placeholder.Render(ctx, w)
	})
}
