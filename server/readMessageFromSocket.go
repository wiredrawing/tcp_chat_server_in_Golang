package server

import (
	"net"
	"time"
)

const BufferMaxSize int = 32

// ReadMessageFromSocket ソケットからメッセージを読み取る
// 引数に指定したソケットからメッセージを読み取る処理を別関数に切り出す
func ReadMessageFromSocket(socket net.Conn, timeout int) string {
	// タイムアウトの時間を指定する
	duration := time.Duration(timeout)
	_ = socket.SetReadDeadline(time.Now().Add(duration * time.Second))
	//fmt.Printf("名前の読みとり開始 ----------------------------->")
	// 長さ0の空っぽの配列を確保する場合は以下のように記述
	buffer := make([]byte, BufferMaxSize)
	gotBuffer := make([]byte, 0, 1024)
	var size int
	var err error
	for {
		// 読み取るまでブロックしてしまうため
		// 最初にタイムアウトを設定しておく
		size, err = socket.Read(buffer)
		if err != nil {
			//fmt.Printf("読み取りサイズ: %d\n", size)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				//log.Print("タイムアウトを延長")
				//log.Printf("タイムアウトエラーの型チェック: %T\n", err)
				_ = socket.SetReadDeadline(time.Now().Add(duration * time.Second))
			}
		}
		// 受け取ったbufferをまとめる
		buffer = buffer[:size]
		gotBuffer = append(gotBuffer, buffer...)
		if size < BufferMaxSize {
			break
		}
	}

	//fmt.Printf("名前の読みとり終了 ----------------------------->")
	//fmt.Printf("名前の読みとり結果: %v\n", string(gotBuffer))
	return string(gotBuffer)
}
