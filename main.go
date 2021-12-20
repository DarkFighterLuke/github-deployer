package main

import (
	"github.com/gorilla/mux"
	"githubDeployer/utils"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var configPath = "conf.yml"
var conf *utils.Config

func main() {
	var err error
	conf, err = utils.GetConfig(configPath)
	if err != nil {
		log.Println(err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc(conf.Server.PayloadEndpoint, payloadHandler).Methods("POST")
	err = http.ListenAndServe(conf.Server.Host+":"+conf.Server.Port, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	temp := strings.Split(string(body), ",")
	temp = strings.Split(temp[0], ":")
	branch := strings.Split(strings.Replace(temp[1], "\"", "", -1), "/")[2]
	log.Println("Received event on branch " + branch)
	if branch == conf.Repository.Branch {
		log.Println("Running " + conf.Script.Path + "...")

		err := exec.Command("chmod", "+x", conf.Script.Path).Run()
		if err != nil {
			log.Println(err)
			return
		}
		err = exec.Command("/bin/bash", conf.Script.Path).Run()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Execution of " + conf.Script.Path + " finished.")
	}
}
