package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

var printf = fmt.Printf

var clientList = make(map[net.Addr]net.Conn)

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
		address := connection.RemoteAddr()

		// 接続元クライアント管理
		var isExist = false
		for addr := range clientList {
			// 既にクライアントリストに登録されているか確認
			if addr.String() == address.String() {
				isExist = true
				break
			}
		}
		if isExist == false {
			clientList[address] = connection
		}
		printf("%v\n", clientList)
		//connection.SetReadDeadline(time.Now().Add(5 * time.Second))

		go handlingConnection(connection)
	}
}

func handlingConnection(c net.Conn) {

	printf("%v", c.RemoteAddr())
	var fromAddress string
	c.Write([]byte("--------->接続を開始しました!\n: "))
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
		//printf("%v\n", stackBuffer)
		//printf("recievedMessage => <%s>\n", recievedMessage)
		//printf("recievedMessage => <%v>\n", recievedMessage)
		//panic("!!!!!!!!!!!!!!")
		if recievedMessage == "exit" || recievedMessage == "end" {
			printf("クライアントリストから削除")
		}
		if recievedMessage == "exit\n" || recievedMessage == "end" {
			printf("クライアントリストから削除")
			c.Write([]byte("Goodbye\n"))
			err := c.Close()
			if err != nil {
				panic(err)
			}
		} else {
			// Print the request to the console
			printf("接続元アドレス: %s\n", fromAddress)
			printf("受信したデータ: %s\n", string(stackBuffer))
		}
		thanksBuffer := []byte("Thanks for your message\n")
		if _, err := c.Write(thanksBuffer); err != nil {
			panic(err)
		} else {
			printf("メッセージを受信しました")
		}
		//}

	}
}
