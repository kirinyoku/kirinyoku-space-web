import { useEffect, useState } from "react";
import {
  fetchLanguages,
  fetchPostsWithFilters,
  fetchTags,
  Post,
} from "./api/api";
import Sidebar from "./components/sidebar/Sidebar";
import Navbar from "./components/navbar/Navbar";
import Pagination from "./components/post/Pagination";
import PostListing from "./components/post/PostListing";
import SearchBar from "./components/search/SearchBar";

const App: React.FC = () => {
  const [selectedTag, setSelectedTag] = useState<string | null>(null);
  const [selectedType, setSelectedType] = useState<string | null>(null);
  const [selectedLanguage, setSelectedLanguage] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [posts, setPosts] = useState<Post[]>([]);
  const [tags, setTags] = useState<string[]>([]);
  const [languages, setLanguages] = useState<string[]>([]);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [loading, setLoading] = useState<boolean>(false);
  const limit = 15;

  useEffect(() => {
    const loadInitialData = async () => {
      try {
        const fetchedTags = await fetchTags();
        setTags(fetchedTags);
        const fetchedLanguages = await fetchLanguages();
        setLanguages(fetchedLanguages);
      } catch (error) {
        console.error("Failed to fetch initial data:", error);
      }
    };
    loadInitialData();
  }, []);

  useEffect(() => {
    const loadPosts = async () => {
      setLoading(true);
      try {
        const response = await fetchPostsWithFilters(
          searchQuery || null,
          selectedTag,
          selectedType,
          selectedLanguage,
          currentPage,
          limit
        );
        setPosts(response.posts);
        const calculatedTotalPages = Math.max(
          1,
          Math.ceil(response.total_count / limit)
        );
        setTotalPages(calculatedTotalPages);
        if (currentPage > calculatedTotalPages) {
          setCurrentPage(1);
        }
      } catch (error) {
        console.error("Failed to fetch posts:", error);
        setPosts([]);
        setTotalPages(1);
      } finally {
        setLoading(false);
      }
    };
    loadPosts();
  }, [searchQuery, selectedTag, selectedType, selectedLanguage, currentPage]);

  const types = ["video", "article", "website"];

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <main className="flex-1 flex flex-col md:flex-row">
        <Sidebar
          tags={tags}
          selectedTag={selectedTag}
          onTagSelect={setSelectedTag}
          types={types}
          selectedType={selectedType}
          onTypeSelect={setSelectedType}
          languages={languages}
          selectedLanguage={selectedLanguage}
          onLanguageSelect={setSelectedLanguage}
        />
        <section className="flex-1 p-4 bg-gray-50">
          <div className="mb-4">
            <SearchBar
              searchQuery={searchQuery}
              onSearchChange={setSearchQuery}
            />
          </div>
          <div className="mb-4">
            {loading ? (
              <div className="text-center text-gray-500">Loading...</div>
            ) : (
              <PostListing posts={posts} />
            )}
          </div>
          <div className="flex justify-center">
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
            />
          </div>
        </section>
      </main>
    </div>
  );
};

export default App;
