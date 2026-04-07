<!--
@component
- Attributes:
	- note

-->

<script lang="ts">
	import { type Note } from '$lib/api/notes';

	const { note }: { note: Note } = $props();
</script>

<button type="button" onclick={() => removeNote(note.id)} disabled={note.pending}>Delete</button>
{#if note.id === editingId}
	<button type="button" onclick={() => saveNote(note.id)}>Save</button>
	<button type="button" onclick={cancelEdit}>X</button> ::
	<input type="text" bind:value={editText} />
{:else}
	<button type="button" onclick={() => editNote(note)} disabled={note.pending}>Edit</button> :: {note.body}
{/if}
{#if note.pending}
	(saving . . .)
{/if}
