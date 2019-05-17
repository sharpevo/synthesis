// Package tml provides functions for TML_Lib, an SDK from Technosoft. Be sure
// to read the offical document before importing this package. Actually, there
// is a more concise and effective implementation, but requires more time for
// testing.
package tml

import (
	"fmt"
	"posam/config"
	"posam/gui/uiutil"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"posam/util/log"
	"reflect"
	"time"
	"tml"
)

var (
	clientMap *concurrentmap.ConcurrentMap
	_MTIMEOUT time.Duration
	_MDELAY   time.Duration

	_CONFIG_MOTION_TIMEOUT        = "tml.motion.timeout"
	_CONFIG_MOTION_DELAY          = "tml.motion.delay"
	_CONFIG_TONPOSOK              = "tml.tonposok"
	_CONFIG_COMPENSATION_BASIC    = "tml.compensation.basic"
	_CONFIG_COMPENSATION_ADVANCED = "tml.compensation.advanced"

	_SET_TONPOSOK          = config.GetBool(_CONFIG_TONPOSOK)
	_COMPENSATION          = config.GetBool(_CONFIG_COMPENSATION_BASIC)
	_COMPENSATION_ADVANCED = config.GetBool(_CONFIG_COMPENSATION_ADVANCED)
)

func init() {
	clientMap = concurrentmap.NewConcurrentMap()
	config.SetDefault(_CONFIG_MOTION_TIMEOUT, 100)
	_MTIMEOUT = time.Duration(config.GetInt(_CONFIG_MOTION_TIMEOUT)) * time.Second
	config.SetDefault(_CONFIG_MOTION_DELAY, 0)
	_MDELAY = time.Duration(config.GetInt(_CONFIG_MOTION_DELAY)) * time.Millisecond
}

func addInstance(client *Client) (*Client, bool) {
	key := string(client.name)
	if c, ok := clientMap.Get(key); ok {
		return c.(*Client), true
	} else {
		clientMap.Set(key, client)
		return client, false
	}
}

// ResetInstance resets instance map. Not tested yet.
func ResetInstance() {
	for item := range clientMap.Iter() {
		client := item.Value.(*Client)
		log.Df("to terminate client: %v", client.name)
		// TODO: terminate client
	}
	clientMap = concurrentmap.NewConcurrentMap()
}

// A client is the abstraction of TML devices connected via RS232. Note that
// the communication channel file descriptor is named as channelDescriptor,
// although it's not required right now for single channel apps.
type Client struct {
	name           string
	baudRate       int
	axisXID        int
	axisXSetupFile string
	axisYID        int
	axisYSetupFile string

	channelDescriptor int
	requestQueue      *blockingqueue.BlockingQueue

	posX float64
	posY float64
	//TODO: speeds
}

// NewClient returns TML device connection instance which been initialized
// and launched at the same time.
func NewClient(
	name string,
	baud int,
	axisXID int,
	axisXSetupFile string,
	axisYID int,
	axisYSetupFile string,
) (*Client, error) {
	client := &Client{
		name:           name,
		baudRate:       baud,
		axisXID:        axisXID,
		axisXSetupFile: axisXSetupFile,
		axisYID:        axisYID,
		axisYSetupFile: axisYSetupFile,
		requestQueue:   blockingqueue.NewBlockingQueue(),
	}
	if c, found := addInstance(client); found {
		return c, fmt.Errorf("client existed")
	}
	go launchClient(client)
	log.Df(
		"client launched(AxisXID: %v, AxisYID: %v)", client.axisXID, client.axisYID)
	return client, nil
}

func (c *Client) connect() (err error) {
	log.Df("Connecting the motor %q...", c.name)

	commType := tml.CHANNEL_RS232
	hostID := 1

	descriptor, err := tml.OpenChannel(
		c.name,
		commType,
		hostID,
		c.baudRate,
	)
	if err != nil {
		return err
	}
	c.channelDescriptor = descriptor

	idxSetup, err := tml.LoadSetup(c.axisXSetupFile)
	if err != nil {
		return err
	}
	if err = tml.SetupAxis(c.axisXID, idxSetup); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.axisXID); err != nil {
		return err
	}
	if err = tml.DriveInitialisation(); err != nil {
		return err
	}
	if err = tml.Power(true); err != nil {
		return err
	}

	idxSetup, err = tml.LoadSetup(c.axisYSetupFile)
	if err != nil {
		return err
	}
	if err = tml.SetupAxis(c.axisYID, idxSetup); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.axisYID); err != nil {
		return err
	}
	if err = tml.DriveInitialisation(); err != nil {
		return err
	}
	if err = tml.Power(true); err != nil {
		return err
	}

	log.D("checking status...")
	var statusx int
	var statusy int
	for i := 0; i < 15; i++ {
		if statusx == 0 {
			err = tml.SelectAxis(c.axisXID)
			if err != nil {
				return err
			}
			err = tml.ReadStatus(3, &statusx)
			if err != nil {
				return err
			}
			statusx = statusx & (1 << 15)
		}

		if statusy == 0 {
			err = tml.SelectAxis(c.axisYID)
			if err != nil {
				return err
			}
			err = tml.ReadStatus(3, &statusy)
			if err != nil {
				return err
			}
			statusy = statusy & (1 << 15)
		}
		if statusx != 0 && statusy != 0 {
			break
		}
		<-time.After(1 * time.Second)
	}
	if statusx == 0 || statusy == 0 {
		return fmt.Errorf("failed to enable power on axes: x(%d) / y(%d)", statusx, statusy)
	}

	log.Df("motor %q is ready", c.name)
	return nil
}

type request struct {
	responsec chan response
	function  string
	arguments []interface{}
}

type response struct {
	err error
}

var launchClient = func(c *Client) {
	log.D("motor client launched")
	c.connect()

	for {
		reqi, err := c.requestQueue.Pop()
		if err != nil {
			log.E("motor client terminated")
			return
		}
		req := reqi.(*request)
		function := reflect.ValueOf(c).MethodByName(req.function)
		args := []reflect.Value{}
		for _, v := range req.arguments {
			args = append(args, reflect.ValueOf(v))
		}
		result := function.Call(args)
		erri := result[0].Interface()
		if erri != nil {
			req.responsec <- response{err: erri.(error)}
			continue
		}
		if err := c.UpdateMotionStatus(); err != nil {
			log.E(err)
		}
		req.responsec <- response{err: nil}
	}
}

// MoveAbsoluteByAxis moves the motor by axis with the following arguments:
//
// - _aid: axis id of X or Y;
//
// - _pos: position to reached expressed in TML position units;
//
// - _spd: slew speed expressed in TML speed units. If the value is zero the
// drive/motor will use the previously value set for speed;
//
// - _acc: acceleration/deceleration rate expressed in TML acceleration
// units. If its value is zero the drive/motor will use the previously value
// set for acceleration;
//
// - _mmt: defines the moment when the motion is started;
//
// - _ref: specifies how the motion reference is computed: from actual values
// of position and speed reference or from actual values of load/motor position
// and speed
func (c *Client) MoveAbsoluteByAxis(
	_aid, _pos, _spd, _acc, _mmt, _ref interface{},
) (err error) {
	aid, ok := _aid.(int)
	if !ok {
		return fmt.Errorf("failed to convert aid %v", _aid)
	}
	pos, spd, acc, mmt, ref, err := parseAbsArgs(_pos, _spd, _acc, _mmt, _ref)
	if err != nil {
		return err
	}
	log.If("moving axis %d to %v...", aid, pos)
	if err = tml.SelectAxis(aid); err != nil {
		return err
	}
	if err = tml.MoveAbsolute(
		tml.CalcPosition(aid, pos),
		tml.CalcSpeed(aid, spd),
		tml.CalcAccel(aid, acc),
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	log.D("done")
	c.CompensateMotion(aid, pos)
	return nil
}

// MoveRelativeByAxis moves the motor by axis with the following arguments:
//
// - _aid: axis id of X or Y;
//
// - _pos: position increment expressed in TML position units;
//
// - _spd: slew speed expressed in TML speed units. If the value is zero the
// drive/motor will use the previously value set for speed;
//
// - _acc: acceleration/deceleration rate expressed in TML acceleration
// units. If its value is zero the drive/motor will use the previously value
// set for acceleration;
//
// - _add: specifies how is computed the position to reach;
//
// - _mmt: defines the moment when the motion is started;
//
// - refa: specifies how the motion reference is computed: from actual values
// of position and speed reference or from actual values of load/motor position
// and speed
func (c *Client) MoveRelativeByAxis(
	_aid, _pos, _spd, _acc, _add, _mmt, _ref interface{},
) (err error) {
	aid, ok := _aid.(int)
	if !ok {
		return fmt.Errorf("failed to convert aid %v", _aid)
	}
	pos, spd, acc, add, mmt, ref, err := parseRelArgs(_pos, _spd, _acc, _add, _mmt, _ref)
	if err != nil {
		return err
	}
	log.If("moving axis %d by %v...", aid, pos)
	if err = tml.SelectAxis(aid); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		tml.CalcPosition(aid, pos),
		tml.CalcSpeed(aid, spd),
		tml.CalcAccel(aid, acc),
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	log.D("done")
	// TODO: compensate
	return nil
}

// MoveAbsByAxis create a MoveAbsoluteByAxis request and send it to the queue.
func (c *Client) MoveAbsByAxis(
	axisID int,
	pos float64,
	speed float64,
	accel float64,
) error {
	req := request{
		responsec: make(chan response),
		function:  "MoveAbsoluteByAxis",
		arguments: []interface{}{
			axisID,
			pos,
			speed,
			accel,
			1,
			1,
		},
	}
	c.requestQueue.Push(&req)
	log.Df(
		"waiting for axis %d response: absolute motion to %v",
		axisID,
		pos,
	)
	resp := <-req.responsec
	if resp.err != nil {
		return resp.err
	}
	return nil
}

// MoveRelByAxis create a MoveRelativeByAxis request and send it to the queue.
func (c *Client) MoveRelByAxis(
	axisID int,
	pos float64,
	speed float64,
	accel float64,
) error {
	req := request{
		responsec: make(chan response),
		function:  "MoveRelativeByAxis",
		arguments: []interface{}{
			axisID,
			pos,
			speed,
			accel,
			true,
			1,
			1,
		},
	}
	c.requestQueue.Push(&req)
	log.Df(
		"waiting for axis %d response: relative motion to %v",
		axisID,
		pos,
	)
	resp := <-req.responsec
	if resp.err != nil {
		return resp.err
	}
	return nil
}

// MoveRelative moves motor with the position increments of x-axis and y-axis
// at the same time.
func (c *Client) MoveRelative(
	_posx, _posy, _spd, _acc, _add, _mmt, _ref interface{},
) (err error) {
	posx, spd, acc, add, mmt, ref, err := parseRelArgs(_posx, _spd, _acc, _add, _mmt, _ref)
	if err != nil {
		return err
	}
	posy, ok := _posy.(float64)
	if !ok {
		return fmt.Errorf("failed to convert posy %v", _posy)
	}
	log.If("moving by (%v,%v)...", posx, posy)
	if err = tml.SelectAxis(c.axisXID); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		tml.CalcPosition(c.axisXID, posx),
		tml.CalcSpeed(c.axisXID, spd),
		tml.CalcAccel(c.axisXID, acc),
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(false, false); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.axisYID); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		tml.CalcPosition(c.axisYID, posy),
		tml.CalcSpeed(c.axisYID, spd),
		tml.CalcAccel(c.axisYID, acc),
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	for xc, tc := false, time.After(_MTIMEOUT); !xc; {
		if _MDELAY != 0 {
			<-time.After(_MDELAY)
		}
		select {
		case <-tc:
			xc = true
			return fmt.Errorf("Aoztech timeout")
		default:
			if err = tml.SelectAxis(c.axisXID); err != nil {
				log.E(err)
			}
			tml.CheckEvent(&xc)
		}
	}
	log.D("done")
	return nil
}

// MoveRel creates a MoveRelative request and send it to the queue.
func (c *Client) MoveRel(
	posx, posy, speed, accel float64,
) error {
	req := request{
		responsec: make(chan response),
		function:  "MoveRelative",
		arguments: []interface{}{
			posx,
			posy,
			speed,
			accel,
			true,
			1,
			1,
		},
	}
	c.requestQueue.Push(&req)
	log.Df(
		"waiting for response: relative motion to (%v, %v)",
		posx,
		posy,
	)
	resp := <-req.responsec
	if resp.err != nil {
		return resp.err
	}
	return nil
}

// MoveAbosolute moves motor to the given x-axis position and y-axis position
// at the same time.
func (c *Client) MoveAbsolute(
	_posx, _posy, _spd, _acc, _mmt, _ref interface{},
) (err error) {
	posx, spd, acc, mmt, ref, err := parseAbsArgs(_posx, _spd, _acc, _mmt, _ref)
	if err != nil {
		return err
	}
	posy, ok := _posy.(float64)
	if !ok {
		return fmt.Errorf("failed to convert posy %v", _posy)
	}
	log.If("moving to (%v,%v)...", posx, posy)
	if err = tml.SelectAxis(c.axisXID); err != nil {
		return err
	}
	if _SET_TONPOSOK {
		log.Df("set TONPOSOK %v", tml.SetIntVariable("TONPOSOK", 100))
	}
	if err = tml.MoveAbsolute(
		tml.CalcPosition(c.axisXID, posx),
		tml.CalcSpeed(c.axisXID, spd),
		tml.CalcAccel(c.axisXID, acc),
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(false, false); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.axisYID); err != nil {
		return err
	}
	if _SET_TONPOSOK {
		log.Df("set TONPOSOK %v", tml.SetIntVariable("TONPOSOK", 100))
	}
	if err = tml.MoveAbsolute(
		tml.CalcPosition(c.axisYID, posy),
		tml.CalcSpeed(c.axisYID, spd),
		tml.CalcAccel(c.axisYID, acc),
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	c.CompensateMotion(c.axisYID, posy)
	if _COMPENSATION && _COMPENSATION_ADVANCED {
		log.I("2nd compensation")
		c.CompensateMotion(c.axisYID, posy)
	}
	for xc, tc := false, time.After(_MTIMEOUT); !xc; {
		if _MDELAY != 0 {
			<-time.After(_MDELAY)
		}
		select {
		case <-tc:
			xc = true
			return fmt.Errorf("Aoztech timeout")
		default:
			if err = tml.SelectAxis(c.axisXID); err != nil {
				log.E(err)
			}
			tml.CheckEvent(&xc)
		}
	}
	log.D("done")
	return nil
}

// MoveAbs creates a MoveAbsolute request and send it to the queue.
func (c *Client) MoveAbs(
	posx, posy, speed, accel float64,
) error {
	req := request{
		responsec: make(chan response),
		function:  "MoveAbsolute",
		arguments: []interface{}{
			posx,
			posy,
			speed,
			accel,
			1,
			1,
		},
	}
	c.requestQueue.Push(&req)
	log.Ef(
		"waiting for response: absolute motion to (%v, %v)",
		posx,
		posy,
	)
	resp := <-req.responsec
	if resp.err != nil {
		return resp.err
	}
	return nil
}

func parseAbsArgs(
	_pos, _spd, _acc, _mmt, _ref interface{},
) (
	pos float64,
	spd float64,
	acc float64,
	mmt int,
	ref int,
	err error,
) {
	var ok bool
	pos, ok = _pos.(float64)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert pos %v", _pos)
	}
	spd, ok = _spd.(float64)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert spd %v", _spd)
	}
	acc, ok = _acc.(float64)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert acc %v", _acc)
	}
	mmt, ok = _mmt.(int)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert mmt %v", _mmt)
	}
	ref, ok = _ref.(int)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert ref %v", _ref)
	}
	return pos, spd, acc, mmt, ref, nil
}

func parseRelArgs(
	_pos, _spd, _acc, _add, _mmt, _ref interface{},
) (
	pos float64,
	spd float64,
	acc float64,
	add bool,
	mmt int,
	ref int,
	err error,
) {
	pos, spd, acc, mmt, ref, err = parseAbsArgs(_pos, _spd, _acc, _mmt, _ref)
	if err != nil {
		return pos, spd, acc, add, mmt, ref, err
	}
	var ok bool
	add, ok = _add.(bool)
	if !ok {
		return pos, spd, acc, add, mmt, ref,
			fmt.Errorf("failed to convert add %v", _add)
	}
	return pos, spd, acc, add, mmt, ref, nil
}

// UpadteMotionStatus updates the position of x-axis and y-axis.
func (c *Client) UpdateMotionStatus() (err error) {
	if c.posX, err = tml.ActualPosition(c.axisXID); err != nil {
		return err
	}
	if c.posY, err = tml.ActualPosition(c.axisYID); err != nil {
		return err
	}
	uiutil.App.UpdateMotorStatusSlot(fmt.Sprintf("Motor: (%v, %v)", c.posX, c.posY))
	return nil
}

// PosX returns the x-axis position of motor
func (c *Client) PosX() float64 {
	return c.posX
}

// PosY returns the y-axis position of motor
func (c *Client) PosY() float64 {
	return c.posY
}

// AxisXID returns the id of x-axis
func (c *Client) AxisXID() int {
	return c.axisXID
}

// AxisYID returns the id of y-axis
func (c *Client) AxisYID() int {
	return c.axisYID
}

// CompensateMotion moves the motor forward or backward according to the actual
// position and the expectation once or twice if _COMPENSATION is specified.
func (c *Client) CompensateMotion(axisID int, target float64) (err error) {
	// TODO: relative compensation
	if !_COMPENSATION {
		return nil
	}
	switch axisID {
	case c.axisYID:
		pos, err := tml.ActualPosition(c.axisYID)
		diffPos := target - pos
		offset := tml.CalcPosition(c.axisYID, diffPos)
		log.Df("compensating axis %d by %v (diff apos %v, actual pos %v)...", c.axisYID, diffPos, offset, pos)
		if err = tml.SelectAxis(c.axisYID); err != nil {
			return err
		}
		if err = tml.MoveRelative(
			tml.CalcPosition(c.axisYID, diffPos),
			tml.CalcSpeed(c.axisYID, 5),
			tml.CalcAccel(c.axisYID, 50),
			true,
			1,
			1,
		); err != nil {
			return err
		}
		if err = tml.SetEventOnMotionComplete(true, false); err != nil {
			return err
		}
		log.D("done")
		return nil
	default:
		return nil
	}
}

// CompensateMotionTPOS compensates the motion by TPOS.
func (c *Client) CompensateMotionTPOS(axisID int, target float64) (err error) {
	if !_COMPENSATION {
		return nil
	}
	switch axisID {
	case c.axisYID:
		pos, err := tml.TargetPosition(c.axisYID)
		diffPos := target - pos
		offset := tml.CalcPosition(c.axisYID, diffPos)
		if diffPos == 0 {
			log.Df("not compensate axis %d by %v (diff tpos %v, actual pos %v)...", c.axisYID, diffPos, offset, pos)
			return nil
		}
		log.Df("compensating axis %d by %v (diff tpos %v, actual pos %v)...", c.axisYID, diffPos, offset, pos)
		if err = tml.SelectAxis(c.axisYID); err != nil {
			return err
		}
		if err = tml.MoveRelative(
			tml.CalcPosition(c.axisYID, diffPos),
			tml.CalcSpeed(c.axisYID, 5),
			tml.CalcAccel(c.axisYID, 50),
			true,
			1,
			1,
		); err != nil {
			return err
		}
		if err = tml.SetEventOnMotionComplete(true, false); err != nil {
			return err
		}
		log.D("done")
		return nil
	default:
		return nil
	}
}
