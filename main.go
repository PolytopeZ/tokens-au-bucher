package main

// C'est juste l'import des libs osef 
import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	// Ce truc c'est pour trouver les infos necessaires de League pour la requete
	port, token, err := findLeague()
	if err != nil {
		fail(err)
	}

	// Ici c'est la requete
	url := fmt.Sprintf("https://127.0.0.1:%s/lol-challenges/v1/update-player-preferences/", port)
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("riot:"+token))

	// Ici {"challengeIds":[]} ca veut dire qu'on vire les tokens des challenges
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(`{"challengeIds":[]}`))
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	// La ca veut dire a LoL, tu te tais et tu prends la requete zbi
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	// Ici on envoie et on lit la réponse
	resp, err := client.Do(req)
	if err != nil {
		fail(err)
	}
	defer resp.Body.Close()

	// La c'est moi qui te dit que c'est ok mon reuf
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("C'est bon c'est fait, appuie sur entree !!!")
	} else {
		fmt.Println("Ah bah ca marche pas.")
	}
	fmt.Scanln()
}

// En gros ce bout de code va trouver le port et le token de League qui sont necessaires
func findLeague() (port, token string, err error) {
	procs, err := process.Processes()
	if err != nil {
		return "", "", err
	}
	portRe := regexp.MustCompile(`--app-port=(\S+?)(?:"|\s|$)`)
	tokenRe := regexp.MustCompile(`--remoting-auth-token=(\S+?)(?:"|\s|$)`)
	for _, p := range procs {
		name, _ := p.Name()
		if name != "LeagueClientUx.exe" {
			continue
		}
		cmd, _ := p.Cmdline()
		pm := portRe.FindStringSubmatch(cmd)
		tm := tokenRe.FindStringSubmatch(cmd)
		if len(pm) > 1 && len(tm) > 1 {
			return pm[1], tm[1], nil
		}
	}
	return "", "", fmt.Errorf("Faut ouvrir le client League of Legends !!!")
}

// Si ca se passe mal ^^
func fail(err error) {
	fmt.Println("Petit bug:", err)
	fmt.Println("Appuie sur entree mon reuf")
	fmt.Scanln()
	os.Exit(1)
}
