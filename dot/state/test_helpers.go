package state

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/trie"

	"github.com/stretchr/testify/require"
)

// branch tree randomly
type testBranch struct {
	hash  common.Hash
	depth int
}

func AddBlocksToState(t *testing.T, blockState *BlockState, depth int) ([]*types.Header, []*types.Header) {
	previousHash := blockState.BestBlockHash()

	branches := []testBranch{}
	r := *rand.New(rand.NewSource(rand.Int63()))

	arrivalTime := uint64(1)
	currentChain := []*types.Header{}
	branchChains := []*types.Header{}

	// create base tree
	for i := 1; i <= depth; i++ {
		block := &types.Block{
			Header: &types.Header{
				ParentHash: previousHash,
				Number:     big.NewInt(int64(i)),
				StateRoot:  trie.EmptyHash,
			},
			Body: &types.Body{},
		}

		currentChain = append(currentChain, block.Header)

		hash := block.Header.Hash()
		err := blockState.AddBlockWithArrivalTime(block, arrivalTime)
		require.Nil(t, err)

		previousHash = hash

		isBranch := r.Intn(2)
		if isBranch == 1 {
			branches = append(branches, testBranch{
				hash:  hash,
				depth: i,
			})
		}

		arrivalTime++
	}

	// create tree branches
	for _, branch := range branches {
		previousHash = branch.hash

		for i := branch.depth; i < depth; i++ {
			block := &types.Block{
				Header: &types.Header{
					ParentHash: previousHash,
					Number:     big.NewInt(int64(i) + 1),
					StateRoot:  trie.EmptyHash,
					Digest:     [][]byte{{byte(i)}},
				},
				Body: &types.Body{},
			}

			branchChains = append(branchChains, block.Header)

			hash := block.Header.Hash()
			err := blockState.AddBlockWithArrivalTime(block, arrivalTime)
			require.Nil(t, err)

			previousHash = hash

			arrivalTime++
		}
	}

	return currentChain, branchChains
}
