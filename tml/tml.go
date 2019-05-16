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
	key := string(client.Name)
	if c, ok := clientMap.Get(key); ok {
		return c.(*Client), true
	} else {
		clientMap.Set(key, client)
		return client, false
	}
}

func ResetInstance() {
	for item := range clientMap.Iter() {
		client := item.Value.(*Client)
		log.Df("terminating client: %v", client.Name)
		//client.Stop()
	}
	clientMap = concurrentmap.NewConcurrentMap()
}

type Client struct {
	Name           string
	BaudRate       int
	AxisXID        int
	AxisXSetupFile string
	AxisYID        int
	AxisYSetupFile string

	ChannelDescriptor int
	RequestQueue      *blockingqueue.BlockingQueue

	PosX float64
	PosY float64
	SpdX float64
	SpdY float64
}

func NewClient(
	name string,
	baud int,
	axisXID int,
	axisXSetupFile string,
	axisYID int,
	axisYSetupFile string,
) (*Client, error) {
	client := &Client{
		Name:           name,
		BaudRate:       baud,
		AxisXID:        axisXID,
		AxisXSetupFile: axisXSetupFile,
		AxisYID:        axisYID,
		AxisYSetupFile: axisYSetupFile,
		RequestQueue:   blockingqueue.NewBlockingQueue(),
	}
	if c, found := addInstance(client); found {
		return c, fmt.Errorf("client existed")
	}
	go launchClient(client)
	log.Df(
		"client launched(AxisXID: %v, AxisYID: %v)", client.AxisXID, client.AxisYID)
	return client, nil
}

func (c *Client) connect() (err error) {
	log.Df("Connecting the motor %q...", c.Name)

	commType := tml.CHANNEL_RS232
	hostID := 1

	descriptor, err := tml.OpenChannel(
		c.Name,
		commType,
		hostID,
		c.BaudRate,
	)
	if err != nil {
		return err
	}
	c.ChannelDescriptor = descriptor

	idxSetup, err := tml.LoadSetup(c.AxisXSetupFile)
	if err != nil {
		return err
	}
	if err = tml.SetupAxis(c.AxisXID, idxSetup); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.AxisXID); err != nil {
		return err
	}
	if err = tml.DriveInitialisation(); err != nil {
		return err
	}
	if err = tml.Power(true); err != nil {
		return err
	}

	idxSetup, err = tml.LoadSetup(c.AxisYSetupFile)
	if err != nil {
		return err
	}
	if err = tml.SetupAxis(c.AxisYID, idxSetup); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.AxisYID); err != nil {
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
			err = tml.SelectAxis(c.AxisXID)
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
			err = tml.SelectAxis(c.AxisYID)
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

	log.Df("motor %q is ready", c.Name)
	return nil
}

type Request struct {
	//AxisID    int
	Responsec chan Response
	Function  string
	Arguments []interface{}
}

type Response struct {
	Error error
}

var launchClient = func(c *Client) {
	log.D("motor client launched")
	c.connect()

	for {
		reqi, err := c.RequestQueue.Pop()
		if err != nil {
			log.E("motor client terminated")
			return
		}
		req := reqi.(*Request)
		function := reflect.ValueOf(c).MethodByName(req.Function)
		args := []reflect.Value{}
		for _, v := range req.Arguments {
			args = append(args, reflect.ValueOf(v))
		}
		result := function.Call(args)
		erri := result[0].Interface()
		if erri != nil {
			req.Responsec <- Response{Error: erri.(error)}
			continue
		}
		if err := c.UpdateMotionStatus(); err != nil {
			log.E(err)
		}
		req.Responsec <- Response{Error: nil}
	}
}

func (c *Client) MoveAbsoluteByAxis(
	aidi interface{},
	posi interface{},
	spdi interface{},
	acci interface{},
	mmti interface{},
	refi interface{},
) (err error) {
	aid, ok := aidi.(int)
	if !ok {
		return fmt.Errorf("failed to convert aid %v", aidi)
	}
	pos, spd, acc, mmt, ref, err := parseAbsArgs(posi, spdi, acci, mmti, refi)
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

func (c *Client) MoveRelativeByAxis(
	aidi interface{},
	posi interface{},
	spdi interface{},
	acci interface{},
	addi interface{},
	mmti interface{},
	refi interface{},
) (err error) {
	aid, ok := aidi.(int)
	if !ok {
		return fmt.Errorf("failed to convert aid %v", aidi)
	}
	pos, spd, acc, add, mmt, ref, err := parseRelArgs(posi, spdi, acci, addi, mmti, refi)
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
	return nil
	//c.CompensateMotion(aid, pos)
}

func (c *Client) MoveAbsByAxis(
	axisID int,
	pos float64,
	speed float64,
	accel float64,
) error {
	req := Request{
		Responsec: make(chan Response),
		Function:  "MoveAbsoluteByAxis",
		Arguments: []interface{}{
			axisID,
			pos,
			speed,
			accel,
			1,
			1,
		},
	}
	c.RequestQueue.Push(&req)
	log.Df(
		"waiting for axis %d response: absolute motion to %v",
		axisID,
		pos,
	)
	resp := <-req.Responsec
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (c *Client) MoveRelByAxis(
	axisID int,
	pos float64,
	speed float64,
	accel float64,
) error {
	req := Request{
		Responsec: make(chan Response),
		Function:  "MoveRelativeByAxis",
		Arguments: []interface{}{
			axisID,
			pos,
			speed,
			accel,
			true,
			1,
			1,
		},
	}
	c.RequestQueue.Push(&req)
	log.Df(
		"waiting for axis %d response: relative motion to %v",
		axisID,
		pos,
	)
	resp := <-req.Responsec
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (c *Client) MoveRelative(
	posxi interface{},
	posyi interface{},
	spdi interface{},
	acci interface{},
	addi interface{},
	mmti interface{},
	refi interface{},
) (err error) {
	posx, spd, acc, add, mmt, ref, err := parseRelArgs(posxi, spdi, acci, addi, mmti, refi)
	if err != nil {
		return err
	}
	posy, ok := posyi.(float64)
	if !ok {
		return fmt.Errorf("failed to convert posy %v", posyi)
	}
	log.If("moving by (%v,%v)...", posx, posy)
	if err = tml.SelectAxis(c.AxisXID); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		tml.CalcPosition(c.AxisXID, posx),
		tml.CalcSpeed(c.AxisXID, spd),
		tml.CalcAccel(c.AxisXID, acc),
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(false, false); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.AxisYID); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		tml.CalcPosition(c.AxisYID, posy),
		tml.CalcSpeed(c.AxisYID, spd),
		tml.CalcAccel(c.AxisYID, acc),
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
			if err = tml.SelectAxis(c.AxisXID); err != nil {
				log.E(err)
			}
			tml.CheckEvent(&xc)
		}
	}
	log.D("done")
	return nil
}

func (c *Client) MoveRel(
	posx float64,
	posy float64,
	speed float64,
	accel float64,
) error {
	req := Request{
		Responsec: make(chan Response),
		Function:  "MoveRelative",
		Arguments: []interface{}{
			posx,
			posy,
			speed,
			accel,
			true,
			1,
			1,
		},
	}
	c.RequestQueue.Push(&req)
	log.Df(
		"waiting for response: relative motion to (%v, %v)",
		posx,
		posy,
	)
	resp := <-req.Responsec
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (c *Client) MoveAbsolute(
	posxi interface{},
	posyi interface{},
	spdi interface{},
	acci interface{},
	mmti interface{},
	refi interface{},
) (err error) {
	posx, spd, acc, mmt, ref, err := parseAbsArgs(posxi, spdi, acci, mmti, refi)
	if err != nil {
		return err
	}
	posy, ok := posyi.(float64)
	if !ok {
		return fmt.Errorf("failed to convert posy %v", posyi)
	}
	log.If("moving to (%v,%v)...", posx, posy)
	if err = tml.SelectAxis(c.AxisXID); err != nil {
		return err
	}
	if _SET_TONPOSOK {
		log.Df("set TONPOSOK %v", tml.SetIntVariable("TONPOSOK", 100))
	}
	if err = tml.MoveAbsolute(
		tml.CalcPosition(c.AxisXID, posx),
		tml.CalcSpeed(c.AxisXID, spd),
		tml.CalcAccel(c.AxisXID, acc),
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(false, false); err != nil {
		return err
	}
	if err = tml.SelectAxis(c.AxisYID); err != nil {
		return err
	}
	if _SET_TONPOSOK {
		log.Df("set TONPOSOK %v", tml.SetIntVariable("TONPOSOK", 100))
	}
	if err = tml.MoveAbsolute(
		tml.CalcPosition(c.AxisYID, posy),
		tml.CalcSpeed(c.AxisYID, spd),
		tml.CalcAccel(c.AxisYID, acc),
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	c.CompensateMotion(c.AxisYID, posy)
	if _COMPENSATION && _COMPENSATION_ADVANCED {
		log.I("2nd compensation")
		c.CompensateMotion(c.AxisYID, posy)
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
			if err = tml.SelectAxis(c.AxisXID); err != nil {
				log.E(err)
			}
			tml.CheckEvent(&xc)
		}
	}
	log.D("done")
	return nil
}

func (c *Client) MoveAbs(
	posx float64,
	posy float64,
	speed float64,
	accel float64,
) error {
	req := Request{
		Responsec: make(chan Response),
		Function:  "MoveAbsolute",
		Arguments: []interface{}{
			posx,
			posy,
			speed,
			accel,
			1,
			1,
		},
	}
	c.RequestQueue.Push(&req)
	log.Ef(
		"waiting for response: absolute motion to (%v, %v)",
		posx,
		posy,
	)
	resp := <-req.Responsec
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func parseAbsArgs(
	posi interface{},
	spdi interface{},
	acci interface{},
	mmti interface{},
	refi interface{},
) (
	pos float64,
	spd float64,
	acc float64,
	mmt int,
	ref int,
	err error,
) {
	var ok bool
	pos, ok = posi.(float64)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert pos %v", posi)
	}
	spd, ok = spdi.(float64)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert spd %v", spdi)
	}
	acc, ok = acci.(float64)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert acc %v", acci)
	}
	mmt, ok = mmti.(int)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert mmt %v", mmti)
	}
	ref, ok = refi.(int)
	if !ok {
		return pos, spd, acc, mmt, ref,
			fmt.Errorf("failed to convert ref %v", refi)
	}
	return pos, spd, acc, mmt, ref, nil
}

func parseRelArgs(
	posi interface{},
	spdi interface{},
	acci interface{},
	addi interface{},
	mmti interface{},
	refi interface{},
) (
	pos float64,
	spd float64,
	acc float64,
	add bool,
	mmt int,
	ref int,
	err error,
) {
	pos, spd, acc, mmt, ref, err = parseAbsArgs(posi, spdi, acci, mmti, refi)
	if err != nil {
		return pos, spd, acc, add, mmt, ref, err
	}
	var ok bool
	add, ok = addi.(bool)
	if !ok {
		return pos, spd, acc, add, mmt, ref,
			fmt.Errorf("failed to convert add %v", addi)
	}
	return pos, spd, acc, add, mmt, ref, nil
}

func (c *Client) UpdateMotionStatus() (err error) {
	if c.PosX, err = tml.ActualPosition(c.AxisXID); err != nil {
		return err
	}
	if c.PosY, err = tml.ActualPosition(c.AxisYID); err != nil {
		return err
	}
	uiutil.App.UpdateMotorStatusSlot(fmt.Sprintf("Motor: (%v, %v)", c.PosX, c.PosY))
	return nil
}

func (c *Client) CompensateMotion(axisID int, target float64) (err error) {
	if !_COMPENSATION {
		return nil
	}
	switch axisID {
	case c.AxisYID:
		pos, err := tml.ActualPosition(c.AxisYID)
		diffPos := target - pos
		offset := tml.CalcPosition(c.AxisYID, diffPos)
		log.Df("compensating axis %d by %v (diff apos %v, actual pos %v)...", c.AxisYID, diffPos, offset, pos)
		if err = tml.SelectAxis(c.AxisYID); err != nil {
			return err
		}
		if err = tml.MoveRelative(
			tml.CalcPosition(c.AxisYID, diffPos),
			tml.CalcSpeed(c.AxisYID, 5),
			tml.CalcAccel(c.AxisYID, 50),
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

func (c *Client) CompensateMotionTPOS(axisID int, target float64) (err error) {
	if !_COMPENSATION {
		return nil
	}
	switch axisID {
	case c.AxisYID:
		pos, err := tml.TargetPosition(c.AxisYID)
		diffPos := target - pos
		offset := tml.CalcPosition(c.AxisYID, diffPos)
		if diffPos == 0 {
			log.Df("not compensate axis %d by %v (diff tpos %v, actual pos %v)...", c.AxisYID, diffPos, offset, pos)
			return nil
		}
		log.Df("compensating axis %d by %v (diff tpos %v, actual pos %v)...", c.AxisYID, diffPos, offset, pos)
		if err = tml.SelectAxis(c.AxisYID); err != nil {
			return err
		}
		if err = tml.MoveRelative(
			tml.CalcPosition(c.AxisYID, diffPos),
			tml.CalcSpeed(c.AxisYID, 5),
			tml.CalcAccel(c.AxisYID, 50),
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
