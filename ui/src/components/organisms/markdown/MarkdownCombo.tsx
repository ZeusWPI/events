import { LoadingOverlay } from "@/components/atoms/LoadingOverlay";
import { useImageCreate } from "@/lib/api/image";
import { useTheme } from "@/lib/hooks/useTheme";
import { getUuid } from "@/lib/utils/utils";
import MDEditor, { commands, MDEditorProps } from "@uiw/react-md-editor";
import { LucideSmilePlus } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import rehypeSanitize from "rehype-sanitize";
import { toast } from "sonner";
import EmojiPicker, { EmojiClickData, Theme } from "emoji-picker-react";

export function MarkdownCombo({ value = "", onChange, ...props }: MDEditorProps) {
  const { theme } = useTheme();

  const [uploading, setUploading] = useState(false);
  const imageCreate = useImageCreate();

  const editorRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  const [emojiOpen, setEmojiOpen] = useState(false);

  const selectionRef = useRef<{ start: number; end: number }>({ start: 0, end: 0 });
  const updateSelection = () => {
    const ta = textareaRef.current;
    if (!ta) return;
    selectionRef.current = {
      start: ta.selectionStart ?? 0,
      end: ta.selectionEnd ?? 0,
    };
  };

  useEffect(() => {
    const root = editorRef.current;
    if (!root) return;

    const ta = root.querySelector<HTMLTextAreaElement>(".w-md-editor-text-input");
    textareaRef.current = ta;

    if (!ta) return;

    const onAny = () => updateSelection();
    ta.addEventListener("keyup", onAny);
    ta.addEventListener("mouseup", onAny);
    ta.addEventListener("select", onAny);
    ta.addEventListener("focus", onAny);

    updateSelection();

    return () => {
      ta.removeEventListener("keyup", onAny);
      ta.removeEventListener("mouseup", onAny);
      ta.removeEventListener("select", onAny);
      ta.removeEventListener("focus", onAny);
    };
  }, []);

  useEffect(() => {
    if (!emojiOpen) return;

    const onDocMouseDown = (e: MouseEvent) => {
      const root = editorRef.current;
      if (!root) return;
      if (!root.contains(e.target as Node)) setEmojiOpen(false);
    };

    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setEmojiOpen(false);
    };

    document.addEventListener("mousedown", onDocMouseDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onDocMouseDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [emojiOpen]);

  useEffect(() => {
    const container = editorRef.current;
    if (!container) return;

    const handlePaste = async (e: ClipboardEvent) => {
      console.log(e);
      // handle pasting text (for auto url)
      const pasteText = e.clipboardData?.getData("text/plain");
      if (pasteText) {
        e.preventDefault();
        console.log(pasteText);
      }

      const allItems = e.clipboardData?.items;
      const items = [...(allItems ?? [])].filter((i) => i.type.startsWith("image/"));
      if (!items.length) return;


      setUploading(true);

      for (const item of items) {
        e.preventDefault();
        const file = item.getAsFile();
        if (!file) continue;

        imageCreate.mutate(
          { name: getUuid(), file },
          {
            onSuccess: (image) => {
              const baseUrl = window.location.origin;
              const inserted = `\n\n![pasted image](${baseUrl}/api/v1/image/${image.id})`;
              onChange?.(`${value}${inserted}`);
            },
            onError: () => toast.error("Failed", { description: "Probably an unsupported file type" }),
            onSettled: () => setUploading(false),
          }
        );
      }
    };

    container.addEventListener("paste", handlePaste);
    return () => container.removeEventListener("paste", handlePaste);
  }, [value, onChange]); // eslint-disable-line react-hooks/exhaustive-deps

  const toggleEmojiSelector = () => setEmojiOpen((v) => !v);

  const handleEmojiClick = (emojiData: EmojiClickData) => {
    const emoji = emojiData.emoji;

    updateSelection();

    const { start, end } = selectionRef.current;

    const nextValue = value.slice(0, start) + emoji + value.slice(end);
    const nextCursor = start + emoji.length;

    onChange?.(nextValue);
    setEmojiOpen(false);

    requestAnimationFrame(() => {
      const ta = textareaRef.current;
      if (!ta) return;
      ta.focus();
      ta.setSelectionRange(nextCursor, nextCursor);
      selectionRef.current = { start: nextCursor, end: nextCursor };
    });
  };

  return (
    <LoadingOverlay loading={uploading}>
      <div ref={editorRef} data-color-mode={theme} className="relative">
        <MDEditor
          value={value}
          height={600}
          onChange={(v) => {
            onChange?.(v);
            requestAnimationFrame(updateSelection);
          }}
          previewOptions={{
            rehypePlugins: [[rehypeSanitize]],
          }}
          extraCommands={[
            commands.group([], {
              name: "Emoji",
              icon: (
                <button
                  type="button"
                  className="mr-1 inline-flex items-center"
                  onMouseDown={(e) => {
                    e.preventDefault();
                    updateSelection();
                  }}
                  onClick={toggleEmojiSelector}
                  aria-label="Insert emoji"
                >
                  <LucideSmilePlus className="w-3 h-3" />
                </button>
              ),
            }),
          ]}
          {...props}
        />

        {emojiOpen && (
          <div className="absolute z-50 top-10 right-2">
            <EmojiPicker
              onEmojiClick={handleEmojiClick}
              autoFocusSearch
              theme={theme === "dark" ? Theme.DARK : Theme.LIGHT}
            />
          </div>
        )}
      </div>
    </LoadingOverlay>
  );
}

