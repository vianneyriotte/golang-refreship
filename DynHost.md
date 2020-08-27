# DynHost

OVH met à disposition un service lui permettant d'associer un sous-domaine à une adresse IP dynamique.

Pour créer un sous domaine dynamique, il faut aller dans "**Domaines**" > "**le_nom_de_domaine**" > "**DynHost**"

Ici, il est alors possible de créer un sous-domaine, ex. **monip.hoplapizza.ovh** et de lui associer un IP lors de la création.

Ensuite, il faut créer un utilisateur, afin de permettre la mise à jour automatisée à partir d'une autre machine. Il faut donc aller dans "**Gérer les accès**" puis "**Créer un identifiant**". On pourra alors définir un identifiant (ex. *hoplapizza.ovh-monidentifiant*), le sous-domaine concerné par les mises à jour (ex. *monip.hoplapizza.ovh*).



## Mise à jour de l'IP

Pour mettre à jour l'IP, il faudra exécuter une requête HTTP GET en fournissant le sous domaine, l'IP, l'utilisateur et le mot de passe.

Si l'on veut mettre à jour le nom de domaine **monip.hoplapizza.ovh** avec l'IP 10.0.0.1 en utilisant l'outil **CURL**:

```sh
curl -X GET 'http://www.ovh.com/nic/update?system=dyndns&hostname=monip.hoplapizza.ovh&myip=10.0.0.1' -u monip.ovh-vianney:mot_de_passe
```



## Mise à jour par code

On peut également écrire un bout de code permettant de requêter de la même manière que CURL.

Pour se faire il faut effectuer un HTTP GERT à l'url http://www.ovh.com/nic/update?system=dyndns&hostname=monip.hoplapizza.ovh&myip=10.0.0.1 et d'ajouter dans les header l'authentification.

```http
Authorization: Basic aG9wbGFwaXp6YS5vdmgtdmlhbm5leTpFbmRya3dhNA==
```

Le hash après le terme 'Basic' est la concaténation de: **utilisateur:mot_de_passe** en **base 64**.

Exemple: 

 ```
BASE64(monip.ovh-vianney:mot_de_passe) = bW9uaXAub3ZoLXZpYW5uZXkt6bW90X2RlX3Bhc3Nl
 ```



## Mise à jour par Box

A voir suivant le type de box, la capacité à paramètrer une url de rafraichissement : Orange ? Free ? 

## Exemple de code en GO LANG

```go
// Build Test:
// 		go build refreship.go && ./refreship vianney.hoplapizza.ovh hoplapizza.ovh-vianney mdp
// Build pour windows: 
// 		GOOS=windows GOARCH=amd64 go build -o refreship.exe refreship.go

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


```

