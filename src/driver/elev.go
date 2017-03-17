package driver

import (
	"fmt"
	"global"
	"time"
)

var lamp_channel_matrix = [global.NUM_FLOORS][global.NUM_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix = [global.NUM_FLOORS][global.NUM_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Set_button_lamp(button global.Button_t, floor global.Floor_t, on_off global.On_off_t) {
	if on_off == global.ON {
		Io_set_bit(lamp_channel_matrix[floor][button])
	} else {
		Io_clear_bit(lamp_channel_matrix[floor][button])
	}
}

func Set_floor_indicator_lamp(floor global.Floor_t) {
	switch {
	case floor == global.FLOOR_1:
		Io_clear_bit(LIGHT_FLOOR_IND1)
		Io_clear_bit(LIGHT_FLOOR_IND2)
	case floor == global.FLOOR_2:
		Io_clear_bit(LIGHT_FLOOR_IND1)
		Io_set_bit(LIGHT_FLOOR_IND2)
	case floor == global.FLOOR_3:
		Io_set_bit(LIGHT_FLOOR_IND1)
		Io_clear_bit(LIGHT_FLOOR_IND2)
	case floor == global.FLOOR_4:
		Io_set_bit(LIGHT_FLOOR_IND1)
		Io_set_bit(LIGHT_FLOOR_IND2)
	}
}

func Set_door_open_lamp(on_off global.On_off_t) {
	if on_off == global.ON {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Open_door() {
	Set_door_open_lamp(global.ON)
	time.Sleep(3 * time.Second)
	Set_door_open_lamp(global.OFF)
}

func Set_stop_lamp(on_off global.On_off_t) {
	if on_off == global.ON {
		Io_set_bit(LIGHT_STOP)
	} else {
		Io_clear_bit(LIGHT_STOP)
	}
}

func Get_floor_sensor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) != 0 {
		return 1
	}
	if Io_read_bit(SENSOR_FLOOR2) != 0 {
		return 2
	}
	if Io_read_bit(SENSOR_FLOOR3) != 0 {
		return 3
	}
	if Io_read_bit(SENSOR_FLOOR4) != 0 {
		return 4
	} else {
		return -1
	}
}

func Get_floor_sensor_signal_floor_t() global.Floor_t {
	if Get_floor_sensor_signal() == 1 {
		return global.FLOOR_1
	}
	if Get_floor_sensor_signal() == 2 {
		return global.FLOOR_2
	}
	if Get_floor_sensor_signal() == 3 {
		return global.FLOOR_3
	}
	if Get_floor_sensor_signal() == 4 {
		return global.FLOOR_4
	} else {
		return global.FLOOR_1
	}
}

func Set_motor_direction(dir global.Motor_direction_t) {
	if dir == global.DIR_STOP {
		Io_write_analog(MOTOR, 0)
	} else if dir == global.DIR_UP {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, global.MOTOR_SPEED)
	} else if dir == global.DIR_DOWN {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, global.MOTOR_SPEED)
	}
}

func Get_button_signal(button global.Button_t, floor global.Floor_t) int {
	return Io_read_bit(button_channel_matrix[floor][button])
}

func Set_all_lamps(on_off global.On_off_t) {
	// Set all button lamps
	Set_button_lamp(global.BUTTON_UP, global.FLOOR_1, on_off)
	Set_button_lamp(global.BUTTON_UP, global.FLOOR_2, on_off)
	Set_button_lamp(global.BUTTON_UP, global.FLOOR_3, on_off)
	Set_button_lamp(global.BUTTON_DOWN, global.FLOOR_2, on_off)
	Set_button_lamp(global.BUTTON_DOWN, global.FLOOR_3, on_off)
	Set_button_lamp(global.BUTTON_DOWN, global.FLOOR_4, on_off)
	Set_button_lamp(global.BUTTON_COMMAND, global.FLOOR_1, on_off)
	Set_button_lamp(global.BUTTON_COMMAND, global.FLOOR_2, on_off)
	Set_button_lamp(global.BUTTON_COMMAND, global.FLOOR_3, on_off)
	Set_button_lamp(global.BUTTON_COMMAND, global.FLOOR_4, on_off)

	// Set door open lamp
	Set_door_open_lamp(on_off)

	// Set stop lamp
	Set_stop_lamp(on_off)
}

func Elevator_to_floor_direct(floor global.Floor_t) {
	switch {
	case floor == global.FLOOR_1:
		Elevator_to_floor_direct_int(1)
	case floor == global.FLOOR_2:
		Elevator_to_floor_direct_int(2)
	case floor == global.FLOOR_3:
		Elevator_to_floor_direct_int(3)
	case floor == global.FLOOR_4:
		Elevator_to_floor_direct_int(4)
	}
	Set_floor_indicator_lamp(floor)
}

func Elevator_to_floor_direct_int(floor int) {
	my_floor := Get_floor_sensor_signal()

	// If the elevator is between two floors
	timer := time.NewTimer(3 * time.Second)
	timeout := false
	go func() {
		<-timer.C
		timeout = true
	}()
	for Get_floor_sensor_signal() == -1 {
		if !timeout {
			Set_motor_direction(global.DIR_UP)
		} else if timeout {
			Set_motor_direction(global.DIR_DOWN)
		}
	}

	// Go to desired floor
	my_floor = Get_floor_sensor_signal()
	if my_floor < floor {
		for Get_floor_sensor_signal() != floor {
			Set_motor_direction(global.DIR_UP)
		}
	} else if my_floor > floor {
		for Get_floor_sensor_signal() != floor {
			Set_motor_direction(global.DIR_DOWN)
		}
	}
	Set_motor_direction(global.DIR_STOP)
}

func Floor_int_to_floor_t(floor_int int) global.Floor_t {
	switch {
	case floor_int == 1:
		return global.FLOOR_1
	case floor_int == 2:
		return global.FLOOR_2
	case floor_int == 3:
		return global.FLOOR_3
	case floor_int == 4:
		return global.FLOOR_4
	}
	return global.FLOOR_1
}

func Floor_t_to_floor_int(floor global.Floor_t) int {
	switch {
	case floor == global.FLOOR_1:
		return 1
	case floor == global.FLOOR_2:
		return 2
	case floor == global.FLOOR_3:
		return 3
	case floor == global.FLOOR_4:
		return 4
	}
	return -1
}

func Elevator_init() {
	fmt.Println("Running elevator initialization.")
	Io_init()

	Set_all_lamps(global.OFF)
	Elevator_to_floor_direct(global.FLOOR_1)

	fmt.Println("Elevator initialization done!")
}
