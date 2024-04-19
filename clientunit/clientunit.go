package clientunit

import "net"

// ClientUnit 接続してきたクライアントを管理する構造体
type ClientUnit struct {
	ClientName string
	Connection net.Conn
}

// FetchClientName クライアント名を取得する
func (c *ClientUnit) FetchClientName() string {
	return c.ClientName
}

func (c *ClientUnit) FetchConnection() net.Conn {
	return c.Connection
}
