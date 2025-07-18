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
import { toast } from "sonner";

const CodeEditor = ({
  user,
  problem,
}: {
  user: CustomUser;
  problem: ProblemType;
}) => {
  const [language, setLanguage] = useState("c_cpp");
  const [code, setCode] = useState("");
  const handleCodeChange = (newValue: string) => {
    setCode(newValue);
  };
  const handleSubmit = async () => {
    if (!user || !problem) {
      console.error("User or problem data is missing");
      return;
    }
    try {
      if (code === "") {
        toast.info("No code to submit");
        return;
      }
      const data = await submitCode(
        problem.id,
        user.id!,
        user.token!,
        code,
        language
      );
      toast.success("Code submitted successfully!");
      console.log("Code submitted successfully:", data);
    } catch (error) {
      console.error("Error submitting code:", error);
      toast.error("Failed to submit code. Please try again.");
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
        value={code}
        onChange={handleCodeChange}
        setOptions={{
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
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
