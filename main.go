package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
)

var printf = fmt.Printf

// 接続してきたクライアントを管理する構造体
type ClientUnit struct {
	clientName string
	connection net.Conn
}

type ClientManager struct {
	clientList map[net.Addr]ClientUnit
}

func (cm *ClientManager) addClient(client ClientUnit) {
	fmt.Printf("---------%v---------------", reflect.TypeOf(client.connection.RemoteAddr()))
	fmt.Printf("---------%v---------------", client.connection.RemoteAddr())
	cm.clientList[client.connection.RemoteAddr()] = client
}
func (cm *ClientManager) removeClient(client ClientUnit) {
	delete(cm.clientList, client.connection.RemoteAddr())
}
func (cm *ClientManager) exists(client ClientUnit) bool {
	if _, ok := cm.clientList[client.connection.RemoteAddr()]; ok {
		return ok
	} else {
		return false
	}
}

// TCPクライアントを管理する構造体を作成
var clientManager = ClientManager{
	clientList: make(map[net.Addr]ClientUnit),
}

func main() {
	_, err := net.ResolveTCPAddr("tcp", ":10080")
	if err != nil {
		log.Fatal(err)
	}
	// Create a listener on TCP port 10080
	listener, err := net.Listen("tcp", ":10080")
	if err != nil {
		panic(err)
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			panic(err)
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		// 新規接続クライアントオブジェクトを作成
		var client = ClientUnit{
			clientName: "",
			connection: connection,
		}

		isExists := clientManager.exists(client)
		if isExists == false {
			clientManager.addClient(client)
		}

		go handlingConnection(client)
	}
}

func handlingConnection(clientUnit ClientUnit) error {
	var c net.Conn = clientUnit.connection
	printf("%v", c.RemoteAddr())
	var fromAddress string
	c.Write([]byte("あなたのお名前を最初に入力してください\n"))
	printf("読み取り開始\n")
	c.SetReadDeadline(time.Now().Add(1 * time.Second))
	for {
		stackBuffer := make([]byte, 0, 1024)
		buffer := make([]byte, 32)
		//printf("待ち受けかいし ")
		for {
			// Read the incoming connection into the buffer

			size, err := c.Read(buffer)
			var opError *net.OpError
			if err != nil {
				if errors.As(err, &opError) == true {
					c.SetReadDeadline(time.Now().Add(1 * time.Second))
					//printf("これはタイムアウトのエラーです")
				} else {
					panic(err)
				}
			}

			//} else {
			//	printf("これはタイムアウトのエラーではありません")
			//	panic(err)
			//	break
			//}
			//if err == io.EOF {
			//	printf("EOF")
			//	break
			//}
			//if err != nil {
			//	printf("%v", reflect.TypeOf(err))
			//	printf("%v", err.Error())
			//
			//}
			//if err != nil {
			//	if err == io.EOF {
			//		break
			//		// io.EOF, etc
			//	} else if err.(*net.OpError).Timeout() {
			//		// no status msgs
			//		// note: TCP keepalive failures also cause this; keepalive is on by default
			//		continue
			//	} else {
			//		panic(err)
			//	}
			//}
			//if err != nil {
			//	//printf("err => %v\n", err)
			//	break
			//}

			//if size == 0 {
			//	printf(string(buffer))
			//	printf("size => %v\n", size)
			//	printf("読み取り完了")
			//	break
			//printf("読み込んだbyteサイズ ->"+
			//	" %v", size)
			//printf("スライシング前%v", buffer)
			buffer = buffer[:size]
			//printf("スライシング後%v", buffer)
			stackBuffer = append(stackBuffer, buffer...)
			if size < 32 {
				break
			}
		}
		//panic(",あああああああああああ")
		//buffer := make([]byte, 512)
		////for {
		//if size, err := c.Read(buffer); err != nil {
		//	printf("size => %v\n", size)
		//	printf("err => %v\n", err)
		//	printf("buffer => %v\n", buffer)
		//	printf("%v", err)
		//	break
		//} else {
		//	buffer = buffer[:size]
		//	printf("buffer => %v\n", buffer)
		//	stackBuffer = append(stackBuffer, buffer...)
		//	if size < 32 {
		//		//break
		//	}
		//}
		////}
		//panic("stop")
		if len(stackBuffer) == 0 {
			//printf("stackBuffer => %v\n", stackBuffer)
			continue
		}
		recievedMessage := string(stackBuffer)
		// ClientUnit構造体のclientNameが空の場合は、クライアント名を登録
		if len(clientUnit.clientName) == 0 {
			// 名前を設定
			clientUnit.clientName = recievedMessage
			c.Write([]byte("ようこそ" + clientUnit.clientName + "さん\n"))

			// 発信者以外のユーザーに接続開始した旨を通知
			for address, value := range clientManager.clientList {
				if address != clientUnit.connection.RemoteAddr() {
					value.connection.Write([]byte(clientUnit.clientName + "さんが入室しました\n"))
				}
			}
			continue
		} else {
			for _, client := range clientManager.clientList {
				if client.connection.RemoteAddr() != clientUnit.connection.RemoteAddr() {
					client.connection.Write([]byte(clientUnit.clientName + "さんが発言しました: " + recievedMessage + "\n"))
				}
			}
		}
		//printf("%v\n", stackBuffer)
		//printf("recievedMessage => <%s>\n", recievedMessage)
		//printf("recievedMessage => <%v>\n", recievedMessage)
		//panic("!!!!!!!!!!!!!!")

		// exit or end を入力されたときはクライアントの接続を切断
		if recievedMessage == "exit" || recievedMessage == "end" {
			printf("クライアントリストから削除")
			c.Write([]byte("ご利用ありがとうございました!!\n"))
			// クライアントリストから削除
			clientManager.removeClient(clientUnit)
			if err := c.Close(); err != nil {
				panic(err)
			}
			return nil
		} else {
			// Print the request to the console
			printf("接続元アドレス: %s\n", fromAddress)
			printf("受信したデータ: %s\n", string(stackBuffer))
			//var reply string = fmt.Sprintf("%sさんが発言しました: %s\n", clientUnit.clientName, string(stackBuffer))
			//c.Write([]byte(reply))
		}
		thanksBuffer := []byte("Thanks for your message\n")
		if _, err := c.Write(thanksBuffer); err != nil {
			panic(err)
		} else {
			printf("メッセージを受信しました")
		}
	}
}
