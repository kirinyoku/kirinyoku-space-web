import Container from "../layout/Container";
import { FaTwitter, FaGithub, FaTelegram } from "react-icons/fa";

const Navbar: React.FC = () => {
  return (
    <nav className="sticky inset-x-0 top-0 z-50 w-full border-b border-b-border h-14 bg-background/50 backdrop-blur-lg shadow-primary/5 shadow-lg transition-all">
      <Container>
        <div className="flex h-14 items-center justify-between border-b border-border">
          <a href="/" className={"flex z-40 font-semibold w-full"}>
            <span className="text-blue-400 text-xl underline">
              kirinyoku's lib
            </span>
          </a>
          <div className="flex justify-end items-center gap-2 sm:gap-4 w-full">
            <ul className="flex items-center gap-4">
              <li>
                <a href="https://x.com/kirinyoku" target="_blank">
                  <FaTwitter size={21} className="text-blue-400" />
                </a>
              </li>
              <li>
                <a href="https://t.me/kirinyoku_space" target="_blank">
                  <FaTelegram size={21} className="text-blue-400" />
                </a>
              </li>
              <li>
                <a href="https://github.com/kirinyoku" target="_blank">
                  <FaGithub size={21} className="text-gray-800" />
                </a>
              </li>
            </ul>
          </div>
        </div>
      </Container>
    </nav>
  );
};

export default Navbar;
