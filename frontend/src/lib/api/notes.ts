export async function getNotes(query?: string) {
    try {
        const url = query
            ? `/api/notes?query=${encodeURIComponent(query)}`
            : `/api/notes`;

        const resp = await fetch(url);
        if (!resp.ok) {
            throw new Error('Something went wrong!');
        }

        return resp.json();
    } catch (error) {
        console.error(error);
    }
}
