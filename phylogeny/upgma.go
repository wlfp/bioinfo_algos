package phylogeny

import "fmt"

type cluster struct {
	name string
	size int
}

type distanceMatrix struct {
	clusters            []cluster
	clusterSimilarities [][]float64
}

func (distanceMatrix *distanceMatrix) findClosestClusters() [2]int {
	var minClusterDistance float64 = distanceMatrix.clusterSimilarities[0][1]
	var closestClusterPair [2]int = [2]int{0, 1}
	for clusterOneIndex, clusterOneSimilarities := range distanceMatrix.clusterSimilarities {
		for otherClusterIndex, otherClusterSimilarity := range clusterOneSimilarities {
			clusterTwoIndex := clusterOneIndex + 1 + otherClusterIndex
			if otherClusterSimilarity < minClusterDistance {
				minClusterDistance = otherClusterSimilarity
				closestClusterPair = [2]int{clusterOneIndex, clusterTwoIndex}
			}
		}
	}
	return closestClusterPair
}

func UPGMA() {
	var distanceMatrix distanceMatrix
	distanceMatrix.clusters = []cluster{{name: "A", size: 1}, {name: "B", size: 1}, {name: "C", size: 1}, {name: "D", size: 1}}
	distanceMatrix.clusterSimilarities = [][]float64{{2, 4, 6}, {4, 6}, {6}}
	fmt.Println(distanceMatrix.findClosestClusters())
}
