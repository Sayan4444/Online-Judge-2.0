import React from "react";
import LoginButton from "../auth/LoginButton";
import Image from "next/image";

const Hero = () => (
  <section className="text-center my-44">
    <Image
      src="/oj.png"
      alt="Online Judge Logo"
      width={200}
      height={200}
      className="mx-auto mb-8"
    />
    <h1 className="text-4xl md:text-5xl font-bold mb-4">
      Welcome to Online Judge
    </h1>
    <p className="text-lg md:text-2xl mb-8">
      Sharpen your coding skills with real-world challenges and contests.
    </p>
    <LoginButton />
    <div className="mt-12 flex flex-wrap justify-center gap-6">
      <p>
        Our one stop platform for competitive programming and coding contests.
      </p>
    </div>
  </section>
);

export default Hero;
