package transaction

type Transaction struct {
	DatasetHash   string
	AlgorithmHash string
	OutputHash    string
}

func NewTransaction(datasetHash, algorithmHash, outputHash string) Transaction {
	return Transaction{
		DatasetHash:   datasetHash,
		AlgorithmHash: algorithmHash,
		OutputHash:    outputHash,
	}
}
