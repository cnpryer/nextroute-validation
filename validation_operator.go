package main

import (
	"context"

	"github.com/nextmv-io/nextroute"
	"github.com/nextmv-io/nextroute/common"
)

func NewValidationUnplanOperator(
	numberOfUnits nextroute.SolveParameter,
) (nextroute.SolveOperatorUnPlan, error) {
	return &validationUnplanOperator{
		SolveOperator: nextroute.NewSolveOperator(
			1.0,
			false,
			nextroute.SolveParameters{numberOfUnits},
		),
	}, nil
}

type validationUnplanOperator struct {
	nextroute.SolveOperator
}

func (o *validationUnplanOperator) NumberOfUnits() nextroute.SolveParameter {
	return o.Parameters()[0]
}

func (o *validationUnplanOperator) Execute(
	ctx context.Context,
	runTimeInformation nextroute.SolveInformation,
) error {
	workSolution := runTimeInformation.
		Solver().
		WorkSolution()

	if workSolution.PlannedPlanUnits().Size() == 0 {
		return nil
	}

	unplanAllInfeasibleVehicles(workSolution)

	return nil
}

func unplanAllInfeasibleVehicles(solution nextroute.Solution) error {
	vehicles := common.Filter(solution.Vehicles(), func(vehicle nextroute.SolutionVehicle) bool {
		return !vehicle.IsEmpty()
	})

	if len(vehicles) == 0 {
		return nil
	}

	for _, vehicle := range vehicles {
		if vehicle.IsEmpty() || vehicle.NumberOfStops() == 0 {
			continue
		}

		for _, constraint := range solution.Model().Constraints() {
			if constraint.(nextroute.SolutionVehicleViolationCheck).DoesVehicleHaveViolations(vehicle) {
				vehicle.Unplan()
			}
		}
	}

	return nil
}
