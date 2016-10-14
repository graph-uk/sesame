package sesameServer

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/graph-uk/sesame/sesameServer/config"
	"github.com/graph-uk/sesame/sesameServer/firewall"
)

type SesameServer struct {
	Config   *config.Config
	firewall *firewall.Firewall
}

func NewSesameServer() (*SesameServer, error) {
	var result SesameServer
	var err error
	result.Config, err = config.LoadConfig()
	if err != nil {
		return nil, err
	}

	err = result.CheckAdminRights()
	if err != nil {
		log.Println("Admin permissions required.")
		return nil, err
	}

	result.firewall = firewall.NewFirewall(result.Config.DurationMinutes)
	return &result, err
}

func (t *SesameServer) ExpiredRulesKiller() {
	for {
		t.firewall.RefreshRules()
		t.firewall.DeleteExpiredRules()
		time.Sleep(time.Minute)
	}
}

func (t *SesameServer) Serve() error {
	log.Println("Sesame ready to work on port: " + strconv.Itoa(t.Config.MasterPort))
	http.HandleFunc("/", t.openPortsHandler)
	if t.Config.URLPostfix != "" {
		http.HandleFunc("/letmein"+t.Config.URLPostfix, t.pageLetMeInHandler)
	}

	err := http.ListenAndServe(":"+strconv.Itoa(t.Config.MasterPort), nil)
	return err
}

func (*SesameServer) CheckAdminRights() error {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
