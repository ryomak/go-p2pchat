package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ryomak/go-p2pchat/control"
	"github.com/ryomak/go-p2pchat/peer"
)

var (
	port   = flag.String("p", "1111", "string flag")
	user   = flag.String("user", "name:192.168.10.1:1111", "string flag")
	myName = flag.String("name", "ryomak", "string flag")
)
var ctrl = control.Control{
	UpdatedText:         make(chan string, 10),
	UpdateUserList:      make(chan []peer.User, 10),
	UpdatedTextFromUser: make(chan string, 10)}

func init() {
	flag.Parse()
	peer.SetMyName(*myName)
}

func main() {
	//if *ip != util.GetMyIP() {
	u := strings.Split(*user, ":")
	go peer.IntroduceMyself(peer.User{
		Name: u[0],
		IP:   u[1],
		Port: u[2],
	})
	//}
	go ctrl.StartControlLoop()
	ctrl.UpdatedText <- "Hello " + peer.GetMyName() + ".\nFor private messages, type the message followed by * and the name of the receiver.\n To leave the conversation type disconnect"
	go peer.RunServer(*port, ctrl.UpdatedText, ctrl.UpdateUserList)
	UserInput()
}

func UserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		ctrl.UpdatedTextFromUser <- text
	}
}
