import copy from "copy-to-clipboard";
import { ClipboardCheckIcon, ClipboardCopyIcon } from "lucide-react";
import { ComponentProps, useEffect, useState } from "react";
import { toast } from "sonner";
import { IconButton } from "./IconButton";
import { TooltipText } from "./TooltipText";

interface Props extends ComponentProps<"button"> {
  text: string;
  tooltip: string;
}

export function Copy({ text, tooltip, ...props }: Props) {
  const [copied, setCopied] = useState(false)

  useEffect(() => {
    if (!copied) return;

    const timer = setTimeout(() => {
      setCopied(false);
    }, 5000);

    return () => clearTimeout(timer);
  }, [copied]);

  const handleCopy = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation()

    const noPrompt = copy(text)

    if (noPrompt) {
      setCopied(true)
      toast.success("Copied to clipboard")
    }
  }

  return (
    <TooltipText text={tooltip}>
      <IconButton onClick={handleCopy} disabled={copied} {...props}>
        {copied
          ? <ClipboardCheckIcon />
          : <ClipboardCopyIcon />
        }
      </IconButton>
    </TooltipText>
  )
}
