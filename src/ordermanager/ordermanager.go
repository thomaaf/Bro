package ordermanager

import (
	"driver"
	"fmt"
	"global"
	"queue"
	"time"
)

func Button_handler(new_order_chan chan queue.Order) {
	fmt.Println("Running: Button handler. ")
	var order queue.Order

	for {
		if global.Lost_network == false {

			if driver.Get_button_signal(global.BUTTON_UP, global.FLOOR_1) != 0 {
				order = queue.Make_new_order(global.BUTTON_UP, global.FLOOR_1, queue.Active, 0)
				new_order_chan <- order
				time.Sleep(1 * time.Second)
			}
			if driver.Get_button_signal(global.BUTTON_UP, global.FLOOR_2) != 0 {
				order = queue.Make_new_order(global.BUTTON_UP, global.FLOOR_2, queue.Active, 0)
				new_order_chan <- order
				time.Sleep(1 * time.Second)

			}
			if driver.Get_button_signal(global.BUTTON_DOWN, global.FLOOR_2) != 0 {
				order = queue.Make_new_order(global.BUTTON_DOWN, global.FLOOR_2, queue.Active, 0)
				new_order_chan <- order
				time.Sleep(1 * time.Second)
			}
			if driver.Get_button_signal(global.BUTTON_UP, global.FLOOR_3) != 0 {
				order = queue.Make_new_order(global.BUTTON_UP, global.FLOOR_3, queue.Active, 0)
				new_order_chan <- order
				time.Sleep(1 * time.Second)
			}
			if driver.Get_button_signal(global.BUTTON_DOWN, global.FLOOR_3) != 0 {
				order = queue.Make_new_order(global.BUTTON_DOWN, global.FLOOR_3, queue.Active, 0)
				new_order_chan <- order
				time.Sleep(1 * time.Second)
			}
			if driver.Get_button_signal(global.BUTTON_DOWN, global.FLOOR_4) != 0 {
				order = queue.Make_new_order(global.BUTTON_DOWN, global.FLOOR_4, queue.Active, 0)
				new_order_chan <- order
				time.Sleep(1 * time.Second)
			}
		}
		if driver.Get_button_signal(global.BUTTON_COMMAND, global.FLOOR_1) != 0 {
			order = queue.Make_new_order(global.BUTTON_COMMAND, global.FLOOR_1, queue.Active, 0)
			new_order_chan <- order
			time.Sleep(1 * time.Second)

		}
		if driver.Get_button_signal(global.BUTTON_COMMAND, global.FLOOR_2) != 0 {
			order = queue.Make_new_order(global.BUTTON_COMMAND, global.FLOOR_2, queue.Active, 0)
			new_order_chan <- order
			time.Sleep(1 * time.Second)

		}
		if driver.Get_button_signal(global.BUTTON_COMMAND, global.FLOOR_3) != 0 {
			order = queue.Make_new_order(global.BUTTON_COMMAND, global.FLOOR_3, queue.Active, 0)
			new_order_chan <- order
			time.Sleep(1 * time.Second)

		}
		if driver.Get_button_signal(global.BUTTON_COMMAND, global.FLOOR_4) != 0 {
			order = queue.Make_new_order(global.BUTTON_COMMAND, global.FLOOR_4, queue.Active, 0)
			new_order_chan <- order
			time.Sleep(1 * time.Second)

		}

	}
}
