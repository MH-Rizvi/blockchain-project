package transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Transaction structure to hold dataset hash, algorithm hash, and output hash
type Transaction struct {
	ID            int
	DatasetHash   string
	AlgorithmHash string
	OutputHash    string
}

var transactionID int

// GenerateDataset generates a random dataset with predefined points and features
func GenerateDataset() [][]float64 {
	numPoints := 10  // 10 data points
	numFeatures := 3 // 3 features per data point (3D points)

	rand.Seed(time.Now().UnixNano())

	dataset := make([][]float64, numPoints)

	for i := 0; i < numPoints; i++ {
		point := make([]float64, numFeatures)
		for j := 0; j < numFeatures; j++ {
			point[j] = rand.Float64() * 100
		}
		dataset[i] = point
	}

	return dataset
}

// EuclideanDistance calculates the Euclidean distance between two points
func EuclideanDistance(p1, p2 []float64) float64 {
	var sum float64
	for i := range p1 {
		diff := p1[i] - p2[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

// KMeans performs K-means clustering on the given dataset
func KMeans(dataset [][]float64, k int) ([][]float64, [][]int) {
	rand.Seed(time.Now().UnixNano())

	centroids := make([][]float64, k)
	assignedClusters := make([][]int, k)

	// Initialize centroids randomly
	for i := 0; i < k; i++ {
		centroids[i] = dataset[rand.Intn(len(dataset))]
	}

	maxIterations := 100
	for i := 0; i < maxIterations; i++ {
		// Clear previous assignments
		for j := 0; j < k; j++ {
			assignedClusters[j] = nil
		}

		// Assign points to nearest centroid
		for _, point := range dataset {
			minDist := math.Inf(1)
			var closestCentroidIdx int
			for idx, centroid := range centroids {
				dist := EuclideanDistance(point, centroid)
				if dist < minDist {
					minDist = dist
					closestCentroidIdx = idx
				}
			}
			assignedClusters[closestCentroidIdx] = append(assignedClusters[closestCentroidIdx], 1)
		}

		// Recalculate centroids
		newCentroids := make([][]float64, k)
		for j := 0; j < k; j++ {
			newCentroids[j] = make([]float64, len(dataset[0]))
			for _, idx := range assignedClusters[j] {
				for f := 0; f < len(dataset[0]); f++ {
					newCentroids[j][f] += dataset[idx][f]
				}
			}
			for f := 0; f < len(dataset[0]); f++ {
				newCentroids[j][f] /= float64(len(assignedClusters[j]))
			}
		}

		// Check for convergence
		converged := true
		for j := 0; j < k; j++ {
			if !equal(newCentroids[j], centroids[j]) {
				converged = false
				break
			}
		}
		if converged {
			break
		}
		centroids = newCentroids
	}

	return centroids, assignedClusters
}

// Helper function to check if two centroids are equal (for convergence)
func equal(c1, c2 []float64) bool {
	for i := range c1 {
		if c1[i] != c2[i] {
			return false
		}
	}
	return true
}

// HashString returns the SHA-256 hash of the input string
func HashString(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// CreateTransaction generates the dataset, runs the algorithm, and creates a transaction
func CreateTransaction(k int) Transaction {

	transactionID++
	// Step 1: Generate the random dataset
	dataset := GenerateDataset()

	// Step 2: Run K-means on the dataset
	centroids, _ := KMeans(dataset, k)

	// Step 3: Generate the hashes
	datasetHash := HashString(fmt.Sprintf("%v", dataset))
	algorithmHash := HashString("KMeans Algorithm")        // Hash of the algorithm
	outputHash := HashString(fmt.Sprintf("%v", centroids)) // Hash of the output

	// Step 4: Return the transaction containing the hashes
	return Transaction{
		ID:            transactionID,
		DatasetHash:   datasetHash,
		AlgorithmHash: algorithmHash,
		OutputHash:    outputHash,
	}
}
