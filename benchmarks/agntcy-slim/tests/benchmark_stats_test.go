package tests

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

const (
	confidenceIntervalAlpha              = 0.05
	defaultMinimumConfidenceIntervalRuns = 20
)

type sampleStats struct {
	Count    int
	Mean     float64
	Variance float64
	StdDev   float64
	CILow    float64
	CIHigh   float64
}

func senderMPSValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.SenderMPS)
	}
	return values
}

func senderMeanLatencyMSValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.SenderMeanLatencyMS)
	}
	return values
}

func senderP50LatencyMSValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.SenderP50LatencyMS)
	}
	return values
}

func senderP99LatencyMSValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.SenderP99LatencyMS)
	}
	return values
}

func observedMPSValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, benchmarkObservedMPS(row))
	}
	return values
}

func senderMessageValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, float64(row.SenderTotalMessages))
	}
	return values
}

func observedMessageValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, benchmarkObservedMessages(row))
	}
	return values
}

func senderCPUPercentValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.SenderCPUPercent)
	}
	return values
}

func responderCPUPercentValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.ResponderCPUPercent)
	}
	return values
}

func nodeCPUPercentValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.NodeCPUPercent)
	}
	return values
}

func totalCPUPercentValues(rows []benchmarkRunResult) []float64 {
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.TotalCPUPercent)
	}
	return values
}

func computeSampleStats(values []float64) sampleStats {
	count := len(values)
	if count == 0 {
		return sampleStats{}
	}

	mean := stat.Mean(values, nil)
	if count == 1 {
		return sampleStats{
			Count:    1,
			Mean:     mean,
			Variance: 0,
			StdDev:   0,
			CILow:    mean,
			CIHigh:   mean,
		}
	}

	variance := 0.0
	stddev := stat.StdDev(values, nil)
	variance = stat.Variance(values, nil)
	standardError := stddev / math.Sqrt(float64(count))
	tDist := distuv.StudentsT{Mu: mean, Sigma: standardError, Nu: float64(count - 1)}
	tailProbability := confidenceIntervalAlpha / 2
	ciLow := tDist.Quantile(tailProbability)
	ciHigh := tDist.Quantile(1 - tailProbability)

	return sampleStats{
		Count:    count,
		Mean:     mean,
		Variance: variance,
		StdDev:   stddev,
		CILow:    math.Max(0, ciLow),
		CIHigh:   ciHigh,
	}
}

func configuredMinimumConfidenceIntervalRuns() int {
	raw := strings.TrimSpace(os.Getenv("MIN_CONFIDENCE_INTERVAL_RUNS"))
	if raw == "" {
		return defaultMinimumConfidenceIntervalRuns
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < 1 {
		return defaultMinimumConfidenceIntervalRuns
	}
	return parsed
}

func confidenceIntervalAvailable(count int) bool {
	return count >= configuredMinimumConfidenceIntervalRuns()
}

func unavailableConfidenceIntervalLabel() string {
	return fmt.Sprintf("n/a (need >=%d runs)", configuredMinimumConfidenceIntervalRuns())
}

func formatCI(stats sampleStats) string {
	if stats.Count == 0 {
		return "-"
	}
	if !confidenceIntervalAvailable(stats.Count) {
		return unavailableConfidenceIntervalLabel()
	}
	return fmt.Sprintf("[%s, %s]", formatFloat(stats.CILow), formatFloat(stats.CIHigh))
}

func formatCIBounds(count int, low float64, high float64) string {
	if count == 0 {
		return "-"
	}
	if !confidenceIntervalAvailable(count) {
		return unavailableConfidenceIntervalLabel()
	}
	return fmt.Sprintf("[%s, %s]", formatFloat(low), formatFloat(high))
}

func formatCISentence(count int, low float64, high float64) string {
	if !confidenceIntervalAvailable(count) {
		return fmt.Sprintf("95%% CI unavailable (need >=%d runs)", configuredMinimumConfidenceIntervalRuns())
	}
	return fmt.Sprintf("95%% CI [%s, %s]", formatFloat(low), formatFloat(high))
}
