interface SearchBarProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
}

const SearchBar: React.FC<SearchBarProps> = ({
  searchQuery,
  onSearchChange,
}) => {
  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onSearchChange(event.target.value);
  };

  return (
    <div className="w-full">
      <input
        type="text"
        value={searchQuery}
        onChange={handleInputChange}
        placeholder="Search posts…"
        className="w-full p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white shadow-primary/5 shadow-lg"
      />
    </div>
  );
};

export default SearchBar;
