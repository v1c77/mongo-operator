package mongo

import (
	"github.com/globalsign/mgo"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
	"sort"
	"strings"
	"fmt"
	"github.com/globalsign/mgo/bson"
)
var logger = logf.Log.WithName("mongocluster.service.mongo")

type Client interface {}


type client struct{}

// Config is the document stored in mongodb that defines the servers in the
// replica set
type Config struct {
	Name            string   `bson:"_id"`
	Version         int      `bson:"version"`
	ProtocolVersion int      `bson:"protocolVersion,omitempty"`
	Members         []Member `bson:"members"`
}

// Status holds data about the status of members of the replica set returned
// from replSetGetStatus
//
// See http://docs.mongodb.org/manual/reference/command/replSetGetStatus/#dbcmd.replSetGetStatus
type Status struct {
	Name    string         `bson:"set"`
	Members []MemberStatus `bson:"members"`
}


// Member holds configuration information for a replica set member.
//
// See http://docs.mongodb.org/manual/reference/replica-configuration/
// for more details
type Member struct {
	// Id is a unique id for a member in a set.
	Id int `bson:"_id"`

	// Address holds the network address of the member,
	// in the form hostname:port.
	Address string `bson:"host"`

	// Arbiter holds whether the member is an arbiter only.
	// This value is optional; it defaults to false.
	Arbiter *bool `bson:"arbiterOnly,omitempty"`

	// BuildIndexes determines whether the mongod builds indexes on this member.
	// This value is optional; it defaults to true.
	BuildIndexes *bool `bson:"buildIndexes,omitempty"`

	// Hidden determines whether the replica set hides this member from
	// the output of IsMaster.
	// This value is optional; it defaults to false.
	Hidden *bool `bson:"hidden,omitempty"`

	// Priority determines eligibility of a member to become primary.
	// This value is optional; it defaults to 1.
	Priority *float64 `bson:"priority,omitempty"`

	// Tags store additional information about a replica member, often used for
	// customizing read preferences and write concern.
	Tags map[string]string `bson:"tags,omitempty"`

	// SlaveDelay describes the number of seconds behind the master that this
	// replica set member should lag rounded up to the nearest second.
	// This value is optional; it defaults to 0.
	SlaveDelay *time.Duration `bson:"slaveDelay,omitempty"`

	// Votes controls the number of votes a server has in a replica set election.
	// This value is optional; it defaults to 1.
	Votes *int `bson:"votes,omitempty"`
}

// MemberState represents the state of a replica set member.
// See http://docs.mongodb.org/manual/reference/replica-states/
type MemberState int

const (
	StartupState = iota
	PrimaryState
	SecondaryState
	RecoveringState
	FatalState
	Startup2State
	UnknownState
	ArbiterState
	DownState
	RollbackState
	ShunnedState
)

var memberStateStrings = []string{
	StartupState:    "STARTUP",
	PrimaryState:    "PRIMARY",
	SecondaryState:  "SECONDARY",
	RecoveringState: "RECOVERING",
	FatalState:      "FATAL",
	Startup2State:   "STARTUP2",
	UnknownState:    "UNKNOWN",
	ArbiterState:    "ARBITER",
	DownState:       "DOWN",
	RollbackState:   "ROLLBACK",
	ShunnedState:    "SHUNNED",
}


// Status holds the status of a replica set member returned from
// replSetGetStatus.
type MemberStatus struct {
	// Id holds the replica set id of the member that the status is describing.
	Id int `bson:"_id"`

	// Address holds address of the member that the status is describing.
	Address string `bson:"name"`

	// Self holds whether this is the status for the member that
	// the session is connected to.
	Self bool `bson:"self"`

	// ErrMsg holds the most recent error or status message received
	// from the member.
	ErrMsg string `bson:"errmsg"`

	// Healthy reports whether the member is up. It is true for the
	// member that the request was made to.
	Healthy bool `bson:"health"`

	// State describes the current state of the member.
	State MemberState `bson:"state"`

	// Uptime describes how long the member has been online.
	Uptime time.Duration `bson:"uptime"`

	// Ping describes the length of time a round-trip packet takes to travel
	// between the remote member and the local instance.  It is zero for the
	// member that the session is connected to.
	Ping time.Duration `bson:"pingMS"`
}

func fmtConfigForLog(config *Config) string {
	memberInfo := make([]string, len(config.Members))
	members := append([]Member(nil), config.Members...)
	sort.SliceStable(members, func(i, j int) bool { return members[i].Id < members[j].Id })
	for i, member := range members {
		voting := "not-voting"
		if member.Votes == nil || *member.Votes > 0 {
			voting = "voting"
		}
		var tags []string
		for key, val := range member.Tags {
			tags = append(tags, fmt.Sprintf("%s:%s", key, val))
		}
		memberInfo[i] = fmt.Sprintf("    {%d %q %v %s},", member.Id, member.Address, strings.Join(tags, ", "), voting)
	}
	return fmt.Sprintf(`{
  Name: %s,
  Version: %d,
  Protocol Version: %d,
  Members: {
%s
  },
}`, config.Name, config.Version, config.ProtocolVersion, strings.Join(memberInfo, "\n"))
}

// formatIPv6AddressWithoutBrackets turns correctly formatted IPv6 addresses
// into the "bad format" (without brackets around the address) that mongo <2.7
// require use.
func formatIPv6AddressWithoutBrackets(address string) string {
	address = strings.Replace(address, "[", "", 1)
	address = strings.Replace(address, "]", "", 1)
	return address
}

// formatIPv6AddressWithBrackets turns the "bad format" IPv6 addresses
// ("<addr>:<port>") that mongo <2.7 uses into correctly format addresses
// ("[<addr>]:<port>").
func formatIPv6AddressWithBrackets(address string) string {
	if strings.Count(address, ":") >= 2 && strings.Count(address, "[") == 0 {
		lastColon := strings.LastIndex(address, ":")
		host := address[:lastColon]
		port := address[lastColon+1:]
		return fmt.Sprintf("[%s]:%s", host, port)
	}
	return address
}
// CurrentConfig returns the Config for the given session's replica set.  If
// there is no current config, the error returned will be mgo.ErrNotFound.
var CurrentConfig = currentConfig

func currentConfig(session *mgo.Session) (*Config, error) {
	cfg := &Config{}
	monotonicSession := session.Clone()
	defer monotonicSession.Close()
	monotonicSession.SetMode(mgo.Monotonic, true)
	err := monotonicSession.DB("local").C("system.replset").Find(nil).One(cfg)
	if err == mgo.ErrNotFound {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("cannot get replset config: %s", err.Error())
	}

	members := make([]Member, len(cfg.Members), len(cfg.Members))
	for index, member := range cfg.Members {
		member.Address = formatIPv6AddressWithBrackets(member.Address)
		members[index] = member
	}
	// Sort the values by Member.Id
	sort.Slice(members, func(i, j int) bool { return members[i].Id < members[j].Id })
	cfg.Members = members
	return cfg, nil
}


func doAttemptInitiate(monotonicSession *mgo.Session, cfg []Config) error {
	var err error
	for _, c := range cfg {
		logger.Info("Initiating replicaset with config: %s",
			fmtConfigForLog(&c))
		if err = monotonicSession.Run(bson.D{{"replSetInitiate", c}},
		nil); err != nil {
			logger.Info(fmt.Sprintf("Unsuccessful attempt to initiate replicaset" +
				": %v", err))
			continue
		}
		return nil
	}
	return err
}


// ======================== mongo client ===============================

func New() Client{
	return &client{}
}

//func (c *Client) GetMongoStatus() (*Config, error) {
//
//}
