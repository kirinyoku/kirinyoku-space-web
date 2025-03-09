interface TypeFilterProps {
  types: string[];
  selectedType: string | null;
  onTypeSelect: (type: string | null) => void;
}

const TypeFilter: React.FC<TypeFilterProps> = ({
  types,
  selectedType,
  onTypeSelect,
}) => {
  const handleTypeClick = (type: string) => {
    onTypeSelect(selectedType === type ? null : type);
  };

  return (
    <div>
      <h3 className="text-md font-medium text-gray-800 mb-2">Type</h3>
      <ul className="flex flex-wrap gap-x-2 space-y-1">
        {types.map((type) => (
          <li key={type}>
            <button
              className={`px-4 py-1 bg-gray-100 rounded-md text-sm cursor-pointer ${
                selectedType === type ? "text-blue-400" : "text-gray-400"
              }`}
              onClick={() => handleTypeClick(type)}
            >
              {type.charAt(0) + type.slice(1)}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TypeFilter;
