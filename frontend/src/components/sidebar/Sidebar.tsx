import HashtagsFilter from "./HashtagsFilter";
import LanguageFilter from "./LanguageFilter";
import TypeFilter from "./TypeFilter";

interface SidebarProps {
  tags: string[];
  selectedTag: string | null;
  onTagSelect: (tag: string | null) => void;

  types: string[];
  selectedType: string | null;
  onTypeSelect: (type: string | null) => void;

  languages: string[];
  selectedLanguage: string | null;
  onLanguageSelect: (language: string | null) => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  tags,
  selectedTag,
  onTagSelect,
  types,
  selectedType,
  onTypeSelect,
  languages,
  selectedLanguage,
  onLanguageSelect,
}) => {
  return (
    <aside className="w-full md:w-80 p-4 shrink-0 shadow-primary/5 shadow-lg z-50">
      <div className="space-y-2">
        <TypeFilter
          types={types}
          selectedType={selectedType}
          onTypeSelect={onTypeSelect}
        />
        <LanguageFilter
          languages={languages}
          selectedLanguage={selectedLanguage}
          onLanguageSelect={onLanguageSelect}
        />
        <HashtagsFilter
          tags={tags}
          selectedTag={selectedTag}
          onTagSelect={onTagSelect}
        />
      </div>
    </aside>
  );
};

export default Sidebar;
