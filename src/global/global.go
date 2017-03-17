package global

const MOTOR_SPEED int = 2800
const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_ELEV = 3
const NUM_ORDER_STATES = 5
const NUM_GLOBAL_ORDERS = 6
const NUM_INTERNAL_ORDERS = 4
const NUM_ORDERS = NUM_GLOBAL_ORDERS + NUM_INTERNAL_ORDERS

type Button_t int

var Lost_network bool

var Is_master bool
var Num_elev_online int
var Iter int

const (
	BUTTON_UP = iota
	BUTTON_DOWN
	BUTTON_COMMAND
)

type Floor_t int

const (
	FLOOR_1 = iota
	FLOOR_2
	FLOOR_3
	FLOOR_4
)

type On_off_t int

const (
	OFF = iota
	ON
)

type Motor_direction_t int

const (
	DIR_DOWN = -1
	DIR_STOP = 0
	DIR_UP   = 1
)


const (
	NONE = iota
	ELEV_1
	ELEV_2
	ELEV_3
)
