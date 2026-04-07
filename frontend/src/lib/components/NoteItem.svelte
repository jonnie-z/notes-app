

<script lang="ts">
	import { type Note } from '$lib/api/notes';

	const {
		note,
		isEditing,
		removeNote,
		saveNote,
		cancelEdit,
		editText,
		editNote
	}: {
		note: Note;
		isEditing: boolean;
		removeNote: (id: number | string) => void;
		saveNote: (id: number | string) => void;
        cancelEdit: () => void;
        editText: string;
        editNote: (note: Note) => void;
	} = $props();
</script>

<button type="button" onclick={() => removeNote(note.id)} disabled={note.pending}>Delete</button>
{#if isEditing}
	<button type="button" onclick={() => saveNote(note.id)}>Save</button>
	<button type="button" onclick={cancelEdit}>X</button> ::
	<input type="text" value={editText} />
{:else}
	<button type="button" onclick={() => editNote(note)} disabled={note.pending}>Edit</button> :: {note.body}
{/if}
{#if note.pending}
	(saving . . .)
{/if}
