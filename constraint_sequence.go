package main

import (
	"github.com/nextmv-io/nextroute"
)

func NewStopSequenceStopConstraint() *stopSequenceStopConstraintImpl {
	return &stopSequenceStopConstraintImpl{}
}

type stopSequenceStopConstraintImpl struct{}

func (c *stopSequenceStopConstraintImpl) EstimateIsViolated(
	move nextroute.SolutionMoveStops,
) (isViolated bool, stopPositionsHint nextroute.StopPositionsHint) {
	checker := newChecker()
	return checker.isMoveInfeasible(move), nextroute.NoPositionsHint()
}

func NewStopSequenceVehicleConstraint() *stopSequenceVehicleConstraintImpl {
	return &stopSequenceVehicleConstraintImpl{}
}

type stopSequenceVehicleConstraintImpl struct{}

func (c *stopSequenceVehicleConstraintImpl) EstimateIsViolated(
	move nextroute.SolutionMoveStops,
) (isViolated bool, stopPositionsHint nextroute.StopPositionsHint) {
	checker := newChecker()
	return checker.isMoveInfeasible(move), nextroute.NoPositionsHint()
}

func (c *stopSequenceVehicleConstraintImpl) DoesVehicleHaveViolations(
	vehicle nextroute.SolutionVehicle,
) bool {
	checker := newChecker()
	return checker.isVehicleInfeasible(vehicle)
}

func newChecker() *stopSequenceChecker {
	return &stopSequenceChecker{}
}

type stopSequenceChecker struct{}

func (c *stopSequenceChecker) isVehicleInfeasible(vehicle nextroute.SolutionVehicle) bool {
	if vehicle.IsEmpty() || vehicle.NumberOfStops() == 0 {
		return false
	}

	stops := vehicle.SolutionStops()
	last := ""
	for _, stop := range stops {
		if stop.IsZero() {
			continue
		}
		stopId := stop.ModelStop().ID()
		if last != "" && stopId < last {
			return true
		}
		last = stopId
	}
	return false
}

func (c *stopSequenceChecker) isMoveInfeasible(move nextroute.SolutionMoveStops) bool {
	generator := nextroute.NewSolutionStopGenerator(move, true, true)
	last := ""
	for stop := generator.Next(); !stop.IsZero(); stop = generator.Next() {
		if stop.IsFirst() || stop.IsLast() {
			continue
		}
		if last != "" && stop.ModelStop().ID() < last {
			return true
		}
		last = stop.ModelStop().ID()
	}
	return false
}
