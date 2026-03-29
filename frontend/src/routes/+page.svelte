<script lang="ts">
	import { onMount } from 'svelte';
	interface Note {
		id: number;
		body: string;
	}

	let note: string = $state('');
	let notes: Note[] = $state([]);

	async function addNote() {
		if (note.trim()) {
			const tempNote = { body: note };

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
					throw new Error('Something went wrong!');
				}
			} catch (error) {
				console.error(error);
			}

			await getNotes();
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

	onMount(getNotes);
</script>

<h1>Notes App</h1>

<input type="text" name="note-entry" id="" bind:value={note} />
<button type="button" onclick={addNote}>Add</button><br />
<br />
{#each notes as note}
	{note.body} <button type="button" onclick={() => deleteNote(note.id)}>Delete</button><br />
{/each}

<!--
Check:
If getNoteIdx(id) returns -1, I send back an error and indicate the note was not found.

-->
