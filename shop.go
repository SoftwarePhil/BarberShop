package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	NUM_OF_BARS   = 3
	NUM_OF_CHAIRS = 5
	HAIR_CUT_TIME = 5000
	NEW_CUSTOMER  = 1000
)

type customerChannel chan *customer

type shop struct {
	waitingRoom []*customer
	full        bool
	customerChannel
}

type barber struct {
	id            int
	time          int
	isCuttingHair bool
	customerChannel
}

type customer struct {
	id int
}

func main() {
	runShop()
}

func newCustomer(id int) *customer {
	a := new(customer)
	a.id = id

	return a
}

func (s *shop) waitingRoomStart(numOfcustomers int, wg *sync.WaitGroup, cc customerChannel) {
	defer wg.Done()
	count := 0
	for i := 0; i < numOfcustomers; i++ {
		newCust := newCustomer(count)
		count++
		select {

		case cc <- newCust:
			fmt.Printf("%s, %d, %s", "Customer ", newCust.id, " has sat down!\n")

		default:
			fmt.Println("customer " + strconv.Itoa(newCust.id) + " had to leave")
		}
		time.Sleep(NEW_CUSTOMER * time.Millisecond)

	}
}

func barberSelect(b []*barber, cc customerChannel) {
	for {
		for i := range b {
			if b[i].isCuttingHair == false {
				newHairCut := <-cc
				b[i].customerChannel <- newHairCut
			}
		}
	}
}

func createBarber(id, time int) *barber {
	a := new(barber)
	a.id = id
	a.time = time
	a.isCuttingHair = false
	a.customerChannel = make(customerChannel)
	return a
}

func (b *barber) cutHair(channel chan *customer) {
	for {
		select {
		case a := <-channel:
			b.isCuttingHair = true
			time.Sleep(HAIR_CUT_TIME * time.Millisecond)
			fmt.Print(strconv.Itoa(b.id) + " haircut for customer " + strconv.Itoa(a.id) + " is done" + "\n")

		default:
			fmt.Println("Barber is sleeping " + strconv.Itoa(b.id))
			time.Sleep(HAIR_CUT_TIME / 4 * time.Millisecond)
		}

		b.isCuttingHair = false
	}
}

func runShop() {
	var wg sync.WaitGroup
	cc := make(chan *customer, NUM_OF_CHAIRS)
	myShop := shop{make([]*customer, NUM_OF_CHAIRS, NUM_OF_CHAIRS), false, cc}
	wg.Add(1)
	go myShop.waitingRoomStart(20, &wg, cc)

	barbers := make([]*barber, NUM_OF_BARS, NUM_OF_BARS)

	for i := range barbers {
		barbers[i] = createBarber(i, HAIR_CUT_TIME)
		go barbers[i].cutHair(cc)
	}

	go barberSelect(barbers, cc)
	wg.Wait()
}
