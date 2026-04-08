

<script lang="ts">
	import { type Note } from '$lib/api/notes';

	let {
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
		saveNote: (id: number | string, body: string) => void;
        cancelEdit: () => void;
        editText: string;
        editNote: (note: Note) => void;
	} = $props();
</script>

<button type="button" onclick={() => removeNote(note.id)} disabled={note.pending}>Delete</button>
{#if isEditing}
	<button type="button" onclick={() => saveNote(note.id, editText)}>Save</button>
	<button type="button" onclick={cancelEdit}>X</button> ::
	<input type="text" bind:value={editText} />
{:else}
	<button type="button" onclick={() => editNote(note)} disabled={note.pending}>Edit</button> :: {note.body}
{/if}
{#if note.pending}
	(saving . . .)
{/if}
