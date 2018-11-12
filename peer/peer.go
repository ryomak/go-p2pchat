package peer

import (
	"encoding/json"
	"net"

	log "github.com/Sirupsen/logrus"

	"github.com/ryomak/go-p2pchat/util"
)

var Dev = true
var usersConnectionsMap = make(map[string]net.Conn)
var usersMap = make(map[string]User)
var myName string
var myPort string

var (
	updateTextChan     chan string
	updateUserListChan chan []User
)

type Message struct {
	Kind  string
	Me    User
	MSG   string
	Users []User
}

type User struct {
	Name string
	IP   string
	Port string
}

func GetMyName() string {
	return myName
}

func SetMyName(name string) {
	myName = name
}

//send message to all peer
func (msg *Message) Send() {
	if Dev {
		log.Println("send")
	}
	log.Println(usersConnectionsMap)
	for user, conn := range usersConnectionsMap {
		if user == myName {
			continue
		}
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}

func (msg *Message) SendToUser(receiver string, updateTextStream chan string) {
	log.Info("sendToUser")
	if _, isExist := usersMap[receiver]; isExist {
		conn := usersConnectionsMap[receiver]
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
	updateTextStream <- (receiver + "is not conn")
}

func RunServer(port string, updateTextCh chan string, updateUserList chan []User) {
	log.Info("starting 'server'")
	myPort = port
	updateTextChan = updateTextCh
	updateUserListChan = updateUserList
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+port)
	if err != nil {
		log.Println("netResolve error")
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Info(err)
			continue
		}
		go receive(conn)
	}
}

func receive(conn net.Conn) {
	log.Println("receive")
	defer conn.Close()
	dec := json.NewDecoder(conn)
	msg := new(Message)
	for {
		if err := dec.Decode(msg); err != nil {
			log.Println(err)
			return
		}
		switch msg.Kind {
		case "CONNECT":
			log.Info("kind = connect")
			if !handleConnect(*msg, conn) {
				return
			}
		case "PRIVATE":
			log.Info("kind = private")
			updateTextChan <- "(private) from " + msg.Me.Name + ": " + msg.MSG
		case "PUBLIC":
			log.Info("kind = publuic")
			updateTextChan <- msg.Me.Name + ": " + msg.MSG
		case "DISCONNECT":
			log.Println("kind = disconnect")
			disconnect(*msg)
		case "HEARTBEAT":
			log.Println("HEARTBEAT")
		case "LIST":
			log.Info("kind = LIST")
			connectToPeers(*msg)
			return
		case "ADD":
			log.Info("kind = ADD")
			addPeer(*msg)
		default:
			log.Info("unknown message type")
		}
	}
}

func handleConnect(msg Message, conn net.Conn) bool {
	users := GetFromUserMap(usersMap)
	users = append(users, User{myName, util.GetMyIP(), myPort})
	response := Message{"LIST", User{}, "", users}
	if _, usernameTaken := usersMap[msg.Me.Name]; usernameTaken {
		response.MSG = "Username already taken, choose another one that is not in the list"
		response.Send()
		return false
	}
	usersMap[msg.Me.Name] = msg.Me
	usersConnectionsMap[msg.Me.Name] = conn
	log.Println(usersConnectionsMap)
	response.SendToUser(msg.Me.Name, updateTextChan)
	return true
}

func addPeer(msg Message) {
	usersMap[msg.Me.Name] = msg.Me
	conn, err := createConnection(msg.Me)
	if err != nil {
		log.Println(err)
		return
	}
	usersConnectionsMap[msg.Me.Name] = conn
	updateUserListChan <- GetFromUserMap(usersMap)
	updateTextChan <- msg.Me.Name + " just joined the chat (from IP: " + msg.Me.IP + ")"
}

func disconnect(msg Message) {
	delete(usersMap, msg.Me.Name)
	delete(usersConnectionsMap, msg.Me.Name)
	updateUserListChan <- GetFromUserMap(usersMap)
	updateTextChan <- msg.Me.Name + " left the chat"
}

func connectToPeers(msg Message) {
	for _, user := range msg.Users {
		conn, err := createConnection(user)
		if err != nil {
			log.Println(err)
			continue
		}
		usersMap[user.Name] = user
		usersConnectionsMap[user.Name] = conn
	}
	updateUserListChan <- GetFromUserMap(usersMap)
	addMessage := Message{"ADD", User{myName, util.GetMyIP(), myPort}, "", make([]User, 0)}
	addMessage.Send()
}

func createConnection(user User) (net.Conn, error) {
	service := user.IP + ":" + user.Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func IntroduceMyself(user User) {
	log.Println("introduceMyself")
	conn, err := createConnection(user)
	if err != nil {
		log.Println(err)
		return
	}
	enc := json.NewEncoder(conn)
	intromessage := Message{"CONNECT", User{myName, util.GetMyIP(), myPort}, "", make([]User, 0)}
	err = enc.Encode(intromessage)
	if err != nil {
		log.Printf("Could not encode msg : %s", err)
	}
	go receive(conn)
}
