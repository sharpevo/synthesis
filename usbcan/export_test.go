package usbcan

var (
	DeviceMap  = deviceMap
	ChannelMap = channelMap
	ClientMap  = clientMap
)

func AddInstance(client *Client) (*Client, bool) {
	return addInstance(client)
}

func (c *Channel) Init() {
	c.init()
}

func (c *Channel) Send() {
	c.send()
}

func (c *Channel) UntilSendable() {
	c.untilSendable()
}

func Instance(key string) (client *Client) {
	return instance(key)
}
