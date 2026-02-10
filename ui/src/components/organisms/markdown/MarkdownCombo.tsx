import { LoadingOverlay } from "@/components/atoms/LoadingOverlay"
import { useImageCreate } from "@/lib/api/image"
import { useTheme } from "@/lib/hooks/useTheme"
import { getUuid } from "@/lib/utils/utils"
import MDEditor, { commands, MDEditorProps } from "@uiw/react-md-editor"
import { LucideSmilePlus } from "lucide-react"
import { useEffect, useMemo, useRef, useState } from "react"
import rehypeSanitize from "rehype-sanitize"
import { toast } from "sonner"
import EmojiPicker, { EmojiClickData, Theme } from "emoji-picker-react"

type Props = MDEditorProps

type Selection = { start: number; end: number }

const isHttpUrl = (text: string) => /^https?:\/\/\S+$/i.test(text.trim())

const canInsertTextWithExecCommand = () => {
  // execCommand is deprecated, but is still supported by most browser and works with the browser's undo stack
  if (typeof document === "undefined") return false
  if (typeof document.execCommand !== "function") return false
  if (typeof document.queryCommandSupported !== "function") return true
  return document.queryCommandSupported("insertText")
}

export function MarkdownCombo({ value = "", onChange, ...props }: Props) {
  const { theme } = useTheme()

  const [uploading, setUploading] = useState(false)
  const [emojiOpen, setEmojiOpen] = useState(false)

  const imageCreate = useImageCreate()

  const editorRef = useRef<HTMLDivElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement | null>(null)

  const selectionRef = useRef<Selection>({ start: 0, end: 0 })

  const updateSelection = () => {
    const ta = textareaRef.current
    if (!ta) return

    selectionRef.current = {
      start: ta.selectionStart ?? 0,
      end: ta.selectionEnd ?? 0,
    }
  }

  const getSelectedText = () => {
    updateSelection()

    const { start, end } = selectionRef.current

    return value.slice(start, end)
  }

  const setSelection = (start: number, end: number) => {
    requestAnimationFrame(() => {
      const ta = textareaRef.current
      if (!ta) return
      ta.focus()
      ta.setSelectionRange(start, end)
      selectionRef.current = { start, end }
    })
  }

  const replaceSelection = (replaceText: string) => {
    updateSelection()

    const { start, end } = selectionRef.current

    const ta = textareaRef.current

    if (!ta) {
      const nextValue = value.slice(0, start) + replaceText + value.slice(end)
      onChange?.(nextValue)
      const nextCursor = start + replaceText.length
      setSelection(nextCursor, nextCursor)

      return
    }

    ta.focus()
    ta.setSelectionRange(start, end)

    const execOk = canInsertTextWithExecCommand() && document.execCommand("insertText", false, replaceText)

    if (execOk) {
      updateSelection()

      return
    }

    ta.setRangeText(replaceText, start, end, "end")
    selectionRef.current = {
      start: start + replaceText.length,
      end: start + replaceText.length,
    }

    // Let react know the input was changed
    try {
      ta.dispatchEvent(
        new InputEvent("input", {
          bubbles: true,
          inputType: "insertText",
          data: replaceText,
        })
      )
    } catch {
      ta.dispatchEvent(new Event("input", { bubbles: true }))
    }
  }

  useEffect(() => {
    const root = editorRef.current
    if (!root) return

    const ta = root.querySelector<HTMLTextAreaElement>(".w-md-editor-text-input")
    textareaRef.current = ta
    if (!ta) return

    const onAny = () => updateSelection()

    ta.addEventListener("keyup", onAny)
    ta.addEventListener("mouseup", onAny)
    ta.addEventListener("select", onAny)
    ta.addEventListener("focus", onAny)

    updateSelection()

    return () => {
      ta.removeEventListener("keyup", onAny)
      ta.removeEventListener("mouseup", onAny)
      ta.removeEventListener("select", onAny)
      ta.removeEventListener("focus", onAny)
    }
  }, [])

  useEffect(() => {
    if (!emojiOpen) return

    const onDocMouseDown = (e: MouseEvent) => {
      const root = editorRef.current
      if (!root) return
      if (!root.contains(e.target as Node)) setEmojiOpen(false)
    }

    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setEmojiOpen(false)
    }

    document.addEventListener("mousedown", onDocMouseDown)
    document.addEventListener("keydown", onKeyDown)

    return () => {
      document.removeEventListener("mousedown", onDocMouseDown)
      document.removeEventListener("keydown", onKeyDown)
    }
  }, [emojiOpen])

  useEffect(() => {
    const container = editorRef.current
    if (!container) return

    const handlePaste = (e: ClipboardEvent) => {
      const pasteText = e.clipboardData?.getData("text/plain") ?? ""
      const selectedText = getSelectedText()

      const imageItems = Array.from(e.clipboardData?.items ?? []).filter((i) => i.type.startsWith("image/"))

      switch (true) {
        case Boolean(selectedText) && isHttpUrl(pasteText): {
          e.preventDefault()
          replaceSelection(`[${selectedText}](${pasteText.trim()})`)

          return
        }

        case imageItems.length > 0: {
          e.preventDefault()
          setUploading(true)

          for (const item of imageItems) {
            const file = item.getAsFile()
            if (!file) continue

            imageCreate.mutate(
              { name: getUuid(), file },
              {
                onSuccess: (image) => {
                  const baseUrl = window.location.origin
                  const alt = getSelectedText() || "pasted image"
                  replaceSelection(`![${alt}](${baseUrl}/api/v1/image/${image.id})`)
                },
                onError: () => toast.error("Failed", { description: "Probably an unsupported file type" }),
                onSettled: () => setUploading(false),
              }
            )
          }

          return
        }

        default:
          return
      }
    }

    container.addEventListener("paste", handlePaste)

    return () => container.removeEventListener("paste", handlePaste)
  }, [value, onChange, imageCreate])

  const emojiTheme = useMemo(() => (theme === "dark" ? Theme.DARK : Theme.LIGHT), [theme])

  const toggleEmojiSelector = () => setEmojiOpen((v) => !v)

  const handleEmojiClick = (emojiData: EmojiClickData) => {
    replaceSelection(emojiData.emoji)
    setEmojiOpen(false)
  }

  return (
    <LoadingOverlay loading={uploading}>
      <div ref={editorRef} data-color-mode={theme} className="relative">
        <MDEditor
          value={value}
          height={600}
          onChange={(v) => {
            onChange?.(v)
            requestAnimationFrame(updateSelection)
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
                    e.preventDefault()
                    updateSelection()
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
            <EmojiPicker onEmojiClick={handleEmojiClick} autoFocusSearch theme={emojiTheme} />
          </div>
        )}
      </div>
    </LoadingOverlay>
  )
}

