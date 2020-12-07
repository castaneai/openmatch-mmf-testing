package openmatch_mmf_testing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"open-match.dev/open-match/pkg/pb"
)

// 2枚チケットを入れたらとりあえずマッチすることを確かめるテスト
func TestRandomMatchmaking(t *testing.T) {
	profile := &pb.MatchProfile{Name: "fake"}

	ticket1 := &pb.Ticket{Id: "test-ticket-1-id"}
	ticket2 := &pb.Ticket{Id: "test-ticket-2-id"}

	poolTickets := map[string][]*pb.Ticket{
		"test-pool": {ticket1, ticket2},
	}

	matches, err := makeMatches(poolTickets, profile)
	assert.NoError(t, err)
	assert.Len(t, matches, 1)

	var matchedTicketIDs []string
	for _, ticket := range matches[0].Tickets {
		matchedTicketIDs = append(matchedTicketIDs, ticket.Id)
	}
	assert.ElementsMatch(t, []string{ticket1.Id, ticket2.Id}, matchedTicketIDs)
}
