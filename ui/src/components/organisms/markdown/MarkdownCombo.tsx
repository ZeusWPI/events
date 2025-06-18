
import { useTheme } from '@/lib/hooks/useTheme';
import MDEditor, { MDEditorProps } from '@uiw/react-md-editor';
import { emojify } from 'node-emoji';
import rehypeSanitize from "rehype-sanitize";

export function MarkdownCombo({ value, ...props }: MDEditorProps) {
  const { theme } = useTheme();

  return (
    <div data-color-mode={theme}>
      <MDEditor
        value={emojify(value ?? "")}
        height={600}
        previewOptions={{
          rehypePlugins: [[rehypeSanitize]],
        }}
        {...props}
      />
    </div>
  )
}
