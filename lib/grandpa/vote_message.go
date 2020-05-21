package grandpa

import (
	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/crypto"
	"github.com/ChainSafe/gossamer/lib/crypto/ed25519"
	"github.com/ChainSafe/gossamer/lib/scale"
)

// CreateVoteMessage returns a signed VoteMessage given a header
func (s *Service) CreateVoteMessage(header *types.Header, kp crypto.Keypair) (*VoteMessage, error) {
	vote := NewVoteFromHeader(header)

	msg, err := scale.Encode(&FullVote{
		stage: s.subround,
		vote:  vote,
		round: s.state.round,
		setID: s.state.setID,
	})
	if err != nil {
		return nil, err
	}

	sig, err := kp.Sign(msg)
	if err != nil {
		return nil, err
	}

	sm := &SignedMessage{
		hash:        vote.hash,
		number:      vote.number,
		signature:   ed25519.NewSignatureBytes(sig),
		authorityID: kp.Public().(*ed25519.PublicKey).AsBytes(),
	}

	return &VoteMessage{
		setID:   s.state.setID,
		round:   s.state.round,
		stage:   s.subround,
		message: sm,
	}, nil
}

// ValidateMessage validates a VoteMessage and adds it to the current votes
// it returns the resulting vote if validated, error otherwise
func (s *Service) ValidateMessage(m *VoteMessage) (*Vote, error) {
	// check for message signature
	pk, err := ed25519.NewPublicKey(m.message.authorityID[:])
	if err != nil {
		return nil, err
	}

	err = validateMessageSignature(pk, m)
	if err != nil {
		return nil, err
	}

	// check that setIDs match
	if m.setID != s.state.setID {
		return nil, ErrSetIDMismatch
	}

	// check for equivocation ie. multiple votes within one subround
	voter, err := s.state.pubkeyToVoter(pk)
	if err != nil {
		return nil, err
	}

	vote := NewVote(m.message.hash, m.message.number)

	equivocated := s.checkForEquivocation(voter, vote)
	if equivocated {
		return nil, ErrEquivocation
	}

	err = s.validateVote(vote)
	if err != nil {
		return nil, err
	}

	s.votes[pk.AsBytes()] = vote

	return vote, nil
}

// checkForEquivocation checks if the vote is an equivocatory vote.
// it returns true if so, false otherwise.
// additionally, if the vote is equivocatory, it updates the service's votes and equivocations.
func (s *Service) checkForEquivocation(voter *Voter, vote *Vote) bool {
	v := voter.key.AsBytes()

	if s.equivocations[v] != nil {
		// if the voter has already equivocated, every vote in that round is an equivocatory vote
		s.equivocations[v] = append(s.equivocations[v], vote)
		return true
	}

	if s.votes[v] != nil {
		// the voter has already voter, all their votes are now equivocatory
		prev := s.votes[v]
		s.equivocations[v] = []*Vote{prev, vote}
		delete(s.votes, v)
		return true
	}

	return false
}

// validateVote checks if the block that is being voted for exists, and that it is a descendant of a
// previously finalized block.
func (s *Service) validateVote(v *Vote) error {
	// check if v.hash corresponds to a valid block
	has, err := s.blockState.HasHeader(v.hash)
	if err != nil {
		return err
	}

	if !has {
		return ErrBlockDoesNotExist
	}

	// check if the block is an eventual descendant of a previously finalized block
	isDescendant, err := s.blockState.IsDescendantOf(s.head, v.hash)
	if err != nil {
		return err
	}

	if !isDescendant {
		return ErrDescendantNotFound
	}

	return nil
}

func validateMessageSignature(pk *ed25519.PublicKey, m *VoteMessage) error {
	msg, err := scale.Encode(&FullVote{
		stage: m.stage,
		vote:  NewVote(m.message.hash, m.message.number),
		round: m.round,
		setID: m.setID,
	})
	if err != nil {
		return err
	}
	ok, err := pk.Verify(msg, m.message.signature[:])
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidSignature
	}

	return nil
}
