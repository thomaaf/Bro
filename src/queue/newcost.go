package queue

import (
	//"driver"
	"fmt"
	"global"
)

func Delegate_order(new_order Order) Order {

	// Elevators_online[i].Elev_ip
	if global.Num_elev_online == 3 {
		global.Iter = global.Iter + 1
		if global.Iter > 2 {
			global.Iter = 0
		}
	} else if global.Num_elev_online == 2 {
		global.Iter = global.Iter + 1
		if global.Iter > 1 {
			global.Iter = 0
		}
	} else if global.Num_elev_online == 1 {
		global.Iter = 0
	}

	fmt.Println("New order is: ", new_order)
	fmt.Println("Iteration nr. : ", global.Iter)
	assigned_elevator_ip := Elevators_online[global.Iter].Elev_ip
	fmt.Println("The assigned to ip is: ", assigned_elevator_ip)

	//assigned_elevator := compare_cost(Elevators_online, new_order)
	//assigned_elevator_ip := assigned_elevator.Elev_ip

	new_order.Assigned_to = assigned_elevator_ip
	fmt.Println("New order is assigned to: ", new_order.Assigned_to)
	fmt.Println("New order in delegate is: ", new_order)
	return new_order
	//Add_new_global_order(new_order)
}
