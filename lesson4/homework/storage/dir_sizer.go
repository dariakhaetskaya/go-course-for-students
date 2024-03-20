package storage

import (
	"context"
	"golang.org/x/sync/errgroup"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: 4,
	}
}

func (a *sizer) getSize(ctx context.Context, d Dir, results chan Result, group *errgroup.Group) error {
	result := Result{}
	defer func() {
		results <- result
	}()

	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return err
	}

	result.Count = int64(len(files))
	for _, file := range files {
		size, err := file.Stat(ctx)
		if err != nil {
			return err
		}
		result.Size += size
	}

	for i := range dirs {
		dir := dirs[i]
		group.Go(func() error {
			return a.getSize(ctx, dir, results, group)
		})
	}

	return nil
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	result := Result{}
	results := make(chan Result)
	defer close(results)

	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(a.maxWorkersCount)
	errGroup.Go(func() error {
		return a.getSize(errGroupCtx, d, results, errGroup)
	})

	errGroupWaitCh := make(chan error)
	go func() {
		errGroupWaitCh <- errGroup.Wait()
		close(errGroupWaitCh)
	}()

	for {
		select {
		case <-ctx.Done():
			return Result{}, ctx.Err()
		case r := <-results:
			result.Size += r.Size
			result.Count += r.Count
		case err := <-errGroupWaitCh:
			return result, err
		}
	}
}
