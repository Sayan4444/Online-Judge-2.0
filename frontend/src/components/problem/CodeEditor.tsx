"use client";
import React, { useState } from "react";
import AceEditor from "react-ace";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "../ui/button";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/mode-c_cpp";
import "ace-builds/src-noconflict/mode-python";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";
import { CustomUser } from "@/app/api/auth/[...nextauth]/options";
import { submitCode } from "@/fetch/submission";

const CodeEditor = ({
  user,
  problem,
}: {
  user: CustomUser;
  problem: ProblemType;
}) => {
  const [language, setLanguage] = useState("c_cpp");
  const onChange = (newValue: string) => {
    console.log("change", newValue);
  };
  const [code, setCode] = useState("");
  const handleCodeChange = (newValue: string) => {
    setCode(newValue);
    onChange(newValue);
  };
  const handleSubmit = async () => {
    if (!user || !problem) {
      console.error("User or problem data is missing");
      return;
    }
    try {
      const data = await submitCode(
        problem.id,
        user.id!,
        user.token!,
        code,
        language
      );
      console.log("Code submitted successfully:", data);
    } catch (error) {
      console.error("Error submitting code:", error);
    }
  };

  return (
    <div className="w-[50%] h-auto">
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline">Select Language: {language}</Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          <DropdownMenuItem onClick={() => setLanguage("java")}>
            Java
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setLanguage("c_cpp")}>
            C++
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setLanguage("python")}>
            Python
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      <AceEditor
        placeholder="Write your code here..."
        mode={language}
        theme="monokai"
        name="editor"
        fontSize={16}
        lineHeight={19}
        showPrintMargin={true}
        showGutter={true}
        highlightActiveLine={true}
        value={``}
        onChange={handleCodeChange}
        setOptions={{
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: false,
          enableSnippets: false,
          enableMobileMenu: false,
          showLineNumbers: true,
          tabSize: 2,
        }}
        height="90%"
        width="100%"
      />
      <Button onClick={handleSubmit}>Submit</Button>
    </div>
  );
};

export default CodeEditor;
