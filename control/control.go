package control

import (
	"fmt"
	"log"
	"strings"

	"github.com/ryomak/go-p2pchat/peer"
)

type Control struct {
	UpdatedText         chan string
	UpdatedTextFromUser chan string
	UpdateUserList      chan []peer.User
	conversationStr     string
	userlist            string
	inputString         string
}

func (ctrl *Control) StartControlLoop() {
	fmt.Println("Running Control loop")
	for {
		select {
		case update := <-ctrl.UpdatedText:
			fmt.Println("received text update", update)
			ctrl.updateText(update)
		case userListChanged := <-ctrl.UpdateUserList:
			fmt.Println("received userListChanged", userListChanged)
			ctrl.updateList(userListChanged)
		case updateFromUser := <-ctrl.UpdatedTextFromUser:
			fmt.Println("<-ctrl.UpdatedTextFromUser", updateFromUser)
			ctrl.handleUserInput(updateFromUser)
		default:
		}
	}
}

func (ctrl *Control) handleUserInput(input string) {
	log.Printf("userInput got message: %s", input)
	whatever := strings.Split(input, "*")
	if input == "disconnect" {
		msg := peer.Message{"DISCONNECT", peer.User{peer.GetMyName(), "", ""}, "", make([]peer.User, 0)}
		msg.Send()
	}
	if len(whatever) > 1 {
		msg := peer.Message{"PRIVATE", peer.User{peer.GetMyName(), "", ""}, whatever[0], make([]peer.User, 0)}
		msg.SendToUser(whatever[1], ctrl.UpdatedTextFromUser)
		ctrl.UpdatedText <- "(private) from " + peer.GetMyName() + ": " + msg.MSG
	} else {
		msg := peer.Message{"Public", peer.User{peer.GetMyName(), "", ""}, whatever[0], make([]peer.User, 0)}
		msg.Send()
		ctrl.UpdatedText <- peer.GetMyName() + ": " + msg.MSG
	}
}

func (ctrl *Control) updateText(toAdd string) {
	ctrl.conversationStr = ctrl.conversationStr + toAdd + "\n" //also keep track of everything in that field
}

func (ctrl *Control) updateList(list []peer.User) {
	ctrl.userlist = ""
	for _, user := range list {
		ctrl.userlist += user.Name + "\n"
	}
}
