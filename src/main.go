package main

import (
	/*
		"driver"
		"global"
		"network"
		"queue"
		"fsm"
		"buttonhandler"*/
	"driver"
	"fmt"
	"fsm"
	"network"
	"ordermanager"
	"queue"
	//"time"
)

func main() {
	driver.Elevator_init()
	driver.Elevator_init()

	new_order_chan := make(chan queue.Order, 10) //Er buffra til 10 for da får alle mulig ebestillinger plass
	updated_order_chan := make(chan queue.Order, 10)
	updated_order_bool_chan := make(chan bool)
	new_order_bool_chan := make(chan bool)

	go fsm.State_handler(new_order_bool_chan, updated_order_bool_chan, updated_order_chan)
	go queue.Order_handler(new_order_bool_chan, new_order_chan, updated_order_chan)
	go ordermanager.Button_handler(new_order_chan)

	// test network
	go network.Network_handler()
	//time.Sleep(4*time.Second)
	go network.Network_sender(new_order_bool_chan, new_order_chan)
	go network.Network_receiver(new_order_bool_chan, new_order_chan)

	fmt.Println("I am ready")

	// Make channels
	/*
		new_order_chan := make(chan queue.Order, 10)
		updated_order_chan := make(chan queue.Order, 10)
		external_order_list_chan := make(chan [global.NUM_EXTERNAL_ORDERS]queue.Order)
		internal_order_list_chan := make(chan [global.NUM_INTERNAL_ORDERS]queue.Order)
		new_order_bool_chan := make(chan bool)
		is_master_chan := make(chan bool)
	*/
	// Run all processes
	/*
		go network.Network_handler(is_master_chan)
		go fsm.State_handler(new_order_bool_chan, updated_order_chan, external_order_list_chan, internal_order_list_chan)
		go queue.Order_handler(new_order_bool_chan, new_order_chan, updated_order_chan, external_order_list_chan, internal_order_list_chan)
		go buttonhandler.Button_handler(new_order_chan)
	*/
	// Keep on running
	for {
	}
}
