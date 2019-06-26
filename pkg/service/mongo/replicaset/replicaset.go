package replicaset

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.smartx.com/mongo-operator/pkg/utils"
	"sort"
	"strings"
	"time"
)

var logger = utils.NewLogger("mongocluster.service.mongo.replicaset")

// Status holds data about the status of members of the replica set returned
// from replSetGetStatus
//
// See http://docs.mongodb.org/manual/reference/command/replSetGetStatus/#dbcmd.replSetGetStatus
type Status struct {
	Name    string         `bson:"set"`
	Members []MemberStatus `bson:"members"`
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

const (
	// maxInitiateAttempts is the maximum number of times to attempt
	// replSetInitiate for each call to Initiate.
	maxInitiateAttempts = 10

	// initiateAttemptDelay is the amount of time to sleep between failed
	// attempts to replSetInitiate.
	initiateAttemptDelay = 100 * time.Millisecond

	// maxInitiateStatusAttempts is the maximum number of attempts
	// to get the replication status after Initiate.
	maxInitiateStatusAttempts = 50

	// initiateAttemptStatusDelay is the amount of time to sleep between failed
	// attempts to replSetGetStatus.
	initiateAttemptStatusDelay = 500 * time.Millisecond
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

var (
	getCurrentStatus = CurrentStatus
	getBuildInfo     = BuildInfo
	attemptInitiate  = doAttemptInitiate
)

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

// Config is the document stored in mongodb that defines the servers in the
// replica set
type Config struct {
	Name            string   `bson:"_id"`
	Version         int      `bson:"version"`
	ProtocolVersion int      `bson:"protocolVersion,omitempty"`
	Members         []Member `bson:"members"`
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

// CurrentStatus returns the status of the replica set for the given session.
func CurrentStatus(session *mgo.Session) (*Status, error) {
	status := &Status{}
	err := session.Run("replSetGetStatus", status)
	if err != nil {
		return nil, fmt.Errorf("cannot get replica set status: %v", err)
	}

	for index, member := range status.Members {
		status.Members[index].Address = formatIPv6AddressWithBrackets(member.Address)
	}
	return status, nil
}

// BuildInfo returns the mongod build info for the given session.
func BuildInfo(session *mgo.Session) (mgo.BuildInfo, error) {
	return session.BuildInfo()
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

// fmtConfigForLog generates a succinct string suitable for debugging what the Members are up to.
// Note that Members will be printed in Id sorted order, regardless of the order in config.Members
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

// doAttemptInitiate will attempt to initiate a mongodb replicaset with each of
// the given configs, returning as soon as one config is successful.
func doAttemptInitiate(monotonicSession *mgo.Session, cfg []Config) error {
	var err error
	for _, c := range cfg {
		logger.Infof("Initiating replicaset with config: %s", fmtConfigForLog(&c))
		if err = monotonicSession.Run(bson.D{{"replSetInitiate", c}}, nil); err != nil {
			logger.Infof("Unsuccessful attempt to initiate replicaset: %v", err)
			continue
		}
		return nil
	}
	return err
}

// Initiate sets up a replica set with the given replica set name with the
// single given member.  It need be called only once for a given mongo replica
// set.  The tags specified will be added as tags on the member that is created
// in the replica set.
//
// Note that you must set DialWithInfo and set Direct = true when dialing into a
// specific non-initiated mongo server.
//
// See http://docs.mongodb.org/manual/reference/method/rs.initiate/ for more
// details.
func Initiate(session *mgo.Session, address, name string, tags map[string]string) error {
	logger.Info("Try to init a Mongo Cluster.")
	monotonicSession := session.Clone()
	defer monotonicSession.Close()
	monotonicSession.SetMode(mgo.Monotonic, true)

	// For mongo 4 and above, we use protocol version 1.
	buildInfo, err := getBuildInfo(monotonicSession)
	if err != nil {
		return err
	}
	protocolVersion := 0
	if buildInfo.VersionAtLeast(4) {
		protocolVersion = 1
	}
	logger.Info("change protocol to suit mongod version", "mongodVersion",
		buildInfo.Version, "protocol", protocolVersion)

	// We don't know mongod's ability to use a correct IPv6 addr format
	// until the server is started, but we need to know before we can start
	// it. Try the older, incorrect format, if the correct format fails.
	cfg := []Config{
		{
			Name:            name,
			Version:         1,
			ProtocolVersion: protocolVersion,
			Members: []Member{{
				Id:      1,
				Address: address,
				Tags:    tags,
			}},
		}, {
			Name:            name,
			Version:         1,
			ProtocolVersion: protocolVersion,
			Members: []Member{{
				Id:      1,
				Address: formatIPv6AddressWithoutBrackets(address),
				Tags:    tags,
			}},
		},
	}

	// Attempt replSetInitiate, with potential retries.
	for i := 0; i < maxInitiateAttempts; i++ {
		monotonicSession.Refresh()
		if err = attemptInitiate(monotonicSession, cfg); err != nil {
			time.Sleep(initiateAttemptDelay)
			continue
		}
		break
	}

	// Wait for replSetInitiate to complete. Even if err != nil,
	// it may be that replSetInitiate is still in progress, so
	// attempt CurrentStatus.
	for i := 0; i < maxInitiateStatusAttempts; i++ {
		monotonicSession.Refresh()
		var status *Status
		status, err = getCurrentStatus(monotonicSession)
		if err != nil {
			logger.Errorf("Initiate: fetching replication status failed: %v", err)
		}
		if err != nil || len(status.Members) == 0 {
			time.Sleep(initiateAttemptStatusDelay)
			continue
		}
		break
	}
	return err
}
