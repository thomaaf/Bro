package network

import (
	"flag"
	"fmt"
	"global"
	//"network/bcast"
	"driver"
	"network/bcast"
	"network/localip"
	"network/peers"
	"queue"
	"strconv"
	"time"
)

const (
	master_port = 20079
	slave_port  = 20179
)

var Local_ip int

type Master_msg struct {
	Address     int
	Global_list [global.NUM_GLOBAL_ORDERS]queue.Order
}

type Slave_msg struct {
	Address       int
	Internal_list [global.NUM_INTERNAL_ORDERS]queue.Order
	External_list [global.NUM_GLOBAL_ORDERS]queue.Order
	Elevator_info queue.Elev_info
}

func Network_handler() {
	var new_info peers.PeerUpdate

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	//localip, _ := localip.LocalIP()

	if id == "" {
		local_ip, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			local_ip = "Disconnected."
		}
		id = fmt.Sprintf(local_ip)

		peer_update_chan := make(chan peers.PeerUpdate)
		peer_transmit_enable_chan := make(chan bool)

		go peers.Transmitter(20243, id, peer_transmit_enable_chan)
		go peers.Receiver(20243, peer_update_chan)
		for {
			select {
			case new_info_catcher := <-peer_update_chan:
				global.Num_elev_online = 0
				new_info = new_info_catcher

				fmt.Printf("Peer update:\n")
				fmt.Printf("  Peers:    %q\n", new_info.Peers)
				fmt.Printf("  New:      %q\n", new_info.New)
				fmt.Printf("  Lost:     %q\n", new_info.Lost)
				for i := 0; i < len(new_info.Peers); i++ {
					str_ip := new_info.Peers[i]
					int_ip, _ := strconv.Atoi(str_ip[12:])
					queue.Elevators_online[i].Elev_ip = int_ip
					global.Num_elev_online = global.Num_elev_online + 1
					fmt.Println("The ip for elev ", i, " is ", int_ip)
					fmt.Println("Num elev online here: ", global.Num_elev_online)
				}
				if len(new_info.Peers) == 2 {
					queue.Elevators_online[2].Elev_ip = -1
				}
				if len(new_info.Peers) == 1 {
					queue.Elevators_online[1].Elev_ip = -1
					queue.Elevators_online[2].Elev_ip = -1
				}
				if len(new_info.Peers) == 0 {
					queue.Elevators_online[1].Elev_ip = -1
					queue.Elevators_online[2].Elev_ip = -1
					queue.Elevators_online[0].Elev_ip = -1

				}
				Local_ip, _ = strconv.Atoi(local_ip[12:])
				queue.My_info.Elev_ip = Local_ip
				Choose_master()

				if global.Is_master {
					var lost_peer bool
					var str_lost_ip string
					var lost_ip int
					// redeleger part 1
					if len(new_info.Lost) == 1 {
						lost_peer = true
						str_lost_ip = new_info.Lost[0]
						lost_ip, _ = strconv.Atoi(str_lost_ip[12:])
					} else if len(new_info.Lost) != 0 {
						fmt.Println("I just want to lose one person at a time, thank you.")
						lost_peer = true
						str_lost_ip = new_info.Lost[0]
						lost_ip, _ = strconv.Atoi(str_lost_ip[12:])
					}
					//--------------
					// redeleger part 2
					if lost_peer {
						for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
							if queue.Global_order_list[i].Assigned_to == lost_ip {
								queue.Global_order_list[i] = queue.Delegate_order(queue.Global_order_list[i])
							}
						}
					}
					//---------------
				}
			}
		}
	}
}
func Choose_master() {

	ip_addresses := queue.Elevators_online
	highest_ip := 0
	for i := 0; i < 2; i++ {
		if ip_addresses[i].Elev_ip > highest_ip {
			highest_ip = ip_addresses[i].Elev_ip
		}
	}
	if Local_ip == highest_ip {
		fmt.Println("I am the master.")
		global.Is_master = true

	} else {
		fmt.Println("I am a slave.")
		global.Is_master = false
	}
	if highest_ip == 0 {
		fmt.Println("I have lost network")
		global.Is_master = false
		global.Lost_network = true
	}
}

func Network_receiver(new_order_bool_chan chan bool, new_order_chan chan queue.Order, new_global_order_bool_chan chan bool) {
	var receive_port int

	slave_receiver := make(chan Master_msg)
	master_receiver := make(chan Slave_msg)

	for {
		if global.Is_master {
			fmt.Println("Connecting as the master.")
			receive_port = slave_port
			go bcast.Receiver(receive_port, master_receiver)

			for {
				select {
				case catch_msg_from_slave := <-master_receiver:
					msg_from_slave := catch_msg_from_slave
					fmt.Println("Master received : ", catch_msg_from_slave)

					Master_msg_handler(msg_from_slave, new_order_bool_chan, new_order_chan, new_global_order_bool_chan)
					//queue.Master_msg_handler(catch_msg_from_slave)
				}
				if global.Is_master == false {
					break
				}
			}
		} else {
			fmt.Println("Connecting as a slave.")
			receive_port = master_port

			go bcast.Receiver(receive_port, slave_receiver)

			for {
				select {
				case catch_msg_from_master := <-slave_receiver:
					msg_from_master := catch_msg_from_master
					Slave_msg_handler(msg_from_master, new_order_bool_chan, new_order_chan, new_global_order_bool_chan)
					fmt.Println("Slave received : ", catch_msg_from_master)
					//queue.Slave_msg_handler(catch_msg_from_master, new_order_bool_chan)
				}
				if global.Is_master == true {
					break
				}
			}
		}
	}
}

func Network_sender(new_order_bool_chan chan bool, new_order_chan chan queue.Order) {
	var broadcast_port int

	master_sender := make(chan Master_msg)
	slave_sender := make(chan Slave_msg)

	for {
		if global.Is_master {
			for {
				fmt.Println("Connecting as the master.")
				broadcast_port = master_port

				go bcast.Transmitter(broadcast_port, master_sender)
				go master_transmit(master_sender)
				time.Sleep(1 * time.Second)
				if global.Is_master == false {
					break
				}
			}
		} else {
			for {
				fmt.Println("Connecting as a slave.")
				broadcast_port = slave_port

				go bcast.Transmitter(broadcast_port, slave_sender)
				go slave_transmit(slave_sender)
				time.Sleep(1 * time.Second)
				if global.Is_master {
					break
				}
			}
		}
	}
}

//FJERNET FOR INNI TRANSMITTENE. Det hjalp
func master_transmit(master_sender chan Master_msg) {
	var master_msg_to_send Master_msg

	master_msg_to_send.Address = Local_ip
	master_msg_to_send.Global_list = queue.Global_order_list
	master_sender <- master_msg_to_send
	fmt.Println("Master sent the global list: ", master_msg_to_send.Global_list)
	fmt.Println("My external list is: ", queue.External_order_list)
	time.Sleep(1 * time.Second)

}

func slave_transmit(slave_sender chan Slave_msg) {
	var slave_msg_to_send Slave_msg

	slave_msg_to_send.Address = Local_ip
	slave_msg_to_send.Internal_list = queue.Internal_order_list
	slave_msg_to_send.External_list = queue.External_order_list
	slave_msg_to_send.Elevator_info = queue.My_info
	fmt.Println("My external list is: ", queue.External_order_list)
	slave_sender <- slave_msg_to_send
	time.Sleep(1 * time.Second)

}

//--
func Master_msg_handler(msg_from_slave Slave_msg, new_order_bool_chan chan bool, new_order_chan chan queue.Order, new_global_order_bool_chan chan bool) {
	var num int
	for i := 0; i < global.Num_elev_online; i++ {
		if queue.Elevators_online[i].Elev_ip == msg_from_slave.Elevator_info.Elev_ip {
			num = i
		}
	}
	//Oppdaterer infoen om den nye heisen
	queue.Elevators_online[num].Elev_last_floor = msg_from_slave.Elevator_info.Elev_last_floor
	queue.Elevators_online[num].Elev_dir = msg_from_slave.Elevator_info.Elev_dir
	queue.Elevators_online[num].Elev_state = msg_from_slave.Elevator_info.Elev_state
	external_order_list := msg_from_slave.External_list
	//internal_order_list := msg_from_slave.Internal_list
	fmt.Println("Master 1")
	//Sjekke etter nye bestillinger
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if external_order_list[i].Assigned_to == 0 && external_order_list[i].Order_state != queue.Finished && external_order_list[i].Order_state != queue.Inactive {
			//If noone ownes it, it must be a new order, and it is not finished or inactive
			new_order := external_order_list[i]
			new_order = queue.Delegate_order(new_order)
			//--
			queue.Add_new_global_order(new_order, new_order_bool_chan, new_order_chan, new_global_order_bool_chan)

		}
	}
	fmt.Println("Master 2")
	//Sjekke etter endrede bestillinger
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		for j := 0; j < global.NUM_GLOBAL_ORDERS; j++ {
			if queue.Global_order_list[i].Button == external_order_list[j].Button && queue.Global_order_list[i].Floor == external_order_list[j].Floor && queue.Global_order_list[i].Order_state != external_order_list[j].Order_state {
				queue.Global_order_list[i].Order_state = external_order_list[i].Order_state
				//Burde også sjekke om noe er satt til Finished og i såfall slette det???

			}
			// Skru på lamper på alle ordre som ikke er inaktiv eller finished
			if queue.Global_order_list[i].Order_state != queue.Inactive && queue.Global_order_list[i].Order_state != queue.Finished {
				driver.Set_button_lamp(queue.Global_order_list[i].Button, queue.Global_order_list[i].Floor, global.ON)
			} else {
				driver.Set_button_lamp(queue.Global_order_list[i].Button, queue.Global_order_list[i].Floor, global.OFF)
			}
		}
	}
	fmt.Println("Master 3")
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if queue.Global_order_list[i].Assigned_to == Local_ip && queue.Global_order_list[i].Order_state != queue.Inactive && queue.Global_order_list[i].Order_state != queue.Finished {
			fmt.Println("Adding new order to external from within master msg loop")
			fmt.Println("My new order bll", queue.Global_order_list[i])
			queue.Add_new_external_order(queue.Global_order_list[i], new_order_bool_chan, new_order_chan, new_global_order_bool_chan) // Bør kanskje kjøres som en go func?? Litt usikker
		}
		// Skru på lamper på alle ordre som ikke er inaktiv eller finished
		if queue.Global_order_list[i].Order_state != queue.Inactive && queue.Global_order_list[i].Order_state != queue.Finished {
			driver.Set_button_lamp(queue.Global_order_list[i].Button, queue.Global_order_list[i].Floor, global.ON)
		} else {
			driver.Set_button_lamp(queue.Global_order_list[i].Button, queue.Global_order_list[i].Floor, global.OFF)
		}
	}
}

func Slave_msg_handler(msg_from_master Master_msg, new_order_bool_chan chan bool, new_order_chan chan queue.Order, new_global_order_bool_chan chan bool) {
	my_ip := Local_ip
	//master_ip := msg_from_master.Address ---- Do we really need this?
	global_order_list := msg_from_master.Global_list
	fmt.Println("Checking the message")
	//Sjekker om noen bestillinger er assigna til seg selv og deretter at den ikke har state Inactive eller Finished
	for i := 0; i < global.NUM_GLOBAL_ORDERS; i++ {
		if global_order_list[i].Assigned_to == my_ip && global_order_list[i].Order_state != queue.Inactive && global_order_list[i].Order_state != queue.Finished {
			queue.Add_new_external_order(global_order_list[i], new_order_bool_chan, new_order_chan, new_global_order_bool_chan) // Bør kanskje kjøres som en go func?? Litt usikker
		}
		// Skru på lamper på alle ordre som ikke er inaktiv eller finished
		if global_order_list[i].Order_state != queue.Inactive && global_order_list[i].Order_state != queue.Finished {
			driver.Set_button_lamp(global_order_list[i].Button, global_order_list[i].Floor, global.ON)
		} else {
			driver.Set_button_lamp(global_order_list[i].Button, global_order_list[i].Floor, global.OFF)
		}
	}

	//queue.Bool_to_new_order_channel(new_order, new_order_bool_chan) // Add_new_external_order gjør vel det?
}
