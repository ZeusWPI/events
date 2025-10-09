
import { useEffect, useRef, useState } from "react";
import { useTheme } from "@/lib/hooks/useTheme";
import MDEditor, { MDEditorProps } from "@uiw/react-md-editor";
import { emojify } from "node-emoji";
import rehypeSanitize from "rehype-sanitize";
import { useImageCreate } from "@/lib/api/image";
import { toast } from "sonner";
import { getUuid } from "@/lib/utils/utils";
import { LoadingOverlay } from "@/components/atoms/LoadingOverlay";

export function MarkdownCombo({ value = "", onChange, ...props }: MDEditorProps) {
  const { theme } = useTheme();

  const [uploading, setUploading] = useState(false)
  const imageCreate = useImageCreate()

  const editorRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const container = editorRef.current;
    if (!container) return;

    const handlePaste = async (e: ClipboardEvent) => {
      const allItems = e.clipboardData?.items;

      const items = [...allItems ?? []].filter(i => i.type.startsWith("image/"))
      if (!items.length) return

      setUploading(true)

      for (const item of items) {
        e.preventDefault();
        const file = item.getAsFile();
        if (file) {
          imageCreate.mutate({ name: getUuid(), file }, {
            onSuccess: (image) => {
              const baseUrl = window.location.origin;
              const newMarkdown = `${value}\n\n![pasted image](${baseUrl}/api/v1/image/${image.id})`;
              onChange?.(newMarkdown);
            },
            onError: () => toast.error("Failed", { description: "Probably an unsupported file type" }),
            onSettled: () => setUploading(false)
          })
        }
      }
    };

    container.addEventListener("paste", handlePaste);
    return () => container.removeEventListener("paste", handlePaste);
  }, [value, onChange]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <LoadingOverlay loading={uploading}>
      <div ref={editorRef} data-color-mode={theme}>
        <MDEditor
          value={emojify(value)}
          height={600}
          onChange={onChange}
          previewOptions={{
            rehypePlugins: [[rehypeSanitize]],
          }}
          {...props}
        />
      </div>
    </LoadingOverlay>
  );
}

