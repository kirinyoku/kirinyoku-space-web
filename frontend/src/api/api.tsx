import axios from "axios";

export interface Post {
  name: string;
  type: string;
  tags: string[];
  url: string;
  timestamp: string;
}

interface PostsResponse {
  posts: Post[];
  total_count: number;
}

const BASE_URL = import.meta.env.VITE_BASE_URL;
const api = axios.create({
  baseURL: BASE_URL,
});

const normalizeResponse = (data: any): PostsResponse => {
  if (!data || typeof data !== "object") {
    return { posts: [], total_count: 0 };
  }
  return {
    posts: Array.isArray(data.posts) ? data.posts : [],
    total_count: typeof data.total_count === "number" ? data.total_count : 0,
  };
};

export const fetchPostsWithFilters = async (
  query: string | null,
  tag: string | null,
  type: string | null,
  language: string | null,
  page: number,
  limit: number
): Promise<PostsResponse> => {
  const params: any = { page, limit };
  if (query) params.search = query;
  if (tag) params.tag = tag;
  if (type) params.type = type;
  if (language) params.language = language;

  const response = await api.get("/posts", { params });
  return normalizeResponse(response.data);
};

export const fetchPosts = async (): Promise<PostsResponse> => {
  const response = await api.get("/posts");
  return normalizeResponse(response.data);
};

export const fetchPostsBySearch = async (
  query: string
): Promise<PostsResponse> => {
  const response = await api.get("/posts", {
    params: { search: query },
  });
  return normalizeResponse(response.data);
};

export const fetchPostsByTag = async (tag: string): Promise<PostsResponse> => {
  const response = await api.get("/posts", { params: { tag } });
  return normalizeResponse(response.data);
};

export const fetchPostsByType = async (
  type: string
): Promise<PostsResponse> => {
  const response = await api.get("/posts", { params: { type } });
  return normalizeResponse(response.data);
};

export const fetchPostsByLanguage = async (
  language: string
): Promise<PostsResponse> => {
  const response = await api.get("/posts", {
    params: { language },
  });
  return normalizeResponse(response.data);
};

export const fetchTags = async (): Promise<string[]> => {
  const response = await api.get("/tags");
  return Array.isArray(response.data) ? response.data : [];
};

export const fetchLanguages = async (): Promise<string[]> => {
  const response = await api.get("/languages");
  return Array.isArray(response.data) ? response.data : [];
};
