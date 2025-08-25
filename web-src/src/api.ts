import type { Post, CreatePostPayload, EditPostPayload } from "./types";

const API_URL = "http://localhost:3001/v1/post";

export async function getPosts(): Promise<Post[]> {
    const res = await fetch(API_URL, {headers: {Accept: "application/json"}});
    if (!res.ok) throw new Error("Failed to fetch posts");
    return res.json();
}

// Update post with optional photo
export async function updatePost(id: number, payload: EditPostPayload, file?: File): Promise<Post> {
    const formData = new FormData();
    formData.append("content", payload.content);
    if (file) formData.append("photo", file);

    const res = await fetch(`${API_URL}/${id}`, {
        method: "PATCH",
        body: formData,
    });

    if (!res.ok) throw new Error("Failed to update post");
    return res.json();
}

export async function deletePost(id: number): Promise<void> {
    const res = await fetch(`${API_URL}/${id}`, {
        method: "DELETE",
        headers: { Accept: "application/json" },
    });
    if (!res.ok) throw new Error("Failed to delete post");
}

export async function createPost(payload: CreatePostPayload): Promise<Post> {
    const res = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify(payload),
    });
    if (!res.ok) throw new Error("Failed to create post");
    return res.json();
}