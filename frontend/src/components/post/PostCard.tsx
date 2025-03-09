import { Post } from "@/api/api";

interface PostCardProps {
  post: Post;
}

const PostCard: React.FC<PostCardProps> = ({ post }) => {
  const handleCardClick = () => {
    window.location.href = post.url;
  };

  const tags = post.tags || [];

  return (
    <div
      className="flex flex-col justify-between min-h-32 p-4 bg-white shadow rounded-lg hover:bg-gray-50 cursor-pointer transition-colors shadow-primary/5 shadow-lg"
      onClick={handleCardClick}
    >
      <h3 className="text-lg font-semibold text-gray-900 mb-2">{post.name}</h3>
      <div className="flex justify-between items-center">
        <div className="flex flex-wrap gap-2">
          {tags.map((tag) => (
            <span
              key={tag}
              className="bg-blue-200 text-blue-800 rounded px-2 py-1 text-xs"
            >
              #{tag}
            </span>
          ))}
        </div>
        <p className="text-sm text-gray-600 font-semibold">{post.type}</p>
      </div>
    </div>
  );
};

export default PostCard;
