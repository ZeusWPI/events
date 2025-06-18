import { useTheme } from "@/lib/hooks/useTheme";
import MDEditor, { MDEditorProps } from "@uiw/react-md-editor";
import { emojify } from "node-emoji";

export function MarkdownViewer({ value }: Pick<MDEditorProps, 'value'>) {
  const { theme } = useTheme();

  return (
    <div data-color-mode={theme}>
      <MDEditor.Markdown
        source={emojify(value ?? "")}
        style={theme === "dark" ? { background: "black" } : undefined}
      />
    </div>
  )
}
