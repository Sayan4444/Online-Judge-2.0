import Image from "next/image";
import React from "react";
import ojImage from "../../../public/oj.png";
import { ThemeToggle } from "../theme/ThemeToggle";

const LandingNavbar = () => {
  return (
    <>
      <div className="flex items-center justify-between p-4 border-b-2">
        <Image src={ojImage} alt="Logo" width={50} height={50} />
        <ThemeToggle/>
      </div>
    </>
  );
};

export default LandingNavbar;
