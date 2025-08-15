import React, { useEffect, useState } from "react";
import { getPosts, createPost, updatePost, deletePost } from "./api";
import type { Post } from "./types";

const path = "http://localhost:3000/public";

const App: React.FC = () => {
    const [posts, setPosts] = useState<Post[]>([]);
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [file, setFile] = useState<File | null>(null);
    const [editId, setEditId] = useState<number | null>(null);

    async function loadPosts() {
        try {
            const data = await getPosts();
            setPosts(data);
        } catch (err) {
            console.error(err);
        }
    }

    useEffect(() => {
        loadPosts();
    }, []);

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        try {
            if (editId) {
                await updatePost(editId, { content }, file || undefined);
                setEditId(null);
            } else {
                await createPost({ title, content });
            }
            setTitle("");
            setContent("");
            setFile(null);
            loadPosts();
        } catch (err) {
            console.error(err);
        }
    }

    async function handleDelete(id: number) {
        try {
            await deletePost(id);
            loadPosts();
        } catch (err) {
            console.error(err);
        }
    }

    function handleEdit(post: Post) {
        setEditId(post.id);
        setTitle(post.title);
        setContent(post.content);
        setFile(null);
    }

    return (
            <div className="max-w-4xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-6 text-center">Posts</h1>

            <form
                onSubmit={handleSubmit}
                className="bg-white p-6 rounded-lg shadow-md mb-8 space-y-4"
            >
                {!editId && (
                    <input
                        type="text"
                        placeholder="Title"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        required
                        className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                )}

                <textarea
                    placeholder="Content"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    required
                    className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />

                {/* File input only when editing */}
                {editId && (
                    <input
                        type="file"
                        accept="image/*"
                        onChange={(e) => {
                            if (e.target.files && e.target.files.length > 0) {
                                setFile(e.target.files[0]);
                            }
                        }}
                        className="w-full"
                    />
                )}

                <div className="flex gap-4">
                    <button
                        type="submit"
                        className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
                    >
                        {editId ? "Update" : "Create"}
                    </button>
                    {editId && (
                        <button
                            type="button"
                            onClick={() => {
                                setEditId(null);
                                setTitle("");
                                setContent("");
                                setFile(null);
                            }}
                            className="bg-gray-300 text-gray-700 px-4 py-2 rounded hover:bg-gray-400 transition"
                        >
                            Cancel
                        </button>
                    )}
                </div>
            </form>

            <div className="space-y-6">
                    {(posts || []).map((post) => (
                        post && post.title && post.content && ( // skip null or invalid posts
                            <div
                                key={post.id}
                                className="bg-white p-4 rounded-lg shadow flex flex-col md:flex-row md:items-start md:gap-6"
                            >
                                <div className="flex-1">
                                    <h3 className="text-xl font-semibold">{post.title}</h3>
                                    <p className="mt-2">{post.content}</p>
                                </div>

                                {post.photo_url && (
                                    <img
                                        src={`${path}${post.photo_url}`}
                                        alt="Post"
                                        className="w-full md:w-48 mt-4 md:mt-0 rounded"
                                    />
                                )}

                                <div className="flex gap-2 mt-4 md:mt-0">
                                    <button
                                        onClick={() => handleEdit(post)}
                                        className="bg-yellow-400 px-3 py-1 rounded hover:bg-yellow-500 transition"
                                    >
                                        Edit
                                    </button>
                                    <button
                                        onClick={() => handleDelete(post.id)}
                                        className="bg-red-500 px-3 py-1 rounded text-white hover:bg-red-600 transition"
                                    >
                                        Delete
                                    </button>
                                </div>
                            </div>
                        )
                    ))}
                </div>

            </div>
    );
};

export default App;
