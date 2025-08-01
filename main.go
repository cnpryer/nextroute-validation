// Package main holds the implementation for the app.
package main

import (
	"context"
	"log"

	"github.com/nextmv-io/nextroute"
	"github.com/nextmv-io/nextroute/check"
	"github.com/nextmv-io/nextroute/factory"
	"github.com/nextmv-io/nextroute/schema"
	"github.com/nextmv-io/sdk/run"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

func main() {
	runner := run.CLI(solver)
	err := runner.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

type options struct {
	Model  factory.Options                `json:"model,omitempty"`
	Solve  nextroute.ParallelSolveOptions `json:"solve,omitempty"`
	Format nextroute.FormatOptions        `json:"format,omitempty"`
	Check  check.Options                  `json:"check,omitempty"`
	Custom customOptions                  `json:"custom,omitempty"`
}

type customOptions struct {
	Validate bool `json:"validate,omitempty"`
}

func solver(
	ctx context.Context,
	input schema.Input,
	options options,
) (runSchema.Output, error) {
	model, err := factory.NewModel(input, options.Model)
	if err != nil {
		return runSchema.Output{}, err
	}

	solveOptions := options.Solve
	var solver nextroute.ParallelSolver

	if options.Custom.Validate {
		solver, err = NewValidationParallelSolver(model)
		if err != nil {
			return runSchema.Output{}, err
		}
		solveOptions.Iterations = 0 // No iterations needed for validation
	} else {
		solver, err = nextroute.NewParallelSolver(model)
		if err != nil {
			return runSchema.Output{}, err
		}
	}

	solutions, err := solver.Solve(ctx, solveOptions)
	if err != nil {
		return runSchema.Output{}, err
	}

	last, err := solutions.Last()
	if err != nil {
		return runSchema.Output{}, err
	}

	output, err := check.Format(
		ctx,
		options,
		options.Check,
		solver,
		last,
	)
	if err != nil {
		return runSchema.Output{}, err
	}
	output.Statistics.Result.Custom = factory.DefaultCustomResultStatistics(last)

	return output, nil
}
