interface HashtagsProps {
  tags: string[];
  selectedTag: string | null;
  onTagSelect: (tag: string | null) => void;
}

const HashtagsFilter: React.FC<HashtagsProps> = ({
  tags,
  selectedTag,
  onTagSelect,
}) => {
  const handleTagClick = (tag: string) => {
    onTagSelect(selectedTag === tag ? null : tag);
  };

  return (
    <div>
      <h3 className="text-md font-medium text-gray-800 mb-2">Hashtags</h3>
      <ul className="flex flex-wrap gap-x-2 space-y-1">
        {/* List of tags */}
        {tags.map((tag) => (
          <li key={tag}>
            <button
              className={`w-full text-left rounded text-sm bg-transparent cursor-pointer ${
                selectedTag === tag
                  ? "text-blue-400"
                  : "text-gray-400 hover:text-blue-400"
              }`}
              onClick={() => handleTagClick(tag)}
            >
              #{tag}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default HashtagsFilter;
