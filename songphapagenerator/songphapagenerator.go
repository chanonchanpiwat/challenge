package songphapagenerator

import (
	"encoding/csv"
	"io"
	"strconv"
	"github.com/chanonchanpiwat/challenge.git/logger"
	"time"
)

type SongPahPa struct {
	Name           string
	AmountSubunits int64
	CCNumber       string
	CVV            string
	ExpMonth       time.Month
	ExpYear        int
}

func SongPahPaParser(record []string) *SongPahPa {

	amount, err := strconv.ParseInt(record[1], 10, 64)
	logger.LogAndExit(err)

	month, err := strconv.Atoi(record[4])
	logger.LogAndExit(err)

	year, err := strconv.Atoi(record[5])
	logger.LogAndExit(err)

	return &SongPahPa{
		Name:           record[0],
		AmountSubunits: amount,
		CCNumber:       record[2],
		CVV:            record[3],
		ExpMonth:       time.Month(month),
		ExpYear:        year,
	}
}

func GenerateSongPhaPaChannel(done <-chan interface{}, csvReader *csv.Reader, sleepMillisecond int) <-chan *SongPahPa {
	ch := make(chan *SongPahPa)
	csvReader.Read()
	go func() {
		defer close(ch)
		for {
			record, err := csvReader.Read()
			if err == io.EOF || len(record) == 0 {
				return
			}

			time.Sleep(time.Duration(sleepMillisecond) * time.Millisecond)

			select {
			case <-done:
				return
			case ch <- SongPahPaParser(record):
			}

		}
	}()

	return ch
}
