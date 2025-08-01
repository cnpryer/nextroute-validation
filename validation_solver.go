package main

import (
	"github.com/nextmv-io/nextroute"
)

func NewValidationParallelSolver(model nextroute.Model) (nextroute.ParallelSolver, error) {
	parallelSolver, err := nextroute.NewSkeletonParallelSolver(model)
	if err != nil {
		return nil, err
	}

	parallelSolver.SetSolverFactory(
		func(
			information nextroute.ParallelSolveInformation,
			solution nextroute.Solution,
		) (nextroute.Solver, error) {
			return NewValidationSolver(model)
		},
	)

	parallelSolver.SetSolveOptionsFactory(nextroute.DefaultSolveOptionsFactory())

	return parallelSolver, nil
}

func NewValidationSolver(model nextroute.Model) (nextroute.Solver, error) {
	solver, err := nextroute.NewSkeletonSolver(model)
	if err != nil {
		return nil, err
	}

	nrPlanUnits := len(model.PlanUnits())
	unplanSolveParameter, err := nextroute.NewSolveParameter(
		nrPlanUnits,
		0,
		0,
		0,
		nrPlanUnits,
		false,
		false,
	)
	if err != nil {
		return nil, err
	}

	unplanOperator, err := NewValidationUnplanOperator(unplanSolveParameter)
	if err != nil {
		return nil, err
	}

	solver.AddSolveOperators(unplanOperator)

	return solver, nil
}
