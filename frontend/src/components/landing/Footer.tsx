import React from "react";

const Footer = () => {
  return (
    <footer className="bg-gray-800 text-white py-2 text-center w-full fixed bottom-0 left-0 z-50">
      <div className="mb-2 flex justify-center space-x-6">
        <a
          href="https://github.com/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="GitHub"
          className="hover:text-gray-400 transition"
        >
          <i className="fab fa-github text-2xl"></i>
        </a>
        <a
          href="https://facebook.com/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="Facebook"
          className="hover:text-gray-400 transition"
        >
          <i className="fab fa-facebook text-2xl"></i>
        </a>
        <a
          href="https://instagram.com/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="Instagram"
          className="hover:text-gray-400 transition"
        >
          <i className="fab fa-instagram text-2xl"></i>
        </a>
        <a
          href="https://discord.gg/"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="Discord"
          className="hover:text-gray-400 transition"
        >
          <i className="fab fa-discord text-2xl"></i>
        </a>
      </div>
      <div>
        Made with <span className="text-red-500">♥</span> by GLUG
      </div>
    </footer>
  );
};

export default Footer;
