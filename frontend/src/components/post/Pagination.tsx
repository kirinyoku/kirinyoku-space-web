import { FaArrowLeft, FaArrowRight } from "react-icons/fa";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
}) => {
  const handlePrevClick = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNextClick = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  return (
    <div className="flex justify-center items-center gap-4 p-4">
      <button
        onClick={handlePrevClick}
        disabled={currentPage === 1}
        className={`text-gray-600 hover:text-blue-500 disabled:text-gray-300 disabled:cursor-not-allowed transition-colors cursor-pointer`}
        aria-label="Previous Page"
      >
        <FaArrowLeft size={20} />
      </button>
      <span className="text-gray-800 font-medium">
        Page {currentPage} of {totalPages}
      </span>
      <button
        onClick={handleNextClick}
        disabled={currentPage === totalPages}
        className={`text-gray-600 hover:text-blue-500 disabled:text-gray-300 disabled:cursor-not-allowed transition-colors cursor-pointer`}
        aria-label="Next Page"
      >
        <FaArrowRight size={20} />
      </button>
    </div>
  );
};

export default Pagination;
