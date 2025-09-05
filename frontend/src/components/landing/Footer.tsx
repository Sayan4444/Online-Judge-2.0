import React from "react";
import { Instagram, Facebook, Github } from "lucide-react";

const Footer = () => {
  return (
    <footer className="py-4 text-center w-full fixed bottom-0 left-0 z-50 border-t-2">
      <div className="mb-2 flex justify-center space-x-6">
        <a
          href="https://github.com/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="GitHub"
          className="hover:text-gray-400 transition"
        >
          <Github className="text-2xl" />
        </a>
        <a
          href="https://facebook.com/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="Facebook"
          className="hover:text-gray-400 transition"
        >
          <Facebook className="text-2xl" />
        </a>
        <a
          href="https://instagram.com/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="Instagram"
          className="hover:text-gray-400 transition"
        >
          <Instagram className="text-2xl" />
        </a>
      </div>
      <div>
        Made with <span className="text-red-500">♥</span> by GLUG
      </div>
    </footer>
  );
};

export default Footer;
