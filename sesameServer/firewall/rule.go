package firewall

import (
	"log"
	"time"
)

type rule struct {
	firewall        *Firewall
	Name            string
	lifetimeMinutes int
}

func NewRule(firewall *Firewall, name string) rule {
	var result rule
	result.Name = name
	result.firewall = firewall
	result.lifetimeMinutes = result.firewall.RuleLifetimeMinutes
	return result
}

func (t *rule) getCreationTime() (time.Time, error) {
	return time.Parse("02.01.2006 15:04 -07", t.Name[len(`sesame `):])
}

func (t *rule) getDeadlineTime() (time.Time, error) {
	CreationTime, err := t.getCreationTime()
	if err != nil {
		log.Println("Warning: incorrect time format in " + t.Name)
		return time.Time{}, err
	}
	deadlineTime := CreationTime
	for i := 0; i < t.lifetimeMinutes; i++ {
		deadlineTime = deadlineTime.Add(time.Minute)
	}
	return deadlineTime, err
}

func (t *rule) isExpired() bool {
	deadlineTime, err := t.getDeadlineTime()
	if err != nil {
		log.Println("Warning: incorrect time format in " + t.Name)
		return false
	}
	if time.Now().After(deadlineTime) {
		return true
	}
	return false
}

func (t *rule) Delete() {
	t.firewall.DeleteRuleByName(t.Name)
}
