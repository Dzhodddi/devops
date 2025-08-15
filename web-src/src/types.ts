export interface Post {
    id: number;
    title: string;
    content: string;
    author_email: string;
    photo_url: string,
}

export interface CreatePostPayload {
    title: string;
    content: string;
}

export interface EditPostPayload {
    content: string;
}
