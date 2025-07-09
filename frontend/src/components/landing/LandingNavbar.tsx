import Image from "next/image";
import React from "react";
import ojImage from "../../../public/oj.png";

const LandingNavbar = () => {
  return (
    <>
      <div className="flex items-center justify-between bg-gray-800 p-4">
        <Image src={ojImage} alt="Logo" width={50} height={50} />
      </div>
    </>
  );
};

export default LandingNavbar;
