package songphapadonationconsumer

import (
	"container/heap"
	"log"
	"github.com/chanonchanpiwat/challenge.git/songphapadonator"
)

type Donator struct {
	Name   string
	Amount int
	topN   int
}

type Result struct {
	SuccessAmount int
	FailedAmount  int
	TopNDonator   []*Donator
	DonationCount int
}

func Consumer(songPhaPaDonationChannel []*songphapadonator.ChargeResponse, topN int) Result {
	var result Result
	h := &DonatorHeap{}
	heap.Init(h)

	for _,rec := range songPhaPaDonationChannel {
		if rec.Charge == nil {
			log.Println("unable to call charge API with error:", rec.Error.Error())
			continue
		}

		result.DonationCount += 1
		amountTransfer := int(rec.Charge.Amount)
		if rec.Charge == nil || !rec.Charge.Paid {
			result.FailedAmount += amountTransfer
		} else {
			result.SuccessAmount += amountTransfer
			donator := Donator{Name: rec.Charge.Card.Name,
				Amount: int(rec.Charge.Amount),
				topN:   topN,
			}
			heap.Push(h, &donator)
		}
	}

	result.TopNDonator = *h
	return result
}


