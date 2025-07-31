package main

import (
	"github.com/nextmv-io/nextroute"
)

func NewStopSequenceStopConstraint() *perStop {
	return &perStop{}
}

type perStop struct{}

func (c *perStop) EstimateIsViolated(
	move nextroute.SolutionMoveStops,
) (isViolated bool, stopPositionsHint nextroute.StopPositionsHint) {
	checker := newChecker()
	return checker.isMoveViolating(move), nextroute.NoPositionsHint()
}

func NewStopSequenceVehicleConstraint() *perVehicle {
	return &perVehicle{}
}

type perVehicle struct{}

func (c *perVehicle) EstimateIsViolated(
	move nextroute.SolutionMoveStops,
) (isViolated bool, stopPositionsHint nextroute.StopPositionsHint) {
	checker := newChecker()
	return checker.isMoveViolating(move), nextroute.NoPositionsHint()
}

func (c *perVehicle) DoesVehicleHaveViolations(
	vehicle nextroute.SolutionVehicle,
) bool {
	checker := newChecker()
	return checker.isVehicleViolating(vehicle)
}

func newChecker() *checker {
	return &checker{}
}

type checker struct{}

func (c *checker) isVehicleViolating(vehicle nextroute.SolutionVehicle) bool {
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

func (c *checker) isMoveViolating(move nextroute.SolutionMoveStops) bool {
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
