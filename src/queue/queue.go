package queue

import (
	"driver"
	"fmt"
	"global"
)

// TODO:
// -- updated order bool chan må lages
// -- updated_order_bool_list_chan
// -- Opdatere delete order til å oppføre seg annerledes om du er master eller slave

var Internal_order_list [global.NUM_INTERNAL_ORDERS]Order
var External_order_list [global.NUM_GLOBAL_ORDERS]Order
var Global_order_list [global.NUM_GLOBAL_ORDERS]Order

var Elevators_online [3]Elev_info

const (
	Inactive = iota
	Active
	Assigned
	Ready
	//Executing
	Finished
)

type Order struct {
	Button      global.Button_t
	Floor       global.Floor_t
	Order_state int
	Assigned_to int
}

type Elev_info struct {
	Elev_ip         int
	Elev_last_floor global.Floor_t
	Elev_dir        global.Motor_direction_t
	Elev_state      int
}

var My_info Elev_info

// Go functions for updating channels //
//Disse fungerer bra:
func Bool_to_new_order_channel(value bool, new_order_bool_chan chan bool) {
	new_order_bool_chan <- value
}

func Bool_to_new_global_order_chan(value bool, new_global_bool_chan chan bool) {
	new_global_bool_chan <- value
}

//Disse bør brukes/er blitt endret på:
/*func Bool_to_updated_order_channel(value bool, updated_order_bool_chan chan bool) {
	updated_order_bool_chan <- value
}

func Bool_to_updated_order_list_channel(value bool, updated_order_bool_list_chan chan bool) {
	updated_order_bool_list_chan <- value
}*/

func Order_to_update_order_chan(order Order, update_order_chan chan Order) {
	update_order_chan <- order
}
func Order_to_new_order_chan(order Order, new_order_chan chan Order) {
	new_order_chan <- order
}

func Order_handler(new_order_bool_chan chan bool, new_order_chan chan Order, update_order_chan chan Order, new_global_order_bool_chan chan bool) {
	fmt.Println("Running: Order handler")
	for {
		select {
		case catch_new_order := <-new_order_chan:
			new_order := catch_new_order
			fmt.Println("Case: A new button is detected. The Order_handler case new order was triggered")
			if new_order.Button == global.BUTTON_COMMAND {
				Add_new_internal_order(new_order, new_order_bool_chan)
				fmt.Println("New order added to Internal order list. The list is now: ", Internal_order_list)

			} else {
				Add_new_external_order(new_order, new_order_bool_chan, new_order_chan, new_global_order_bool_chan)
				fmt.Println("New order added to External order list. The list is now: ", External_order_list)

			}
		case catch_update_order := <-update_order_chan:
			fmt.Println("Case: A updated order is detected. The Order_handler case update order was triggered")
			update_order := catch_update_order

			Update_state(update_order)

			if update_order.Order_state == Finished {
				Delete_order(update_order) // MÅ TILPASSES OM DU ER MASTER ELLER SLAVE

			}

			// fra gomp:
			/*if network.Is_master {
				if update_order.Order_state == Finished {
					for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
						if Internal_order_list[i].Floor == update_order.Floor{
							update_order = Internal_order_list[i]
							driver.Set_button_lamp(update_order.Button, update_order.Floor, global.OFF)
							Delete_order(update_order)
						}
					}
					for i := 0; i < global.NUM_EXTERNAL_ORDERS; i++{
						if External_order_list[i].Floor == update_order.Floor{
							update_order = External_order_list[i]
							driver.Set_button_lamp(update_order.Button, update_order.Floor, global.OFF)
							Delete_order(update_order)
						}
					}
				}
			}*/

		}
	}
}

//-------- Case tar i mot ny ordre. Door_open må sende dette på denne channelen. Update _ state func er ikke tilpassa

func Update_state(update_order Order) {
	fmt.Println("Running Update_state")
	fmt.Println("My update_order: ", update_order)

	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if update_order.Button == Internal_order_list[i].Button && update_order.Floor == Internal_order_list[i].Floor {
			fmt.Println("Update internal order loop.")
			Internal_order_list[i].Order_state = update_order.Order_state

		}
	}
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if update_order.Button == External_order_list[i].Button && update_order.Floor == External_order_list[i].Floor {
			fmt.Println("Update external order loop.")
			External_order_list[i].Order_state = update_order.Order_state
		}
	}
	// THIS IS NEW
	if global.Is_master == true || global.Lost_network == true {
		for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
			if update_order.Button == Global_order_list[i].Button && update_order.Floor == Global_order_list[i].Floor {
				fmt.Println("Update external order loop.")
				Global_order_list[i].Order_state = update_order.Order_state
			}
		}

	}
	//END OF NEW

	if update_order.Order_state == Finished {
		for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
			if update_order.Floor == Internal_order_list[i].Floor {
				fmt.Println("Setting also the internal to finished.")
				Internal_order_list[i].Order_state = update_order.Order_state

			}
		}
		for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
			if update_order.Floor == External_order_list[i].Floor {
				fmt.Println("Setting also the external to finished")
				External_order_list[i].Order_state = update_order.Order_state
			}
		}

	}

}

func Add_new_internal_order(new_order Order, new_order_bool_chan chan bool) {
	new_order_floor := new_order.Floor
	fmt.Println(Internal_order_list)

	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if Internal_order_list[i].Order_state == Inactive || Internal_order_list[i].Order_state == Finished {
			Internal_order_list[i] = new_order
			fmt.Println("New internal order was added!", Internal_order_list[i])
			go Bool_to_new_order_channel(true, new_order_bool_chan)
			driver.Set_button_lamp(new_order.Button, new_order.Floor, global.ON)
			break
		}
		if Internal_order_list[i].Floor == new_order_floor {

			fmt.Println("The order is already in the internal order list.", Internal_order_list[i], "new:", new_order)
			//go Bool_to_new_order_channel(true, new_order_bool_chan)
			break
		}
		if i == global.NUM_INTERNAL_ORDERS-1 {
			fmt.Println("Error: No internal order was added.")
		}
	}
}

func Add_new_external_order(new_order Order, new_order_bool_chan chan bool, new_order_chan chan Order, new_global_order_bool_chan chan bool) {
	fmt.Println("Adder new order fra master")
	fmt.Println("I am the master", global.Is_master)
	new_order_floor := new_order.Floor
	new_order_button := new_order.Button

	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if External_order_list[i].Order_state == Inactive || External_order_list[i].Order_state == Finished {
			External_order_list[i] = new_order
			fmt.Println("New external order was added!")
			fmt.Println(My_info.Elev_ip, "this is my IP yo")
			if new_order.Assigned_to == My_info.Elev_ip {
				fmt.Println("The order was for meeee")
				go Bool_to_new_order_channel(true, new_order_bool_chan)
			}
			if global.Is_master == true && new_order.Assigned_to == 0 {
				fmt.Println("just for fun. delegating order.", global.Is_master)
				new_order = Delegate_order(new_order)
				Add_new_global_order(new_order, new_order_bool_chan, new_order_chan, new_global_order_bool_chan)

			}
			break
		}
		if External_order_list[i].Floor == new_order_floor && External_order_list[i].Button == new_order_button && External_order_list[i].Assigned_to == 0 {
			if new_order.Assigned_to == My_info.Elev_ip {
				External_order_list[i].Assigned_to = new_order.Assigned_to
				Update_state(External_order_list[i])
				// NEW
				go Bool_to_new_order_channel(true, new_order_bool_chan)
				//
			}
		} else if External_order_list[i].Floor == new_order_floor && External_order_list[i].Button == new_order_button {
			fmt.Println("The order is already in the external order list.", External_order_list[i])
			//go Order_to_new_order_chan(new_order, new_order_chan)
			break
		}
		if i == global.NUM_GLOBAL_ORDERS-1 {
			fmt.Println("Error: No external order was added.")
		}
	}

}

//--
func Add_new_global_order(new_order Order, new_order_bool_chan chan bool, new_order_chan chan Order, new_global_order_bool_chan chan bool) {
	new_order_floor := new_order.Floor
	new_order_button := new_order.Button

	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if Global_order_list[i].Order_state == Inactive {
			Global_order_list[i] = new_order
			fmt.Println("New external order was added!")
			//--
			go Bool_to_new_global_order_chan(true, new_global_order_bool_chan)
			//go Bool_to_new_order_channel(true, new_order_bool_chan)
			driver.Set_button_lamp(new_order.Button, new_order.Floor, global.ON)
			break
		}
		if Global_order_list[i].Floor == new_order_floor && Global_order_list[i].Button == new_order_button {
			fmt.Println("The order is already in the global order list.", Global_order_list[i])
			//go Order_to_new_order_chan(new_order, new_order_chan)
			break
		}
		if i == global.NUM_GLOBAL_ORDERS-1 {
			fmt.Println("Error: No external order was added.")
		}
	}

}

// THIS IS NEW
func Delete_order(updated_order Order) {
	fmt.Println("Inside Delete_order func")
	Delete_external_order()
	Delete_internal_order()
	if global.Is_master == true {
		Delete_global_order()
	}

}

func Delete_global_order() {
	clean_order := Make_new_order(global.BUTTON_UP, global.FLOOR_1, Inactive, 0)

	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if Global_order_list[i].Order_state == Finished {
			fmt.Println("An external order is marked finished.")
			for j := i; j < global.NUM_ORDERS; j++ {
				if j < global.NUM_GLOBAL_ORDERS-1 {
					Global_order_list[j] = Global_order_list[j+1]
				} else if j == global.NUM_GLOBAL_ORDERS-1 {
					Global_order_list[j] = clean_order
				}
			}
		}
	}
	fmt.Println("External order deleted, external order list is now: ", External_order_list)

}

//END OF NEW

func Delete_external_order() {
	clean_order := Make_new_order(global.BUTTON_UP, global.FLOOR_1, Inactive, global.NONE)

	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if External_order_list[i].Order_state == Finished {
			fmt.Println("An external order is marked finished.")
			for j := i; j < global.NUM_ORDERS; j++ {
				if j < global.NUM_GLOBAL_ORDERS-1 {
					External_order_list[j] = External_order_list[j+1]
				} else if j == global.NUM_GLOBAL_ORDERS-1 {
					External_order_list[j] = clean_order
				}
			}
		}
	}
	fmt.Println("External order deleted, external order list is now: ", External_order_list)

}

func Delete_internal_order() {
	clean_order := Make_new_order(global.BUTTON_UP, global.FLOOR_1, Inactive, global.NONE)

	for i := 0; i < global.NUM_INTERNAL_ORDERS; i++ {
		if Internal_order_list[i].Order_state == Finished {
			fmt.Println("An internal order is marked finished.")
			for j := i; j < global.NUM_INTERNAL_ORDERS; j++ {
				if j < global.NUM_INTERNAL_ORDERS-1 {
					Internal_order_list[j] = Internal_order_list[j+1]
				} else if j == global.NUM_INTERNAL_ORDERS-1 {
					Internal_order_list[j] = clean_order
				}
			}
		}
	}
	fmt.Println("Internal order deleted, internal order list is now: ", Internal_order_list)
}

func Make_new_order(button global.Button_t, floor global.Floor_t, order_state int, assigned_to int) Order {
	var new_order Order

	new_order.Button = button
	new_order.Floor = floor
	new_order.Order_state = order_state
	new_order.Assigned_to = assigned_to

	return new_order
}
