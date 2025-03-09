interface LanguageFilterProps {
  languages: string[];
  selectedLanguage: string | null;
  onLanguageSelect: (language: string | null) => void;
}

const LanguageFilter: React.FC<LanguageFilterProps> = ({
  languages,
  selectedLanguage,
  onLanguageSelect,
}) => {
  const handleLanguageClick = (language: string) => {
    onLanguageSelect(selectedLanguage === language ? null : language);
  };

  return (
    <div>
      <h3 className="text-md font-medium text-gray-800 mb-2">Language</h3>
      <ul className="flex flex-wrap gap-x-2 space-y-1">
        {languages.map((language) => (
          <li key={language}>
            <button
              className={`px-4 py-1 bg-gray-100 rounded-md text-sm cursor-pointer ${
                selectedLanguage === language
                  ? "text-blue-400"
                  : "text-gray-400"
              }`}
              onClick={() => handleLanguageClick(language)}
            >
              {language}
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default LanguageFilter;
