package app

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type (
	App struct {
		InputFilePath string
		state         AppState
	}

	AppState struct {
		PlacesQuantity  int
		OpenTime        int // in minutes
		CloseTime       int // in minutes
		HourCost        int
		ClientsAtPlaces map[string]int
		Places          map[int]*Place
		ClientsQueue    []string // FIFO queue
	}

	Place struct {
		IsBusy         bool
		ReservedAtTime int // in minutes
		BusyTotalTime  int // in minutes
		Revenue        int
	}
)

func (a *App) Run() {
	// Open input file
	f, err := os.Open(a.InputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Scan first 3 lines
	for i := 0; i < 3; i++ {
		scanner.Scan()
		line := scanner.Text()
		switch i {
		case 0: // places quantity
			a.state.PlacesQuantity, err = strconv.Atoi(line)
			if err != nil {
				inputFormatError(line, err)
			}
		case 1: // open/close time
			splitted := strings.Split(line, " ")
			a.state.OpenTime, err = parseTime(splitted[0])
			if err != nil {
				inputFormatError(line, err)
			}
			a.state.CloseTime, err = parseTime(splitted[1])
			if err != nil {
				inputFormatError(line, err)
			}
		case 2: // hour cost
			a.state.HourCost, err = strconv.Atoi(line)
			if err != nil {
				inputFormatError(line, err)
			}
		}
	}

	// State init
	a.state.ClientsAtPlaces = make(map[string]int)
	a.state.Places = make(map[int]*Place)
	for i := 1; i <= a.state.PlacesQuantity; i++ {
		a.state.Places[i] = &Place{}
	}

	// Printing open time
	fmt.Println(convertMinutesToStringTime(a.state.OpenTime))

	// Scan events
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		splitted := strings.Split(line, " ")
		eventTime, err := parseTime(splitted[0])
		if err != nil {
			inputFormatError(line, err)
		}
		eventID := splitted[1]
		clientName := splitted[2]

		switch eventID {
		case "1":
			if eventTime < a.state.OpenTime {
				a.generateError(eventTime, ErrNotOpenYet)
				break
			}
			_, inClub := a.state.ClientsAtPlaces[clientName]
			if inClub {
				a.generateError(eventTime, ErrYouShallNotPass)
				break
			}
			a.state.ClientsAtPlaces[clientName] = -1
		case "2":
			wantedPlaceNum, err := strconv.Atoi(splitted[3])
			if err != nil {
				inputFormatError(line, err)
			}
			currentPlaceNum, inClub := a.state.ClientsAtPlaces[clientName]
			if !inClub {
				a.generateError(eventTime, ErrClientUnknown)
				break
			}
			place := a.state.Places[wantedPlaceNum]
			if place.IsBusy {
				a.generateError(eventTime, ErrPlaceIsBusy)
				break
			}
			if currentPlaceNum != -1 { // case with place changing
				a.clientLeavePlace(eventTime, clientName, currentPlaceNum)
			}
			a.clientSatAtPlace(eventTime, clientName, wantedPlaceNum)
		case "3":
			isAnyPlaceAvailable := false
			for _, place := range a.state.Places {
				if !place.IsBusy {
					isAnyPlaceAvailable = true
					break
				}
			}
			if isAnyPlaceAvailable {
				a.generateError(eventTime, ErrICanWaitNoLonger)
				break
			}
			if len(a.state.ClientsQueue) >= a.state.PlacesQuantity {
				a.generateClientLeave(eventTime, clientName)
				break
			}
			a.state.ClientsQueue = append(a.state.ClientsQueue, clientName)
		case "4":
			placeNum, inClub := a.state.ClientsAtPlaces[clientName]
			if !inClub {
				a.generateError(eventTime, ErrClientUnknown)
				break
			}
			a.clientLeavePlace(eventTime, clientName, placeNum)

			if len(a.state.ClientsQueue) > 0 {
				queueClientName := a.state.ClientsQueue[0]      // peek
				a.state.ClientsQueue = a.state.ClientsQueue[1:] // deque

				a.clientSatAtPlace(eventTime, queueClientName, placeNum)
				a.generateClientSatAtPlace(eventTime, queueClientName, placeNum)
			}
		default:
			log.Fatalf("unknown event ID: %s", eventID)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error on scanning line: %s", err.Error())
	}

	// Late clients
	if len(a.state.ClientsAtPlaces) > 0 {
		var lateClients []string
		for clientName := range a.state.ClientsAtPlaces {
			lateClients = append(lateClients, clientName)
		}
		sort.Slice(lateClients, func(i, j int) bool { // sort clients alphabetically
			return lateClients[i] < lateClients[j]
		})
		for _, clientName := range lateClients {
			placeNum := a.state.ClientsAtPlaces[clientName]
			place := a.state.Places[placeNum]
			if place != nil {
				placeTime := a.state.CloseTime - place.ReservedAtTime
				place.BusyTotalTime += placeTime
				place.Revenue += a.calcPlaceRevenue(placeTime)
			}
			a.generateClientLeave(a.state.CloseTime, clientName)
		}
	}

	// Printing close time
	fmt.Println(convertMinutesToStringTime(a.state.CloseTime))

	places := make([]Place, a.state.PlacesQuantity)
	for num, place := range a.state.Places {
		places[num-1] = *place
	}
	for i, place := range places {
		fmt.Printf("%d %d %s\n", i+1, place.Revenue, convertMinutesToStringTime(place.BusyTotalTime))
	}
}

func (a *App) clientSatAtPlace(time int, clientName string, placeNum int) {
	a.state.ClientsAtPlaces[clientName] = placeNum
	place := a.state.Places[placeNum]
	place.IsBusy = true
	place.ReservedAtTime = time
}

func (a *App) clientLeavePlace(time int, clientName string, placeNum int) {
	if placeNum < 1 {
		log.Fatal("cant leave from place with number < 1")
	}
	delete(a.state.ClientsAtPlaces, clientName)
	place := a.state.Places[placeNum]
	place.IsBusy = false
	placeTime := time - place.ReservedAtTime
	place.BusyTotalTime += placeTime
	place.Revenue += a.calcPlaceRevenue(placeTime)
}

func (a *App) calcPlaceRevenue(placeTime int) int {
	return int(float64(a.state.HourCost) * math.Ceil(float64(placeTime)/60))
}

func (a *App) generateClientLeave(time int, name string) {
	fmt.Printf("%s 11 %s\n", convertMinutesToStringTime(time), name)
}

func (a *App) generateClientSatAtPlace(time int, name string, placeNum int) {
	fmt.Printf("%s 12 %s %d\n", convertMinutesToStringTime(time), name, placeNum)
}

func (a *App) generateError(time int, error error) {
	fmt.Printf("%s 13 %s\n", convertMinutesToStringTime(time), error.Error())
}
