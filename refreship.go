package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"encoding/base64"
	"net/http"
	"time"
	"os"
)

type Result struct {
	Ip string `json:"ip"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {

	// Vérification du nombre d'arguments
	arglength := len(os.Args)
	if( arglength != 4){
		fmt.Println(`
Erreur:
	Veuillez fournir le nom de domaine, l'identifiant et le mot de passe en paramètre de lancement!
Exemple: 
	./refreship nom-de-domaine identfiant mot-de-passe
		`)
		os.Exit(2)
	}

	// Récupération de l'adresse IP internet
    result := Result{}
	r, err := myClient.Get("https://api.ipify.org?format=json")
    if err != nil {
        fmt.Println(err.Error()) 
    }
    defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	sbody := string(body);
	err = json.Unmarshal([]byte(sbody), &result)
	if err != nil {
		fmt.Println(err.Error()) 
		os.Exit(2)
    }
	println("Votre adresse IP internet est: " +result.Ip)

	// Récupération des paramètres d'appel (dns, identifiant, mot de passe)
	dns := os.Args[1]
	user := os.Args[2]
	password := os.Args[3]
	urlref := "http://www.ovh.com/nic/update?system=dyndns&hostname=" + dns + "&myip=" + result.Ip

	req, err := http.NewRequest("GET", urlref, nil)
	req.Header.Add("Authorization","Basic " + basicAuth(user,password)) 
	r, err = myClient.Do(req) 
    if err != nil {
        fmt.Println(err.Error()) 
    }
	defer r.Body.Close()
	body, err = ioutil.ReadAll(r.Body)
	sbody = string(body)
	println(sbody)
}
