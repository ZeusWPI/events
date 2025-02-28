interface Props {
  text: string;
}

export function DividerText({ text }: Props) {
  return (
    <div className=" flex py-5 items-center">
      <div className="grow border-t" />
      <span className="shrink mx-4 text-muted-foreground">{text}</span>
      <div className="grow border-t" />
    </div>
  );
}
