import React from "react";
import LoginButton from "../auth/LoginButton";
import Image from "next/image";

const Hero = () => (
  <section className="py-16 bg-gray-50 text-center">
    <Image
      src="/oj.png"
      alt="Online Judge Logo"
      width={200}
      height={200}
      className="mx-auto mb-8"
    />
    <h1 className="text-4xl md:text-5xl font-bold mb-4 text-gray-900">
      Welcome to Online Judge
    </h1>
    <p className="text-lg md:text-xl text-gray-600 mb-8">
      Sharpen your coding skills with real-world challenges, contests, and a
      vibrant community.
    </p>
    <LoginButton />
    <div className="mt-12 flex flex-wrap justify-center gap-6">
      <div className="bg-white p-6 rounded-xl shadow-md min-w-[220px]">
        <h3 className="text-blue-600 font-semibold mb-2 text-lg">
          Practice Problems
        </h3>
        <p className="text-gray-800">
          Solve hundreds of curated coding problems.
        </p>
      </div>
      <div className="bg-white p-6 rounded-xl shadow-md min-w-[220px]">
        <h3 className="text-blue-600 font-semibold mb-2 text-lg">Contests</h3>
        <p className="text-gray-800">
          Compete in regular coding contests and climb the leaderboard.
        </p>
      </div>
      <div className="bg-white p-6 rounded-xl shadow-md min-w-[220px]">
        <h3 className="text-blue-600 font-semibold mb-2 text-lg">Community</h3>
        <p className="text-gray-800">
          Join discussions, ask questions, and help others.
        </p>
      </div>
    </div>
  </section>
);

export default Hero;
