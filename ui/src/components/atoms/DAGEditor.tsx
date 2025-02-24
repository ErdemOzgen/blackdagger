import React, { useEffect, useRef } from 'react';
import MonacoEditor, { loader } from '@monaco-editor/react';
import * as monaco from 'monaco-editor';
import { configureMonacoYaml } from 'monaco-yaml';

configureMonacoYaml(monaco, {
  enableSchemaRequest: true,
  hover: true,
  completion: true,
  validate: true,
  format: true,
  schemas: [
    {
      uri: 'https://raw.githubusercontent.com/ErdemOzgen/blackdagger/refs/heads/main/schemas/dag.schema.json',
      fileMatch: ['*'], 
    },
  ],
});

loader.config({ monaco });

type Props = {
  value: string;
  onChange: (value?: string) => void;
};

function DAGEditor({ value, onChange }: Props) {
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);

  useEffect(() => {
    return () => {
      editorRef.current?.dispose();
    };
  }, []);

  const editorDidMount = (editor: monaco.editor.IStandaloneCodeEditor) => {
    editorRef.current = editor;

    setTimeout(() => {
      editor.getAction('editor.action.formatDocument')?.run();
    }, 100);
  };
  return (
    <MonacoEditor
      height="120vh"
      value={value}
      onChange={onChange}
      language="yaml"
      theme="vs-dark"
      onMount={editorDidMount}
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
