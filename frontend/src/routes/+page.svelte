<script lang="ts">
	import { onMount } from 'svelte';

	import NoteItem from '$lib/components/NoteItem.svelte';
	import { type Note, getNotes, createNote, updateNote, deleteNote } from '$lib/api/notes';

	const MIN_LENGTH = 3;

	let inputRef: HTMLInputElement;
	let note: string = $state('');
	let query: string = $state('');
	let notes: Note[] = $state([]);
	let editingId: number | string | null = $state(null);
	let editText = $state('');
	let informationalText = $state('');
	let isLoading = $state(false);
	let page = $state('1');
	let pageSize = $state('10');

	async function addNote() {
		if (note.trim()) {
			try {
				const tempId = 'temp-' + Date.now();
				const tempNote = { id: tempId, body: note, pending: true };

				notes = [...notes, tempNote];
				note = '';
				let createdNote: Note;

				createdNote = await createNote(tempNote);

				if (!createdNote) {
					notes = notes.filter((n) => n.id !== tempId);
					note = tempNote.body;
				}

				notes = notes.map((n) => (n.id === tempId ? createdNote : n));
			} catch (error) {
				console.log(error);
			}
		}

		inputRef.focus();
	}

	async function refreshNotes() {
		query = '';

		try {
			notes = await getNotes('', pageSize, page);
		} catch (error) {
			console.error(error);
		}
	}

	async function removeNote(id: number | string) {
		try {
			await deleteNote(id);
		} catch (error) {
			console.log(error);
		}

		notes = notes.filter((n) => n.id !== id);
	}

	function editNote(note: Note) {
		if (editingId === null) {
			editingId = note.id;
			editText = note.body;
		} else {
			informationalText = 'You need to save or cancel editing before editing another note!';
		}
	}

	async function saveNote(id: number | string, body: string) {
		try {
			const tempNote = { id: id, body: body };

			const updatedNote = await updateNote(tempNote);

			notes = notes.map((n) => {
				if (n.id === id) {
					n.body = updatedNote.body;
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

	let searchTimeout: ReturnType<typeof setTimeout> | null = null;

	async function searchNotes() {
		if (searchTimeout) {
			clearTimeout(searchTimeout);
		}

		searchTimeout = setTimeout(executeSearch, 400);
	}

	async function executeSearch() {
		isLoading = true;

		if (query.length == 0 || query.length >= MIN_LENGTH) {
			try {
				notes = await getNotes(query, pageSize, page);
			} catch (error) {
				console.error(error);
			} finally {
				isLoading = false;
			}
		}
	}

	onMount(async () => {
		isLoading = true;

		try {
			notes = await getNotes('', pageSize, page);
		} catch (error) {

		} finally {
			isLoading = false;
		}
		});
</script>

<h1>Notes App</h1>

<input bind:this={inputRef} type="text" name="note-entry" bind:value={note} />
<button type="button" onclick={addNote}>Add</button>
<button type="button" onclick={refreshNotes}>Refresh</button><br />
<input type="text" name="search-query" bind:value={query} oninput={searchNotes} /><br />
<label for="pageSize">Page Size:</label>
<select name="pageSize" id="pageSize" bind:value={pageSize}>
	<option value="1">1</option>
	<option value="5">5</option>
	<option value="10">10</option>
	<option value="20">20</option>
</select>
<label for="page">Page:</label>
<select name="page" id="page" bind:value={page}>
	<option value="1">1</option>
	<option value="2">2</option>
	<option value="3">3</option>
	<option value="4">4</option>
	<option value="5">5</option>
</select>
<br />
<br />

{#if isLoading}
	<p>Loading . . .</p>
{:else if notes.length > 0}
	{#each notes as note}
		<NoteItem
			{note}
			isEditing={note.id === editingId}
			{removeNote}
			{saveNote}
			{cancelEdit}
			{editText}
			{editNote}
		></NoteItem><br />
	{/each}
{:else}
	NO NOTES FOUND
{/if}

<br />
<h5>{informationalText}</h5>

<!-- 

-->
