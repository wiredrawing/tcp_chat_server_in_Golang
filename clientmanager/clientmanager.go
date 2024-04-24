package clientmanager

import (
	"errors"
	"net"
	"wiredrawing/go/socket-application/clientunit"
)

type ClientManager struct {
	ClientList map[net.Addr]*clientunit.ClientUnit
}

func (cm *ClientManager) AddClient(refClient *clientunit.ClientUnit) (bool, error) {
	var address = refClient.Connection.RemoteAddr()
	if _, ok := cm.ClientList[address]; ok != true {
		cm.ClientList[address] = refClient
		return true, nil
	}
	return false, nil
}

func (cm *ClientManager) addClient(refClient *clientunit.ClientUnit) (bool, error) {
	var address = refClient.Connection.RemoteAddr()
	if _, ok := cm.ClientList[address]; ok != true {
		cm.ClientList[address] = refClient
		return true, nil
	}
	return false, errors.New("クライアントは既に存在しています")
}
func (cm *ClientManager) RemoveClient(client *clientunit.ClientUnit) {
	delete(cm.ClientList, client.Connection.RemoteAddr())
}
func (cm *ClientManager) Exists(client *clientunit.ClientUnit) bool {
	if _, ok := cm.ClientList[client.Connection.RemoteAddr()]; ok {
		return ok
	} else {
		return false
	}
}
