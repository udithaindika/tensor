/*
Copyright 2017 Pearson, Inc.

Licensed under the Apache License, Version 2.0 (the "LICENSE"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pearsonappeng/tensor/db"
	"github.com/pearsonappeng/tensor/models/common"

	"github.com/pearsonappeng/tensor/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/mgo.v2/bson"
	"syscall"
)

func main() {
	if util.InteractiveSetup {
		os.Exit(doSetup())
	}
}

func doSetup() int {
	logrus.Info("Checking database connectivity.. Please be patient.")

	if err := db.Connect(); err != nil {
		logrus.Fatal("\n Cannot connect to database!\n" + err.Error())
	}

	stdin := bufio.NewReader(os.Stdin)

	user := common.User{
		ID:              bson.NewObjectId(),
		Username:        "admin",
		IsSystemAuditor: false,
		IsSuperUser:     true,
		Created:         time.Now(),
	}
	// username is optional (default admin)
	username := readNewline("\n > Username (optional, default `admin`): ", stdin)
	if username != "" {
		user.Username = strings.ToLower(username)
	}

	var ouser common.User
	err := db.Users().Find(bson.M{"username": user.Username}).One(&ouser)

	if err == nil {
		// user already exists
		fmt.Printf("\n Welcome back, %v! (a user with this username/email is already set up..)\n\n", ouser.Username)
	} else {
		user.Email = readNewline("\n > Email: ", stdin)
		if user.Email == "" {
			logrus.Fatal("\n Email is required\n")
			return 1
		}
		user.Email = strings.ToLower(user.Email)

		user.FirstName = readNewline(" > First Name: ", stdin)
		if user.FirstName == "" {
			logrus.Fatal("\n First Name is required\n")
			return 1
		}

		user.LastName = readNewline(" > Last Name: ", stdin)
		if user.LastName == "" {
			logrus.Fatal("\n First Lasti is required\n")
			return 1
		}

		fmt.Print(" > Password: ")
		pwd, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil && string(pwd) == "" {
			logrus.Fatal("\n Password is required\n")
			return 1
		}

		fmt.Print("\n > Confirm password: ")
		cpwd, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil && string(cpwd) == "" {
			logrus.Fatal("\n Confirm password is required\n")
			return 1
		}

		if string(pwd) != string(cpwd) {
			logrus.Fatal("\n Password do not match\n")
			return 1
		}

		user.Password = string(pwd)

		pwdHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
		user.Password = string(pwdHash)

		if err := db.Users().Insert(user); err != nil {
			fmt.Printf(" Failed to create. If you already setup a user, you can disregard this error.\n %v\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("\n You are all setup %v!\n", ouser.FirstName+" "+ouser.LastName)
	}
	fmt.Printf(" You can login with `%v`.\n", user.Username)

	return 0
}

func readNewline(pre string, stdin *bufio.Reader) string {
	fmt.Print(pre)

	str, _ := stdin.ReadString('\n')
	str = strings.Replace(str, "\n", "", -1)

	return str
}
