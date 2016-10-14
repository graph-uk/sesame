package sesameServer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (t *SesameServer) openPortsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		IP := t.extractIP(r)
		if IP == "" {
			log.Println("Warning: Cannot parse IP")
		} else {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Error: cannot read body")
			} else {
				if !t.checkFormToken(&body) && !t.checkJSONToken(&body) {
					log.Println("Warning: token is not verified")
					w.Write([]byte("Ciao!\n"))
				} else {
					err = t.firewall.OpenSlavePortsForIP(IP, t.Config.SlavePorts)
					if err == nil {
						w.Write([]byte("Tervetuloa!\n"))
					}
					log.Println("Request received: IP=" + r.RemoteAddr + " ARR_IP=" + r.Header.Get(t.Config.ARR_IP_HeaderField))
				}
			}
		}
	}
}

func (t *SesameServer) extractIP(r *http.Request) string {
	IP := ""
	StandardClientIP := r.RemoteAddr[:strings.Index(r.RemoteAddr, `:`)]
	ARR_API := r.Header.Get(t.Config.ARR_IP_HeaderField)
	if ARR_API != "" {
		ARR_API = ARR_API[:strings.Index(ARR_API, `:`)]
	}
	if ARR_API == "" {
		IP = StandardClientIP
	} else {
		IP = ARR_API
	}
	return IP
}

func (t *SesameServer) checkFormToken(body *[]byte) bool {
	if string(*body) == "Token="+t.Config.Token {
		return true
	}
	return false
}

func (t *SesameServer) checkJSONToken(body *[]byte) bool {
	type openPortsRequest struct {
		Token string
	}
	var request openPortsRequest

	err := json.Unmarshal(*body, &request)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if request.Token == t.Config.Token {
		return true
	}
	return false
}
