package openmatch_mmf_testing

import (
	"fmt"
	"time"

	"open-match.dev/open-match/pkg/matchfunction"
	"open-match.dev/open-match/pkg/pb"
)

type MatchFunctionService struct {
	qsc pb.QueryServiceClient
}

func NewMatchFunctionService(qsc pb.QueryServiceClient) *MatchFunctionService {
	return &MatchFunctionService{qsc: qsc}
}

func (s *MatchFunctionService) Run(request *pb.RunRequest, stream pb.MatchFunction_RunServer) error {
	poolTickets, err := matchfunction.QueryPools(stream.Context(), s.qsc, request.Profile.Pools)
	if err != nil {
		return fmt.Errorf("failed to query pools: %+v", err)
	}
	matches, err := makeMatches(poolTickets, request.Profile)
	for _, match := range matches {
		if err := stream.Send(&pb.RunResponse{Proposal: match}); err != nil {
			return fmt.Errorf("failed to send proposal: %+v", err)
		}
	}
	return nil
}

func makeMatches(poolTickets map[string][]*pb.Ticket, profile *pb.MatchProfile) ([]*pb.Match, error) {
	tickets := map[string]*pb.Ticket{}
	for _, pool := range poolTickets {
		for _, ticket := range pool {
			tickets[ticket.GetId()] = ticket
		}
	}

	var matches []*pb.Match

	t := time.Now().Format("2006-01-02T15:04:05.00")

	thisMatch := make([]*pb.Ticket, 0, 2)
	matchNum := 0

	for _, ticket := range tickets {
		thisMatch = append(thisMatch, ticket)

		if len(thisMatch) >= 2 {
			matches = append(matches, &pb.Match{
				MatchId:       fmt.Sprintf("profile-%s-time-%s-num-%d", profile.Name, t, matchNum),
				MatchProfile:  profile.Name,
				MatchFunction: "demo-match-function",
				Tickets:       thisMatch,
			})

			thisMatch = make([]*pb.Ticket, 0, 2)
			matchNum++
		}
	}

	return matches, nil
}
