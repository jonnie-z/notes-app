export interface Note {
    id: number | string;
    body: string;
    pending?: boolean;
}

async function handleResponse<T>(resp: Response): Promise<T> {
    if (!resp.ok) {
        const text = await resp.text();
        throw new Error(text || "Request failed");
    }

    return resp.json();
}

export async function getNotes(query?: string, pageSize?: string, page?: string): Promise<Note[]> {
    let url = query
        ? `/api/notes?query=${encodeURIComponent(query)}`
        : `/api/notes`;
    
    if (pageSize) {
        if (url.includes('?')) {
            url = url  + `&pageSize=${pageSize}`
        } else {
            url = url  + `?pageSize=${pageSize}`
        }
    }

    if (page) {
        if (url.includes('?')) {
            url = url  + `&page=${page}`
        } else {
            url = url  + `?page=${page}`
        }
    }

    const resp = await fetch(url);
    return handleResponse(resp);
}

export async function createNote(tempNote: Note): Promise<Note> {
    console.log(tempNote);
    console.log(JSON.stringify(tempNote));
    const resp = await fetch('api/notes', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(tempNote)
    });

    return handleResponse(resp);
}

export async function updateNote(tempNote: Note): Promise<Note> {
    const resp = await fetch(`api/notes/${tempNote.id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ body: tempNote.body })
    });

    if (!resp.ok) {
        throw new Error('Something went wrong!');
    }

    return handleResponse(resp);
}

export async function deleteNote(id: number | string) {
    const resp = await fetch(`api/notes/${id}`, {
        method: 'DELETE'
    });

    return handleResponse(resp);
}