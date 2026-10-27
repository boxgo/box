package kafka

import (
	"crypto/sha512"
	"hash"

	"github.com/xdg/scram"
)

// sha512Fcn is the SHA-512 hash generator missing from github.com/xdg/scram.
var sha512Fcn scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

// XDGSCRAMClient implements the sarama SCRAMClient interface
// using github.com/xdg/scram.
type XDGSCRAMClient struct {
	convo    *scram.ClientConversation
	hashFcn  scram.HashGeneratorFcn
	authzID  string
	anchored bool
}

// Begin prepares the client for the SCRAM exchange with the server.
func (c *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	if c.hashFcn == nil {
		c.hashFcn = scram.SHA256
	}

	client, err := c.hashFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}

	c.convo = client.NewConversation()
	c.authzID = authzID
	c.anchored = true

	return nil
}

// Step takes a string provided from a server (or just an empty string for the
// very first step) and attempts to advance the SCRAM conversation.
func (c *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	return c.convo.Step(challenge)
}

// Done returns true when the SCRAM conversation is complete.
func (c *XDGSCRAMClient) Done() bool {
	return c.convo.Done()
}

// defaultSCRAMClientGenerator returns the default SCRAM client generator
// matching the configured mechanism (SCRAM-SHA-256 or SCRAM-SHA-512).
func (s SASLConfig) defaultSCRAMClientGenerator() SCRAMClient {
	var hashFcn scram.HashGeneratorFcn

	switch s.Mechanism {
	case SASLTypeSCRAMSHA512:
		hashFcn = sha512Fcn
	default:
		hashFcn = scram.SHA256
	}

	return &XDGSCRAMClient{
		hashFcn: hashFcn,
	}
}
