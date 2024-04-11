package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var printf = fmt.Printf

// 接続してきたクライアントを管理する構造体
type ClientUnit struct {
	clientName string
	connection net.Conn
}

type ClientManager struct {
	clientList map[net.Addr]*ClientUnit
}

func (cm *ClientManager) addClient(refClient *ClientUnit) (bool, error) {
	var address net.Addr = refClient.connection.RemoteAddr()
	if _, ok := cm.clientList[address]; ok != true {
		cm.clientList[address] = refClient
		return true, nil
	}
	return false, errors.New("クライアントは既に存在しています")
}
func (cm *ClientManager) removeClient(client *ClientUnit) {
	delete(cm.clientList, client.connection.RemoteAddr())
}
func (cm *ClientManager) exists(client *ClientUnit) bool {
	if _, ok := cm.clientList[client.connection.RemoteAddr()]; ok {
		return ok
	} else {
		return false
	}
}

// TCPクライアントを管理する構造体を作成
var clientManager = ClientManager{
	clientList: make(map[net.Addr]*ClientUnit),
}

func fetchReceiveBufferFromServer(connection *net.TCPConn) {
	const ByteSize = 1024
	for {
		// 読み取り開始
		//fmt.Println("TCPサーバーからの読み取り開始---")
		readBytes := make([]byte, ByteSize)
		_ = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
		for {
			buffer := make([]byte, ByteSize)
			size, err := connection.Read(buffer)
			if err != nil {
				// timeoutエラーを検出
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					//fmt.Printf("Timeout時間を延慶")
					_ = connection.SetReadDeadline(time.Now().Add(5 * time.Second))
				} else {
					panic(err)
				}
			}
			// 読み取ったバイト数のみスライスする
			buffer = buffer[:size]
			readBytes = append(readBytes, buffer...)
			if (size < ByteSize) && (size > 0) {
				break
			}
		}
		fmt.Printf("サーバーからのレスポンス-- %v", string(readBytes))
	}
	connection.Write([]byte("TCPクライアントから接続 -------------------------->"))
}

func main() {
	// コマンドライン引数を取得
	asClient := flag.Bool("client", false, "TCPクライアントとして実行する場合")
	ipAddress := flag.String("address", "127.0.0.1", "接続先あるいはListen先のIPアドレス")
	portNumber := flag.Int("port", 10080, "接続先あるいはListen先のポート番号")
	flag.Parse()

	if (*asClient) == true {
		// 読み取らせるマックスbyte数

		fmt.Println("TCPクライアントの起動---")
		// TCPクライアントとして起動させた場合
		addr := &net.TCPAddr{
			IP:   net.ParseIP(*ipAddress),
			Port: *portNumber,
		}
		connection, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			panic(err)
		}

		go fetchReceiveBufferFromServer(connection)

		for {
			var prompt = " >>> "
			fmt.Printf(prompt)
			scanner := bufio.NewScanner(os.Stdin)
			result := scanner.Scan()
			if result == true {
				var input string = scanner.Text()
				// 入力したテキストをサーバーへ送信
				if _, err := connection.Write([]byte(input)); err != nil {
					panic(err)
				}
			}
		}

	} else {
		tcp, err := net.ResolveTCPAddr("tcp", ":10080")
		tcp.IP = net.ParseIP(*ipAddress)
		tcp.Port = *portNumber
		if err != nil {
			log.Fatal(err)
		}
		// Create a listener on TCP port 10080
		listener, err := net.ListenTCP("tcp", tcp)
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
			var client *ClientUnit = new(ClientUnit)
			client.clientName = ""
			client.connection = connection

			isExists := clientManager.exists(client)
			fmt.Printf("isExists => %v\n", isExists)
			if isExists == false {
				clientManager.addClient(client)
			}
			fmt.Printf("クライアントリストclientManager: %v\n", clientManager)
			fmt.Printf("クライアントリストclientManager.clientList: %v\n", clientManager.clientList)
			go handlingConnection(client)
		}
	}
}

func handlingConnection(clientUnit *ClientUnit) error {
	var c net.Conn = clientUnit.connection
	printf("%v", c.RemoteAddr())
	c.Write([]byte("<< あなたのお名前を最初に入力してください >>\n"))
	c.Write([]byte("<< ユーザー一覧を表示する場合は、usersと入力してください >>\n"))
	printf("読み取り開始\n")
	if err := c.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
		panic(err)
	}
	for {
		stackBuffer := make([]byte, 0, 1024)
		buffer := make([]byte, 32)
		//printf("待ち受けかいし ")
		for {

			size, err := c.Read(buffer)
			//var opError *net.OpError
			if err != nil {
				//if errors.Is(err, &net.OpError{}) {
				//	printf("IsメソッドがTrueを返す場合")
				//}
				//e := errors.Unwrap(err)
				//printf("e ===> %v ]]]]\n", e)
				//printf("エラー発生 %v ===================== \n", reflect.TypeOf(e))
				//fmt.Printf("errors.Unwrap(err) => %v >>> \n", errors.Unwrap(err))
				//fmt.Printf("reflect.TypeOf(err) => %v >>> \n", reflect.TypeOf(err))
				//fmt.Printf("reflect.TypeOf(err).Comparable() ==> %v >>> \n", reflect.TypeOf(err).Comparable())
				var ope = new(net.OpError)
				if ok := errors.As(err, &(ope)); ok {
				}
				if _, ok := err.(*net.OpError); ok {
					//fmt.Printf(" エラー時 %v\n", reflect.TypeOf(err))
					// 型アサーション成功時<エラー型>が指定したエラーの場合
					//fmt.Printf(" タイムアウトエラー時 %v\n", err)
					if err := c.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
						panic(err)
					}
				}
			}

			buffer = buffer[:size]
			//printf("スライシング後%v", buffer)
			stackBuffer = append(stackBuffer, buffer...)
			if size < 32 {
				break
			}
		}

		if len(stackBuffer) == 0 {
			continue
		}
		recievedMessage := string(stackBuffer)

		if recievedMessage == "users" {
			users := make([]string, len(clientManager.clientList))
			for address, value := range clientManager.clientList {
				s := fmt.Sprintf("[%s]:%s", address, value.clientName)
				users = append(users, s)
			}
			fmt.Printf("接続中のユーザー: %v\n", users)
			connectedUsers := strings.Join(users, "\n")
			fmt.Printf("接続中のユーザー: %v\n", connectedUsers)
			c.Write([]byte("接続中のユーザー: " + connectedUsers + "\n"))
			continue
		}
		// ClientUnit構造体のclientNameが空の場合は、クライアント名を登録
		if len((*clientUnit).clientName) == 0 {
			// 名前を設定
			(*clientUnit).clientName = recievedMessage
			c.Write([]byte("ようこそ" + clientUnit.clientName + "さん\n"))

			// 発信者以外のユーザーに接続開始した旨を通知
			for address, value := range clientManager.clientList {
				if address != clientUnit.connection.RemoteAddr() {
					value.connection.Write([]byte(clientUnit.clientName + "さんが入室しました\n"))
				}
			}
			continue
		} else {
			for address, client := range clientManager.clientList {
				if address != clientUnit.connection.RemoteAddr() {
					formattedMessage := colorWrapping("33", clientUnit.clientName+"さんが発言しました: "+recievedMessage+"\n")
					client.connection.Write([]byte(formattedMessage))
				}
			}
		}
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
			printf("%s: %s\n", clientUnit.clientName, string(stackBuffer))
			//printf("接続元アドレス: %s\n", fromAddress)
			//printf("受信したデータ: %s\n", string(stackBuffer))
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

func colorWrapping(colorCode string, text string) string {
	return "\033[" + colorCode + "m" + text + "\033[0m"
}
