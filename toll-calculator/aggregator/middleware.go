package main

import (
	"time"

	"github.com/Stiffjobs/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

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
