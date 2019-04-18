package usbcan

var (
	DeviceMap  = deviceMap
	ChannelMap = channelMap
	ClientMap  = clientMap
)

func AddInstance(client *Client) (*Client, bool) {
	return addInstance(client)
}

func (c *Channel) Send() {
	c.send()
}
