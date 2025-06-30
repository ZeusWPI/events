
import { ComponentProps, useEffect, useState } from "react";
import { Indeterminate } from "./Indeterminate";

interface Props extends ComponentProps<'img'> {
  file?: File;
  isLoading: boolean;
}

export function FileImg({ file, isLoading, ...props }: Props) {
  const [url, setUrl] = useState<string | null>(null);

  useEffect(() => {
    if (!file) {
      setUrl(null)
      return
    }

    if (file) {
      const objectUrl = URL.createObjectURL(file);
      setUrl(objectUrl);

      return () => {
        URL.revokeObjectURL(objectUrl);
      };
    }
  }, [file]);

  if (isLoading) {
    return <Indeterminate />;
  }

  if (!file || !url) {
    return null;
  }

  return <img src={url} {...props} />;
}

