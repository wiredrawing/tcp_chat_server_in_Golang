package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"go-socket/clientmanager"
	"go-socket/clientunit"
	"go-socket/server"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var printf = fmt.Printf

//type ClientManager struct {
//	clientList map[net.Addr]*clientunit.ClientUnit
//}

// TCPクライアントを管理する構造体を作成
var clientManager = clientmanager.ClientManager{
	ClientList: make(map[net.Addr]*clientunit.ClientUnit),
}

// ログイン可能なユーザーリスト
var userList = make([]map[string]string, 16)

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
				var netError = new(net.Error)
				if errors.As(err, netError) && (*netError).Timeout() {
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

	// ログイン可能なユーザーリストを定義しておく
	userList = append(userList, map[string]string{"user": "user1", "password": "password1"})
	userList = append(userList, map[string]string{"user": "user2", "password": "password2"})
	userList = append(userList, map[string]string{"user": "user3", "password": "password3"})

	// ログファイルを作成する
	logFile, err := os.OpenFile("./log.dat", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("ログファイルが作成できませんでした %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// コマンドライン引数を取得
	asClient := flag.Bool("client", false, "TCPクライアントとして実行する場合")
	ipAddress := flag.String("address", "127.0.0.1", "接続先あるいはListen先のIPアドレス")
	portNumber := flag.Int("port", 10080, "接続先あるいはListen先のポート番号")
	flag.Parse()

	if (*asClient) == true {
		// 読み取らせるマックスbyte数

		fmt.Println("TCPクライアントの起動---")

		// IPを直接指定させる
		fmt.Println("接続先のIPアドレスを入力してください...")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if (*scanner).Scan() {
				var connectToIp = scanner.Text()
				// IPアドレスの形式チェック
				var ip = net.ParseIP(connectToIp)
				if (ip.String() == connectToIp) && (ip.To4() != nil) {
					fmt.Printf("接続先を[%s]に確定しました。\n", connectToIp)
					*ipAddress = connectToIp
					break
				} else {
					fmt.Println("IPアドレスの形式が不正です。再度入力してください...")
					continue
				}
			}
		}
		fmt.Println("接続先のポート番号を入力してください...")
		for {
			// 標準入力から文字列を取得する
			scanner := bufio.NewScanner(os.Stdin)
			var result = scanner.Scan()
			if result == true {
				var input = scanner.Text()
				if ip, err := strconv.Atoi(input); err == nil {
					if (ip >= 1024) && (ip <= 65535) {
						fmt.Printf("ポート番号を[%d]に確定しました。\n", ip)
						*portNumber = ip
						break
					}
				} else {
					fmt.Printf("エラー: %v\n", err)
				}
				fmt.Println("妥当なポート番号を入力して下さい")
			}
		}
		// TCPクライアントとして起動させた場合
		addr := &net.TCPAddr{
			IP:   net.ParseIP(*ipAddress),
			Port: *portNumber,
		}

		connection, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			fmt.Printf("TCPサーバー[%s:%s]接続エラー: %v\n", *ipAddress, *portNumber, err)
			panic(err)
		}

		go fetchReceiveBufferFromServer(connection)

		for {
			var prompt = " >>> "
			fmt.Printf(prompt)
			scanner := bufio.NewScanner(os.Stdin)
			result := scanner.Scan()
			if result == true {
				var input = scanner.Text()
				// 入力したテキストをサーバーへ送信
				if _, err := connection.Write([]byte(input)); err != nil {
					panic(err)
				}
			}
		}

	} else {
		fmt.Println("> Starting TCP server...")
		fmt.Println("> Started to accept the connection from the client...")

		fmt.Println("ListenするIPアドレスを入力してください...")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			var result = scanner.Scan()
			if result == true {
				var input = scanner.Text()
				var ip = net.ParseIP(input)
				fmt.Println(ip.String())
				if (ip.String() == input) && (ip.To4() != nil) {
					fmt.Printf("%sをListenします\n", input)
					*ipAddress = input
					break
				} else {
					fmt.Println("IPアドレスの形式が不正です。再度入力してください...")
					continue
				}
			}
		}

		fmt.Println("Listenするポート番号を入力してください...")
		for {
			scanner := bufio.NewScanner(os.Stdin)
			var result = scanner.Scan()
			if result == true {
				// String from console.
				var input = scanner.Text()
				if port, err := strconv.Atoi(input); err == nil {
					*portNumber = port
					fmt.Printf("%dをListenします\n", port)
					break
				} else {
					fmt.Printf("妥当なポート番号を入力してください\n")
					fmt.Printf("エラー: %v\n", err)
				}
			}
		}

		hostName := fmt.Sprintf("%s:%d", *ipAddress, *portNumber)

		tcp, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", *ipAddress, *portNumber))
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

		var serverMessagenger = false
		for {
			connection, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			// 新規接続クライアントオブジェクトを作成
			var client = new(clientunit.ClientUnit)
			client.ClientName = ""
			client.Connection = connection
			fmt.Printf("client オブジェクト -> %v", client)

			isExists := clientManager.Exists(client)
			fmt.Printf("isExists => %v\n", isExists)
			if isExists == false {
				clientManager.AddClient(client)
			}
			fmt.Printf("クライアントリストclientManager: %v\n", clientManager)
			fmt.Printf("クライアントリストclientManager.clientList: %v\n", clientManager.ClientList)
			// サーバーからの送信処理は一度だけ実行
			if serverMessagenger == false {
				fmt.Printf("サーバーからのメッセージ送信処理を開始します\n")
				go SendMessageFromServer(clientManager)
				serverMessagenger = true
			}
			go handlingConnection(client, hostName)
		}
	}
}

// authenticate ユーザー名とパスワードを認証する
func authenticate(userName string, password string) bool {
	var isAuthenticated = false
	for _, user := range userList {
		if user["user"] == userName && user["password"] == password {
			isAuthenticated = true
			break
		}
	}
	return isAuthenticated
}

func login(c net.Conn) map[string]string {
	var userInfo map[string]string = make(map[string]string)
	var messageToClient string = ""
	messageToClient = "username: "
	_, _ = c.Write([]byte(messageToClient))
	var userName = ""
	for {
		userName = server.ReadMessageFromSocket(c, 2)
		if len(userName) > 0 {
			break
		}
	}
	userInfo["user"] = userName
	var password = ""
	messageToClient = "password: "
	_, _ = c.Write([]byte(messageToClient))
	for {
		password = server.ReadMessageFromSocket(c, 2)
		if len(password) > 0 {
			break
		}
	}
	userInfo["password"] = password
	return userInfo
}
func handlingConnection(clientUnit *clientunit.ClientUnit, hostName string) error {
	var messageToClient string
	var c = clientUnit.Connection
	printf("クライアント接続元情報[%v]\n", c.RemoteAddr())
	messageToClient = fmt.Sprintf("[%s]ようこそ\n", hostName)
	_, _ = c.Write([]byte(messageToClient))

	// 認証処理を行う ----------------------------------------------------
	for {
		userInfo := login(c)
		ok := authenticate(userInfo["user"], userInfo["password"])
		if ok == true {
			_, _ = c.Write([]byte("認証処理に成功しました\n"))
			break
		}
		_, _ = c.Write([]byte("認証処理に失敗しました."))
		continue
		//_ = c.Close()
		//return errors.New("認証処理に失敗しました")
	}

	var socketName string
	for {
		socketName = server.ReadMessageFromSocket(c, 2)
		if len(socketName) > 0 {
			// クライアント名を登録
			fmt.Printf("クライアント[%s]さんが入室しました.\n", socketName)
			fmt.Printf("socketName => %v\n", socketName)
			messageToClient = fmt.Sprintf("こんにちわ[%s]さん 楽しんでね!\n", socketName)
			(*clientUnit).ClientName = socketName
			for key, value := range clientManager.ClientList {
				if key == clientUnit.Connection.RemoteAddr() {
					// 発言者本人へのリプライ
					_, _ = value.Connection.Write([]byte(messageToClient))
				} else {
					// 本人以外へのリプライ
					messageToClient = fmt.Sprintf("[%s]さんが入室しました\n", socketName)
					_, _ = value.Connection.Write([]byte(messageToClient))
				}
			}
			_, _ = c.Write([]byte(messageToClient))
			break
		}
	}
	c.Write([]byte("<< ユーザー一覧を表示する場合は、usersと入力してください >>\n"))
	printf("読み取り開始\n")
	for {
		recievedMessage := server.ReadMessageFromSocket(c, 2)
		if len(recievedMessage) == 0 {
			continue
		}

		if recievedMessage == "users" {
			users := make([]string, len(clientManager.ClientList))
			for address, value := range clientManager.ClientList {
				s := fmt.Sprintf("[%s]:%s", address, value.ClientName)
				users = append(users, s)
			}
			fmt.Printf("接続中のユーザー: %v\n", users)
			connectedUsers := strings.Join(users, "\n")
			fmt.Printf("接続中のユーザー: %v\n", connectedUsers)
			c.Write([]byte("接続中のユーザー: " + connectedUsers + "\n"))
			continue
		}
		// ClientUnit構造体のclientNameが空の場合は、クライアント名を登録
		if len((*clientUnit).ClientName) > 0 {
			for address, client := range clientManager.ClientList {
				if address != clientUnit.Connection.RemoteAddr() {
					formattedMessage := colorWrapping("33", clientUnit.ClientName+"さんが発言しました: "+recievedMessage+"\n")
					client.Connection.Write([]byte(formattedMessage))
				}
			}
		}
		// exit or end を入力されたときはクライアントの接続を切断
		if recievedMessage == "exit" || recievedMessage == "end" {
			printf("クライアントリストから削除")
			c.Write([]byte("ご利用ありがとうございました!!\n"))
			// クライアントリストから削除
			clientManager.RemoveClient(clientUnit)
			if err := c.Close(); err != nil {
				panic(err)
			}
			return nil
		} else {
			// Print the request to the console
			printf("%s: %s\n", clientUnit.ClientName, recievedMessage)
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

// SendMessageFromServer この関数はサーバーからクライアントへメッセージを送信する
func SendMessageFromServer(clientManager clientmanager.ClientManager) {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("メッセージを入力してください on Server: ")
		if input.Scan() {
			var message = input.Text()
			for _, value := range clientManager.ClientList {
				_, _ = value.Connection.Write([]byte(message))
				fmt.Printf("メッセージ: %v\n", message)
			}
		}
	}
}
