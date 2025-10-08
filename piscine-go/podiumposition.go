package piscine

func PodiumPosition(podium [][]string) [][]string {
	// Sort the podium slice based on the first element of each inner slice
	// which represents the finishing position
	for i := 0; i < len(podium); i++ {
		for j := i + 1; j < len(podium); j++ {
			// Compare the finishing positions (assuming the first element is the position)
			if podium[i][0] > podium[j][0] {
				// Swap the slices if the current slice has a higher position
				podium[i], podium[j] = podium[j], podium[i]
			}
		}
	}
	return podium
}
