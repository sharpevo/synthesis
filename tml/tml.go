package tml

import (
	"fmt"
	"log"
	"posam/util/blockingqueue"
	"posam/util/concurrentmap"
	"reflect"
	"time"
	"tml"
)

var clientMap *concurrentmap.ConcurrentMap

func init() {
	clientMap = concurrentmap.NewConcurrentMap()
}

func Instance(key string) *Client {
	if clienti, ok := clientMap.Get(key); ok {
		return clienti.(*Client)
	}
	return nil
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
		log.Println("terminating client: ", client.Name)
		//client.Stop()
	}
	clientMap = concurrentmap.NewConcurrentMap()
}

type Clienter interface {
	connect() error
	MoveRelative(int, int, float64, float64, float64) error
	MoveAbsolute(int, int, float64, float64, float64) error
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
	go client.launch()
	log.Println(">>> client: ", client.AxisXID, client.AxisYID)
	return client, nil
}

func (c *Client) connect() (err error) {
	log.Println("Connecting the motor %q...", c.Name)

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

	log.Println("checking status...")
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
		time.Sleep(1 * time.Second)
	}
	if statusx == 0 || statusy == 0 {
		return fmt.Errorf("failed to enable power on axes: x(%d) / y(%d)", statusx, statusy)
	}

	log.Printf("motor %q is ready", c.Name)
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

func (c *Client) launch() {
	log.Println("motor client launched")
	c.connect()
	for {
		reqi, err := c.RequestQueue.Pop()
		if err != nil {
			log.Println("motor client terminated")
			return
		}
		req := reqi.(*Request)
		//err = c.checkAxisReady(req.AxisID)
		//if err != nil {
		//req.Responsec <- Response{Error: err}
		//continue
		//}
		//c.execute(req)
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
		req.Responsec <- Response{Error: nil}
		if err := c.UpdateMotionStatus(); err != nil {
			log.Println(err)
		}
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
	fmt.Printf("moving axis %d to %v...", aid, pos)
	if err = tml.SelectAxis(aid); err != nil {
		return err
	}
	if err = tml.MoveAbsolute(
		pos,
		spd,
		acc,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	fmt.Printf("done\n")
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
	fmt.Printf("moving axis %d by %v...", aid, pos)
	if err = tml.SelectAxis(aid); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		pos,
		spd,
		acc,
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	fmt.Printf("done\n")
	return nil
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
	log.Printf(
		"waiting for axis %d response: absolute motion to %v\n",
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
	log.Printf(
		"waiting for axis %d response: relative motion to %v\n",
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
	fmt.Printf("moving by (%v,%v)...", posx, posy)
	if err = tml.SelectAxis(c.AxisXID); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		posx,
		spd,
		acc,
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(false, false); err != nil {
		return err
	}
	xcompleted := false
	if err = tml.SelectAxis(c.AxisYID); err != nil {
		return err
	}
	if err = tml.MoveRelative(
		posy,
		spd,
		acc,
		add,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	for {
		time.Sleep(200 * time.Millisecond)
		if err = tml.SelectAxis(c.AxisXID); err != nil {
			log.Println(err)
		}
		tml.CheckEvent(&xcompleted)
		if xcompleted {
			break
		}
	}
	fmt.Println("done")
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
	log.Printf(
		"waiting for response: relative motion to (%v, %v)\n",
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
	fmt.Printf("moving to (%v,%v)...", posx, posy)
	if err = tml.SelectAxis(c.AxisXID); err != nil {
		return err
	}
	if err = tml.MoveAbsolute(
		posx,
		spd,
		acc,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(false, false); err != nil {
		return err
	}
	xcompleted := false
	if err = tml.SelectAxis(c.AxisYID); err != nil {
		return err
	}
	if err = tml.MoveAbsolute(
		posy,
		spd,
		acc,
		mmt,
		ref,
	); err != nil {
		return err
	}
	if err = tml.SetEventOnMotionComplete(true, false); err != nil {
		return err
	}
	for {
		time.Sleep(200 * time.Millisecond)
		if err = tml.SelectAxis(c.AxisXID); err != nil {
			log.Println(err)
		}
		tml.CheckEvent(&xcompleted)
		if xcompleted {
			break
		}
	}
	fmt.Println("done")
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
	log.Printf(
		"waiting for response: absolute motion to (%v, %v)\n",
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
	if err = tml.SelectAxis(c.AxisXID); err != nil {
		return err
	}
	var posx float64
	if err = tml.GetLongVariable("APOS", &posx); err != nil {
		return err
	}
	c.PosX = tml.ParsePosition(posx)
	if err = tml.SelectAxis(c.AxisYID); err != nil {
		return err
	}
	var posy float64
	if err = tml.GetLongVariable("APOS", &posy); err != nil {
		return err
	}
	c.PosY = tml.ParsePosition(posy)
	return nil
}
