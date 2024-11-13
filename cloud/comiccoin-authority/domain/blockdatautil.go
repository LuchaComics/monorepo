package domain

// func GetBlockNumberToHashMap(blockData []BlockData) map[uint64]string {
// 	blockNumberToHash := make(map[uint64]string)
// 	for _, block := range blockData {
// 		blockNumberToHash[block.Header.GetNumber()] = block.Hash
// 	}
// 	sortedBlockNumbers := make([]uint64, 0, len(blockNumberToHash))
// 	for blockNumber := range blockNumberToHash {
// 		sortedBlockNumbers = append(sortedBlockNumbers, blockNumber)
// 	}
// 	sort.Slice(sortedBlockNumbers, func(i, j int) bool {
// 		return sortedBlockNumbers[i] < sortedBlockNumbers[j]
// 	})
// 	sortedBlockNumberToHash := make(map[uint64]string)
// 	for _, blockNumber := range sortedBlockNumbers {
// 		sortedBlockNumberToHash[blockNumber] = blockNumberToHash[blockNumber]
// 	}
// 	return sortedBlockNumberToHash
// }
