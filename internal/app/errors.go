package app

import (
	"errors"
	"log"
)

var (
	ErrYouShallNotPass  = errors.New("YouShallNotPass")
	ErrNotOpenYet       = errors.New("NotOpenYet")
	ErrPlaceIsBusy      = errors.New("PlaceIsBusy")
	ErrClientUnknown    = errors.New("ClientUnknown")
	ErrICanWaitNoLonger = errors.New("ICanWaitNoLonger!")
)

func inputFormatError(line string, err error) {
	log.Fatalf("Error at line: %s\nError: %s", line, err.Error())
}
