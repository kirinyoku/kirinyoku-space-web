import { Post } from "@/api/api";
import PostCard from "./PostCard";

interface PostListingProps {
  posts: Post[];
}

const PostListing: React.FC<PostListingProps> = ({ posts }) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {posts.map((post) => (
        <PostCard key={post.url} post={post} />
      ))}
    </div>
  );
};

export default PostListing;
