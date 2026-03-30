<script lang="ts">
	import { onMount } from 'svelte';
	interface Note {
		id: number;
		body: string;
	}

	let note: string = $state('');
	let notes: Note[] = $state([]);
	let editingId: number | null = $state(null);
	let editText = $state('');
    let informationalText = $state('');

	async function addNote() {
		if (note.trim()) {
			const tempNote = { body: note };

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
					throw new Error('Something went wrong!');
				}
			} catch (error) {
				console.error(error);
			}

			await getNotes();

			note = '';
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

	async function deleteNote(id: number) {
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

		await getNotes();
	}

	function editNote(note: Note) {
		if (editingId === null) {
			editingId = note.id;
			editText = note.body;
		} else {
            informationalText = 'You need to save or cancel editing before editing another note!';
        }
	}

	async function saveNote(id: number) {
		try {
			const tempNote = { body: editText };

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
		} catch (error) {
			console.error(error);
		}

		editText = '';
		editingId = null;
        informationalText = '';

		await getNotes();
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
<button type="button" onclick={addNote}>Add</button><br />
<br />
{#each notes as note}
	<button type="button" onclick={() => deleteNote(note.id)}>Delete</button>
	{#if note.id === editingId}
		<button type="button" onclick={() => saveNote(note.id)}>Save</button>
		<button type="button" onclick={cancelEdit}>X</button> ::
		<input type="text" bind:value={editText} />
	{:else}
		<button type="button" onclick={() => editNote(note)}>Edit</button> :: {note.body}
	{/if}
	<br />
{/each}
<br />
<h5>{informationalText}</h5>

<!--
Quick check: A.
Question answer: I"m leaning towards A, because I feel like B could potentially cause some state/race condition issues, but I don"t know if I"m being overly cautious.
-->
