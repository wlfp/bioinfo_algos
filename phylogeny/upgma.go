package phylogeny

import (
	"fmt"
	"math"
	"slices"
)

type cluster struct {
	name string
	size int
}

type distanceMatrix struct {
	clusters            []cluster
	deadClusterIndices  []int
	clusterSimilarities [][]float64
}

func (distanceMatrix *distanceMatrix) findClosestClusters() [2]int {
	var minClusterDistance float64 = math.MaxFloat64
	var closestClusterPair [2]int = [2]int{-1, -1}
	for clusterOneIndex, clusterOneSimilarities := range distanceMatrix.clusterSimilarities {
		if slices.Contains(distanceMatrix.deadClusterIndices, clusterOneIndex) {
			continue
		}
		for otherClusterIndex, otherClusterSimilarity := range clusterOneSimilarities {
			clusterTwoIndex := clusterOneIndex + 1 + otherClusterIndex
			if slices.Contains(distanceMatrix.deadClusterIndices, clusterTwoIndex) {
				continue
			}
			if otherClusterSimilarity < minClusterDistance {
				minClusterDistance = otherClusterSimilarity
				closestClusterPair = [2]int{clusterOneIndex, clusterTwoIndex}
			}
		}
	}
	return closestClusterPair
}

func (distanceMatrix *distanceMatrix) mergeClosestClusters(closestClusters [2]int) {
	var newCluster cluster
	newCluster.name = distanceMatrix.clusters[closestClusters[0]].name + distanceMatrix.clusters[closestClusters[1]].name
	newCluster.size = distanceMatrix.clusters[closestClusters[0]].size + distanceMatrix.clusters[closestClusters[1]].size
	distanceMatrix.deadClusterIndices = append(distanceMatrix.deadClusterIndices, closestClusters[0])
	distanceMatrix.deadClusterIndices = append(distanceMatrix.deadClusterIndices, closestClusters[1])
	distanceMatrix.clusters = append(distanceMatrix.clusters, newCluster)

	distanceMatrix.clusterSimilarities = append(distanceMatrix.clusterSimilarities, []float64{})

	distanceMatrix.updateDistanceMatrix(closestClusters)
	fmt.Println(distanceMatrix)
}

// memberIndices contains the two clusters that just combined.
func (distanceMatrix *distanceMatrix) updateDistanceMatrix(memberIndices [2]int) {
	for clusterSimilaritiesIndex := range distanceMatrix.clusterSimilarities {
		if slices.Contains(distanceMatrix.deadClusterIndices, clusterSimilaritiesIndex) {
			continue
		}
		distanceMatrix.clusterSimilarities[clusterSimilaritiesIndex] = append(distanceMatrix.clusterSimilarities[clusterSimilaritiesIndex], distanceMatrix.computeNewClusterAverageDistance(memberIndices[0], memberIndices[1], clusterSimilaritiesIndex))
	}
}

func (distanceMatrix *distanceMatrix) computeNewClusterAverageDistance(clusterOneIndex, clusterTwoIndex, oldClusterIndex int) float64 {
	var similarityToOne, similarityToTwo float64
	if clusterOneIndex < oldClusterIndex {
		similarityToOne = distanceMatrix.clusterSimilarities[clusterOneIndex][oldClusterIndex-clusterOneIndex-1]
	} else {
		similarityToOne = distanceMatrix.clusterSimilarities[oldClusterIndex][clusterOneIndex-oldClusterIndex-1]
	}
	if clusterTwoIndex < oldClusterIndex {
		similarityToTwo = distanceMatrix.clusterSimilarities[clusterTwoIndex][oldClusterIndex-clusterTwoIndex-1]
	} else {
		similarityToTwo = distanceMatrix.clusterSimilarities[oldClusterIndex][clusterTwoIndex-oldClusterIndex-1]
	}
	return (similarityToOne + similarityToTwo) / float64(distanceMatrix.clusters[clusterOneIndex].size+distanceMatrix.clusters[clusterTwoIndex].size)
}

func (distanceMatrix *distanceMatrix) upgma() {
	for len(distanceMatrix.clusters)-len(distanceMatrix.deadClusterIndices) > 1 {
		closestClusters := distanceMatrix.findClosestClusters()
		distanceMatrix.mergeClosestClusters(closestClusters)
	}
}

func UPGMA() {
	var distanceMatrix distanceMatrix
	distanceMatrix.clusters = []cluster{{name: "A", size: 1}, {name: "B", size: 1}, {name: "C", size: 1}, {name: "D", size: 1}}
	distanceMatrix.clusterSimilarities = [][]float64{{2, 4, 6}, {4, 6}, {6}}
	distanceMatrix.upgma()
}
