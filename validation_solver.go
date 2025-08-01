package main

import (
	"github.com/nextmv-io/nextroute"
)

func NewValidationParallelSolver(model nextroute.Model) (nextroute.ParallelSolver, error) {
	parallelSolver, err := nextroute.NewSkeletonParallelSolver(model)
	if err != nil {
		return nil, err
	}

	// Set the solver factory for the parallel solver. This factory is used to
	// create new solver instances for each cycle. The information contains data
	// about the current cycle and which solver of the n solvers is being
	// created. The solution is the best solution of the previous cycle (and
	// globally best).
	//
	// In this example we create identical solvers with custom operators, but
	// you can also create different solvers with different operators. There
	// is a random component in the operators, so the solvers will behave
	// differently.
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
