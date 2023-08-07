package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NUMBER_OF_ORDERS = 10

var ordersMade, ordersFailed, total int

type order struct {
	id int
	message string
	success bool
}

type producer struct {
	data chan order
	quit chan chan error
}

func (p *producer) Close() error {

	ch := make(chan error)
	p.quit <- ch
	
	return <-ch
}

func buildOrder(orderNumber int) *order {

	var message string
	success := false
	orderNumber++
	
	if orderNumber <= NUMBER_OF_ORDERS {
		
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d! \n", orderNumber)
		
		generator := rand.Intn(12) + 1
		
		if generator < 5 { ordersFailed++ } else { ordersMade++ }
		total++
		
		fmt.Printf("Making order #%d. It will take %d seconds....\n", orderNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)
		
		if generator <= 2 {
			message = fmt.Sprintf("*** We ran out of service #%d!", orderNumber)
		} else if generator <= 4 {
			message = fmt.Sprintf("*** The order #%d failed", orderNumber)
		} else {
			success = true
			message = fmt.Sprintf("Order #%d is ready!", orderNumber)
		}
		
		return &order {
			id: orderNumber,
			message: message,
			success: success,
		}
	}
	
	return &order {
		id: orderNumber,
		success: success,
	}
}

func builder(producer *producer) {

	var i int = 0
	
	for {
		currentOrder := buildOrder(i)
    i = currentOrder.id

    select {

    case producer.data <- *currentOrder:
    case quitChan := <-producer.quit:
      
      close(quitChan)
      close(producer.data)

      return
    }
	}
}

func consumer(producer *producer) {

  for order := range producer.data {

    if order.id <= NUMBER_OF_ORDERS {

      if order.success {

        color.Green(order.message)
        color.Green("Order #%d was successful", order.id)
      } else {

        color.Red("Something went wroing with order #%d", order.id)
      }
    } else {

      color.Cyan("Stop generating orders")
      err := producer.Close()

      if err != nil {

        color.Red("Error closing channel")
      }
    }
  }

  
}



func main() {

	rand.Seed(time.Now().UnixNano())

	color.Cyan("Generating orders")

	producerInstance := producer{
		data: make(chan order),
		quit: make(chan chan error),
	}
	
	go builder(&producerInstance)

  consumer(&producerInstance)

  color.Cyan("A total of %d orders were created. %d failed and %d were success", total, ordersFailed, ordersMade)
}
