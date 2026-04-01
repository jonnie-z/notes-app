<script lang="ts">
	import { onMount } from 'svelte';
	interface Note {
		id: number | string;
		text: string;
		pending?: boolean;
	}

	let note: string = $state('');
	let notes: Note[] = $state([]);
	let editingId: number | string | null = $state(null);
	let editText = $state('');
	let informationalText = $state('');

	async function addNote() {
		if (note.trim()) {
			const tempId = 'temp-' + Date.now();
			const tempNote = { id: tempId, text: note, pending: true };

			notes = [...notes, tempNote];
			note = '';
			const json = JSON.stringify(tempNote);
			try {
				const resp = await fetch('api/notes', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: json
				});

				if (!resp.ok) {
					notes = notes.filter((n) => n.id !== tempId);
					note = tempNote.text;
					throw new Error('Something went wrong!');
				}

				const createdNote = await resp.json();

				notes = notes.map((n) => (n.id === tempId ? createdNote : n));
			} catch (error) {
				console.error(error);
			}
		}
	}

	async function getNotes() {
		try {
			const resp = await fetch('/api/notes');

			if (!resp.ok) {
				throw new Error('Something went wrong!');
			}

			notes = await resp.json();
		} catch (error) {
			console.error(error);
		}
	}

	async function refreshNotes() {
		getNotes();
	}

	async function deleteNote(id: number | string) {
		try {
			const resp = await fetch(`api/notes/${id}`, {
				method: 'DELETE'
			});

			if (!resp.ok) {
				throw new Error('Something went wrong!');
			}
		} catch (error) {
			console.error(error);
		}

		notes = notes.filter((n) => n.id !== id);
	}

	function editNote(note: Note) {
		if (editingId === null) {
			editingId = note.id;
			editText = note.text;
		} else {
			informationalText = 'You need to save or cancel editing before editing another note!';
		}
	}

	async function saveNote(id: number | string) {
		try {
			const tempNote = { text: editText };

			const json = JSON.stringify(tempNote);
			const resp = await fetch(`api/notes/${id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: json
			});

			if (!resp.ok) {
				throw new Error('Something went wrong!');
			}

			const updatedNote: Note = await resp.json();

			notes = notes.map((n) => {
				if (n.id === id) {
					n.text = updatedNote.text;
					return n;
				} else {
					return n;
				}
			});
		} catch (error) {
			console.error(error);
		}

		editText = '';
		editingId = null;
		informationalText = '';
	}

	function cancelEdit() {
		editText = '';
		editingId = null;
		informationalText = '';
	}

	onMount(getNotes);
</script>

<h1>Notes App</h1>

<input type="text" name="note-entry" bind:value={note} />
<button type="button" onclick={addNote}>Add</button>
<button type="button" onclick={refreshNotes}>Refresh</button><br />
<br />
{#each notes as note}
	<button type="button" onclick={() => deleteNote(note.id)} disabled={note.pending}>Delete</button>
	{#if note.id === editingId}
		<button type="button" onclick={() => saveNote(note.id)}>Save</button>
		<button type="button" onclick={cancelEdit}>X</button> ::
		<input type="text" bind:value={editText} />
	{:else}
		<button type="button" onclick={() => editNote(note)} disabled={note.pending}>Edit</button> :: {note.text}
	{/if}
	{#if note.pending}
		(saving . . .)
	{/if}
	<br />
{/each}
<br />
<h5>{informationalText}</h5>


<!-- 

-->