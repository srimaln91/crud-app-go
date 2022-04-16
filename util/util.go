package util

import "github.com/srimaln91/crud-app-go/core/entities"

func ChunkEventSlice(slice []entities.Event, chunkSize int) [][]entities.Event {
	var chunks [][]entities.Event
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
