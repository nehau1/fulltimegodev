package main

import (
	"fmt"
	"time"

	"github.com/Stiffjobs/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyCalc prometheus.Histogram
	reqLatencyAgg  prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) Aggregator {
	fmt.Print("hello")
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})

	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calculator_request_counter",
		Name:      "calculate",
	})

	reqAggLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqCalcLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "calculator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:           next,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyCalc: reqCalcLatency,
		reqLatencyAgg:  reqAggLatency,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"err":  err,
			}).Error("CalculateInvoice")
		} else {
			logrus.WithFields(logrus.Fields{
				"took":          time.Since(start),
				"totalAmount":   inv.TotalAmount,
				"totalDistance": inv.TotalDistance,
				"obuID":         inv.OBUID,
			}).Info("CalculateInvoice")
		}
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)

	return
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}
