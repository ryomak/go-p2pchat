# go-p2pchat
## Usage 
```
go get -u github.com/ryomak/go-p2pchat

make run
```
### User1
start user1 with setting user1 name and open port
```go run chat.go -name "user1" -p "1111" ```

### User2
start user2 and introduce user2 to user1
```go run chat.go -name "user2" -p "1112" -user "user1@(IP of user1):(Port of user1)" ```

### User3
start user3 and introduce user3 to user1
connecting other node automatically
```go run chat.go -name "user3" -p "1112" -user "user1@(IP of user1):(Port of user1)" ```

## TODO
- [ ] public mode
- [ ] node list => all
- [ ] control
