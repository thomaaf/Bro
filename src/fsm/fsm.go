package fsm

import (
	"driver"
	"fmt"
	"global"
	"network"
	"os"
	"queue"
	"time"
)

// TODO:
// -- Legge inn en go func med timer som sjekker om man er stuck
// -- Event idle reagerer ikke på updated orders, så viktig at heisene pusher new_order_chan ikke bare ved knappetrykk, men også hvis master assigner en ordre til dem
// -- Event idle sjekker nå external list. MÅ sjekke global list. Skal kun gjøre ting når de er blitt assigna til deg i global list
// -- Sjekker ikke om man skal stoppe når man kjører forbi etasjer

// Maries todo:
// - sette endring av states inn i event funksjonene
// - ha endring av states som global variabel, kan slette lokal (tror ikke vi bruker den noen plass, dobbeltsjekk)
// - set button lamp off bør settes inn i update state funksjonen
// - atm: når to bestillinger i samme floor så blir den stuck i door open - bare et par ganger

// Elevator states
const (
	Idle int = iota
	Moving
	Door_open
	Stuck
)

var Elev_state int

//var Dir global.Motor_direction_t /// la til global variabel Dir (direction)
var current_order queue.Order
var Empty_list bool

func State_handler(new_order_bool_chan chan bool, updated_order_bool_chan chan bool, update_order_chan chan queue.Order, new_global_order_bool_chan chan bool) {
	fmt.Println("Running: State handler")
	Elev_state = Idle
	queue.My_info.Elev_state = Elev_state
	for i := 0; i < global.Num_elev_online; i++ {
		if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
			queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
		}
	}

	for {
		switch Elev_state {
		case Idle:
			event_idle(new_order_bool_chan, new_global_order_bool_chan)
			Elev_state = Moving
			queue.My_info.Elev_state = Elev_state
			for i := 0; i < global.Num_elev_online; i++ {
				if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
					queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
				}
			}
			//elev_state = Moving
		case Moving:
			event_moving(update_order_chan)
			//Elev_state = Door_open
			//elev_state = Door_open
		case Door_open:
			event_door_open(update_order_chan)
			//Elev_state = Idle
			//elev_state = Idle
			time.Sleep(1 * time.Second)
		case Stuck:
			os.Exit(0)

		}
	}
}

func event_idle(new_order_bool_chan chan bool, new_global_order_bool_chan chan bool) {
	fmt.Println("Running event: Idle.")
	order_exist := false
	fmt.Println("Current order is before", current_order)
	fmt.Println("Num elev online in Idle hehe: ", global.Num_elev_online)

	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if queue.Internal_order_list[i].Order_state != queue.Inactive {
			if queue.Internal_order_list[i].Order_state != queue.Finished {
				current_order = queue.Internal_order_list[i]
				fmt.Println("The current order is internal---------------------------------------")
				order_exist = true
				Empty_list = false
				break
			}
		}

	}

	// THIS IS NEW
	if global.Lost_network == true {
		for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
			if queue.Global_order_list[i].Order_state != queue.Inactive {
				if queue.Global_order_list[i].Order_state != queue.Finished {
					fmt.Println(queue.Global_order_list[i].Assigned_to == network.Local_ip)
					current_order = queue.Global_order_list[i]
					order_exist = true
					Empty_list = false
					break
				}
			}
		}
	}

	// END NEW
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if queue.External_order_list[i].Order_state != queue.Inactive {
			if queue.External_order_list[i].Order_state != queue.Finished {
				if queue.External_order_list[i].Assigned_to == network.Local_ip || global.Num_elev_online == 0 {
					fmt.Println(queue.External_order_list[i].Assigned_to == network.Local_ip)
					current_order = queue.External_order_list[i]
					order_exist = true
					Empty_list = false
					break
				}
			}
		}
	}
	fmt.Println("Current order is after", current_order)
	Empty_list = true
	if order_exist == false {
		select {
		case <-new_order_bool_chan:
			var this_order queue.Order
			new_order_Assigned_to_me := false
			//fmt.Println("Got new order bool ", catch_new_order_bool, " in Idle.")
			//fmt.Println("Now checking for orders that needs to be done inside the select case in event_idle")
			fmt.Println("Current order is", current_order)
			for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
				if queue.Internal_order_list[i].Order_state != queue.Inactive {
					if queue.Internal_order_list[i].Order_state != queue.Finished {

						fmt.Println("The current order is internal---------------------------------------")
						//current_order = queue.Internal_order_list[i]
						this_order = queue.Internal_order_list[i]
						current_order = this_order
						fmt.Println("This order is: ", this_order)
						new_order_Assigned_to_me = true

						break
					}
				}
				if new_order_Assigned_to_me == true {
					Empty_list = false
					break
				}

			}
			for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
				if queue.External_order_list[i].Order_state != queue.Inactive {
					if queue.External_order_list[i].Order_state != queue.Finished {
						if queue.External_order_list[i].Assigned_to == network.Local_ip || global.Num_elev_online == 0 {
							fmt.Println(queue.External_order_list[i].Assigned_to, "is equal to", network.Local_ip, "-----------------", "network.Num_elev_onlibe", global.Num_elev_online)
							current_order = queue.External_order_list[i]
							new_order_Assigned_to_me = true
						}
					}
				}

				// THIS IS NEW
				if global.Lost_network == true || global.Num_elev_online == 1 {
					for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
						if queue.Global_order_list[i].Order_state != queue.Inactive {
							if queue.Global_order_list[i].Order_state != queue.Finished {
								fmt.Println("Taking all the global orders since i'm alone")
								current_order = queue.Global_order_list[i]
								order_exist = true
								Empty_list = false
								break
							}
						}
					}
				}

			}
			// END NEW

			if new_order_Assigned_to_me == true {
				fmt.Println("The order is assigned to me, I'll take it.")
				Empty_list = false
				break
			}

		case <-new_global_order_bool_chan:
			var this_order queue.Order
			new_order_Assigned_to_me := false
			//fmt.Println("Got new order bool ", catch_new_order_bool, " in Idle.")
			//fmt.Println("Now checking for orders that needs to be done inside the select case in event_idle")
			fmt.Println("Current order is", current_order)
			for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
				if queue.Internal_order_list[i].Order_state != queue.Inactive {
					if queue.Internal_order_list[i].Order_state != queue.Finished {

						fmt.Println("The current order is internal---------------------------------------")
						//current_order = queue.Internal_order_list[i]
						this_order = queue.Internal_order_list[i]
						current_order = this_order
						fmt.Println("This order is: ", this_order)
						new_order_Assigned_to_me = true

						break
					}
				}
				if new_order_Assigned_to_me == true {
					Empty_list = false
					break
				}

			}
			for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
				if queue.External_order_list[i].Order_state != queue.Inactive {
					if queue.External_order_list[i].Order_state != queue.Finished {
						if queue.External_order_list[i].Assigned_to == network.Local_ip || global.Num_elev_online == 0 {
							fmt.Println(queue.External_order_list[i].Assigned_to, "is equal to", network.Local_ip, "-----------------", "network.Num_elev_onlibe", global.Num_elev_online)
							current_order = queue.External_order_list[i]
							new_order_Assigned_to_me = true
						}
					}
				}

				// THIS IS NEW
				if global.Lost_network == true || global.Num_elev_online == 1 {
					for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
						if queue.Global_order_list[i].Order_state != queue.Inactive {
							if queue.Global_order_list[i].Order_state != queue.Finished {
								fmt.Println("Taking all the global orders since i'm alone")
								current_order = queue.Global_order_list[i]
								order_exist = true
								Empty_list = false
								break
							}
						}
					}
				}

			}
			// END NEW

			if new_order_Assigned_to_me == true {
				fmt.Println("The order is assigned to me, I'll take it.")
				Empty_list = false
				break
			}

		}
	}
}

// Elev_state = Moving // <- sette global state inne i funksjonen

func event_moving(update_order_chan chan queue.Order) {
	fmt.Println("Running event: Moving.")
	elevator_to_floor(update_order_chan)
	/*if Elev_state != Idle {
		Elev_state = Door_open
	}*/

}

func event_door_open(update_order_chan chan queue.Order) {
	fmt.Println("Running event: Door open.")

	// Open door
	driver.Open_door()
	fmt.Println("Door opened.")
	// - set button lamp off bør settes inn i update state funksjonen
	driver.Set_button_lamp(current_order.Button, current_order.Floor, global.OFF) //-- can be moved to before open door
	fmt.Println("Door open lamp set on.")

	// Set order state to finished
	current_order.Order_state = queue.Finished
	fmt.Println("Current order state set to finished.")
	// ---- hmhmhmhmmh
	go queue.Order_to_update_order_chan(current_order, update_order_chan)
	fmt.Println("Order sent on updated order chan.")

	Elev_state = Idle // <- sette global state inne i funksjonen
	queue.My_info.Elev_state = Elev_state
	for i := 0; i < global.Num_elev_online; i++ {
		if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
			queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
		}
	}
}

func elevator_to_floor(update_order_chan chan queue.Order) {
	// Check if the elevator is between two floors
	between_two_floors_timer := time.NewTimer(3 * time.Second)
	timeout_between_floors := false
	go func() {
		<-between_two_floors_timer.C
		timeout_between_floors = true
	}()
	for driver.Get_floor_sensor_signal() == -1 {
		if !timeout_between_floors {
			//Dir = global.DIR_UP
			queue.My_info.Elev_dir = global.DIR_UP
			for i := 0; i < global.Num_elev_online; i++ {
				if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
					queue.Elevators_online[i].Elev_dir = queue.My_info.Elev_dir
				}
			}
			driver.Set_motor_direction(global.DIR_UP)
		} else if timeout_between_floors {
			//Dir = global.DIR_DOWN
			queue.My_info.Elev_dir = global.DIR_DOWN
			for i := 0; i < global.Num_elev_online; i++ {
				if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
					queue.Elevators_online[i].Elev_dir = queue.My_info.Elev_dir
				}
			}
			driver.Set_motor_direction(global.DIR_DOWN)
		}
	}

	check_if_stuck_timer := time.NewTimer(15 * time.Second)
	timeout := false
	go func() {
		<-check_if_stuck_timer.C
		timeout = true
	}()

	// Go to desired floor
	current_floor_int := driver.Get_floor_sensor_signal()
	current_floor := driver.Floor_int_to_floor_t(current_floor_int)
	floor_int := driver.Floor_t_to_floor_int(current_order.Floor)
	fmt.Println("Current floor int: ", current_floor_int, ", floor int: ", floor_int, ", current floor: ", current_floor)

	if current_floor_int < floor_int {
		fmt.Println("Going up.")
		//Dir = global.DIR_UP
		queue.My_info.Elev_dir = global.DIR_UP
		for i := 0; i < global.Num_elev_online; i++ {
			if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
				queue.Elevators_online[i].Elev_dir = queue.My_info.Elev_dir
			}
		}
		driver.Set_motor_direction(global.DIR_UP)

		for driver.Get_floor_sensor_signal() != floor_int {
			current_floor = driver.Floor_int_to_floor_t(driver.Get_floor_sensor_signal())

			// When arriving at any floor, check for order
			if driver.Get_floor_sensor_signal() != -1 {
				//fmt.Println("Floor sensor signal is not equal to minus one.")
				this_floor := driver.Floor_int_to_floor_t(driver.Get_floor_sensor_signal())
				queue.My_info.Elev_last_floor = this_floor
				for i := 0; i < global.Num_elev_online; i++ {
					if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
						queue.Elevators_online[i].Elev_last_floor = queue.My_info.Elev_last_floor
					}
				}
				driver.Set_floor_indicator_lamp(this_floor)
				//pick_up_order_on_the_way(current_floor, order_list, updated_order_chan, current_order)
				//time.Sleep(10 * time.Millisecond)
				//is_order := stop_if_order_in_floor(current_floor, update_order_chan)
				//fmt.Println("Checking floor if any order")
				is_order := check_if_order_in_floor(this_floor)
				//fmt.Println("We have an order in this floor.")
				//fmt.Println("Is order: ", is_order)
				if is_order {
					//fmt.Println("Setting state to door open.")
					Elev_state = Door_open
					queue.My_info.Elev_state = Elev_state
					for i := 0; i < global.Num_elev_online; i++ {
						if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
							queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
						}
					}
					break
				}
				time.Sleep(10 * time.Millisecond)
				check_if_stuck_timer.Reset(15 * time.Second)
			} else if timeout {
				Elev_state = Stuck
				queue.My_info.Elev_state = Elev_state
				for i := 0; i < global.Num_elev_online; i++ {
					if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
						queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
					}
				}
				break
			}
		}

	} else if current_floor_int > floor_int {
		fmt.Println("Going down.")
		//Dir = global.DIR_DOWN
		queue.My_info.Elev_dir = global.DIR_DOWN
		for i := 0; i < global.Num_elev_online; i++ {
			if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
				queue.Elevators_online[i].Elev_dir = queue.My_info.Elev_dir
			}
		}
		driver.Set_motor_direction(global.DIR_DOWN)

		for driver.Get_floor_sensor_signal() != floor_int {
			current_floor = driver.Floor_int_to_floor_t(driver.Get_floor_sensor_signal())

			// When we arrive at any floor, check for order
			if driver.Get_floor_sensor_signal() != -1 {
				this_floor := driver.Floor_int_to_floor_t(driver.Get_floor_sensor_signal())
				queue.My_info.Elev_last_floor = this_floor
				for i := 0; i < global.Num_elev_online; i++ {
					if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
						queue.Elevators_online[i].Elev_last_floor = queue.My_info.Elev_last_floor
					}
				}
				driver.Set_floor_indicator_lamp(this_floor)
				//pick_up_order_on_the_way(current_floor, order_list, updated_order_chan, current_order)
				//time.Sleep(10 * time.Millisecond)
				//is_order := stop_if_order_in_floor(current_floor, update_order_chan)
				is_order := check_if_order_in_floor(this_floor)
				if is_order {
					Elev_state = Door_open
					queue.My_info.Elev_state = Elev_state
					for i := 0; i < global.Num_elev_online; i++ {
						if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
							queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
						}
					}
					break
				}
				time.Sleep(10 * time.Millisecond)
				check_if_stuck_timer.Reset(15 * time.Second)
			} else if timeout {
				Elev_state = Stuck
				queue.My_info.Elev_state = Elev_state
				for i := 0; i < global.Num_elev_online; i++ {
					if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
						queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
					}
				}
				break
			}
		}
	} else {
		if current_order.Order_state == queue.Inactive || current_order.Order_state == queue.Finished {
			Elev_state = Idle
			queue.My_info.Elev_state = Elev_state
			for i := 0; i < global.Num_elev_online; i++ {
				if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
					queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
				}
			}
		} else {
			Elev_state = Door_open
			queue.My_info.Elev_state = Elev_state
			for i := 0; i < global.Num_elev_online; i++ {
				if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
					queue.Elevators_online[i].Elev_state = queue.My_info.Elev_state
				}
			}
			// was idle before, but it should open the door if it's in the same floor
			// must delete the order somewhere, atm it's a loop
		}
	}

	// Stop when at desired floor
	//Dir = global.DIR_STOP
	queue.My_info.Elev_dir = global.DIR_STOP
	for i := 0; i < global.Num_elev_online; i++ {
		if queue.My_info.Elev_ip == queue.Elevators_online[i].Elev_ip {
			queue.Elevators_online[i].Elev_dir = queue.My_info.Elev_dir
		}
	}
	driver.Set_motor_direction(global.DIR_STOP)
}

/*func stop_if_order_in_floor(floor global.Floor_t, update_order_chan chan queue.Order) bool {
	is_order_in_floor := check_if_order_in_floor(floor)
	if is_order_in_floor {
		driver.Set_motor_direction(global.DIR_STOP)
		event_door_open(update_order_chan)
	}
	return is_order_in_floor
}*/

func check_if_order_in_floor(floor global.Floor_t) bool {
	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if queue.Internal_order_list[i].Floor == floor {
			if queue.Internal_order_list[i].Order_state != queue.Inactive && queue.Internal_order_list[i].Order_state != queue.Finished {
				current_order = queue.Internal_order_list[i]
				//fmt.Println("The current order is: ", current_order, " , and the queue element is: ", queue.Internal_order_list[i])
				return true
			}
		}
	}
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if queue.External_order_list[i].Floor == floor {
			if queue.External_order_list[i].Button == global.BUTTON_UP && queue.My_info.Elev_dir == global.DIR_DOWN {
				return false
			} else if queue.External_order_list[i].Button == global.BUTTON_DOWN && queue.My_info.Elev_dir == global.DIR_UP {
				return false
			}
			if queue.External_order_list[i].Order_state != queue.Inactive && queue.External_order_list[i].Order_state != queue.Finished {
				current_order = queue.External_order_list[i]
				return true
			}
		}
	}
	return false
}
