package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type int    `json:"type"`
	Mess string `json:"mess"`
}

type MessageOneUser struct {
	Type int        `json:"type"`
	Mess ReturnUser `json:"mess"`
	XR   float64    `json:"xr"`
	YR   float64    `json:"yr"`
}

type ReturnUser struct {
	Name  string `json:"name"`
	Nhash string `json:"nhash"`
}

type ReturnUsers struct {
	Type  int          `json:"type"`
	Users []ReturnUser `json:"mess"`
}

type User struct {
	c          *websocket.Conn
	name       string
	nhash      string
	enemynhash string
	mux        sync.Mutex
}

type Users struct {
	users map[string](*User)
	mux   sync.Mutex
}

const port string = "3000"

var upgrader = websocket.Upgrader{}
var hosts Users
var ALLUSERS Users

func NEWHOST(user *User) {
	tin, _ := json.Marshal(ReturnUsers{1, append([]ReturnUser{}, ReturnUser{user.name, user.nhash})})
	users := GETUSERS()
	for i := 0; i < len(users); i++ {
		(*users[i]).mux.Lock()
		if (*users[i]).nhash != user.nhash {
			(*users[i]).c.WriteMessage(1, tin)
		}
		(*users[i]).mux.Unlock()
	}
	hosts.mux.Lock()
	hosts.users[(*user).nhash] = user
	hosts.mux.Unlock()
}

func GETHOST(nhash string) (*User, bool) {
	hosts.mux.Lock()
	defer hosts.mux.Unlock()
	user, ok := hosts.users[nhash]
	return user, ok
}

func GETHOSTS() [](*User) {
	hosts.mux.Lock()
	defer hosts.mux.Unlock()
	var users [](*User)
	for _, user := range hosts.users {
		users = append(users, user)
	}
	return users
}

func DELHOST(name string) {
	users := GETUSERS()
	tin, _ := json.Marshal(Message{4, name})
	for i := 0; i < len(users); i++ {
		(*users[i]).mux.Lock()
		if (*users[i]).nhash != name {
			(*users[i]).c.WriteMessage(1, tin)
		}
		(*users[i]).mux.Unlock()
	}
	hosts.mux.Lock()
	delete(hosts.users, name)
	hosts.mux.Unlock()
}

func NEWUSER(user *User) {
	ALLUSERS.mux.Lock()
	ALLUSERS.users[(*user).nhash] = user
	ALLUSERS.mux.Unlock()
}

func GETUSERS() [](*User) {
	ALLUSERS.mux.Lock()
	defer ALLUSERS.mux.Unlock()
	var users [](*User)
	for _, user := range ALLUSERS.users {
		users = append(users, user)
	}
	return users
}

func GETUSER(nhash string) (user User, ok bool) {
	ALLUSERS.mux.Lock()
	link, ok := ALLUSERS.users[nhash]
	if ok {
		user = (*link)
	}
	ALLUSERS.mux.Unlock()
	return user, ok
}

func DELUSER(name string) {
	ALLUSERS.mux.Lock()
	delete(ALLUSERS.users, name)
	ALLUSERS.mux.Unlock()
}

func getHash(in []byte) string {
	hasher := sha1.New()
	hasher.Write(in)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func main() {
	log.Println("Сервер запущен")
	hosts.users = make(map[string](*User))
	ALLUSERS.users = make(map[string](*User))
	log.Println("Регистрация файлов сайта..")
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/speaker", speaker)
	log.Println("Файлы зарегестрированы")
	log.Println("Начало прослушивания порта", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func speaker(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Клиент:", r.RemoteAddr+". Ошибка:", err)
		return
	}
	defer c.Close()
	log.Println("Клиент:", r.RemoteAddr+". Подключен.")
	var cuser User
	cuser.c = c
	ishost := false
	islogined := false
	for {
		messType, in, err := c.ReadMessage()
		if err != nil {
			log.Println("Клиент:", r.RemoteAddr+". Возможная ошибка (норма: 1000, 1001, 1005, 1006):", err)
			break
		} else {
			if messType == 1 {
				var mess Message
				err = json.Unmarshal(in, &mess)
				if err != nil {
					log.Panic(err)
					break
				}
				if mess.Type == 1 {
					cuser.name = mess.Mess
					newnhash := getHash([]byte(mess.Mess + time.Now().String()))
					_, ok := GETUSER(newnhash)
					for ok {
						newnhash = getHash([]byte(mess.Mess + time.Now().String()))
						_, ok = GETUSER(newnhash)
					}
					cuser.nhash = newnhash
					var returnUsers ReturnUsers
					returnUsers.Type = 1
					users := GETHOSTS()
					for i := 0; i < len(users); i++ {
						(*users[i]).mux.Lock()
						returnUsers.Users = append(returnUsers.Users, ReturnUser{(*users[i]).name, (*users[i]).nhash})
						(*users[i]).mux.Unlock()
					}
					in, _ = json.Marshal(returnUsers)
					c.WriteMessage(messType, in)
					log.Println("New Name:", cuser.name, cuser.nhash)
					islogined = true
					NEWUSER(&cuser)
					continue
				} else if islogined {
					if mess.Type == 2 {
						euser, ok := GETHOST(mess.Mess)
						if ok {
							DELHOST((*euser).nhash)
							(*euser).mux.Lock()
							cuser.mux.Lock()
							cuser.enemynhash = (*euser).nhash
							(*euser).enemynhash = cuser.nhash
							tin, _ := json.Marshal(MessageOneUser{2, ReturnUser{cuser.name, cuser.nhash}, rand.Float64(), rand.Float64()})
							in, _ = json.Marshal(MessageOneUser{2, ReturnUser{(*euser).name, (*euser).nhash}, rand.Float64(), rand.Float64()})
							cuser.c.WriteMessage(messType, in)
							(*euser).c.WriteMessage(messType, tin)
							(*euser).mux.Unlock()
							cuser.mux.Unlock()
							continue
						} else {
							mess.Mess = "neok"
						}
					} else if mess.Type == 3 {
						NEWHOST(&cuser)
						mess.Type = 3
						mess.Mess = "ok"
					} else if mess.Type == 5 {
						cuser.mux.Lock()
						enemy, ok := GETUSER(cuser.enemynhash)
						enemy.mux.Lock()
						if ok {
							in, _ = json.Marshal(Message{5, ""})
							enemy.c.WriteMessage(1, in)
						}
						DELUSER(cuser.nhash)
						DELUSER(enemy.nhash)
						enemy.mux.Unlock()
						cuser.mux.Unlock()
						continue
					} else if mess.Type == 6 {
						cuser.mux.Lock()
						enemy, ok := GETUSER(cuser.enemynhash)
						enemy.mux.Lock()
						if ok {
							in, _ = json.Marshal(Message{6, ""})
							enemy.c.WriteMessage(1, in)
						}
						enemy.mux.Unlock()
						cuser.mux.Unlock()
						continue
					}
				} else {
					log.Println("Error")
				}
				log.Println("not continued")
				in, _ = json.Marshal(mess)
				c.WriteMessage(1, in)
				continue
			} else if messType == 2 {
				log.Println("Неожиданная битовая последовательность")
				continue
			} else {
				log.Println("Клиент:", r.RemoteAddr, "ошибка обработки.")
				continue
			}
		}
	}
	cuser.mux.Lock()
	_, ishost = GETHOST(cuser.nhash)
	_, islogined = GETUSER(cuser.nhash)
	cuser.mux.Unlock()
	if ishost {
		DELHOST(cuser.nhash)
	}
	if islogined {
		DELUSER(cuser.nhash)
	}
	log.Println("Клиент:", r.RemoteAddr, "отключен.")
	return
}
