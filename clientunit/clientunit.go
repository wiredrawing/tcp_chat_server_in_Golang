package clientunit

import "net"

// ClientUnit 接続してきたクライアントを管理する構造体
type ClientUnit struct {
	ClientName string
	Connection net.Conn
}

// FetchClientAddress クライアントのアドレスを取得する
func (c *ClientUnit) FetchClientAddress() net.Addr {
	return c.Connection.RemoteAddr()
}

// FetchClientName クライアント名を取得する
func (c *ClientUnit) FetchClientName() string {
	return c.ClientName
}

// FetchConnection クライアントの接続情報を取得する
func (c *ClientUnit) FetchConnection() net.Conn {
	return c.Connection
}
