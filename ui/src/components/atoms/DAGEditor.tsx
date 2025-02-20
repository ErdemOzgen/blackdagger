import React from 'react';
import MonacoEditor from 'react-monaco-editor';

type Props = {
  value: string;
  onChange: (value?: string) => void;
};

function DAGEditor({ value, onChange }: Props) {
  return (
    <MonacoEditor
      height="120vh"
      value={value}
      onChange={onChange}
      language="yaml"
      theme="vs-dark"
      options={{
        automaticLayout: true,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        quickSuggestions: { other: true, comments: false, strings: true },
        formatOnType: true,
        formatOnPaste: true,
        renderValidationDecorations: 'on',
        lineNumbers: 'on',
        glyphMargin: true,
      }}
    />
  );
}

export default DAGEditor;
