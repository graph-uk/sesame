package firewall

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Firewall struct {
	Rules               []rule
	RuleLifetimeMinutes int
}

func NewFirewall(RuleLifetimeMinutes int) *Firewall {
	var result Firewall
	result.RuleLifetimeMinutes = RuleLifetimeMinutes
	result.RefreshRules()
	return &result
}

func (t *Firewall) RefreshRules() {
	t.Rules = t.Rules[:0] // delete all rules
	sesameRules := t.getSesameRulesNames()
	for _, curRuleName := range *sesameRules {
		t.Rules = append(t.Rules, NewRule(t, curRuleName))
	}
}

func (t *Firewall) DeleteExpiredRules() {
	for _, curRule := range t.Rules {
		if curRule.isExpired() {
			log.Println("Rule is expired: " + curRule.Name)
			curRule.Delete()
		}
	}
}

func (t *Firewall) DeleteRuleByName(ruleName string) {
	// get full output of firewall show
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", `name=`+ruleName+``)
	err := cmd.Run()
	var out, outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	if err != nil {
		log.Println("Warning: cannot delete rule " + ruleName + "\r\n" + err.Error() + "\r\n" + out.String() + "\r\n" + outErr.String())
	}
}

func (t *Firewall) getSesameRulesNames() *[]string {
	var result []string
	allRulesArray := t.getAllRulesNames()

	// filter needed
	for _, curRuleName := range *allRulesArray {
		match, _ := regexp.MatchString("sesame .*", curRuleName)
		if match {
			result = append(result, curRuleName)
		}
	}

	return &result
}

func (t *Firewall) getAllRulesNames() *[]string {
	var result []string

	// get full output of "Firewall show"
	cmd := exec.Command("netsh", "advFirewall", "Firewall", "show", "rule", "name=all")
	var out, outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Cannot get Firewall rules")
		fmt.Println("netsh Firewall says:")
		fmt.Println(outErr.String())
		os.Exit(1)
	}
	outStr := out.String()

	// extract all rule names
	outputStringArr := strings.Split(outStr, "\n")
	for _, curLine := range outputStringArr {
		if strings.Contains(curLine, `Rule Name:`) {
			ruleName := strings.TrimSpace(curLine[10:])
			result = append(result, ruleName)
		}
	}

	return &result
}

func (t *Firewall) OpenSlavePortsForIP(ip, ports string) error {
	timenow := time.Now().Format("02.01.2006 15:04 -07")
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name=sesame "+timenow,
		"protocol=tcp",
		"localport="+ports,
		"dir=in",
		"remoteip="+ip,
		"action=allow")
	var out, outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		log.Println("Cannot create rule. Check permissions.")
		log.Println("netsh firewall says:")
		log.Println(outErr.String())
		return err
	}
	return nil
}
