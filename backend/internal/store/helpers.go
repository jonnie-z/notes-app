package store

func removeNoteAtIndex(notes []Note, idx int) []Note {
	return append(notes[:idx], notes[idx+1:]...)
}

func getNoteIdx(notes []Note, id int) int {
	for i, note := range notes {
		if note.ID == id {
			return i
		}
	}

	return -1
}